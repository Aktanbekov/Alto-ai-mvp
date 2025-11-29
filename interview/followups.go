package interview

// Optional generic mapping if you want AI to specify a type like "clarify_home_ties".
var FollowupByType = map[string][]string{
	"clarify_purpose": {
		"q1f_clarify_purpose",
	},
	"clarify_university": {
		"q2f_university_exact",
	},
	"clarify_financial": {
		"q5f_finance_clarify",
		"q6f_finance_detail",
	},
	"clarify_home_ties": {
		"q7f_home_country_career",
		"q8f_ties_detail",
	},
}

// Utility: has this followup been asked already in this session.
func hasAskedQuestion(s *Session, questionID string) bool {
	for _, ans := range s.Answers {
		if ans.QuestionID == questionID {
			return true
		}
	}
	return false
}

// DecideNextQuestion uses AI eval plus graph rules to select the next question id.
func DecideNextQuestion(current Question, s *Session, eval *EvalResult) string {
	if eval != nil && eval.NeedsFollowup {
		if next := pickFollowupQuestion(current, eval.SuggestedFollowup, s); next != "" {
			return next
		}
	}
	return current.NextID
}

func pickFollowupQuestion(current Question, followupType string, s *Session) string {
	if followupType == "" {
		return ""
	}

	// Allowed followups for this question
	allowed := map[string]bool{}
	for _, id := range current.FollowupCandidates {
		allowed[id] = true
	}

	candidates, ok := FollowupByType[followupType]
	if !ok {
		return ""
	}

	for _, id := range candidates {
		if !allowed[id] {
			continue
		}
		if hasAskedQuestion(s, id) {
			continue
		}
		return id
	}

	return ""
}
