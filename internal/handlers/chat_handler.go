package handlers

import (
	"altoai_mvp/logic"
	"altoai_mvp/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequest struct {
	Messages []logic.Message `json:"messages"`
}

type ChatResponse struct {
	Content string `json:"content"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Messages) == 0 {
		response.Error(c, http.StatusBadRequest, "messages cannot be empty")
		return
	}

	responseText, err := logic.GetGPTResponse(req.Messages)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get GPT response: "+err.Error())
		return
	}

	response.OK(c, ChatResponse{Content: responseText})
}


