package converter

func IsValidImageEncoding(encoding string) bool {
	if encoding == "jpg" || encoding == "png" {
		return true
	}
	return false
}
