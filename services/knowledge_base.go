package services

import (
	"fmt"
	"go-wiki/config"
	"os"
	"path/filepath"
	"strings"
)

// 查找功能
func Search(keywords string) ([]string, error) {
	var results []string
	c := config.Init()
	kwList := ParseKeywords(keywords)
	if len(kwList) == 0 {
		return results, nil
	}

	// 1. 搜索wiki目录下的.md文件
	err := filepath.Walk(c.WikiPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := strings.ToLower(string(content))
		matched := true
		for _, kw := range kwList {
			if !strings.Contains(text, kw) {
				matched = false
				break
			}
		}
		if matched {
			summary := ExtractSummary(string(content))
			excerpt := ExtractExcerpt(text, kwList[0], string(content))
			results = append(results, fmt.Sprintf("来源: %s\n摘要: %s\n%s", path, summary, excerpt))
		}
		return nil
	})
	if err != nil {
		fmt.Println("Search wiki error: %v", err)
	}

	// 2. 搜索raw目录下的.txt文件（转换后的原始文档）
	rawPath := filepath.Join(filepath.Dir(c.WikiPath), "..", "raw")
	err = filepath.Walk(rawPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".txt") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := strings.ToLower(string(content))
		matched := true
		for _, kw := range kwList {
			if !strings.Contains(text, kw) {
				matched = false
				break
			}
		}
		if matched {
			// 原始文档直接返回完整内容
			results = append(results, fmt.Sprintf("来源(原始文档): %s\n%s", path, string(content)))
		}
		return nil
	})
	if err != nil {
		fmt.Println("Search raw error: %v", err)
	}

	return results, nil
}
func ParseKeywords(keywords string) []string {
	raw := strings.Fields(keywords)
	var output []string
	for _, kw := range raw {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw != "" {
			output = append(output, kw)
		}
	}
	return output
}

// extractExcerpt 提取关键词所在章节的完整内容
func ExtractExcerpt(lowerText, keywords, originalText string) string {
	idx := strings.Index(lowerText, keywords)
	if idx < 0 {
		return ""
	}

	// 找到关键词所在的位置，向前找到章节标题（## 开头的行）
	// 向后找到下一个章节标题或文件结尾
	lines := strings.Split(originalText, "\n")
	startLine := 0
	endLine := len(lines)

	// 找到关键词所在的行号
	currentPos := 0
	keywordLine := 0
	for i, line := range lines {
		if currentPos+len(line) >= idx {
			keywordLine = i
			break
		}
		currentPos += len(line) + 1 // +1 for newline
	}

	// 向前找到最近的章节标题（## 开头）
	for i := keywordLine; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") {
			startLine = i
			break
		}
	}

	// 向后找到下一个章节标题（## 开头，且不是子标题###）
	for i := keywordLine + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") {
			endLine = i
			break
		}
	}

	// 提取章节内容
	sectionLines := lines[startLine:endLine]
	return strings.Join(sectionLines, "\n")
}

// extractSummary 提取 frontmatter 中的 summary 字段
func ExtractSummary(content string) string {
	// 简单提取 summary: "xxx" 行
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "summary:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "summary:"))
		}
	}
	return "（无摘要）"
}
