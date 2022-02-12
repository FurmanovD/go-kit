package filesys

import (
	"os"
	gopath "path"
	"runtime"
)

func GetUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func GetFullPath(path string) string {

	// empty path
	if path == "" || path == "." {
		wd, _ := os.Getwd()
		return wd
	}

	// full path
	if runtime.GOOS == "windows" {
		// "D:..." case:
		if len(path) >= 2 && path[2] == ':' {
			if (path[0] >= 'a' && path[0] <= 'z') ||
				(path[0] >= 'A' && path[0] <= 'Z') {
				return path
			}
		}
	} else if path[0] == '/' || path[0] == '\\' {
		// "/..." case
		return path
	}

	// relative to a home directory:
	if path[0] == '~' {
		return gopath.Join(GetUserHomeDir(), path[1:])
	}

	// relative
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return gopath.Join(wd, path)
}
