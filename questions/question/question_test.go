package question

import (
	"errors"
	"os"
	"testing"
)

func TestLoadQuestionsMissingFile(t *testing.T) {
	_, err := LoadQuestions("non_existent_file.json")
	if err == nil {
		t.Fatal("Expected an error for missing file, but got nil")
	}

	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("Expected error type *os.PathError, got %T", err)
	}

	if pathErr.Op != "open" {
		t.Errorf("Expected PathError.Op to be 'open', got '%s'", pathErr.Op)
	}
	if pathErr.Path != "non_existent_file.json" {
		t.Errorf("Expected PathError.Path to be 'non_existent_file.json', got '%s'", pathErr.Path)
	}
}

func TestLoadQuestionsMalformedJSON(t *testing.T) {
	malformedJSON := []byte(`{"this is": "not valid JSON"`)
	tmpfile, err := os.CreateTemp("", "malformed*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(malformedJSON); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadQuestions(tmpfile.Name())
	if err == nil {
		t.Error("Expected an error for malformed JSON, but got nil")
	}
}

func TestLoadQuestionsValidJSON(t *testing.T) {
	validJSON := []byte(`[
		{
			"question": "What is the capital of France?",
			"options": ["London", "Berlin", "Paris", "Madrid"],
			"correct_answer": 2
		}
	]`)
	tmpfile, err := os.CreateTemp("", "valid*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(validJSON); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	questions, err := LoadQuestions(tmpfile.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(questions))
	}
	if questions[0].QuestionText != "What is the capital of France?" {
		t.Errorf("Unexpected question text: %s", questions[0].QuestionText)
	}
}
