package utils

import "os"

// JudgeFileExist judge file or dir exist
func JudgeFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

type Element interface {
	~int | ~int32 | ~int64 | ~string
}

func ElementInSlice[E Element](elements []E, target E) bool {
	for _, element := range elements {
		if element == target {
			return true
		}
	}
	return false
}
