package subscribers

import (
	"context"
	"github.com/cshong0618/haruka/post/pkg/usecase"
	"github.com/nats-io/nats.go"
)

type CreatePostCommandSubscriber struct {
	c *nats.EncodedConn
	postService *usecase.PostService
}

func NewCreatePostCommandSubscriber(
	c *nats.EncodedConn,
	postService *usecase.PostService) *CreatePostCommandSubscriber {
	return &CreatePostCommandSubscriber{c: c, postService: postService}
}

type CreatePostOutput struct {
	OK    bool   `json:"ok"`
	Error *Error `json:"error"`
}

type CreatePostInput struct {
	OwnerID string `json:"ownerId"`
	Content string `json:"content"`
}

func (subscriber *CreatePostCommandSubscriber) CreatePost(
	subject, reply string, input *CreatePostInput) {
	err := subscriber.postService.CreatePost(
		context.Background(), input.OwnerID, input.Content)
	
	if err.HasError() {
		subscriber.c.Publish(reply, CreatePostOutput{
			OK:    false,
			Error: &Error{Reason: err.Error()},
		})
	} else {
		subscriber.c.Publish(reply, CreatePostOutput{
			OK: true,
		})
	}
}