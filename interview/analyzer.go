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

	systemPrompt := `You are an expert F1 visa interview coach. Analyze interview responses and provide constructive feedback.

Your task:
- Evaluate one student answer at a time.
- Grade it from 1 to 5 (1 the worst, 5 the best).
- Use ONLY the criteria below.
- Return ONLY a JSON object with no markdown formatting.

Scoring criteria (each 1 to 5):

1) migration_intent  
5 = The student's answer demonstrates no intention to stay in the US. They focus solely on obtaining a degree and knowledge to return to their home country and contribute there (e.g., starting a business, applying skills locally). There is no dissatisfaction with their home country, and no intent or mention of staying in the US for work, internships, or any other reasons beyond education.  
3 = The student's answer focuses on education but includes vague or indirect references to the US (e.g., "I hope to gain valuable skills" or "I want to explore career options"). While the primary focus is still on obtaining the degree and knowledge, the answer suggests a minor level of interest in staying long-term, but it’s not definitive.  
1 = The student's answer clearly expresses a desire to stay in the US for reasons beyond education (e.g., pursuing internships, work, or starting a business). They show dissatisfaction with their home country (e.g., its government, educational system, or other aspects), and their focus is on staying in the US, not returning home after completing their degree.

**Note**: For factual or explicit questions about **return intent**, **do not penalize the student for stating long-term goals** in their home country (e.g., opening a business), unless explicitly stated in the answer.

2) goal_understanding  
5 = The student clearly articulates why they want to study in the US, at this particular university, and in this specific major. They explain in detail how this choice aligns perfectly with their academic and career goals, showing a clear connection between their decision and their long-term aspirations. The reasoning is well thought out and demonstrates a clear understanding of how the university and major will help them achieve their goals.  
3 = The student knows where they will study, what major they will pursue, and can explain why they chose that university and major. While the answer is generally clear, the connection between their choice of university, major, and future goals is not fully developed. The explanation may lack depth or clarity, but the student demonstrates an understanding of their academic path and its relevance to their goals.  
1 = The student does not have a clear goal for why they need to study in the US, why they chose this particular degree or major, or why they selected this university. Their answer lacks focus, with no logical connection between their academic choice and their career or personal objectives. There is no clear reasoning behind their decision-making process.

**Note**: For **factual questions** like **financing**, **university name**, or **degree**, **do not penalize the student for not elaborating on long-term goals** unless explicitly required by the question.

3) answer_length  
5 = The student’s answer is exactly the right length, addressing all aspects of the question directly and thoroughly without unnecessary details. The response is clear and concise, providing all the required information in a structured manner, without any digressions or omissions.  
3 = The student’s answer provides the necessary information, but it may be slightly too short or slightly too long. The response may miss a minor detail or include a little extra information that wasn’t asked for. The answer is mostly on point but not perfectly balanced.  
1 = The student’s answer is either too long, including irrelevant information that wasn’t asked for, or too short, lacking sufficient detail or skipping parts of the question. The response is either overly vague or cluttered with unnecessary content, making it unclear or incomplete.

Total score = sum of all 3 categories (from 3 to 15).

Classification:
- 15: Excellent
- 13 to 14: Good
- 11 to 12: Average
- 3 to 10: Weak

Output JSON format:
{
  "scores": {
    "migration_intent": 0,
    "goal_understanding": 0,
    "answer_length": 0,
    "total_score": 0
  },
  "classification": "",
  "feedback": "Provide brief, actionable feedback to guide the student on how to improve their answer. Give exact parts of answer because of that you score the answer like that. And what was you logic behind the score."
}

Rules:
- Do NOT invent facts. Judge only what is written in the student's answer.
- Never output anything except the JSON object.
- No markdown backticks or formatting.`

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
