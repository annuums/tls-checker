package email

import (
  "fmt"
  "strings"
  "time"

  mail "github.com/xhit/go-simple-mail/v2"
)

// EmailClient holds SMTP configuration for sending emails.
type EmailClient struct {
  Enabled  bool
  To       string
  Subject  string
  SMTPHost string
  SMTPPort int
  SMTPUser string
  SMTPPass string
}

// NewEmailClient creates a new EmailClient with provided settings.
func NewEmailClient(enabled bool, to, subject, host string, port int, user, pass string) *EmailClient {

  return &EmailClient{
    Enabled:  enabled,
    To:       to,
    Subject:  subject,
    SMTPHost: host,
    SMTPPort: port,
    SMTPUser: user,
    SMTPPass: pass,
  }
}

// SendEmail sends an email with the specified subject and body.
func (ec *EmailClient) SendEmail(subject, body string) error {
  if !ec.Enabled {
    return nil
  }

  // SMTP 클라이언트 설정
  server := mail.NewSMTPClient()
  server.Host = ec.SMTPHost
  server.Port = ec.SMTPPort
  server.Username = ec.SMTPUser
  server.Password = ec.SMTPPass
  server.Encryption = mail.EncryptionSTARTTLS
  server.KeepAlive = false
  server.ConnectTimeout = 10 * time.Second
  server.SendTimeout = 10 * time.Second

  // SMTP 서버 연결
  smtpClient, err := server.Connect()
  if err != nil {
    return fmt.Errorf("failed to connect to SMTP server: %w", err)
  }

  // 이메일 메시지 구성
  msg := mail.NewMSG()
  msg.SetFrom(ec.SMTPUser).
    AddTo(strings.Split(ec.To, ",")...).
    SetSubject(subject)
  msg.SetBody(mail.TextPlain, body)

  // 이메일 전송
  if err = msg.Send(smtpClient); err != nil {
    return fmt.Errorf("failed to send email: %w", err)
  }

  return nil
}
