package hls

import (
	"github.com/rubpy/crawly"
	"github.com/rubpy/crawly/cclient"
	"github.com/rubpy/crawly/csync"
)

//////////////////////////////////////////////////

type Crawler struct {
	crawly.Crawler

	client cclient.Client

	settings csync.Value[CrawlerSettings]
}

func NewCrawler(opts ...ConfigOption) (*Crawler, error) {
	var cfg config

	for _, opt := range opts {
		opt(&cfg)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	cr, err := buildCrawlerFromConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return cr, nil
}
