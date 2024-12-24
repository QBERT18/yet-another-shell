package yetcommand

import (
	"os"
	"path/filepath"
	"regexp"
)

func SearchProgramInPath(program string, paths []string) string {
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

func CustomSplit(input string) []string {
	re := regexp.MustCompile(`exit 0|[^\s]+`)
	return re.FindAllString(input, -1)
}
