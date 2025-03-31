package util

import "regexp"

func IsValidRgb(color string) bool {

  regex := regexp.MustCompile(`^#([A-Fa-f0-9]{3,6})$`)
  return regex.MatchString(color)
}
