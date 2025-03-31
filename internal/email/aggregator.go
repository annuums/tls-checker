package email

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/anuums/tls-checker/internal/tlscheck"
)

type groupedAlert struct {
	Scope         []string
	Hostnames     []string
	Expiry        time.Time
	RemainingDays int
}

func groupByDNSScope(certs []*tlscheck.ValidateCert) map[string]*groupedAlert {

	groups := make(map[string]*groupedAlert)

	for _, vc := range certs {

		if vc == nil || vc.Cert == nil {

			continue
		}

		scope := append([]string{}, vc.Cert.DNSNames...)
		sort.Strings(scope)
		key := strings.Join(scope, ",")

		if group, exists := groups[key]; exists {

			if !contains(group.Hostnames, vc.Hostname) {

				group.Hostnames = append(group.Hostnames, vc.Hostname)
			}
		} else {

			now := time.Now()
			remaining := int(vc.Cert.NotAfter.Sub(now).Hours() / 24)

			groups[key] = &groupedAlert{
				Scope:         scope,
				Hostnames:     []string{vc.Hostname},
				Expiry:        vc.Cert.NotAfter,
				RemainingDays: remaining,
			}
		}
	}

	return groups
}

func contains(slice []string, item string) bool {

	for _, s := range slice {

		if s == item {

			return true
		}
	}

	return false
}

func SendGroupedAlerts(certificates []*tlscheck.ValidateCert, ec *EmailClient) error {

	log.Println("Sending Email Alert...")

	groups := groupByDNSScope(certificates)

	log.Printf("  group by dns sope :: %v\n", groups)

	for _, group := range groups {

		// 이메일 제목에 기본 제목과 스코프 표시
		subject := fmt.Sprintf("%s [%s]", ec.Subject, strings.Join(group.Scope, " "))

		log.Printf("  subject: %v\n", subject)

		// 이메일 본문 구성
		body := "The following hostnames share the same TLS certificate scope and are nearing expiration:\n\n"
		body += fmt.Sprintf("Certificate Scope: %s\n", strings.Join(group.Scope, ", "))
		body += fmt.Sprintf("Hostnames: %s\n", strings.Join(group.Hostnames, ", "))
		body += fmt.Sprintf("Expires At: %s\n", group.Expiry.Format(time.RFC1123))
		body += fmt.Sprintf("Remaining Days: %d\n", group.RemainingDays)
		body += "\nPlease take appropriate action."

		log.Printf("  body :: %v\n", body)

		// 이메일 전송
		if err := ec.SendEmail(subject, body); err != nil {

			return fmt.Errorf("failed to send email for group %s: %w", strings.Join(group.Scope, ","), err)
		}

		log.Printf("  mail sent for :: %v\n", group)
	}

	return nil
}
