package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joefazee/go-api-tdd/pkg/common"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"net/http"
	"time"
)

func (s *server) createUser(ctx *gin.Context) {

	req := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid fields",
		})
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := s.store.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *server) userLogin(ctx *gin.Context) {

	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid fields",
		})
		return
	}

	user, err := s.store.FindUserByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}

	err = common.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login",
		})
		return
	}

	token, err := s.jwt.CreateToken(*user, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to create token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token.Token,
	})
}
