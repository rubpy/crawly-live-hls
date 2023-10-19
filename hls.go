package hls

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

//////////////////////////////////////////////////

type M3UPlaylist struct {
	HasStartTag   bool
	HasVersionTag bool

	HasError               bool
	HasIndependentSegments bool
	HasStreamInf           bool
	HasIFrameStreamInf     bool
	HasMedia               bool
}

func decodeM3ULine(playlist *M3UPlaylist, line string) (err error) {
	if playlist == nil {
		return errors.New("playlist cannot be a nil pointer")
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	switch {
	case line == "#EXTM3U":
		playlist.HasStartTag = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"):
		playlist.HasVersionTag = true

	case line == "#EXT-X-INDEPENDENT-SEGMENTS":
		playlist.HasIndependentSegments = true

	case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
		playlist.HasMedia = true
	case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
		playlist.HasStreamInf = true
	case strings.HasPrefix(line, "#EXT-X-I-FRAME-STREAM-INF:"):
		playlist.HasIFrameStreamInf = true
	}

	return
}

func decodeM3U(reader io.Reader) (playlist M3UPlaylist, err error) {
	if reader == nil {
		err = io.ErrUnexpectedEOF
		return
	}

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return
	}

	var eof bool
	for !eof {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}

		if len(line) < 1 || line == "\r" || line == "\n" {
			continue
		}

		err = decodeM3ULine(&playlist, line)
		if err != nil {
			break
		}
	}

	return
}

//////////////////////////////////////////////////

var (
	InvalidStreamURL = errors.New("invalid stream URL")

	HLSUnexpectedResponse = errors.New("unexpected HLS response")
)

func IsValidStreamURL(streamURL string) (valid bool) {
	if streamURL == "" {
		return false
	}

	u, err := url.Parse(streamURL)
	if err != nil || u.Scheme == "" {
		return false
	}

	return true
}

func (cr *Crawler) CheckHLS(ctx context.Context, streamURL string) (available bool, err error) {
	if streamURL == "" || !IsValidStreamURL(streamURL) {
		err = InvalidStreamURL
		return
	}

	if cr.client == nil {
		err = NilClient
		return
	}

	if ctx == nil {
		ctx = context.Background()
	} else {
		if err = ctx.Err(); err != nil {
			return
		}
	}

	resp, err := cr.client.Request(ctx, "GET", streamURL, nil, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		playlist, err := decodeM3U(resp.Body)
		if err != nil {
			return false, fmt.Errorf("decodeM3U: %w", err)
		}

		if playlist.HasStartTag && !playlist.HasError && (playlist.HasMedia || playlist.HasStreamInf || playlist.HasIFrameStreamInf) {
			available = true
		}
	}

	return
}
