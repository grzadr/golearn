package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/grzadr/golearn/questions/question"
)

func main() {
	questionFile := filepath.Join("data", "questions.json")
	questions, err := question.LoadQuestions(questionFile)
	if err != nil {
		log.Fatal("Error loading questions:", err)
	}

	fmt.Println("Quiz loaded successfully!")
	fmt.Printf("Number of questions: %d\n", len(questions))

	// We'll add the quiz logic here in the next step
}
