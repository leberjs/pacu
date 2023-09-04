package internal

import (
	"strings"
)

func assignCred(s string) string {
	ss := strings.Split(s, "=")
	return strings.TrimSpace(ss[1])
}
