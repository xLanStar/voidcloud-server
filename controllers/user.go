package controllers

import (
	"log"
	"net/http"
	"voidcloud-server/internal/exception"
	"voidcloud-server/internal/user"
	"voidcloud-server/services/auth"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var data struct {
		Name     string `json:"name"`
		Account  string `json:"account"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.BindJSON(&data)

	// 檢查是否已存在相同的帳號
	originalUser := user.GetUserByAccount(data.Account)

	if originalUser != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	// 建立資料
	user := user.CreateUser(data.Name, data.Account, data.Email, data.Password)

	// 產生 Token
	token, err := auth.GenerateToken(data.Account, data.Password)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func Login(c *gin.Context) {
	var data struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	c.BindJSON(&data)

	// 檢查帳號密碼是否有誤
	user := user.GetUserByAccount(data.Account)

	log.Println(user)

	if user == nil || data.Password != user.Password {
		c.Status(http.StatusUnauthorized)
		return
	}

	// 產生 Token
	token, err := auth.GenerateToken(data.Account, data.Password)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  string(token),
		"user":   user,
		"folder": user.RootFolder,
	})
}

func Validate(c *gin.Context) {
	user, ok := c.Get("user")
	if ok && user != nil {
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": exception.UNKNOWN,
		})
	}
}
