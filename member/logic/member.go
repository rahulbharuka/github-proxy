package logic

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// ListAllMembers list all members of an org and return a list sorted by number of followers.
func (h *handlerImpl) ListAllMembers(ctx *gin.Context) {
	org := ctx.Param("org")
	users, err := h.github.ListAllMembers(ctx, org)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	sort.Slice(users, func(i, j int) bool { return users[i].Followers > users[j].Followers })
	ctx.JSON(http.StatusOK, users)
}
