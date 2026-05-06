package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-wiki/config"
	"io"
	"net/http"
	"os"
	"strings"
)

type LlmRequst struct {
	Model     string     `json:"model"`
	MaxTokens int        `json:"max_tokens"`
	Messages  []Messages `json:"messages"`
}
type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type LlmResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Organize 结合wiki搜索结果，调用LLM整理输出（通用）
func Organize(keyword string, wikiResults []string) (string, error) {
	c := config.Init()
	apiKey := os.Getenv("OPENROUTER_APIKEY")
	wikiStr := JoinResults(wikiResults)
	prompt := fmt.Sprintf(`你是一个资深运维专家，专注于云中间件运维。       
                                                                                
  用户查询关键词：%s                                                            
                                                                                
  以下是知识库中与查询相关的文档内容：
  %s                                                                            
                                                                                
  请完成以下任务：                                                              
  1. 从知识库内容中提取与"%s"直接相关的知识点                                   
  2. 按逻辑分类整理（如：定义、排查步骤、处理方法、注意事项等）                 
  3. 如果知识库信息不足，请基于中间件运维经验补充                               
  4. 给出清晰的结构化输出，使用Markdown格式                                     
  5. 标注每个知识点的来源文档`, keyword, wikiStr, keyword)
	reqBody := LlmRequst{
		Model:     c.LlmModel,
		MaxTokens: 3000,
		Messages:  []Messages{{Role: "user", Content: prompt}},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("marshal failed", err)
		return "", err
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", c.LlmUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create request failed", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	//client := &http.Client{Timeout: 60 * time.Second}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body failed:%w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("LLM API error: %w", body)
	}
	var result LlmResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal failed:%w", err)
	}
	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "未获取到LLM响应", nil
}
func JoinResults(results []string) string {
	if len(results) == 0 {
		return "（知识库未找到相关内容）"
	}
	return strings.Join(results, "\n\n---\n\n")
}
