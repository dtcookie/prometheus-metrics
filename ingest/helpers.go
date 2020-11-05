package main

import (
	"flag"
	"os"
	"strings"
	"unicode"
)

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func pempty(s *string) bool {
	if s == nil {
		return true
	}
	return empty(*s)
}
func uscore(s string) string {
	var result string
	for _, r := range s {
		if unicode.IsSpace(r) {
			result = result + "_"
		} else {
			result = result + string(r)
		}
	}
	return result
}

func quote(s string) string {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return s
	}
	return "\"" + s + "\""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func setFlag(source *string, target *string) {
	if pempty(source) {
		return
	}
	*target = *source
}

func readFlag(key string, def string, desc string, target *string) {
	ptr := flag.String(key, def, desc)
	if ptr == nil {
		return
	}
	s := *ptr
	if len(strings.TrimSpace(s)) == 0 {
		return
	}
	*target = s
}
