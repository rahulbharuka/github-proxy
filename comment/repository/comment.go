package repository

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-pg/pg"
	"github.com/rahulbharuka/github-proxy/comment/storage"
)

var (
	// initOnce protects the following
	initCommentRepoOnce  sync.Once
	singletonCommentRepo *commentRepoImpl

	// ErrNoData ...
	ErrNoData = errors.New("no comments for given org")
)

// Comment is a storage object for comment table.
type Comment struct {
	ID        uint64    `json:"id"`
	Org       string    `json:"org"`
	Author    string    `json:"author"`
	Comment   string    `json:"comment"`
	IsDeleted bool      `json:"is_deleted" pg:",use_zero"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// String ...
func (c Comment) String() string {
	return fmt.Sprintf("Comment<%d %s %s %s %t %v %v>", c.ID, c.Org, c.Author, c.Comment, c.IsDeleted, c.CreatedAt, c.UpdatedAt)
}

// commentRepoImpl ...
type commentRepoImpl struct {
	db *pg.DB
}

// CommentRepo implements following methods.
// go:generate mockery -inpkg -case underscore -name CommentRepo
type CommentRepo interface {
	ListAll(ctx context.Context, org string) ([]Comment, error)
	Save(ctx context.Context, c *Comment) error
	DeleteAll(ctx context.Context, org string) error
}

// NewCommentRepo returns the CommentRepo handler.
func NewCommentRepo() CommentRepo {
	initCommentRepoOnce.Do(func() {
		singletonCommentRepo = &commentRepoImpl{
			db: storage.NewDBHandler(),
		}
	})
	return singletonCommentRepo
}

// ListAll lists all comments
func (r *commentRepoImpl) ListAll(ctx context.Context, org string) ([]Comment, error) {
	var comments []Comment
	err := r.db.Model(&comments).Where("org=? and is_deleted=?", org, false).Select()
	if err != nil {
		log.Printf("ERROR: failed to list comments for org %v, err: %v", org, err)
		return nil, err
	}

	return comments, nil
}

// Save saves the comment in table.
func (r *commentRepoImpl) Save(ctx context.Context, c *Comment) error {
	currentTime := time.Now()
	c.CreatedAt = currentTime
	c.UpdatedAt = currentTime

	err := r.db.Insert(c)
	if err != nil {
		log.Printf("ERROR: failed to save comment %+v, err: %v", c, err)
		return err
	}
	return nil
}

// DeleteAll marks all record for given org as deleted.
func (r *commentRepoImpl) DeleteAll(ctx context.Context, org string) error {
	c := &Comment{
		Org:       org,
		IsDeleted: true,
		UpdatedAt: time.Now(),
	}
	resp, err := r.db.Model(c).Set("is_deleted=?is_deleted, updated_at=?updated_at").Where("org=?org and is_deleted=?", false).Update()

	if err != nil {
		log.Printf("ERROR: failed to delete comments for org %v, err: %v", org, err)
		return err
	}

	rowsAffected := resp.RowsAffected()
	if rowsAffected <= 0 {
		log.Printf("INFO: %v rows affected while deleting comments for org %v, err: %v", rowsAffected, org, err)
		return ErrNoData
	}

	return nil
}
