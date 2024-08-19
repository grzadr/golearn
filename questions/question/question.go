package question

import (
	"encoding/json"
	"os"
)

type Question struct {
	QuestionText  string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correct_answer"`
}

func LoadQuestions(filename string) ([]Question, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var questions []Question
	err = json.Unmarshal(data, &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}
