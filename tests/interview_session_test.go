package tests

import (
	"testing"
	"time"
	"altoai_mvp/interview"
)

func TestNewSession(t *testing.T) {
	// Load questions first
	err := interview.LoadQuestions("../interview/questions.json")
	if err != nil {
		t.Fatalf("LoadQuestions failed: %v", err)
	}
	
	session := interview.NewSession("user123")
	
	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}
	if session.UserID != "user123" {
		t.Errorf("Expected UserID to be 'user123', got %s", session.UserID)
	}
	if len(session.SelectedQuestions) == 0 {
		t.Error("Session should have selected questions")
	}
	if session.CurrentQuestion == "" {
		t.Error("Session should have a current question")
	}
	if session.QuestionIndex != 0 {
		t.Errorf("Expected QuestionIndex to be 0, got %d", session.QuestionIndex)
	}
	if session.Status != interview.SessionStatusActive {
		t.Errorf("Expected Status to be %s, got %s", interview.SessionStatusActive, session.Status)
	}
	if len(session.Answers) != 0 {
		t.Error("New session should have no answers")
	}
}

func TestSaveAndGetSession(t *testing.T) {
	session := interview.NewSession("test-user")
	sessionID := session.ID
	
	interview.SaveSession(session)
	
	retrieved, ok := interview.GetSession(sessionID)
	if !ok {
		t.Error("Session should be retrievable after saving")
	}
	if retrieved.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, retrieved.ID)
	}
	if retrieved.UserID != "test-user" {
		t.Errorf("Expected UserID 'test-user', got %s", retrieved.UserID)
	}
}

func TestGetNonExistentSession(t *testing.T) {
	_, ok := interview.GetSession("non-existent-id")
	if ok {
		t.Error("Non-existent session should return false")
	}
}

func TestSessionAnswers(t *testing.T) {
	session := interview.NewSession("test-user")
	
	answer := interview.Answer{
		QuestionID:   "q1",
		QuestionText: "Test question",
		Text:         "Test answer",
		CreatedAt:    time.Now(),
	}
	
	session.Answers = append(session.Answers, answer)
	interview.SaveSession(session)
	
	retrieved, _ := interview.GetSession(session.ID)
	if len(retrieved.Answers) != 1 {
		t.Errorf("Expected 1 answer, got %d", len(retrieved.Answers))
	}
	if retrieved.Answers[0].Text != "Test answer" {
		t.Errorf("Expected answer text 'Test answer', got %s", retrieved.Answers[0].Text)
	}
}

