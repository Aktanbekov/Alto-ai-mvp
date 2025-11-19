package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Question represents an interview question with follow-ups and evaluation criteria
type Question struct {
	Text       string
	Category   string
	FollowUps  []string
	GoodFlags  []string // Keywords that indicate a good answer
	RedFlags   []string // Keywords that indicate red flags
}

// UserProfile stores information about the applicant
type UserProfile struct {
	VisaType      string
	PurposeTravel string
	HasFamily     bool
	Answers       map[string]string
}

// InterviewSession manages the practice session
type InterviewSession struct {
	Profile   UserProfile
	Questions []Question
	Score     int
	Feedback  []string
}

func main() {
	fmt.Println("=== Visa Interview Practice App ===\n")
	
	session := NewInterviewSession()
	session.CollectProfile()
	session.ConductInterview()
	session.ShowResults()
}

func NewInterviewSession() *InterviewSession {
	return &InterviewSession{
		Profile:   UserProfile{Answers: make(map[string]string)},
		Questions: initializeQuestions(),
		Feedback:  []string{},
	}
}

func initializeQuestions() []Question {
	return []Question{
		{
			Text:      "What is the purpose of your trip to the United States?",
			Category:  "purpose",
			FollowUps: []string{"How long do you plan to stay?", "Where will you be staying?"},
			GoodFlags: []string{"tourism", "visit", "business", "conference", "specific"},
			RedFlags:  []string{"maybe", "not sure", "find job", "stay permanently"},
		},
		{
			Text:      "Do you have family or friends in the United States?",
			Category:  "ties",
			FollowUps: []string{"What is their immigration status?", "When did you last see them?"},
			GoodFlags: []string{"yes", "no", "citizen", "specific"},
			RedFlags:  []string{"illegal", "overstayed", "don't know"},
		},
		{
			Text:      "What do you do for work?",
			Category:  "financial",
			FollowUps: []string{"How long have you worked there?", "What is your salary?"},
			GoodFlags: []string{"years", "company", "manager", "stable"},
			RedFlags:  []string{"unemployed", "just started", "looking for"},
		},
		{
			Text:      "Why will you return to your home country?",
			Category:  "ties",
			FollowUps: []string{"Do you own property?", "What family do you have there?"},
			GoodFlags: []string{"job", "family", "property", "business", "children"},
			RedFlags:  []string{"nothing", "don't know", "maybe won't"},
		},
	}
}

func (s *InterviewSession) CollectProfile() {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("What type of visa are you applying for? (B1/B2, F1, etc.): ")
	visaType, _ := reader.ReadString('\n')
	s.Profile.VisaType = strings.TrimSpace(visaType)
	
	fmt.Print("Purpose of travel (tourism/business/study): ")
	purpose, _ := reader.ReadString('\n')
	s.Profile.PurposeTravel = strings.TrimSpace(purpose)
	
	fmt.Print("Do you have family in the US? (yes/no): ")
	family, _ := reader.ReadString('\n')
	s.Profile.HasFamily = strings.TrimSpace(strings.ToLower(family)) == "yes"
	
	fmt.Println("\n--- Profile created! Starting interview practice ---\n")
}

func (s *InterviewSession) ConductInterview() {
	reader := bufio.NewReader(os.Stdin)
	
	// Select relevant questions based on profile
	selectedQuestions := s.selectQuestions()
	
	for i, q := range selectedQuestions {
		fmt.Printf("Question %d: %s\n", i+1, q.Text)
		fmt.Print("Your answer: ")
		
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		s.Profile.Answers[q.Category] = answer
		
		// Evaluate answer
		feedback := s.evaluateAnswer(answer, q)
		s.Feedback = append(s.Feedback, feedback)
		fmt.Printf("💬 %s\n\n", feedback)
		
		// Ask follow-up based on answer
		if len(q.FollowUps) > 0 && len(answer) > 10 {
			followUp := s.selectFollowUp(q, answer)
			if followUp != "" {
				fmt.Printf("Follow-up: %s\n", followUp)
				fmt.Print("Your answer: ")
				followUpAnswer, _ := reader.ReadString('\n')
				followUpAnswer = strings.TrimSpace(followUpAnswer)
				
				followUpFeedback := s.evaluateFollowUp(followUpAnswer)
				s.Feedback = append(s.Feedback, followUpFeedback)
				fmt.Printf("💬 %s\n\n", followUpFeedback)
			}
		}
	}
}

