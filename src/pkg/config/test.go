package config

import (
	"insomnia/src/pkg/utils"
	"os"
	"path"
	"strings"
)

func cfgInit() {
	var isTest bool
	for _, v := range os.Args {
		if strings.HasPrefix(v, "-test.") {
			isTest = true
			break
		}
	}
	if !isTest && strings.Contains(os.Args[0], "___Test") || strings.Contains(os.Args[0], ".test") {
		isTest = true
	}
	if isTest || !InDocker {
		// locate app.yaml, set go test directory to the directory where app.yaml is located
		pwd, _ := os.Getwd()
		parentDir := pwd
		var appDir string
		for {
			if utils.JudgeFileExist(path.Join(parentDir, "app.yaml")) {
				appDir = parentDir
				break
			}
			temp := path.Dir(parentDir)
			if parentDir == temp {
				break
			}
			parentDir = temp
		}
		if appDir != "" {
			_ = os.Chdir(appDir)
		}
	}
}
