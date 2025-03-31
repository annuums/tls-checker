package util

import "testing"

func TestIsValidRgb(t *testing.T) {

  tests := map[string]bool{
    "#FFF":    true,
    "#ffffff": true,
    "#abc123": true,
    "123456":  false,
    "#ggg":    false,
    "#12345g": false,
  }

  for input, expected := range tests {

    result := IsValidRgb(input)
    if result != expected {

      t.Errorf("IsValidRgb(%q) = %v; want %v", input, result, expected)
    }
  }
}
