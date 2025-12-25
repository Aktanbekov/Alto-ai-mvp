package tests

import (
	"testing"
	"altoai_mvp/interview"
)

func TestLoadQuestions(t *testing.T) {
	// Test loading questions from JSON file
	err := interview.LoadQuestions("../interview/questions.json")
	if err != nil {
		t.Fatalf("LoadQuestions failed: %v", err)
	}
	
	if len(interview.QuestionsByCategory) == 0 {
		t.Error("Questions should be loaded")
	}
	
	// Check that categories exist
	requiredCategories := []string{
		"Purpose of Study",
		"Academic Background",
		"University Choice",
		"Financial Capability",
		"Family/Sponsor Info",
		"Post-Graduation Plans",
		"Immigration Intent",
	}
	
	for _, category := range requiredCategories {
		questions, ok := interview.QuestionsByCategory[category]
		if !ok {
			t.Errorf("Category %s should exist", category)
		}
		if len(questions) == 0 {
			t.Errorf("Category %s should have questions", category)
		}
	}
}

func TestQuestionSelection(t *testing.T) {
	err := interview.LoadQuestions("../interview/questions.json")
	if err != nil {
		t.Fatalf("LoadQuestions failed: %v", err)
	}
	
	selected := interview.SelectQuestionsForSession()
	if len(selected) == 0 {
		t.Error("Should select questions for session")
	}
	
	// Check that we have the right number of questions
	expectedTotal := 2 + 2 + 2 + 2 + 1 + 2 + 1 // 12 total
	if len(selected) != expectedTotal {
		t.Errorf("Expected %d questions, got %d", expectedTotal, len(selected))
	}
	
	// Check that all selected questions have valid structure
	for _, q := range selected {
		if q.ID == "" {
			t.Error("Selected question should have ID")
		}
		if q.Text == "" {
			t.Error("Selected question should have text")
		}
		if q.Category == "" {
			t.Error("Selected question should have category")
		}
	}
}

