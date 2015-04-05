package main

import (
	"github.com/ryanlower/setting"
)

type config struct {
	Port                string `env:"PORT"`
	AllowedContentTypes string `env:"ALLOWED_CONTENT_TYPE_REGEX" default:"^image/"` // uncompiled regex
}

// Load config from environment
func (c *config) load() {
	setting.Load(c)
}
