package config

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Config struct {
	GCS           string
	GCSPublicBase string
	gcsBucket     string
	gcsBasePath   string
	gcsBaseURL    string
}

func (c *Config) Build() error {
	gcsURL, err := url.Parse(c.GCS)
	if err != nil {
		return errors.Wrapf(err, "failed to parse gcsURL: %v", gcsURL)
	}
	if gcsURL.Scheme != "gs" || gcsURL.Host == "" {
		return errors.Errorf("invalid gcsURL: %v", gcsURL)
	}
	if c.GCSPublicBase == "" {
		c.GCSPublicBase = "https://storage.googleapis.com"
	}
	c.gcsBucket = gcsURL.Host
	c.gcsBasePath = strings.TrimPrefix(gcsURL.Path, "/")

	gcsBasePath, err := url.JoinPath(c.GCSPublicBase, c.gcsBucket, c.gcsBasePath)
	if err != nil {
		return errors.Wrapf(err, "failed to join paths: %v %v %v", c.GCSPublicBase, c.gcsBucket, c.gcsBasePath)
	}
	c.gcsBaseURL = gcsBasePath
	return nil
}

func (c *Config) GCSBucketName() string {
	return c.gcsBucket
}

func (c *Config) GCSBasePath() string {
	return c.gcsBasePath
}

func (c *Config) GCSBaseURL() string {
	return c.gcsBaseURL
}
