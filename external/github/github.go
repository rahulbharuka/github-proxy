package github

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/google/go-github/github"
)

var (
	// initOnce protects the following
	initOnce         sync.Once
	singletonHandler *handlerImpl
)

// User is a model for a Git user.
type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

// Handler in an interface for interacting with Github API v3.
// go:generate mockery -inpkg -case underscore -name Handler
type Handler interface {
	IsValidOrg(ctx context.Context, org string) (bool, error)
	IsMember(ctx context.Context, org, user string) (bool, error)
	ListAllMembers(ctx context.Context, org string) ([]*User, error)
}

type handlerImpl struct {
	client *github.Client
}

// GetHandler initializes the Github client and return the handler.
func GetHandler() Handler {
	initOnce.Do(func() {
		singletonHandler = &handlerImpl{
			client: github.NewClient(nil),
		}
	})
	return singletonHandler
}

// IsValidOrg checks whether the organization exists in Github.
func (h *handlerImpl) IsValidOrg(ctx context.Context, org string) (bool, error) {
	_, resp, err := h.client.Organizations.Get(ctx, org)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		log.Printf("ERROR: failed to validate org %v from Github, err: %v", org, err)
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}

// IsMember checks whether the user is a public member of specified org in Github.
func (h *handlerImpl) IsMember(ctx context.Context, org, user string) (bool, error) {
	isMember, _, err := h.client.Organizations.IsPublicMember(ctx, org, user)
	if err != nil {
		log.Printf("ERROR: failed to validate org %v and member %v from Github, err: %v", org, user, err)
		return false, err
	}
	return isMember, nil
}

// ListAllMembers fetch all members of specified Github org and return slice of *User.
func (h *handlerImpl) ListAllMembers(ctx context.Context, org string) ([]*User, error) {
	allMembers := []*github.User{}
	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		members, resp, err := h.client.Organizations.ListMembers(ctx, org, opt)
		if err != nil {
			log.Printf("ERROR: failed to fetch org %v members from Github, err: %v", org, err)
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("ERROR: failed to fetch org %v members from Github, status: %v", org, resp.Status)
			return nil, errors.New(resp.Status)
		}

		allMembers = append(allMembers, members...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	users := []*User{}
	for _, member := range allMembers {
		user, resp, err := h.client.Users.Get(ctx, member.GetLogin())
		if err != nil {
			log.Printf("ERROR: failed to fetch user %v from Github, err: %v", member.GetLogin(), err)
			return nil, err
		}
		if resp.StatusCode == http.StatusOK {
			users = append(users,
				&User{
					Login:     user.GetLogin(),
					AvatarURL: user.GetAvatarURL(),
					Followers: user.GetFollowers(),
					Following: user.GetFollowing(),
				})
		}
	}
	return users, nil
}
