package usecase

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func BuildCache() error {
	cacheTime := time.Now()
	// check whether files will changed
	fileInfo, err := os.Stat("path")
	if err != nil {
		log.Fatal(err)
	}

	modifyTime := fileInfo.ModTime()

	// if not changed, return
	if modifyTime.Before(cacheTime) {
		return nil
	}
	// build cache

	// record the cache time
	return nil
}

func QueryByTag(tag string) error {
	if tagFiles, ok, _ := readCache(); ok {
		if files, ok := tagFiles[tag]; !ok {
			return fmt.Errorf("tag %s not found", tag)
		} else {
			fmt.Println(files)
			return nil
		}
	}

	rootPath := "./"

	filesPath := []string{}

	// Walk through all files and append path which contains '.go'
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

		// if not set tag, continue to next file
		if len(lines) < 3 {
			continue
		}

		thirdLine := strings.TrimSpace(lines[2])
		if strings.HasPrefix(thirdLine, "/* @tag:") && strings.HasSuffix(thirdLine, "*/") {
			tagsLine := strings.TrimPrefix(thirdLine, "/* @tag:")
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

	fmt.Println(tagFiles["abc"])

	return nil
}

func QueryByExpression(expression string) ([]string, error) {
	return []string{}, errors.New("not implemented")
}

func ListTags() error {
	return errors.New("not implemented")
}

func readCache() (map[string][]string, bool, error) {
	return map[string][]string{}, false, errors.New("not implemented")
}
