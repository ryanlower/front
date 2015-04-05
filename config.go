package main

import (
	"github.com/ryanlower/setting"
)

type config struct {
	Port                string `env:"PORT"`
	AllowedContentTypes string `env:"ALLOWED_CONTENT_TYPE_REGEX" default:"^image/"` // uncompiled regex
	AWS                 struct {
		AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
		Region          string `env:"AWS_REGION" default:"us-east-1"`
	}
	S3 struct {
		Bucket string `env:"S3_BUCKET"`
	}
}

// Load config from environment
func (c *config) load() {
	setting.Load(c)
}
