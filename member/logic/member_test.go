package logic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/github-proxy/external/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListAllMembers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	githubMock := &github.MockHandler{}
	h := &handlerImpl{
		github: githubMock,
	}

	t.Run("happy-path", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("ListAllMembers", mock.Anything, mock.Anything).Return([]*github.User{}, nil).Once()
		h.ListAllMembers(ctx)
		assert.Equal(t, http.StatusOK, respWriter.Code)
	})

	t.Run("github-api-error", func(t *testing.T) {
		respWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respWriter)
		ctx.Params = []gin.Param{gin.Param{Key: "org", Value: "github"}}

		githubMock.On("ListAllMembers", mock.Anything, mock.Anything).Return(nil, errors.New("some github error")).Once()
		h.ListAllMembers(ctx)
		assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	})
}
