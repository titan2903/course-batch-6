package domain

import (
	"errors"
	"time"
)

type Exercise struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type Question struct {
	ID            int       `json:"id"`
	ExerciseID    int       `json:"excercise_id"`
	Body          string    `json:"body"`
	OptionA       string    `json:"option_a"`
	OptionB       string    `json:"option_b"`
	OptionC       string    `json:"option_c"`
	OptionD       string    `json:"option_d"`
	CorrectAnswer string    `json:"correct_answer"`
	Score         int       `json:"score"`
	CreatorID     int       `json:"creator_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Answer struct {
	ID         int       `json:"id"`
	ExerciseID int       `json:"exercise_id"`
	QuestionID int       `json:"question_id"`
	UserID     int       `json:"user_id"`
	Answer     string    `json:"answer"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ExerciseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type QuestionRequest struct {
	Body          string `json:"body"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectAnswer string `json:"correct_answer"`
}

type AnswerRequest struct {
	Answer string `json:"answer"`
}

func NewExercise(title, description string) (exerciseRequest *ExerciseRequest, err error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if description == "" {
		return nil, errors.New("description is required")
	}

	exerciseRequest = &ExerciseRequest{
		Title:       title,
		Description: description,
	}

	return exerciseRequest, err
}

func NewQuestion(body, optionA, optionB, optionC, optionD, correctAnswer string) (questionRequest *QuestionRequest, err error) {
	if body == "" {
		return nil, errors.New("body is required")
	}

	if optionA == "" || optionB == "" || optionC == "" || optionD == "" {
		return nil, errors.New("have to choose one of the multiple choices")
	}

	if correctAnswer == "" {
		return nil, errors.New("correctAnswer is required")
	}

	questionRequest = &QuestionRequest{
		Body:          body,
		OptionA:       optionA,
		OptionB:       optionB,
		OptionC:       optionC,
		OptionD:       optionD,
		CorrectAnswer: correctAnswer,
	}

	return questionRequest, err
}

func NewAnswer(answer string) (answerRequest *AnswerRequest, err error) {
	if answer == "" {
		return nil, errors.New("answer is required")
	}

	answerRequest = &AnswerRequest{
		Answer: answer,
	}

	return answerRequest, err
}
