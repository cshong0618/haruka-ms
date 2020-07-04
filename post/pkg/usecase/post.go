package usecase

import (
	"context"
	errors2 "errors"
	"github.com/cshong0618/haruka/post/pkg/domain/post"
	"github.com/cshong0618/haruka/post/pkg/utils/audit"
	"github.com/cshong0618/haruka/post/pkg/utils/errors"
)

type PostService struct {
	postRepository post.Repository
}

func NewPostService(postRepository post.Repository) *PostService {
	return &PostService{postRepository: postRepository}
}

func (service *PostService) CreatePost(
	ctx context.Context,
	ownerID string,
	content string) (werr errors.WrappedError) {
	defer audit.Start(ctx, "CreatePost", "ownerID", ownerID).Capture(&werr).End()

	err := service.postRepository.CreatePost(ctx, ownerID, content)
	werr.SetError(err)

	return
}

func (service *PostService) UpdatePost(
	ctx context.Context,
	postID string,
	ownerID string,
	content string) (werr errors.WrappedError) {
	defer audit.Start(ctx, "UpdatePost",
		"postID", postID,
		"ownerID", ownerID).Capture(&werr).End()
	postOwner, err := service.postRepository.PostOwner(ctx, postID)

	if err != nil {
		werr.SetError(err)
		return
	}

	if ownerID != postOwner {
		werr.SetError(errors2.New("post does not belong to this user"))
		return
	}

	err = service.postRepository.UpdatePost(ctx, postID, content)
	werr.SetError(err)
	return
}

func (service *PostService) LikePost(
	ctx context.Context,
	postID string,
	likerID string) error {
	return service.postRepository.UpdatePostLike(ctx, postID, likerID)
}

func (service *PostService) CreateComment(
	ctx context.Context,
	postID string,
	ownerID string,
	content string) error {
	return service.postRepository.CreateComment(ctx, postID, ownerID, content)
}

func (service *PostService) LikeComment(
	ctx context.Context,
	postID string,
	commentID string,
	likerID string) error {
	return service.postRepository.UpdateCommentLike(ctx, postID, commentID, likerID)
}
