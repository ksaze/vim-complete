package stale

import "os"

func IsStale(out string, srcs ...string) bool {
	outInfo, err := os.Stat(out)
	if err != nil {
		return true
	}
	for _, s := range srcs {
		srcInfo, err := os.Stat(s)
		if err != nil {
			return true
		}
		if srcInfo.ModTime().After(outInfo.ModTime()) {
			return true
		}
	}
	return false
}
