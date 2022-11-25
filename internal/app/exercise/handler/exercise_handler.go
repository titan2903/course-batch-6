package handler

import (
	"course-batch-6/internal/app/domain"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseHandler struct {
	db *gorm.DB
}

func NewExerciseHandler(db *gorm.DB) *ExerciseHandler {
	return &ExerciseHandler{
		db: db,
	}
}

func (eh ExerciseHandler) GetExerciseByID(c *gin.Context) {
	idString := c.Param("exerciseId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var exercise domain.Exercise
	err = eh.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	c.JSON(http.StatusOK, exercise)
}

func (eh ExerciseHandler) GetScore(c *gin.Context) {
	idString := c.Param("exerciseId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var exercise domain.Exercise
	err = eh.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	userID := c.Request.Context().Value("user_id").(int)

	var answers []domain.Answer
	err = eh.db.Where("exercise_id = ? AND user_id = ?", id, userID).Find(&answers).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "not answere yet",
		})
		return
	}

	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score Score
	wg := new(sync.WaitGroup)
	for _, question := range exercise.Questions {
		wg.Add(1)
		go func(question domain.Question) {
			defer wg.Done()
			if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
				score.Inc(question.Score)
			}
		}(question)
	}

	wg.Wait()

	c.JSON(http.StatusOK, map[string]int{
		"score": score.total,
	})
}

func (eh ExerciseHandler) CreateExercise(c *gin.Context) {
	var exerciseBody domain.ExerciseRequest
	if err := c.ShouldBind(&exerciseBody); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid body",
		})
		return
	}

	exerciseRequest, err := domain.NewExercise(exerciseBody.Title, exerciseBody.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	exercise := &domain.Exercise{
		Title:       exerciseRequest.Title,
		Description: exerciseRequest.Description,
	}

	if err := eh.db.Create(exercise).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}

func (eh ExerciseHandler) CreateQuestion(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
		c.Abort()
		return
	}

	auths := strings.Split(authHeader, " ")
	if len(auths) != 2 {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
		c.Abort()
		return
	}

	var user domain.User
	data, err := user.DecryptJWT(auths[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
		c.Abort()
		return
	}

	userID := int(data["user_id"].(float64))

	idString := c.Param("exerciseId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var exercise domain.Exercise
	err = eh.db.Where("id = ?", id).Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	var questionBody domain.QuestionRequest
	if err := c.ShouldBind(&questionBody); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid body",
		})
		return
	}

	querstionRequest, err := domain.NewQuestion(questionBody.Body, questionBody.OptionA, questionBody.OptionB, questionBody.OptionC, questionBody.OptionD, questionBody.CorrectAnswer)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	question := &domain.Question{
		CreatorID:     userID,
		ExerciseID:    id,
		Body:          querstionRequest.Body,
		OptionA:       querstionRequest.OptionA,
		OptionB:       querstionRequest.OptionB,
		OptionC:       querstionRequest.OptionC,
		OptionD:       querstionRequest.OptionD,
		CorrectAnswer: querstionRequest.CorrectAnswer,
		Score:         10,
	}

	if err := eh.db.Create(question).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}

func (eh ExerciseHandler) CreateAnswer(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
		c.Abort()
		return
	}

	auths := strings.Split(authHeader, " ")
	if len(auths) != 2 {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
		c.Abort()
		return
	}

	var user domain.User
	data, err := user.DecryptJWT(auths[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
		c.Abort()
		return
	}

	userID := int(data["user_id"].(float64))

	idStringExercise := c.Param("exerciseId")
	exerciseID, err := strconv.Atoi(idStringExercise)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var exercise domain.Exercise
	err = eh.db.Where("id = ?", exerciseID).Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	idStringQuestion := c.Param("questionId")
	questionID, err := strconv.Atoi(idStringQuestion)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var question domain.Question
	err = eh.db.Where("id = ?", questionID).Take(&question).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "question not found",
		})
		return
	}

	var answerBody domain.AnswerRequest
	if err := c.ShouldBind(&answerBody); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid body",
		})
		return
	}

	answerRequest, err := domain.NewAnswer(answerBody.Answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	answer := &domain.Answer{
		ExerciseID: exerciseID,
		QuestionID: questionID,
		UserID:     userID,
		Answer:     answerRequest.Answer,
	}

	if err := eh.db.Create(answer).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}

type Score struct {
	total int
	mu    sync.Mutex
}

func (s *Score) Inc(value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total += value
}
