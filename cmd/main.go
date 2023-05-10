package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootPath := "./"

	filesPath := []string{}

	// Walk through all files and append path which contains .go
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			filesPath = append(filesPath, path)
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(filesPath)

	// Add file name to tag key
	tagFiles := map[string][]string{}
	for _, path := range filesPath {
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		fileContent := string(content)
		lines := strings.Split(fileContent, "\n")
		if len(lines) >= 2 {
			secondLine := strings.TrimSpace(lines[1])
			if strings.HasPrefix(secondLine, "/* $tag:") && strings.HasSuffix(secondLine, "*/") {
				tagsLine := strings.TrimPrefix(secondLine, "/* $tag:")
				tagsLine = strings.TrimSuffix(tagsLine, "*/")
				tagsLine = strings.TrimSpace(tagsLine)
				tags := strings.Split(tagsLine, ",")
				for _, tag := range tags {
					if _, ok := tagFiles[tag]; !ok {
						tagFiles[tag] = []string{}
					}
					tagFiles[tag] = append(tagFiles[tag], path)
				}
			}
		}
	}

	fmt.Println(tagFiles)

}
