package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func NextDateHandler(c *gin.Context) {
	nowURL := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")
	log.Println(nowURL, date, repeat, c.Request.URL.String())
	now, err := time.Parse("20060102", nowURL)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := NextDate(now, date, repeat)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusOK, resp)
}

// Обработчик для добалвения задания
func postTaskHandler(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := checkTask(&task)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	id, err := postTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": id})
}

// Обработчик для получения ближайших задач
func getRecentTasksHandler(c *gin.Context) {
	var (
		err   error
		tasks Tasks
	)
	search := c.Query("search")
	if search == "" {
		tasks, err = getRecentTasks()
	} else {
		tasks, err = getTasksBySearch(search)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Обработчик для получения задачи по ID
func getTaskByIDHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrNoID.Error()})
		return
	}
	task, err := getTaskByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrTaskNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Обработчик для обновления задачи
func updateTaskHandler(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := checkTask(&task)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	err = updateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{})
}

// Обработчик для выполненных задач
func doneTaskHandeler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrNoID.Error()})
		return
	}
	task, err := getTaskByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrTaskNotFound.Error()})
		return
	}
	if task.Repeat == "" {
		err := deleteTask(task.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = updateTask(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusAccepted, gin.H{})
}

// Обработчик для удаления задач
func deleteTaskHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrNoID.Error()})
		return
	}
	err := deleteTask(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{})
}

// Обработчик для аутентификации
func authHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Password != os.Getenv("TODO_PASSWORD") {
		c.JSON(http.StatusForbidden, gin.H{"error": ErrWrongPassword.Error()})
		return
	}
	payload := sha256.Sum256([]byte(user.Password))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"HSPW": payload})
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var resp TokenResponse
	resp.Token = signedToken
	log.Println("JWT token: " + resp.Token)
	c.JSON(http.StatusAccepted, resp)
}

// Middleware для аутентификации
func auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string // JWT-токен из куки
			// получаем куку
			cookie, err := c.Request.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET")), nil
			})
			if err != nil {
				log.Println(err, jwtToken)
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			res, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.AbortWithError(http.StatusInternalServerError, errors.New("failed to typecast to jwt.MapCalims"))
				return
			}
			HSPW := res["HSPW"]
			if fmt.Sprintf("%v", HSPW) == fmt.Sprintf("%v", sha256.Sum256([]byte(pass))) {
				valid = true
			}
			if !valid {
				// возвращаем ошибку авторизации 401
				c.AbortWithError(http.StatusUnauthorized, errors.New("authentification required"))
				return
			}
		}
		c.Next()
	})
}
