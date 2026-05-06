package controller

import (
	"go-wiki/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	Keyword string `json:"keyword"` // 查询关键词，如 "ISR副本同步"
}
type QueryResponse struct {
	Keyword     string   `json:"keyword"`
	WikiSources []string `json:"wiki_sources"` // 命中的wiki来源
	LLMAnswer   string   `json:"llm_answer"`   // 大模型整理的答案
}

func QueryHandler(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}
	// 通用搜索
	wikiResults, err := services.Search(req.Keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 调用 LLM 整理
	answer, err := services.Organize(req.Keyword, wikiResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, QueryResponse{
		Keyword:     req.Keyword,
		WikiSources: wikiResults,
		LLMAnswer:   answer,
	})
}
