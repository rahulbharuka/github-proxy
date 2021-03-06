package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/github-proxy/comment/repository"
	"github.com/rahulbharuka/github-proxy/external/github"
)

// Handler is the logic handler interface
type Handler interface {
	PostComment(ctx *gin.Context)
	ListAllComments(ctx *gin.Context)
	DeleteAllComments(ctx *gin.Context)
}

// handlerImpl is a implementation of Handler interface
type handlerImpl struct {
	commentRepo repository.CommentRepo
	github      github.Handler
}

// GetHandler initializes and returns the logic layer handler.
func GetHandler() Handler {
	return &handlerImpl{
		commentRepo: repository.NewCommentRepo(),
		github:      github.GetHandler(),
	}
}

// handlerError is a helper function to return JSON error.
func handlerError(ctx *gin.Context, errCode int, err error) {
	ctx.JSON(errCode, gin.H{
		"message": err.Error(),
	})
}
