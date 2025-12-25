package interview

import "time"

// Question represents one node in your interview graph.
type Question struct {
	ID                 string   `json:"id"`                  // e.g. "q1_purpose"
	Category           string   `json:"category"`            // e.g. "Purpose of Study"
	Text               string   `json:"text"`                // full question text
	NextID             string   `json:"next_id"`             // linear next question in normal flow
	FollowupCandidates []string `json:"followup_candidates"` // allowed followups from this node
	Tags               []string `json:"tags"`                // semantic tags: ["purpose", "intent", "risk"]
}

// Answer is one student response.
type Answer struct {
	QuestionID   string    `json:"question_id"`
	QuestionText string    `json:"question_text"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"created_at"`
	// Optional: store AI eval snapshot per answer for analytics
	Eval *EvalResult `json:"eval,omitempty"`
	// New grading system analysis
	Analysis *AnalysisResponse `json:"analysis,omitempty"`
}

// Scores are cumulative across the entire session.
type Scores struct {
	Academic       int `json:"academic"`
	Financial      int `json:"financial"`
	IntentToReturn int `json:"intent_to_return"`
	OverallRisk    int `json:"overall_risk"`
}

// SessionStatus allows you to extend states later (paused, aborted, etc.).
type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusFinished SessionStatus = "finished"
	SessionStatusAborted  SessionStatus = "aborted"
)

// Session holds the state of one full interview attempt.
type Session struct {
	ID                string        `json:"id"`
	UserID            string        `json:"user_id,omitempty"`  // if you later have accounts
	CurrentQuestion   string        `json:"current_question"`   // question ID
	SelectedQuestions []Question    `json:"selected_questions"` // questions selected for this session
	QuestionIndex     int           `json:"question_index"`     // current question index in SelectedQuestions
	Answers           []Answer      `json:"answers"`
	Scores            Scores        `json:"scores"`
	Status            SessionStatus `json:"status"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	// Session summary for completed interviews
	Summary *SessionSummary `json:"summary,omitempty"`
}

// AnalysisScores represents the 3–15 grading system for a single answer
type AnalysisScores struct {
	MigrationIntent   int `json:"migration_intent"`   // 1–5
	GoalUnderstanding int `json:"goal_understanding"` // 1–5
	AnswerLength      int `json:"answer_length"`      // 1–5
	TotalScore        int `json:"total_score"`        // 3–15 (sum of all criteria)
}

// FeedbackByCriterion contains feedback for each scoring criterion
type FeedbackByCriterion struct {
	MigrationIntent   string `json:"migration_intent"`
	GoalUnderstanding string `json:"goal_understanding"`
	AnswerLength      string `json:"answer_length"`
}

// StructuredFeedback contains detailed feedback in the new format
type StructuredFeedback struct {
	Overall      string              `json:"overall"`
	ByCriterion  FeedbackByCriterion `json:"by_criterion"`
	Improvements []string            `json:"improvements"`
}

// AnalysisResponse contains detailed analysis of a single answer (new grading system)
type AnalysisResponse struct {
	Scores         AnalysisScores     `json:"scores"`
	Classification string             `json:"classification"` // Excellent, Good, Average, Weak
	Feedback       StructuredFeedback `json:"feedback"`       // Structured feedback with overall, by_criterion, and improvements
}

// AnalysisRecord stores a complete analysis record
type AnalysisRecord struct {
	ID        string           `json:"id"`
	SessionID string           `json:"sessionId,omitempty"`
	Question  string           `json:"question"`
	Answer    string           `json:"answer"`
	Analysis  AnalysisResponse `json:"analysis"`
	CreatedAt time.Time        `json:"createdAt"`
}

// SessionSummary provides overall assessment of a completed interview session
type SessionSummary struct {
	SessionID      string    `json:"sessionId"`
	TotalQuestions int       `json:"totalQuestions"`
	AverageScore   float64   `json:"averageScore"`
	OverallGrade   string    `json:"overallGrade"`
	StrongAreas    []string  `json:"strongAreas"`
	WeakAreas      []string  `json:"weakAreas"`
	CommonRedFlags []string  `json:"commonRedFlags"`
	Recommendation string    `json:"recommendation"`
	CompletedAt    time.Time `json:"completedAt"`
}
