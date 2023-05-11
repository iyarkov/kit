package support

import "os"

func SkipIntegrations() bool {
	for _, arg := range os.Args {
		if arg == "integration" {
			return false
		}
	}
	return true
}
