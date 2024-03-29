package mesh

import (
	"os"
	"strings"
)

func Internal(path string) string {
	if strings.HasPrefix(path, "http") {
		if strings.HasPrefix(path, "https") {
			path = strings.TrimPrefix(path, "https://")
		} else {
			path = strings.TrimPrefix(path, "http://")
		}
	}
	if os.Getenv("USE_HTTPS") == "false" {
		return "http://" + path
	} else {
		return "https://" + path
	}
}
