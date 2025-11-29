package interview

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

const FirstQuestionID = "q1_purpose"

var (
	sessions   = make(map[string]*Session)
	sessionsMu sync.RWMutex
)

func NewSession(userID string) *Session {
	now := time.Now()
	return &Session{
		ID:              uuid.NewString(),
		UserID:          userID,
		CurrentQuestion: FirstQuestionID,
		Answers:         []Answer{},
		Scores:          Scores{},
		Status:          SessionStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
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
