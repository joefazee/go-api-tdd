package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"net/http"
)

func (s *server) createPost(ctx *gin.Context) {

	authUser := s.getAuthUser(ctx)

	req := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Title == "" || req.Body == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid fields",
		})
		return
	}

	_, err := s.store.CreatePost(&domain.Post{
		UserID: authUser.ID,
		Title:  req.Title,
		Body:   req.Body,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error creating post " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   "",
		"message": "post created",
	})

}

func (s *server) listUserPosts(ctx *gin.Context) {

	authUser := s.getAuthUser(ctx)

	posts, err := s.store.GetAllUserPosts(authUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error creating post " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": "",
		"posts": posts,
	})
}
