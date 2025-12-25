package interview

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// QuestionsByCategory stores questions organized by category
var QuestionsByCategory map[string][]string

// QuestionSelectionRules defines how many questions to ask from each category
var QuestionSelectionRules = map[string]int{
	"Purpose of Study":       2,
	"Academic Background":    2,
	"University Choice":      2,
	"Financial Capability":   2,
	"Family/Sponsor Info":    1,
	"Post-Graduation Plans":  2,
	"Immigration Intent":     1,
}

// CategoryOrder defines the order in which categories should be asked
var CategoryOrder = []string{
	"Purpose of Study",
	"Academic Background",
	"University Choice",
	"Financial Capability",
	"Family/Sponsor Info",
	"Post-Graduation Plans",
	"Immigration Intent",
}

func LoadQuestions(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read questions file: %w", err)
	}

	var categories map[string][]string
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("unmarshal questions: %w", err)
	}

	QuestionsByCategory = make(map[string][]string)
	for category, questions := range categories {
		QuestionsByCategory[category] = questions
	}

	// Validate that all required categories exist
	for category := range QuestionSelectionRules {
		if _, ok := QuestionsByCategory[category]; !ok {
			return fmt.Errorf("required category '%s' not found in questions file", category)
		}
	}

	return nil
}

// SelectQuestionsForSession selects questions according to the rules
func SelectQuestionsForSession() []Question {
	var selectedQuestions []Question
	rand.Seed(time.Now().UnixNano())

	for _, category := range CategoryOrder {
		count, ok := QuestionSelectionRules[category]
		if !ok {
			continue
		}

		questions, ok := QuestionsByCategory[category]
		if !ok || len(questions) == 0 {
			continue
		}

		// Select random questions from this category
		available := make([]string, len(questions))
		copy(available, questions)
		
		// Shuffle and take the required count
		rand.Shuffle(len(available), func(i, j int) {
			available[i], available[j] = available[j], available[i]
		})

		// Take up to 'count' questions
		toTake := count
		if toTake > len(available) {
			toTake = len(available)
		}

		for i := 0; i < toTake; i++ {
			questionID := fmt.Sprintf("q%d_%s", len(selectedQuestions)+1, sanitizeCategory(category))
			selectedQuestions = append(selectedQuestions, Question{
				ID:       questionID,
				Category: category,
				Text:     available[i],
			})
		}
	}

	return selectedQuestions
}

// sanitizeCategory converts category name to a valid ID suffix
func sanitizeCategory(category string) string {
	// Simple sanitization - replace spaces and special chars
	result := ""
	for _, char := range category {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			result += string(char)
		} else if char == ' ' || char == '/' {
			result += "_"
		}
	}
	return result
}
