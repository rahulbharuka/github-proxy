package logic

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rahulbharuka/github-proxy/comment/repository"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/github-proxy/external/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListAllComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	githubMock := &github.MockHandler{}
	commentRepoMock := &repository.MockCommentRepo{}
	h := &handlerImpl{
		github:      githubMock,
		commentRepo: commentRepoMock,
	}

	t.Run("happy-path", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("IsValidOrg", mock.Anything, mock.Anything).Return(true, nil).Once()
		commentRepoMock.On("ListAll", mock.Anything, mock.Anything).Return([]repository.Comment{}, nil).Once()
		h.ListAllComments(ctx)
		assert.Equal(t, http.StatusOK, respWriter.Code)
	})

	t.Run("github-api-err", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("IsValidOrg", mock.Anything, mock.Anything).Return(false, errors.New("some github error")).Once()
		h.ListAllComments(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})

	t.Run("repo-err", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("IsValidOrg", mock.Anything, mock.Anything).Return(true, nil).Once()
		commentRepoMock.On("ListAll", mock.Anything, mock.Anything).Return([]repository.Comment{}, errors.New("some repo error")).Once()
		h.ListAllComments(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})
}

func TestDeleteAllComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	githubMock := &github.MockHandler{}
	commentRepoMock := &repository.MockCommentRepo{}
	h := &handlerImpl{
		github:      githubMock,
		commentRepo: commentRepoMock,
	}

	t.Run("happy-path", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("IsValidOrg", mock.Anything, mock.Anything).Return(true, nil).Once()
		commentRepoMock.On("DeleteAll", mock.Anything, mock.Anything).Return(nil).Once()
		h.DeleteAllComments(ctx)
		assert.Equal(t, http.StatusOK, respWriter.Code)
	})

	t.Run("github-api-err", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("IsValidOrg", mock.Anything, mock.Anything).Return(false, errors.New("some github error")).Once()
		h.DeleteAllComments(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})

	t.Run("repo-err", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("IsValidOrg", mock.Anything, mock.Anything).Return(true, nil).Once()
		commentRepoMock.On("DeleteAll", mock.Anything, mock.Anything).Return(errors.New("some repo error")).Once()
		h.DeleteAllComments(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})
}

func TestPostComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	githubMock := &github.MockHandler{}
	commentRepoMock := &repository.MockCommentRepo{}
	h := &handlerImpl{
		github:      githubMock,
		commentRepo: commentRepoMock,
	}

	t.Run("happy-path", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}
		jsonBody := `{"author":"awesome.user","comment":"test comment"}`
		body := ioutil.NopCloser(bytes.NewReader([]byte(jsonBody)))
		ctx.Request = &http.Request{Body: body}

		githubMock.On("IsMember", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
		commentRepoMock.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
		h.PostComment(ctx)
		assert.Equal(t, http.StatusOK, respWriter.Code)
	})

	t.Run("bad-request-body", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		jsonBody := `{"author":"awesome.user","comment":"test comment"}}`
		body := ioutil.NopCloser(bytes.NewReader([]byte(jsonBody)))
		ctx.Request = &http.Request{Body: body}

		h.PostComment(ctx)
		assert.Equal(t, http.StatusBadRequest, respWriter.Code)
	})

	t.Run("github-api-error", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		jsonBody := `{"author":"awesome.user","comment":"test comment"}`
		body := ioutil.NopCloser(bytes.NewReader([]byte(jsonBody)))
		ctx.Request = &http.Request{Body: body}

		githubMock.On("IsMember", mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("some github error")).Once()
		h.PostComment(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})

	t.Run("save-error", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}
		jsonBody := `{"author":"awesome.user","comment":"test comment"}`
		body := ioutil.NopCloser(bytes.NewReader([]byte(jsonBody)))
		ctx.Request = &http.Request{Body: body}

		githubMock.On("IsMember", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
		commentRepoMock.On("Save", mock.Anything, mock.Anything).Return(errors.New("save error")).Once()
		h.PostComment(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})

}
