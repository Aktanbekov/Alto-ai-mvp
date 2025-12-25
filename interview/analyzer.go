package interview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// VisaAnalyzer handles AI-powered analysis of visa interview answers
type VisaAnalyzer struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
	// Cache the system prompt to avoid regenerating it
	systemPrompt string
}

// NewVisaAnalyzer creates a new VisaAnalyzer instance
func NewVisaAnalyzer(apiKey string) *VisaAnalyzer {
	if apiKey == "" {
		// Try to get from environment
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			apiKey = os.Getenv("GPT_API_KEY")
		}
	}

	systemPrompt := `You are an F1 visa interview grading engine.

You grade ONE student answer at a time.

INPUT YOU WILL RECEIVE:
- question: the F1 visa interview question asked by the officer
- answer: the student's answer text

YOUR TASK:
1) Score the answer in 3 criteria, each from 1 to 5:
   - migration_intent
   - goal_understanding
   - answer_length

2) Compute total_score = migration_intent + goal_understanding + answer_length.

3) Set classification based on total_score:
   - 15 => "Excellent"
   - 13–14 => "Good"
   - 11–12 => "Average"
   - 3–10 => "Weak"

4) Provide structured feedback:
   - "overall": 1–3 sentences summarizing the quality of the answer.
   - "by_criterion": short explanations for each score.
   - "improvements": 1–3 concrete, actionable suggestions.

GENERAL RULES:
- Grade ONLY based on what is written in the answer.
- Do NOT invent facts or assume things that are not clearly stated.
- Focus on the content and structure of the answer, not English accent or minor grammar mistakes.
- Student style can be simple and direct. Do not penalize just for not sounding academic.

DETECT QUESTION TYPE FIRST:
Some questions are about goals and intent (for example: “Why do you want to study in the US?”, “What are your plans after graduation?”, “Do you plan to work in the US?”).
Other questions are factual (for example: “Who is sponsoring your studies?”, “Do you have any gaps?”, “What is your university name?”).

Use the full rubric for goal or intent questions.
For factual questions, do NOT penalize the student for not talking about long-term goals unless the question clearly asks for them.

SCORING RULES:

1) migration_intent (1 to 5)

This criterion measures how clearly the answer avoids migration risk and shows intention to return home.

Give HIGH scores when:
- The student clearly plans to return to their home country.
- They have specific plans in their home country, such as opening a business or working there.
- They mention strong ties to home country (family, career, long term projects).
- They do NOT show interest in staying in the US to work or live long term.

Give LOW scores when:
- They talk about staying in the US long term, getting a job there, using OPT/CPT for long term career, or living in the US.
- They complain about or speak negatively about their home country (government, economy, education) as a reason to leave.
- They mention backup plans to study or move to multiple foreign countries (for example “If not the US, I will go to Canada, Switzerland, or somewhere else”) in a way that sounds like migration focus.

Score definitions:

5 = The answer clearly shows strong intent to return home. The student focuses on education in the US and then returning to their home country to work or start something (for example, opening a school or business). They speak positively or neutrally about their home country and do not mention working or staying in the US.

3 = The answer is mostly education focused and does not clearly say they will stay in the US. There may be small or vague references to the US or “opportunities” there, but no strong or explicit plan to stay long term. There is some uncertainty, but no clear migration risk.

1 = The answer clearly suggests migration risk. They openly talk about wanting to work, live, or stay in the US after study, or they speak negatively about their home country as a reason to leave. They might also mention multiple foreign country options in a way that looks like they want to leave their home country permanently.

Special case:
- If the question is purely factual and not about plans or intent (for example, “Who is paying your tuition?”, “Do you have any gaps?”), and the answer does not mention anything risky, give migration_intent = 5 by default.

2) goal_understanding (1 to 5)

This criterion measures how clearly the student connects their studies to their future goals.

Use this mainly when the question asks about:
- “Why US?”
- “Why this university?”
- “Why this major?”
- “Plans after graduation”
- “Future goals”

Give HIGH scores when:
- The answer clearly explains why they chose the US, this university, and this major.
- They connect the program and major to long-term goals, especially in their home country.
- The reasoning is logical and believable, not just memorized.

Give LOW scores when:
- They have no clear goals.
- They cannot explain why they chose the program, university, or country.
- There is no link between their studies and future plans.

Score definitions:

5 = Very clear and logical understanding of goals. The student explains why this major and university fit their future plan, especially in their home country. The answer mentions concrete goals (for example, opening a programming school in their country, contributing to a specific industry).

3 = The student knows their major and university and has some goals, but the explanation is generic or not very detailed. The connection between studies and goals exists but is not very strong or specific.

1 = The student does not show clear goals. They do not explain why they chose this major or university. There is no logical connection between their study plan and their future.

Special rule for factual questions:
If the question is factual (for example “Who is sponsoring you?”, “Do you have any gaps?”, “What is your university name?”) and the answer correctly gives the needed information, give goal_understanding = 5 as long as the answer is clear and appropriate. Do NOT penalize them for not mentioning future goals if the question did not ask for that.

3) answer_length (1 to 5)

This criterion measures whether the answer length fits the question.

Do NOT expect long answers for every question. Short and direct answers can be perfect, especially for simple factual questions.

Give HIGH scores when:
- The answer clearly and directly answers the question.
- It includes enough detail for understanding, but not unnecessary stories.
- For simple factual questions, a short, direct answer is fine.

Give LOW scores when:
- The answer is extremely short and misses important information.
- The answer is very long and goes off-topic with irrelevant details.

Score definitions:

5 = The length fits the question well. For complex questions (goals, plans, “why US”), the answer has 2–5 sentences with some details. For simple factual questions, the answer can be short and direct but still complete.

3 = Slightly too short or slightly too long, but the main information is there. The answer may miss a minor detail or include a little extra information that was not needed.

1 = Clearly too short or too long. Too short = the answer feels incomplete or vague. Too long = the answer contains a lot of irrelevant content and becomes unclear.

OUTPUT FORMAT:

You must return ONLY a JSON object with this exact structure and field names:

{
  "scores": {
    "migration_intent": 0,
    "goal_understanding": 0,
    "answer_length": 0,
    "total_score": 0
  },
  "classification": "",
  "feedback": {
    "overall": "",
    "by_criterion": {
      "migration_intent": "",
      "goal_understanding": "",
      "answer_length": ""
    },
    "improvements": []
  }
}

FILLING THE FIELDS:

- "scores.migration_intent": integer 1 to 5
- "scores.goal_understanding": integer 1 to 5
- "scores.answer_length": integer 1 to 5
- "scores.total_score": integer = sum of the three scores
- "classification": one of "Excellent", "Good", "Average", "Weak"
- "feedback.overall": 1–3 sentence summary
- "feedback.by_criterion.migration_intent": 1–2 sentences explaining why you gave that score, quoting short parts of the answer when helpful.
- "feedback.by_criterion.goal_understanding": 1–2 sentences explaining that score.
- "feedback.by_criterion.answer_length": 1–2 sentences explaining that score.
- "feedback.improvements": an array of 1–3 short strings with actionable tips on how to improve the answer next time.

OUTPUT RULES:
- Output JSON only, no markdown, no backticks.
- Do not add any other keys.
- Do not include explanations outside the JSON.
`

	return &VisaAnalyzer{
		apiKey:       apiKey,
		apiURL:       "https://api.openai.com/v1/chat/completions",
		systemPrompt: systemPrompt,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// AnalyzeAnswer analyzes a single answer and returns detailed feedback
func (va *VisaAnalyzer) AnalyzeAnswer(question, answer string) (*AnalysisResponse, error) {
	if va.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// Build session messages with system prompt (only once)
	sessionMessages := []GPTMessage{
		{
			Role:    "system",
			Content: va.systemPrompt,
		},
	}

	return va.callGPTAPI(sessionMessages, question, answer)
}

// AnalyzeAnswerWithSession analyzes an answer with full session context
// The system prompt is sent only once, then we append conversation history
func (va *VisaAnalyzer) AnalyzeAnswerWithSession(session *Session, question, answer string) (*AnalysisResponse, error) {
	if va.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// Start with system prompt (sent once per API call, but contains all rules)
	sessionMessages := []GPTMessage{
		{
			Role:    "system",
			Content: va.systemPrompt,
		},
	}

	// Add previous Q&A pairs from the session for context
	// These messages don't repeat the rules, just the conversation
	for _, prevAnswer := range session.Answers {
		sessionMessages = append(sessionMessages, GPTMessage{
			Role:    "user",
			Content: fmt.Sprintf("Question: %s\nStudent's Answer: %s", prevAnswer.QuestionText, prevAnswer.Text),
		})

		// Add assistant response if analysis exists
		if prevAnswer.Analysis != nil {
			analysisJSON, err := json.Marshal(prevAnswer.Analysis)
			if err == nil {
				sessionMessages = append(sessionMessages, GPTMessage{
					Role:    "assistant",
					Content: string(analysisJSON),
				})
			}
		}
	}

	return va.callGPTAPI(sessionMessages, question, answer)
}

// GetSessionMessages builds the full conversation history for a session
// This can be useful if you want to inspect what's being sent to the API
func (va *VisaAnalyzer) GetSessionMessages(session *Session) []GPTMessage {
	messages := []GPTMessage{
		{
			Role:    "system",
			Content: va.systemPrompt,
		},
	}

	for _, prevAnswer := range session.Answers {
		messages = append(messages, GPTMessage{
			Role:    "user",
			Content: fmt.Sprintf("Question: %s\nStudent's Answer: %s", prevAnswer.QuestionText, prevAnswer.Text),
		})

		if prevAnswer.Analysis != nil {
			analysisJSON, err := json.Marshal(prevAnswer.Analysis)
			if err == nil {
				messages = append(messages, GPTMessage{
					Role:    "assistant",
					Content: string(analysisJSON),
				})
			}
		}
	}

	return messages
}

// GenerateSessionSummary generates a summary from multiple analysis records
func (va *VisaAnalyzer) GenerateSessionSummary(analyses []AnalysisRecord) (*SessionSummary, error) {
	if len(analyses) == 0 {
		return nil, fmt.Errorf("no analyses provided")
	}

	totalScore := 0
	for _, record := range analyses {
		totalScore += record.Analysis.Scores.TotalScore
	}

	avgScore := float64(totalScore) / float64(len(analyses))

	return &SessionSummary{
		TotalQuestions: len(analyses),
		AverageScore:   avgScore,
		OverallGrade:   getGradeFromScore(int(avgScore)),
		StrongAreas:    extractCommonStrengths(analyses),
		WeakAreas:      extractCommonWeaknesses(analyses),
		CommonRedFlags: extractCommonRedFlags(analyses),
		Recommendation: generateRecommendation(avgScore, analyses),
		CompletedAt:    time.Now(),
	}, nil
}

// GPTMessage represents a message in the GPT conversation
type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (va *VisaAnalyzer) callGPTAPI(sessionMessages []GPTMessage, question, answer string) (*AnalysisResponse, error) {
	type GPTRequest struct {
		Model       string       `json:"model"`
		MaxTokens   int          `json:"max_tokens"`
		Messages    []GPTMessage `json:"messages"`
		Temperature float64      `json:"temperature"`
	}

	type GPTChoice struct {
		Message GPTMessage `json:"message"`
	}

	type GPTResponse struct {
		Choices []GPTChoice `json:"choices"`
	}

	// Add the new question and answer to the session messages
	// This is just the Q&A content, not the rules
	sessionMessages = append(sessionMessages, GPTMessage{
		Role:    "user",
		Content: fmt.Sprintf("Question: %s\nStudent's Answer: %s", question, answer),
	})

	gptReq := GPTRequest{
		Model:       "gpt-3.5-turbo",
		MaxTokens:   1000,
		Temperature: 0.3,
		Messages:    sessionMessages,
	}

	reqBody, err := json.Marshal(gptReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", va.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+va.apiKey)

	resp, err := va.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var gptResp GPTResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(gptResp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	content := gptResp.Choices[0].Message.Content
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var analysis AnalysisResponse
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}

	return &analysis, nil
}

// Helper functions for session summary generation

func getGradeFromScore(score int) string {
	switch {
	case score == 15:
		return "A" // Excellent: 15
	case score >= 13:
		return "B" // Good: 13-14
	case score >= 11:
		return "C" // Average: 11-12
	default:
		return "D" // Weak: 3-10
	}
}

func generateRecommendation(avgScore float64, analyses []AnalysisRecord) string {
	if avgScore >= 15 {
		return "Excellent performance! You're well-prepared. Focus on maintaining confidence and natural delivery during the actual interview."
	} else if avgScore >= 13 {
		return "Good foundation. Review the specific feedback for each answer and practice the improved versions. Focus on being more specific and confident in your responses."
	} else if avgScore >= 11 {
		return "You need more practice. Focus on providing specific examples, showing strong ties to your home country, and demonstrating clear post-graduation plans."
	}
	return "Significant improvement needed. Consider working with an advisor to strengthen your answers. Focus on clarity, specificity, and addressing visa officer concerns about immigrant intent."
}

// ScoreToPercentage converts 3-15 scale to 0-100 percentage for display
func ScoreToPercentage(score int) float64 {
	if score < 3 {
		score = 3
	}
	if score > 15 {
		score = 15
	}
	// Formula: ((score - 3) / 12) * 100
	return float64(score-3) * (100.0 / 12.0)
}

func extractCommonStrengths(analyses []AnalysisRecord) []string {
	criteriaScores := make(map[string]int)
	criteriaCount := make(map[string]int)

	for _, record := range analyses {
		scores := record.Analysis.Scores

		if scores.MigrationIntent >= 4 {
			criteriaScores["migration_intent"] += scores.MigrationIntent
			criteriaCount["migration_intent"]++
		}
		if scores.GoalUnderstanding >= 4 {
			criteriaScores["goal_understanding"] += scores.GoalUnderstanding
			criteriaCount["goal_understanding"]++
		}
		if scores.AnswerLength >= 4 {
			criteriaScores["answer_length"] += scores.AnswerLength
			criteriaCount["answer_length"]++
		}
	}

	var strengths []string
	for criterion, count := range criteriaCount {
		if count >= len(analyses)/2 {
			strengths = append(strengths, formatCriterionName(criterion))
		}
	}

	return strengths
}

func extractCommonWeaknesses(analyses []AnalysisRecord) []string {
	criteriaScores := make(map[string]int)

	for _, record := range analyses {
		scores := record.Analysis.Scores

		if scores.MigrationIntent <= 3 {
			criteriaScores["migration_intent"]++
		}
		if scores.GoalUnderstanding <= 3 {
			criteriaScores["goal_understanding"]++
		}
		if scores.AnswerLength <= 3 {
			criteriaScores["answer_length"]++
		}
	}

	var weaknesses []string
	for criterion, count := range criteriaScores {
		if count >= len(analyses)/2 {
			weaknesses = append(weaknesses, formatCriterionName(criterion))
		}
	}

	return weaknesses
}

func extractCommonRedFlags(analyses []AnalysisRecord) []string {
	flagMap := make(map[string]bool)

	for _, record := range analyses {
		scores := record.Analysis.Scores

		if scores.MigrationIntent <= 2 {
			flagMap["Shows potential immigration intent"] = true
		}
		if scores.GoalUnderstanding <= 2 {
			flagMap["Unclear academic goals"] = true
		}
		if scores.AnswerLength <= 2 {
			flagMap["Poor answer structure or length"] = true
		}
	}

	var flags []string
	for flag := range flagMap {
		flags = append(flags, flag)
	}

	return flags
}

func formatCriterionName(criterion string) string {
	switch criterion {
	case "migration_intent":
		return "No immigration intent"
	case "goal_understanding":
		return "Clear understanding of academic goals"
	case "answer_length":
		return "Appropriate answer length"
	default:
		return criterion
	}
}
