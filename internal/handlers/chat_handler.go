package handlers

import (
	"altoai_mvp/interview"
	"altoai_mvp/pkg/response"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequest struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	SessionID string `json:"session_id,omitempty"` // Optional: for continuing existing interview
}

type ChatResponse struct {
	Content         string                      `json:"content"`                    // The question text or completion message
	SessionID       string                      `json:"session_id,omitempty"`       // Session ID for client to track
	QuestionID      string                      `json:"question_id,omitempty"`      // Current question ID
	Finished        bool                        `json:"finished"`                   // Whether interview is complete
	Scores          *interview.Scores           `json:"scores,omitempty"`           // Current risk scores
	IsNewSession    bool                        `json:"is_new_session,omitempty"`   // Whether this is a new session
	Analysis        *interview.AnalysisResponse `json:"analysis,omitempty"`         // Detailed analysis of the answer
	Grade           string                      `json:"grade,omitempty"`            // Letter grade (A-F) for the answer
	Suggestions     []string                    `json:"suggestions,omitempty"`      // Improvement suggestions
	ImprovedVersion string                      `json:"improved_version,omitempty"` // Suggested improved answer
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get or create session
	var session *interview.Session
	var isNewSession bool

	if req.SessionID != "" {
		// Try to retrieve existing session
		if s, ok := interview.GetSession(req.SessionID); ok {
			session = s
			isNewSession = false
		} else {
			// Session not found, create new one
			session = interview.NewSession("")
			interview.SaveSession(session)
			isNewSession = true
		}
	} else {
		// No session ID provided, create new session
		session = interview.NewSession("")
		interview.SaveSession(session)
		isNewSession = true
	}

	// If session is finished, return completion message
	if session.Status == interview.SessionStatusFinished {
		completionMsg := buildCompletionMessage(session)
		response.OK(c, ChatResponse{
			Content:      completionMsg,
			SessionID:    session.ID,
			Finished:     true,
			Scores:       &session.Scores,
			IsNewSession: false,
		})
		return
	}

	// If this is a new session, return the first question
	if isNewSession {
		currentQ, ok := interview.Questions[session.CurrentQuestion]
		if !ok {
			response.Error(c, http.StatusInternalServerError, "initial question not found")
			return
		}

		response.OK(c, ChatResponse{
			Content:      currentQ.Text,
			SessionID:    session.ID,
			QuestionID:   currentQ.ID,
			Finished:     false,
			IsNewSession: isNewSession,
		})
		return
	}

	// If no messages provided, return current question
	if len(req.Messages) == 0 {
		currentQ, ok := interview.Questions[session.CurrentQuestion]
		if !ok {
			response.Error(c, http.StatusInternalServerError, "current question not found")
			return
		}

		response.OK(c, ChatResponse{
			Content:    currentQ.Text,
			SessionID:  session.ID,
			QuestionID: currentQ.ID,
			Finished:   false,
			Scores:     &session.Scores,
		})
		return
	}

	// Find the last user message
	var lastUserMessage string
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			lastUserMessage = req.Messages[i].Content
			break
		}
	}

	if lastUserMessage == "" {
		response.Error(c, http.StatusBadRequest, "no user message found")
		return
	}

	// Get current question
	currentQ, ok := interview.Questions[session.CurrentQuestion]
	if !ok {
		session.Status = interview.SessionStatusFinished
		interview.SaveSession(session)
		response.Error(c, http.StatusInternalServerError, "current question not found")
		return
	}

	// Check if we've already answered this question (prevent duplicate processing)
	alreadyAnswered := false
	for _, ans := range session.Answers {
		if ans.QuestionID == currentQ.ID {
			alreadyAnswered = true
			break
		}
	}

	// If we've already answered this question, just return the next question
	if alreadyAnswered {
		// Decide next question (use nil eval since we're not processing a new answer)
		nextID := currentQ.NextID
		if nextID == "" || nextID == "end" {
			session.Status = interview.SessionStatusFinished

			// Generate session summary before completing
			summary, err := interview.GenerateSessionSummary(session)
			if err == nil && summary != nil {
				session.Summary = summary
			}

			interview.SaveSession(session)

			completionMsg := buildCompletionMessage(session)
			response.OK(c, ChatResponse{
				Content:   completionMsg,
				SessionID: session.ID,
				Finished:  true,
				Scores:    &session.Scores,
			})
			return
		}

		// Update session with next question
		session.CurrentQuestion = nextID
		interview.SaveSession(session)

		nextQ := interview.Questions[nextID]
		response.OK(c, ChatResponse{
			Content:    nextQ.Text,
			SessionID:  session.ID,
			QuestionID: nextQ.ID,
			Finished:   false,
			Scores:     &session.Scores,
		})
		return
	}

	// Record the answer
	answer := interview.Answer{
		QuestionID:   currentQ.ID,
		QuestionText: currentQ.Text,
		Text:         lastUserMessage,
		CreatedAt:    time.Now(),
	}

	// Call new analyzer for detailed feedback
	analysis, err := interview.AnalyzeAnswer(currentQ, lastUserMessage)
	if err != nil {
		// Log error for debugging
		log.Printf("Error analyzing answer: %v", err)
		// Continue without analysis (graceful degradation)
		analysis = nil
	} else if analysis != nil {
		log.Printf("Analysis successful: Classification=%s, TotalScore=%d",
			analysis.Classification, analysis.Scores.TotalScore)
	}

	// Attach analysis to answer
	if analysis != nil {
		answer.Analysis = analysis
		// Also create EvalResult for backward compatibility with scoring system
		eval := interview.ConvertAnalysisToEval(analysis, currentQ)
		answer.Eval = eval
		// Update scores using the converted eval
		interview.ApplyEval(session, eval)
	}

	session.Answers = append(session.Answers, answer)

	// Check if we've reached 12 questions limit
	if len(session.Answers) >= 12 {
		session.Status = interview.SessionStatusFinished

		// Generate session summary before completing
		summary, err := interview.GenerateSessionSummary(session)
		if err == nil && summary != nil {
			session.Summary = summary
		}

		interview.SaveSession(session)

		completionMsg := buildCompletionMessage(session)
		response.OK(c, ChatResponse{
			Content:   completionMsg,
			SessionID: session.ID,
			Finished:  true,
			Scores:    &session.Scores,
			Analysis:  analysis,
			Grade:     getGradeFromAnalysis(analysis),
		})
		return
	}

	// Decide next question (use eval from analysis if available)
	var eval *interview.EvalResult
	if analysis != nil {
		eval = interview.ConvertAnalysisToEval(analysis, currentQ)
	}
	nextID := interview.DecideNextQuestion(currentQ, session, eval)
	if nextID == "" || nextID == "end" {
		session.Status = interview.SessionStatusFinished

		// Generate session summary before completing
		summary, err := interview.GenerateSessionSummary(session)
		if err == nil && summary != nil {
			session.Summary = summary
		}

		interview.SaveSession(session)

		completionMsg := buildCompletionMessage(session)
		response.OK(c, ChatResponse{
			Content:   completionMsg,
			SessionID: session.ID,
			Finished:  true,
			Scores:    &session.Scores,
			Analysis:  analysis,
			Grade:     getGradeFromAnalysis(analysis),
		})
		return
	}

	// Update session with next question
	session.CurrentQuestion = nextID
	interview.SaveSession(session)

	nextQ := interview.Questions[nextID]

	response.OK(c, ChatResponse{
		Content:         nextQ.Text,
		SessionID:       session.ID,
		QuestionID:      nextQ.ID,
		Finished:        false,
		Scores:          &session.Scores,
		Analysis:        analysis,
		Grade:           getGradeFromAnalysis(analysis),
		Suggestions:     getSuggestionsFromAnalysis(analysis),
		ImprovedVersion: getImprovedVersionFromAnalysis(analysis),
	})
}

