package main

import (
  "github.com/anuums/tls-checker/internal/config"
  "github.com/anuums/tls-checker/internal/email"
  "github.com/anuums/tls-checker/internal/tlscheck"
  "log"
  "strconv"

  "github.com/anuums/tls-checker/internal/slack"
)

func sendSlack() {

}

func main() {

  log.Println("Staring TLS Checker...")
  log.Printf("  Check period :: %v\n", config.RuntimeConf.CheckPeriod)
  log.Printf("  Hostnaems:: %v\n", config.RuntimeConf.Hostnames)

  slackClient := slack.NewSlackClient(config.RuntimeConf.SlackConfig.Token)
  tlsChecker := tlscheck.NewChecker(config.RuntimeConf.CheckPeriod)

  validateCerts := make([]*tlscheck.ValidateCert, 0)
  for _, hostname := range config.RuntimeConf.Hostnames {

    vc := tlsChecker.Check(hostname)
    if vc != nil {

      validateCerts = append(validateCerts, vc)
    }
  }

  for _, vc := range validateCerts {

    if vc.ShouldAlert {

      if config.RuntimeConf.SlackConfig.Enabled {

        err := slackClient.SendAlert(vc)
        if err != nil {

          log.Printf("failed to send notification :: %v\n", err)
        }
      }
    }
  }

  if config.RuntimeConf.EmailConfig.Enabled {

    port, err := strconv.Atoi(config.RuntimeConf.EmailConfig.SMTPPort)
    if err != nil {

      log.Fatalf("failed to parse SMTP port: %v\n", err)
    }

    emailClient := email.NewEmailClient(
      config.RuntimeConf.EmailConfig.Enabled,
      config.RuntimeConf.EmailConfig.To,
      config.RuntimeConf.EmailConfig.Subject,
      config.RuntimeConf.EmailConfig.SMTPHost,
      port,
      config.RuntimeConf.EmailConfig.SMTPUser,
      config.RuntimeConf.EmailConfig.SMTPPass,
    )

    if err := email.SendGroupedAlerts(validateCerts, emailClient); err != nil {

      log.Printf("failed to send email alerts :: %v\n", err)
    }
  }
}
