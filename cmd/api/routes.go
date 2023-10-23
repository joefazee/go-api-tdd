package main

import "github.com/gin-gonic/gin"

func (s *server) setupRoutes() {

	mux := gin.Default()

	v1 := mux.Group("/api/v1")

	v1.POST("/users/create", s.createUser)
	v1.POST("/users/login", s.userLogin)

	blog := v1.Group("/blog")
	blog.Use(s.applyAuthentication())
	blog.POST("/create", s.createPost)
	blog.GET("/list", s.listUserPosts)

	s.router = mux
}
