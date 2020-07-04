package post

import "context"

type Repository interface {
	CreatePost(ctx context.Context, ownerID string, content string) error
	UpdatePost(ctx context.Context, postID string, content string) error
	UpdatePostLike(ctx context.Context, postID string, likerID string) error

	CreateComment(ctx context.Context, postID string, ownerID string, content string) error
	UpdateCommentLike(ctx context.Context, postID string, commentID string, likerID string) error

	PostExists(ctx context.Context, postID string) (bool, error)
	PostOwner(ctx context.Context, postID string) (string, error)

	// Aggregates
	GetPostByID(ctx context.Context, postID string) (Post, error)
	GetPostsByOwner(ctx context.Context, ownerID string) ([]Post, error)
}