func (s *InterviewSession) selectQuestions() []Question {
	// Pattern: Select questions based on user profile
	selected := []Question{}
	
	for _, q := range s.Questions {
		// Always ask purpose and ties questions
		if q.Category == "purpose" || q.Category == "ties" {
			selected = append(selected, q)
		}
		// Ask financial questions for tourist visas
		if q.Category == "financial" && strings.Contains(strings.ToLower(s.Profile.VisaType), "b") {
			selected = append(selected, q)
		}
	}
	
	return selected
}

func (s *InterviewSession) evaluateAnswer(answer string, q Question) string {
	answerLower := strings.ToLower(answer)
	
	// Check for red flags first
	for _, flag := range q.RedFlags {
		if strings.Contains(answerLower, strings.ToLower(flag)) {
			s.Score -= 10
			return fmt.Sprintf("⚠️ RED FLAG: Avoid mentioning '%s'. This raises concerns about your intentions.", flag)
		}
	}
	
	// Check for good indicators
	goodCount := 0
	for _, flag := range q.GoodFlags {
		if strings.Contains(answerLower, strings.ToLower(flag)) {
			goodCount++
		}
	}
	
	// Evaluate answer quality
	if len(answer) < 10 {
		return "❌ Too brief. Provide more specific details to build credibility."
	}
	
	if goodCount >= 2 {
		s.Score += 20
		return "✅ Good answer! Specific and confident."
	} else if goodCount == 1 {
		s.Score += 10
		return "👍 Decent answer, but could be more specific with dates, names, or details."
	}
	
	s.Score += 5
	return "⚡ Okay, but try to be more concrete. Use specific examples."
}

func (s *InterviewSession) selectFollowUp(q Question, answer string) string {
	// Pattern: Select follow-up based on answer content
	if len(q.FollowUps) == 0 {
		return ""
	}
	
	answerLower := strings.ToLower(answer)
	
	// Smart follow-up selection
	if strings.Contains(answerLower, "family") || strings.Contains(answerLower, "friend") {
		for _, fu := range q.FollowUps {
			if strings.Contains(strings.ToLower(fu), "status") || strings.Contains(strings.ToLower(fu), "see them") {
				return fu
			}
		}
	}
	
	if strings.Contains(answerLower, "stay") || strings.Contains(answerLower, "visit") {
		for _, fu := range q.FollowUps {
			if strings.Contains(strings.ToLower(fu), "long") {
				return fu
			}
		}
	}
	
	// Default: return first follow-up
	return q.FollowUps[0]
}

func (s *InterviewSession) evaluateFollowUp(answer string) string {
	if len(answer) < 5 {
		return "❌ Too vague. Consular officers want specific information."
	}
	
	// Pattern: Check for specificity indicators
	hasNumbers := strings.ContainsAny(answer, "0123456789")
	hasSpecificWords := strings.Contains(strings.ToLower(answer), "years") ||
		strings.Contains(strings.ToLower(answer), "months") ||
		strings.Contains(strings.ToLower(answer), "citizen")
	
	if hasNumbers || hasSpecificWords {
		s.Score += 15
		return "✅ Excellent! Specific details strengthen your case."
	}
	
	s.Score += 5
	return "👍 Acceptable, but more specificity would help."
}

func (s *InterviewSession) ShowResults() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("INTERVIEW PRACTICE COMPLETE")
	fmt.Println(strings.Repeat("=", 50))
	
	fmt.Printf("\n📊 Your Score: %d/100\n\n", s.Score)
	
	if s.Score >= 70 {
		fmt.Println("🎉 Great job! You're well-prepared for your interview.")
	} else if s.Score >= 40 {
		fmt.Println("⚡ Good start, but practice more to improve confidence.")
	} else {
		fmt.Println("⚠️ Needs improvement. Focus on being more specific and confident.")
	}
	
	fmt.Println("\n📝 Key Takeaways:")
	fmt.Println("1. Be specific: Use dates, numbers, and concrete details")
	fmt.Println("2. Show ties: Emphasize reasons to return home")
	fmt.Println("3. Be consistent: Don't contradict yourself")
	fmt.Println("4. Be confident: Hesitation raises red flags")
	fmt.Println("5. Be honest: Don't exaggerate or lie")
}