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

func QueryByTag(path, tag string) error {
	if tagFiles, ok, _ := readCache(); ok {
		if files, ok := tagFiles[tag]; !ok {
			return fmt.Errorf("tag %s not found", tag)
		} else {
			fmt.Println(files)
			return nil
		}
	}

	// Walk through all files and append path which contains '.go'
	filesPath, err := getFilesPath(path)
	if err != nil {
		return err
	}

	// find the file which contains target tag
	files := []string{}
	for _, path := range filesPath {
		fileTags, hasTags, err := extractTags(path)
		if err != nil {
			return err
		}
		if hasTags {
			for _, fileTag := range fileTags {
				if fileTag == tag {
					files = append(files, path)
				}
			}
		}
	}

	fmt.Println(files)

	return nil
}

func QueryByExpression(path, expression string) error {
	// '', -
	// 'divide and conquer' tree
	expression = strings.TrimPrefix(expression, " ")
	expression = strings.TrimSuffix(expression, " ")

	tokens := []string{}

	temp := []byte{}
	for i := 0; i < len(expression); i++ {
		c := expression[i]
		switch c {
		case '\'':
			i++
			for i < len(expression) {
				if expression[i] != '\'' {
					temp = append(temp, expression[i])
				} else {
					tokens = append(tokens, string(temp))
					temp = temp[:0]
					i++
					break
				}
				i++
			}
		case ' ':
			tokens = append(tokens, string(temp))
			temp = temp[:0]
		default:
			temp = append(temp, c)
		}
	}
	tokens = append(tokens, string(temp))

	for _, t := range tokens {
		fmt.Printf("%s,", t)
	}
	positiveTags, negativeTags := []string{}, []string{}
	for _, tag := range tokens {
		if tag[0] == '-' {
			negativeTags = append(negativeTags, tag[1:])
		} else {
			positiveTags = append(positiveTags, tag)
		}
	}

	// Walk through all files and append path which contains '.go'
	filesPath, err := getFilesPath(path)
	if err != nil {
		return err
	}

	// find the file which corresponds to expression
	files := []string{}
	for _, path := range filesPath {
		fileTags, hasTags, err := extractTags(path)
		if err != nil {
			return err
		}
		if hasTags {
			isMatch := true
			for _, pTag := range positiveTags {
				contain := false
				for _, fileTag := range fileTags {
					if fileTag == pTag {
						contain = true
						break
					}
				}
				if !contain {
					isMatch = false
					break
				}
			}
			if !isMatch {
				continue
			}

			for _, nTag := range negativeTags {
				contain := false
				for _, fileTag := range fileTags {
					if fileTag == nTag {
						contain = true
						break
					}
				}
				if contain {
					isMatch = false
					break
				}
			}
			if !isMatch {
				continue
			}

			files = append(files, path)
		}
	}

	fmt.Println(files)

	return nil
}

func ListTags(path string) error {
	if tagFiles, ok, _ := readCache(); ok {
		tags := []string{}
		for key := range tagFiles {
			tags = append(tags, key)
		}
		fmt.Println(tags)
		return nil
	}

	// Walk through all files and append path which contains '.go'
	filesPath, err := getFilesPath(path)
	if err != nil {
		fmt.Println(err)
	}

	// Add file name to tag key
	tagsMap := map[string]bool{}
	for _, path := range filesPath {
		tags, hasTags, err := extractTags(path)
		if err != nil {
			return err
		}

		if hasTags {
			for _, tag := range tags {
				tagsMap[tag] = true
			}
		}
	}

	tags := []string{}
	for key := range tagsMap {
		tags = append(tags, key)
	}

	for _, tag := range tags {
		fmt.Printf("%s,", tag)
	}

	fmt.Println()
	return nil
}

func writeTagsToFile(tags []string, path string) {

}

func readCache() (map[string][]string, bool, error) {
	return map[string][]string{}, false, errors.New("not implemented")
}

func getFilesPath(path string) ([]string, error) {
	filesPath := []string{}

	// Walk through all files and append path which contains '.go'
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			filesPath = append(filesPath, path)
		}

		return nil
	})

	return filesPath, err
}

func getTagFiles(filesPath []string) (map[string][]string, error) {
	tagFiles := map[string][]string{}
	for _, path := range filesPath {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fileContent := string(content)
		lines := strings.Split(fileContent, "\n")

		// if not set tag, continue to next file
		if len(lines) < 3 {
			continue
		}

		thirdLine := strings.TrimSpace(lines[2])
		if strings.HasPrefix(thirdLine, "/* @tags:") && strings.HasSuffix(thirdLine, "*/") {
			tagsLine := strings.TrimPrefix(thirdLine, "/* @tags:")
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

	return tagFiles, nil
}

func extractTags(path string) ([]string, bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	fileContent := string(content)
	lines := strings.Split(fileContent, "\n")

	// if not set tag, continue to next file
	if len(lines) < 3 {
		return nil, false, nil
	}

	thirdLine := strings.TrimSpace(lines[2])
	if strings.HasPrefix(thirdLine, "/* @tags:") && strings.HasSuffix(thirdLine, "*/") {
		tagsLine := strings.TrimPrefix(thirdLine, "/* @tags:")
		tagsLine = strings.TrimSuffix(tagsLine, "*/")
		tagsLine = strings.TrimSpace(tagsLine)
		tags := strings.Split(tagsLine, ",")
		return tags, true, nil
	}

	return nil, false, nil
}

// Tag Statistics/Analytics: Providing insights and analytics related to tag usage, such as popular tags, tag frequency, or tag-based trends.
// Tag Synonyms: Supporting the ability to define synonyms for tags, so that searching or filtering with a synonym will retrieve items with the associated tag.
// Tag Relationships: Allowing users to define relationships between tags, such as parent-child relationships, synonyms, or related tags, to enhance searchability and discoverability.
