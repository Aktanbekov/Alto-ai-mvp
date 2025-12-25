package tests

import (
	"encoding/json"
	"testing"
	"altoai_mvp/interview"
)

func TestAnalysisScoresJSON(t *testing.T) {
	scores := interview.AnalysisScores{
		MigrationIntent:  5,
		GoalUnderstanding: 4,
		AnswerLength:     5,
		TotalScore:       14,
	}

	jsonData, err := json.Marshal(scores)
	if err != nil {
		t.Fatalf("Failed to marshal AnalysisScores: %v", err)
	}

	var unmarshaled interview.AnalysisScores
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AnalysisScores: %v", err)
	}

	if unmarshaled.MigrationIntent != scores.MigrationIntent {
		t.Errorf("MigrationIntent mismatch: got %d, want %d", unmarshaled.MigrationIntent, scores.MigrationIntent)
	}
	if unmarshaled.GoalUnderstanding != scores.GoalUnderstanding {
		t.Errorf("GoalUnderstanding mismatch: got %d, want %d", unmarshaled.GoalUnderstanding, scores.GoalUnderstanding)
	}
	if unmarshaled.AnswerLength != scores.AnswerLength {
		t.Errorf("AnswerLength mismatch: got %d, want %d", unmarshaled.AnswerLength, scores.AnswerLength)
	}
	if unmarshaled.TotalScore != scores.TotalScore {
		t.Errorf("TotalScore mismatch: got %d, want %d", unmarshaled.TotalScore, scores.TotalScore)
	}
}

func TestStructuredFeedbackJSON(t *testing.T) {
	feedback := interview.StructuredFeedback{
		Overall: "Good answer overall",
		ByCriterion: interview.FeedbackByCriterion{
			MigrationIntent:  "No migration intent shown",
			GoalUnderstanding: "Clear goals",
			AnswerLength:     "Appropriate length",
		},
		Improvements: []string{"Be more specific", "Add examples"},
	}

	jsonData, err := json.Marshal(feedback)
	if err != nil {
		t.Fatalf("Failed to marshal StructuredFeedback: %v", err)
	}

	var unmarshaled interview.StructuredFeedback
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal StructuredFeedback: %v", err)
	}

	if unmarshaled.Overall != feedback.Overall {
		t.Errorf("Overall mismatch: got %s, want %s", unmarshaled.Overall, feedback.Overall)
	}
	if len(unmarshaled.Improvements) != len(feedback.Improvements) {
		t.Errorf("Improvements length mismatch: got %d, want %d", len(unmarshaled.Improvements), len(feedback.Improvements))
	}
}

func TestAnalysisResponseJSON(t *testing.T) {
	response := interview.AnalysisResponse{
		Scores: interview.AnalysisScores{
			MigrationIntent:  5,
			GoalUnderstanding: 4,
			AnswerLength:     5,
			TotalScore:       14,
		},
		Classification: "Good",
		Feedback: interview.StructuredFeedback{
			Overall: "Good answer",
			ByCriterion: interview.FeedbackByCriterion{
				MigrationIntent:  "Good",
				GoalUnderstanding: "Good",
				AnswerLength:     "Good",
			},
			Improvements: []string{"Improve clarity"},
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal AnalysisResponse: %v", err)
	}

	var unmarshaled interview.AnalysisResponse
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AnalysisResponse: %v", err)
	}

	if unmarshaled.Classification != response.Classification {
		t.Errorf("Classification mismatch: got %s, want %s", unmarshaled.Classification, response.Classification)
	}
	if unmarshaled.Scores.TotalScore != response.Scores.TotalScore {
		t.Errorf("TotalScore mismatch: got %d, want %d", unmarshaled.Scores.TotalScore, response.Scores.TotalScore)
	}
}


