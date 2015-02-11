package main

import (
	"path/filepath"
	"strings"

	"github.com/zerklabs/auburn/utils"
)

func getOssecConfPath(p string) string {
	if strings.Contains(strings.ToLower(p), "ossec.conf") {
		p = filepath.Dir(p)
	}

	if p != "" && !utils.StrInArray(p, paths) {
		newPaths := make([]string, 0)
		newPaths = append(newPaths, p)
		for _, v := range paths {
			newPaths = append(newPaths, v)
		}

		paths = newPaths
	}

	for _, v := range paths {
		if utils.FileExists(filepath.Join(v, "ossec.conf")) {
			return filepath.Join(v, "ossec.conf")
		}
	}

	return ""
}

func getClientKeysPath(p string) string {
	if strings.Contains(strings.ToLower(p), "client.keys") {
		p = filepath.Dir(p)
	}

	if p != "" && !utils.StrInArray(p, paths) {
		newPaths := make([]string, 0)
		newPaths = append(newPaths, p)
		for _, v := range paths {
			newPaths = append(newPaths, v)
		}

		paths = newPaths
	}

	for _, v := range paths {
		if utils.FileExists(filepath.Join(v, "client.keys")) {
			return filepath.Join(v, "client.keys")
		}
	}

	return ""
}
