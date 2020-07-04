package subscribers

import (
	"context"
	"github.com/cshong0618/haruka/post/pkg/usecase"
	"github.com/nats-io/nats.go"
)

type UpdatePostCommandSubscriber struct {
	c           *nats.EncodedConn
	postService *usecase.PostService
}

func NewUpdatePostCommandSubscriber(
	c *nats.EncodedConn,
	postService *usecase.PostService) *UpdatePostCommandSubscriber {
	return &UpdatePostCommandSubscriber{c: c, postService: postService}
}

type UpdatePostOutput struct {
	OK    bool   `json:"ok"`
	Error *Error `json:"error"`
}

type UpdatePostInput struct {
	PostID  string `json:"postId"`
	OwnerID string `json:"ownerId"`
	Content string `json:"content"`
}

func (subscriber *UpdatePostCommandSubscriber) UpdatePost(
	subject, reply string, input *UpdatePostInput) {
	err := subscriber.postService.UpdatePost(context.Background(),
		input.PostID, input.OwnerID, input.Content)

	if err.HasError() {
		subscriber.c.Publish(reply, UpdatePostOutput{
			OK:    false,
			Error: &Error{Reason: err.Error()},
		})
	} else {
		subscriber.c.Publish(reply, UpdatePostOutput{
			OK: true,
		})
	}
}