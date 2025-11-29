package interview

import (
	"encoding/json"
	"fmt"
	"os"
)

var Questions map[string]Question

func LoadQuestions(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read questions file: %w", err)
	}

	var list []Question
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("unmarshal questions: %w", err)
	}

	Questions = make(map[string]Question, len(list))
	for _, q := range list {
		if q.ID == "" {
			return fmt.Errorf("question with empty id: %+v", q)
		}
		Questions[q.ID] = q
	}

	// Simple check for missing next_id targets
	for _, q := range Questions {
		if q.NextID != "" && q.NextID != "end" {
			if _, ok := Questions[q.NextID]; !ok {
				return fmt.Errorf("question %s has next_id %s which does not exist", q.ID, q.NextID)
			}
		}
	}

	return nil
}
