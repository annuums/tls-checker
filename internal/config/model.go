package config

import "github.com/anuums/tls-checker/internal/tlscheck"

type groupedAlert struct {
	Scope     string
	Hostnames []string
	Cert      *tlscheck.ValidateCert
}
