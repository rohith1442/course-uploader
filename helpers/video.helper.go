package helpers

func IsValidVideoExtension(extensionType string, extensions []string) bool {
	for _, validExt := range extensions {
		if validExt == extensionType {
			return true
		}
	}
	return false
}
