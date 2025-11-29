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
	return &VisaAnalyzer{
		apiKey: apiKey,
		apiURL: "https://api.openai.com/v1/chat/completions",
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

	prompt := fmt.Sprintf(`You are an expert F1 visa interview coach. Analyze this interview response and provide constructive feedback.

Question: %s

Student's Answer: %s 

Your task:
- Evaluate one student answer at a time.
- Grade it from 1 to 5 (1 the worst, 5 the best).
- Use ONLY the criteria below.
- Return ONLY a JSON object with no markdown formatting.

Scoring criteria (each 1 to 5):

1) goal_understanding  
5 = student clearly understands the academic goal and reason for study  
3 = somewhat understands the goal  
1 = no clear goal or confusion

2) logical_mindset  
5 = answer is logical and fits a realistic study plan  
3 = partly logical with small issues  
1 = illogical, confused, or contradictory

3) no_migration_intent  
5 = clearly no immigration intent  
3 = slightly unclear  
1 = suggests wanting to stay in the US, work long term, or shows signs of intent

4) no_hate_to_home_country  
5 = respects home country, neutral or positive tone  
3 = slightly negative or weak respect  
1 = shows hate or rejection of home country

5) answer_quality  
5 = direct, clear, enough detail, no unnecessary long story  
3 = a bit too short or slightly too long  
1 = very unclear, too short, or includes irrelevant long information

Total score = sum of all 5 categories (from 5 to 25).

Classification:
- 22 to 25: Excellent
- 17 to 21: Good
- 12 to 16: Average
- 8 to 11: Weak
- 5 to 7: Poor

Output JSON format:
{
  "scores": {
    "goal_understanding": 0,
    "logical_mindset": 0,
    "no_migration_intent": 0,
    "no_hate_to_home_country": 0,
    "answer_quality": 0,
    "total_score": 0
  },
  "classification": "",
  "feedback": "1 to 2 short sentences with suggestions for improvement."
}

Rules:
- Do NOT invent facts. Judge only what is written.
- Never output anything except the JSON object.
- No markdown backticks or formatting.`, question, answer)

	return va.callGPTAPI(prompt)
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

func (va *VisaAnalyzer) callGPTAPI(prompt string) (*AnalysisResponse, error) {
	type GPTMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

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

	gptReq := GPTRequest{
		Model:       "gpt-4o-mini", // Better quality than 3.5-turbo
		MaxTokens:   1000,
		Temperature: 0.3, // Lower for more consistent grading
		Messages: []GPTMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
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

// Helper functions for session summary generation (FIXED for 5-25 scale)

func getGradeFromScore(score int) string {
	// Score range: 5-25
	switch {
	case score >= 22:
		return "A" // Excellent: 22-25
	case score >= 17:
		return "B" // Good: 17-21
	case score >= 12:
		return "C" // Average: 12-16
	case score >= 8:
		return "D" // Weak: 8-11
	default:
		return "F" // Poor: 5-7
	}
}

func generateRecommendation(avgScore float64, analyses []AnalysisRecord) string {
	// avgScore range: 5.0-25.0
	if avgScore >= 22 {
		return "Excellent performance! You're well-prepared. Focus on maintaining confidence and natural delivery during the actual interview."
	} else if avgScore >= 17 {
		return "Good foundation. Review the specific feedback for each answer and practice the improved versions. Focus on being more specific and confident in your responses."
	} else if avgScore >= 12 {
		return "You need more practice. Focus on providing specific examples, showing strong ties to your home country, and demonstrating clear post-graduation plans."
	} else if avgScore >= 8 {
		return "Significant improvement needed. Consider working with an advisor to strengthen your answers. Focus on clarity, specificity, and addressing visa officer concerns about immigrant intent."
	}
	return "Major revision required. Your answers show fundamental issues with clarity, goals, or immigration intent. Please work with a visa consultant before attempting the interview."
}

// ScoreToPercentage converts 5-25 scale to 0-100 percentage for display
func ScoreToPercentage(score int) float64 {
	if score < 5 {
		score = 5
	}
	if score > 25 {
		score = 25
	}
	// Formula: ((score - 5) / 20) * 100
	return float64(score-5) * 5.0
}

func extractCommonStrengths(analyses []AnalysisRecord) []string {
	// Count high-scoring criteria across all analyses
	criteriaScores := make(map[string]int)
	criteriaCount := make(map[string]int)
	
	for _, record := range analyses {
		scores := record.Analysis.Scores
		
		// Track each criterion
		if scores.GoalUnderstanding >= 4 {
			criteriaScores["goal_understanding"] += scores.GoalUnderstanding
			criteriaCount["goal_understanding"]++
		}
		if scores.LogicalMindset >= 4 {
			criteriaScores["logical_mindset"] += scores.LogicalMindset
			criteriaCount["logical_mindset"]++
		}
		if scores.NoMigrationIntent >= 4 {
			criteriaScores["no_migration_intent"] += scores.NoMigrationIntent
			criteriaCount["no_migration_intent"]++
		}
		if scores.NoHateToHomeCountry >= 4 {
			criteriaScores["no_hate_to_home_country"] += scores.NoHateToHomeCountry
			criteriaCount["no_hate_to_home_country"]++
		}
		if scores.AnswerQuality >= 4 {
			criteriaScores["answer_quality"] += scores.AnswerQuality
			criteriaCount["answer_quality"]++
		}
	}

	var strengths []string
	for criterion, count := range criteriaCount {
		if count >= len(analyses)/2 { // Appears in at least half with high score
			strengths = append(strengths, formatCriterionName(criterion))
		}
	}
	
	return strengths
}

func extractCommonWeaknesses(analyses []AnalysisRecord) []string {
	// Count low-scoring criteria across all analyses
	criteriaScores := make(map[string]int)
	
	for _, record := range analyses {
		scores := record.Analysis.Scores
		
		// Track each criterion that's weak (score <= 3)
		if scores.GoalUnderstanding <= 3 {
			criteriaScores["goal_understanding"]++
		}
		if scores.LogicalMindset <= 3 {
			criteriaScores["logical_mindset"]++
		}
		if scores.NoMigrationIntent <= 3 {
			criteriaScores["no_migration_intent"]++
		}
		if scores.NoHateToHomeCountry <= 3 {
			criteriaScores["no_hate_to_home_country"]++
		}
		if scores.AnswerQuality <= 3 {
			criteriaScores["answer_quality"]++
		}
	}

	var weaknesses []string
	for criterion, count := range criteriaScores {
		if count >= len(analyses)/2 { // Appears in at least half with low score
			weaknesses = append(weaknesses, formatCriterionName(criterion))
		}
	}
	
	return weaknesses
}

func extractCommonRedFlags(analyses []AnalysisRecord) []string {
	flagMap := make(map[string]bool)
	
	for _, record := range analyses {
		scores := record.Analysis.Scores
		
		// Critical red flags
		if scores.NoMigrationIntent <= 2 {
			flagMap["Shows potential immigration intent"] = true
		}
		if scores.NoHateToHomeCountry <= 2 {
			flagMap["Negative attitude towards home country"] = true
		}
		if scores.GoalUnderstanding <= 2 {
			flagMap["Unclear academic goals"] = true
		}
		if scores.LogicalMindset <= 2 {
			flagMap["Illogical or contradictory statements"] = true
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
	case "goal_understanding":
		return "Clear understanding of academic goals"
	case "logical_mindset":
		return "Logical and coherent responses"
	case "no_migration_intent":
		return "Strong ties to home country"
	case "no_hate_to_home_country":
		return "Positive attitude about home country"
	case "answer_quality":
		return "Well-structured and clear answers"
	default:
		return criterion
	}
}