package scripts

import (
	"os"
)

func getCacheDir() string {
	var odxCacheDir string
	if xdgCacheHome, ok := os.LookupEnv("XDG_CACHE_HOME"); ok {
		odxCacheDir = xdgCacheHome
	} else {
		// TODO: do we care about this error?
		userCacheDir, _ := os.UserCacheDir()
		odxCacheDir = userCacheDir
		return odxCacheDir
	}
	odxCacheDir = odxCacheDir + string(os.PathSeparator) + "odx"
	if err := os.MkdirAll(odxCacheDir, 0744); err != nil && !os.IsExist(err) {
		return ""
	}
	return odxCacheDir
}
