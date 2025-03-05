package yetcommand

import (
	"os"
	"path/filepath"
	"strings"
)

func SearchProgramInPath(program string) string {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, dir := range paths {
		fullPath := filepath.Join(dir, program)
		if FileExistsAndExecutable(fullPath) {
			return fullPath
		}
	}
	return ""
}

func FileExistsAndExecutable(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir() && (info.Mode()&0111 != 0)
}
