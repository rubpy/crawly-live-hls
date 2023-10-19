package hls

import (
	"strconv"
	"strings"

	"github.com/rubpy/crawly"
)

//////////////////////////////////////////////////

type Handle struct {
	Type  HandleType
	Value string
}

func (h Handle) Valid() bool {
	return h.Type != 0 && h.Value != ""
}

func (h Handle) Equal(handle crawly.Handle) bool {
	if hh, ok := handle.(Handle); ok {
		return hh.Type == h.Type && hh.Value == h.Value
	}

	return false
}

func (h Handle) String() string {
	var s strings.Builder
	s.WriteRune('{')
	s.WriteString(h.Type.String())
	s.WriteString(":")
	s.WriteString(strconv.Quote(h.Value))
	s.WriteRune('}')

	return s.String()
}

type HandleType uint

const (
	HandleStreamURL HandleType = (iota + 1)
)

func (ht HandleType) String() string {
	switch ht {
	case HandleStreamURL:
		return "StreamURL"
	}

	return ""
}

//////////////////////////////////////////////////

func StreamURL(streamURL string) Handle {
	return Handle{HandleStreamURL, streamURL}
}
