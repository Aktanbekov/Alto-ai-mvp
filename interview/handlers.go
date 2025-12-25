package interview

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateSessionResponse struct {
	SessionID    string `json:"session_id"`
	QuestionID   string `json:"question_id"`
	QuestionText string `json:"question_text"`
}

func CreateSessionHandler(c *gin.Context) {
	// In future you can extract userID from auth
	s := NewSession("")
	SaveSession(s)

	if len(s.SelectedQuestions) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no questions selected for session"})
		return
	}

	firstQ := s.SelectedQuestions[0]
	resp := CreateSessionResponse{
		SessionID:    s.ID,
		QuestionID:   firstQ.ID,
		QuestionText: firstQ.Text,
	}

	c.JSON(http.StatusOK, resp)
}

type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type SubmitAnswerResponse struct {
	NextQuestionID   string `json:"next_question_id,omitempty"`
	NextQuestionText string `json:"next_question_text,omitempty"`
	Finished         bool   `json:"finished"`
	Scores           Scores `json:"scores"`
	CurrentQuestion  string `json:"current_question,omitempty"`
}

func SubmitAnswerHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing session id"})
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	s, ok := GetSession(sessionID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	if s.Status != SessionStatusActive {
		resp := SubmitAnswerResponse{
			Finished:        true,
			Scores:          s.Scores,
			CurrentQuestion: s.CurrentQuestion,
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Find current question in selected questions
	var currentQ *Question
	for i, q := range s.SelectedQuestions {
		if q.ID == s.CurrentQuestion {
			currentQ = &s.SelectedQuestions[i]
			break
		}
	}

	if currentQ == nil {
		s.Status = SessionStatusFinished
		SaveSession(s)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "current question not found"})
		return
	}

	// Record the answer
	answer := Answer{
		QuestionID:   req.QuestionID,
		QuestionText: currentQ.Text,
		Text:         req.Answer,
		CreatedAt:    time.Now(),
	}
	s.Answers = append(s.Answers, answer)

	// Call AI evaluator with session context
	eval, err := CallLLM(s, *currentQ, req.Answer)
	if err != nil {
		// You can log and degrade gracefully to rule based only
		eval = nil
	}
	// Attach eval to last answer for analytics
	if eval != nil {
		s.Answers[len(s.Answers)-1].Eval = eval
	}

	// Update scores
	ApplyEval(s, eval)

	// Move to next question
	s.QuestionIndex++
	if s.QuestionIndex >= len(s.SelectedQuestions) {
		// All questions answered
		s.Status = SessionStatusFinished
		SaveSession(s)
		resp := SubmitAnswerResponse{
			Finished:        true,
			Scores:          s.Scores,
			CurrentQuestion: "",
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Set next question
	nextQ := s.SelectedQuestions[s.QuestionIndex]
	s.CurrentQuestion = nextQ.ID
	SaveSession(s)

	resp := SubmitAnswerResponse{
		NextQuestionID:   nextQ.ID,
		NextQuestionText: nextQ.Text,
		Finished:         false,
		Scores:           s.Scores,
		CurrentQuestion:  nextQ.ID,
	}

	c.JSON(http.StatusOK, resp)
}
