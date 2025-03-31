package tlscheck

import "testing"

func TestCheck(t *testing.T) {

  checker := NewChecker(365)
  vc := checker.Check("google.com")

  if vc == nil {

    t.Fatal("Expected a valid cert, got nil")
  }

  if !vc.ShouldAlert {

    t.Logf("No alert needed, expires in more than threshold")
  } else {

    t.Logf("ALERT: expires within threshold")
  }
}
