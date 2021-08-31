package mesh

import "os"

func Internal(path string) string {
	if os.Getenv("USE_HTTPS") == "false" {
		return "http://" + path
	} else {
		return "https://" + path
	}
}
