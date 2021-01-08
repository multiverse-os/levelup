package validations

import (
	"regexp"
)

const RestrictedNameChars = `[a-zA-Z0-9][a-zA-Z0-9_.-]`

var RestrictedNamePattern = regexp.MustCompile(`^` + RestrictedNameChars + `+$`)
