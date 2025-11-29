package interview

import "fmt"

// ScoreDelta indicates how one answer should adjust the session scores.
type ScoreDelta struct {
	Academic       int `json:"academic"`
	Financial      int `json:"financial"`
	IntentToReturn int `json:"intent_to_return"`
	OverallRisk    int `json:"overall_risk"`
}

// EvalResult is what the AI returns for a single answer.
type EvalResult struct {
	Quality            int        `json:"quality"`                 // 0–10
	Clarity            int        `json:"clarity"`                 // 0–10
	Confidence         int        `json:"confidence"`              // 0–10
	Flags              []string   `json:"flags"`                   // ["mentions_work_in_us", ...]
	IntentToReturnRisk int        `json:"intent_to_return_risk"`   // 0–10
	SuggestedFollowup  string     `json:"suggested_followup_type"` // "clarify_home_ties", etc.
	NeedsFollowup      bool       `json:"needs_followup"`
	ScoreDelta         ScoreDelta `json:"score_delta"`
}

func ApplyEval(s *Session, eval *EvalResult) {
	if eval == nil {
		return
	}
	s.Scores.Academic += eval.ScoreDelta.Academic
	s.Scores.Financial += eval.ScoreDelta.Financial
	s.Scores.IntentToReturn += eval.ScoreDelta.IntentToReturn
	s.Scores.OverallRisk += eval.ScoreDelta.OverallRisk
}

// ApplyAnalysis applies an AnalysisResponse to the session
// This converts the analysis to EvalResult and applies it, maintaining backward compatibility
func ApplyAnalysis(s *Session, analysis *AnalysisResponse, q Question) {
	if analysis == nil {
		return
	}
	eval := ConvertAnalysisToEval(analysis, q)
	ApplyEval(s, eval)
}

// GenerateSessionSummary generates a summary from all answers in the session
func GenerateSessionSummary(s *Session) (*SessionSummary, error) {
	if len(s.Answers) == 0 {
		return nil, fmt.Errorf("no answers in session")
	}

	// Convert answers to analysis records
	analyses := make([]AnalysisRecord, 0, len(s.Answers))
	for _, answer := range s.Answers {
		if answer.Analysis != nil {
			analyses = append(analyses, AnalysisRecord{
				ID:        fmt.Sprintf("analysis_%s_%d", s.ID, len(analyses)),
				SessionID: s.ID,
				Question:  answer.QuestionText,
				Answer:    answer.Text,
				Analysis:  *answer.Analysis,
				CreatedAt: answer.CreatedAt,
			})
		}
	}

	if len(analyses) == 0 {
		return nil, fmt.Errorf("no analyses found in session")
	}

	analyzer := GetAnalyzer()
	if analyzer == nil {
		return nil, fmt.Errorf("analyzer not initialized")
	}

	summary, err := analyzer.GenerateSessionSummary(analyses)
	if err != nil {
		return nil, err
	}

	summary.SessionID = s.ID
	return summary, nil
}
