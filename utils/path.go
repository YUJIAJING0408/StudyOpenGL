package utils

import "os"

func Relative2FullPath(relativePath string) (fullPath string, err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return workDir + "\\" + relativePath, nil
}
