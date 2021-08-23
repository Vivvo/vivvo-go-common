package mesh

func internal(path string) string {
	if os.Getenv("USE_HTTPS") == "true" {
		return "https://" + path
	} else {
		return "http://" + path
	}
}
