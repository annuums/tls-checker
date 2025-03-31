package config

import (
  "github.com/anuums/tls-checker/internal/util"
  "os"
  "strconv"
  "strings"
)

var RuntimeConf = Config{}

type SlackConfig struct {
  Enabled   bool `default:"false"`
  Token     string
  ChannelID string
  ColorBar  string
}

type EmailConfig struct {
  Enabled  bool `default:"false"`
  To       string
  Subject  string
  SMTPHost string
  SMTPPort string
  SMTPUser string
  SMTPPass string
}

type Config struct {
  Hostnames   []string
  CheckPeriod int
  Timezone    string
  SlackConfig
  EmailConfig
}

func init() {
  RuntimeConf = LoadConfigFromEnv()
}

func LoadConfigFromEnv() Config {

  checkPeriod := 14
  if str := os.Getenv("TLS_EXPIRATION_CHECK_PERIOD"); str != "" {

    if i, err := strconv.Atoi(str); err == nil {
      checkPeriod = i
    }
  }

  tz := os.Getenv("TIMEZONE")
  if tz == "" {

    tz = "Etc/UTC"
  }

  color := os.Getenv("SLACK_MESSAGE_COLOR_BAR")
  if !util.IsValidRgb(color) {

    color = ""
  }

  config := Config{
    Timezone:    tz,
    CheckPeriod: checkPeriod,
    Hostnames:   strings.Split(os.Getenv("HOSTNAMES"), ","),
    // Slack Configuration
    SlackConfig: SlackConfig{
      Token:     os.Getenv("SLACK_TOKEN"),
      ChannelID: os.Getenv("SLACK_CHANNEL_ID"),
      ColorBar:  color,
    },
    // Email SMTP Configuration
    EmailConfig: EmailConfig{
      To:       os.Getenv("EMAIL_TO"),
      Subject:  os.Getenv("EMAIL_SUBJECT"),
      SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
      SMTPPort: os.Getenv("EMAIL_SMTP_PORT"),
      SMTPUser: os.Getenv("EMAIL_SMTP_USER"),
      SMTPPass: os.Getenv("EMAIL_SMTP_PASS"),
    },
  }

  if config.EmailConfig.Subject == "" {

    config.EmailConfig.Subject = "[TLS Checker] TLS Certificate is About to Expire Soon"
  }

  if config.SlackConfig.Token != "" && config.SlackConfig.ChannelID != "" {

    config.SlackConfig.Enabled = true
  }

  if config.EmailConfig.To != "" &&
    config.EmailConfig.Subject != "" && config.EmailConfig.SMTPHost != "" &&
    config.EmailConfig.SMTPPort != "" && config.EmailConfig.SMTPUser != "" && config.EmailConfig.SMTPPass != "" {

    config.EmailConfig.Enabled = true
  }

  return config
}
