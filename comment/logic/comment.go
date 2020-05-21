package logic

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/github-proxy/comment/model"
	"github.com/rahulbharuka/github-proxy/comment/repository"
)

// PostComment posts a comment for the org
func (h *handlerImpl) PostComment(ctx *gin.Context) {
	org := ctx.Param("org")

	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("ERROR: failed to read request body, err: %v", err)
		handlerError(ctx, http.StatusBadRequest, err)
		return
	}
	c := &repository.Comment{}
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Printf("ERROR: failed to unmarshal request body, err: %v", err)
		handlerError(ctx, http.StatusBadRequest, err)
		return
	}
	c.Org = org

	isValid, err := h.github.IsMember(ctx, c.Org, c.Author)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isValid {
		log.Printf("INFO: user %v is not a member of org %v", c.Author, c.Org)
		handlerError(ctx, http.StatusNotFound, errors.New("user is not a member of specified org"))
		return
	}

	err = h.commentRepo.Save(ctx, c)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "added comment")
}

// ListAllComments fetches all comments for an org.
func (h *handlerImpl) ListAllComments(ctx *gin.Context) {
	org := ctx.Param("org")

	isValid, err := h.github.IsValidOrg(ctx, org)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isValid {
		log.Printf("INFO: %v is not a valid Github org", org)
		handlerError(ctx, http.StatusNotFound, errors.New("Specified org does not exist"))
		return
	}

	comments, err := h.commentRepo.ListAll(ctx, org)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]*model.Comment, len(comments))
	for i, c := range comments {
		resp[i] = &model.Comment{
			Author:    c.Author,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteAllComments soft deletes all comments for an org.
func (h *handlerImpl) DeleteAllComments(ctx *gin.Context) {
	org := ctx.Param("org")
	isValid, err := h.github.IsValidOrg(ctx, org)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isValid {
		log.Printf("INFO: %v is not a valid Github org", org)
		handlerError(ctx, http.StatusNotFound, errors.New("Specified org does not exist"))
		return
	}

	err = h.commentRepo.DeleteAll(ctx, org)
	if err == repository.ErrNoData {
		handlerError(ctx, http.StatusNoContent, err)
		return
	}

	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "deleted all comments !")
}
