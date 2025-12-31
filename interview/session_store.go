package interview

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	sessions   = make(map[string]*Session)
	sessionsMu sync.RWMutex
)

func NewSession(userID string) *Session {
	return NewSessionWithLevel(userID, "")
}

func NewSessionWithLevel(userID string, level string) *Session {
	now := time.Now()
	
	// Select questions for this session based on level
	selectedQuestions := SelectQuestionsForSession(level)
	
	session := &Session{
		ID:               uuid.NewString(),
		UserID:           userID,
		SelectedQuestions: selectedQuestions,
		QuestionIndex:    0,
		Answers:          []Answer{},
		Scores:           Scores{},
		Status:           SessionStatusActive,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	
	// Set current question to first selected question
	if len(selectedQuestions) > 0 {
		session.CurrentQuestion = selectedQuestions[0].ID
	}
	
	return session
}

func SaveSession(s *Session) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s.UpdatedAt = time.Now()
	sessions[s.ID] = s
}

func GetSession(id string) (*Session, bool) {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	s, ok := sessions[id]
	return s, ok
}
