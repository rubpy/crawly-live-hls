package hls

import (
	"github.com/rubpy/crawly"
)

//////////////////////////////////////////////////

type CrawlerSettings struct {
	crawly.CrawlerSettings
}

var DefaultSettings = CrawlerSettings{
	CrawlerSettings: crawly.DefaultCrawlerSettings,
}

//////////////////////////////////////////////////

func (cr *Crawler) loadSettings() CrawlerSettings {
	return cr.settings.Load()
}

func (cr *Crawler) setSettings(settings CrawlerSettings) {
	cr.settings.Store(settings)
	crawly.SetCrawlerSettings(&cr.Crawler, settings.CrawlerSettings)
}

func (cr *Crawler) Settings() CrawlerSettings {
	return cr.loadSettings()
}

func (cr *Crawler) SetSettings(settings CrawlerSettings) {
	cr.setSettings(settings)
}
