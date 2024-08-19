package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/grzadr/golearn/questions/question"
)

const (
	timeLimit   = 30 // seconds per question
	hintPenalty = 5  // seconds deducted for using a hint
)

func main() {
	questionFile := filepath.Join("data", "questions.json")
	questions, err := question.LoadQuestions(questionFile)
	if err != nil {
		log.Fatal("Error loading questions:", err)
	}

	// Randomize question order
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	fmt.Println("Welcome to the Enhanced Quiz Game!")
	fmt.Printf("You will be asked %d questions. You have %d seconds for each question.\n", len(questions), timeLimit)
	fmt.Println("Type 'hint' to get a hint (but you'll lose 5 seconds!)")
	fmt.Println("Press Enter to start the quiz...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	score := 0
	for i, q := range questions {
		if askQuestion(i+1, q) {
			score++
		}
	}

	fmt.Printf("\nQuiz completed! Your score: %d out of %d\n", score, len(questions))
}

func askQuestion(num int, q question.Question) bool {
	fmt.Printf("\nQuestion %d: %s\n", num, q.QuestionText)
	for i, option := range q.Options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	answerCh := make(chan string)

	go func() {
		defer close(answerCh) // Ensure the channel is closed when the goroutine exits
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Your answer (enter the number or 'hint'): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading input: %v", err)
				return
			}
			answerCh <- strings.TrimSpace(input)

			// Check if the timer has expired
			select {
			case <-timer.C:
				return // Exit the goroutine if the timer has expired
			default:
				// Continue if the timer hasn't expired
			}
		}
	}()

	var correct bool
	remainingTime := timeLimit

	for {
		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			return false
		case input, ok := <-answerCh:
			if !ok {
				// Channel closed, should not happen unless there's an error
				return false
			}
			if input == "hint" {
				if remainingTime > hintPenalty {
					remainingTime -= hintPenalty
					timer.Reset(time.Duration(remainingTime) * time.Second)
					fmt.Printf("Hint: The answer is not %d\n", getIncorrectOption(q))
					continue
				} else {
					fmt.Println("Not enough time for a hint!")
					continue
				}
			}

			answerNum, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input. Please enter a number or 'hint'.")
				continue
			}

			answer := answerNum - 1
			correct = (answer == q.CorrectAnswer)
			timer.Stop()

			if correct {
				fmt.Println("Correct!")
			} else {
				fmt.Printf("Incorrect. The correct answer was: %s\n", q.Options[q.CorrectAnswer])
			}
			return correct
		}
	}
}

func getIncorrectOption(q question.Question) int {
	for {
		option := rand.Intn(len(q.Options))
		if option != q.CorrectAnswer {
			return option + 1
		}
	}
}
