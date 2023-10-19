package hls

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rubpy/crawly"
	"github.com/rubpy/crawly/clog"
)

//////////////////////////////////////////////////

type EntityData struct {
	Live bool `json:"live"`
}

func (cr *Crawler) entityHandler(ctx context.Context, entity *crawly.Entity, result *crawly.TrackingResult) error {
	handle, ok := entity.Handle.(Handle)
	if !ok || !handle.Valid() {
		return crawly.InvalidHandle
	}

	data, _ := entity.Data.(EntityData)
	defer func() {
		entity.Data = data
	}()

	if handle.Type != HandleStreamURL {
		return crawly.InvalidHandle
	}

	{
		lp := clog.Params{
			Message: "checkHLS",
			Level:   slog.LevelDebug,

			Values: clog.ParamGroup{
				"streamURL": handle.Value,
			},
		}

		available, err := cr.CheckHLS(ctx, handle.Value)
		if err == nil {
			data.Live = available

			lp.Set("available", available)
		} else {
			err = fmt.Errorf("CheckHLS: %w", err)
		}

		lp.Err = err
		cr.Log(ctx, lp)

		if err != nil {
			return err
		}
	}

	return nil
}