// buildCompletionMessage creates a completion message based on session summary or scores
func buildCompletionMessage(session *interview.Session) string {
	// Use session summary if available (new grading system)
	if session.Summary != nil {
		return "Thank you for completing the interview practice session! " +
			"Your overall grade is: " + session.Summary.OverallGrade + " (Average Score: " +
			fmt.Sprintf("%.1f", session.Summary.AverageScore) + "). " +
			session.Summary.Recommendation + " " +
			"Good luck with your visa interview!"
	}

	// Fallback to old scoring system
	scores := session.Scores
	totalRisk := scores.Academic + scores.Financial + scores.IntentToReturn + scores.OverallRisk
	avgRisk := float64(totalRisk) / 4.0

	var assessment string
	if avgRisk < 25 {
		assessment = "excellent"
	} else if avgRisk < 50 {
		assessment = "good"
	} else if avgRisk < 75 {
		assessment = "moderate"
	} else {
		assessment = "needs improvement"
	}

	return "Thank you for completing the interview practice session! " +
		"Your overall assessment is: " + assessment + ". " +
		"Keep practicing to improve your answers and confidence. " +
		"Good luck with your visa interview!"
}

// Helper functions to extract data from analysis
func getGradeFromAnalysis(analysis *interview.AnalysisResponse) string {
	if analysis == nil {
		return ""
	}
	// Map classification to a rough letter grade for backward compatibility
	switch strings.ToLower(analysis.Classification) {
	case "excellent":
		return "A"
	case "good":
		return "B"
	case "average":
		return "C"
	case "weak":
		return "D"
	case "poor":
		return "F"
	default:
		return ""
	}
}

func getSuggestionsFromAnalysis(analysis *interview.AnalysisResponse) []string {
	if analysis == nil {
		return nil
	}
	if strings.TrimSpace(analysis.Feedback) == "" {
		return nil
	}
	// Use feedback as a single suggestion entry
	return []string{analysis.Feedback}
}

func getImprovedVersionFromAnalysis(analysis *interview.AnalysisResponse) string {
	if analysis == nil {
		return ""
	}
	// Improved version no longer provided in new analysis format
	return ""
}
