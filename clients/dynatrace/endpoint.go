package dynatrace

import (
	"fmt"
	"strings"
)

type endpoint interface {
	Format(string) string
}

type saasEndpoint struct {
	EnvironmentID string
}

func (se *saasEndpoint) Format(path string) string {
	return fmt.Sprintf("https://%s.live.dynatrace.com%s", se.EnvironmentID, path)
}

type managedEndpoint struct {
	BaseURL       string
	EnvironmentID string
}

func (me *managedEndpoint) Format(path string) string {
	base := me.BaseURL
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}
	return fmt.Sprintf("%se/%s", me.EnvironmentID, path)
}
