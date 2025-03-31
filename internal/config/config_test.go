package config

import (
  "os"
  "testing"
)

func TestLoadConfigFromEnv(t *testing.T) {

  os.Setenv("TLS_EXPIRATION_CHECK_PERIOD", "10")
  os.Setenv("TIMEZONE", "Asia/Seoul")
  os.Setenv("HOSTNAMES", "google.com,example.com")
  os.Setenv("SLACK_TOKEN", "fake-token")
  os.Setenv("SLACK_CHANNEL_ID", "fake-channel")
  os.Setenv("SLACK_MESSAGE_COLOR_BAR", "#123456")

  cfg := LoadConfigFromEnv()

  if cfg.CheckPeriod != 10 {

    t.Errorf("expected CheckPeriod 10, got %d", cfg.CheckPeriod)
  }

  if len(cfg.Hostnames) != 2 {

    t.Errorf("expected 2 hostnames, got %d", len(cfg.Hostnames))
  }

  if cfg.SlackConfig.ColorBar != "#123456" {

    t.Errorf("invalid color bar: %s", cfg.SlackConfig.ColorBar)
  }
}
