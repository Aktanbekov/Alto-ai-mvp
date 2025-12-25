package tests

import (
	"testing"
	"time"
	"altoai_mvp/interview"
)

func TestGenerateSessionSummary(t *testing.T) {
	session := interview.NewSession("test-user")
	
	// Add some answers with analyses
	session.Answers = []interview.Answer{
		{
			QuestionID:   "q1",
			QuestionText: "Question 1",
			Text:         "Answer 1",
			CreatedAt:    time.Now(),
			Analysis: &interview.AnalysisResponse{
				Scores: interview.AnalysisScores{
					MigrationIntent:  5,
					GoalUnderstanding: 4,
					AnswerLength:     5,
					TotalScore:       14,
				},
				Classification: "Good",
			},
		},
		{
			QuestionID:   "q2",
			QuestionText: "Question 2",
			Text:         "Answer 2",
			CreatedAt:    time.Now(),
			Analysis: &interview.AnalysisResponse{
				Scores: interview.AnalysisScores{
					MigrationIntent:  4,
					GoalUnderstanding: 5,
					AnswerLength:     4,
					TotalScore:       13,
				},
				Classification: "Good",
			},
		},
	}

	summary, err := interview.GenerateSessionSummary(session)
	if err != nil {
		t.Fatalf("GenerateSessionSummary failed: %v", err)
	}

	if summary.TotalQuestions != 2 {
		t.Errorf("Expected 2 questions, got %d", summary.TotalQuestions)
	}
	if summary.AverageScore != 13.5 {
		t.Errorf("Expected average score 13.5, got %.2f", summary.AverageScore)
	}
	if summary.OverallGrade != "B" {
		t.Errorf("Expected grade B, got %s", summary.OverallGrade)
	}
}

func TestGenerateSessionSummaryEmptySession(t *testing.T) {
	session := interview.NewSession("test-user")
	
	_, err := interview.GenerateSessionSummary(session)
	if err == nil {
		t.Error("Should return error for empty session")
	}
}


