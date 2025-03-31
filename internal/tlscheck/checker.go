package tlscheck

import (
  "crypto/tls"
  "crypto/x509"
  "log"
  "strings"
  "time"
)

type ValidateCert struct {
  Hostname    string
  ShouldAlert bool
  Cert        *x509.Certificate
}

type Checker struct {
  threshold float64
}

func NewChecker(days int) *Checker {

  return &Checker{threshold: float64(days)}
}

func (c *Checker) Check(hostname string) *ValidateCert {

  if !strings.Contains(hostname, ":") {

    hostname += ":443"
  }

  conn, err := tls.Dial("tcp", hostname, &tls.Config{})
  if err != nil {

    log.Printf("[tls] connection failed: %v", err)
    return nil
  }
  defer conn.Close()

  certs := conn.ConnectionState().PeerCertificates
  now := time.Now()
  for _, cert := range certs {

    if now.Before(cert.NotAfter) {

      daysLeft := cert.NotAfter.Sub(now).Hours() / 24
      log.Printf("[tls] %s is valid until %s", hostname, cert.NotAfter.Format(time.RFC1123))

      if daysLeft < c.threshold {

        return &ValidateCert{
          Hostname:    hostname,
          ShouldAlert: true,
          Cert:        cert,
        }
      }
      break
    }
  }

  return nil
}
