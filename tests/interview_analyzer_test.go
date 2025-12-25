package tests

import (
	"testing"
	"altoai_mvp/interview"
)

// Note: getGradeFromScore is not exported, so we test through ScoreToPercentage
// which uses the grade calculation internally

func TestScoreToPercentage(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected float64
	}{
		{"Minimum score - 3", 3, 0.0},
		{"Maximum score - 15", 15, 100.0},
		{"Middle score - 9", 9, 50.0},
		{"Good score - 13", 13, 83.33},
		{"Average score - 11", 11, 66.67},
		{"Below minimum - 0", 0, 0.0},
		{"Above maximum - 20", 20, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interview.ScoreToPercentage(tt.score)
			// Allow small floating point differences
			if result < tt.expected-0.1 || result > tt.expected+0.1 {
				t.Errorf("ScoreToPercentage(%d) = %.2f, want approximately %.2f", tt.score, result, tt.expected)
			}
		})
	}
}

func TestNewVisaAnalyzer(t *testing.T) {
	// Test with empty API key (should use environment)
	analyzer := interview.NewVisaAnalyzer("")
	if analyzer == nil {
		t.Error("NewVisaAnalyzer returned nil")
	}

	// Test with provided API key
	analyzer2 := interview.NewVisaAnalyzer("test-key")
	if analyzer2 == nil {
		t.Error("NewVisaAnalyzer returned nil")
	}
}

