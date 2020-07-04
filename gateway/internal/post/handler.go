package post

import (
	"errors"
	"github.com/cshong0618/haruka/gateway/internal"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"time"
)

type PostHandler struct {
	c *nats.EncodedConn
}

func NewPostHandler(nc *nats.Conn) *PostHandler {
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil
	}

	return &PostHandler{c: c}
}

func (handler *PostHandler) CreatePost(c echo.Context) error {
	var input CreatePostInput
	if err := c.Bind(&input); err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(400, response)
		return err
	}

	postInput := CreatePostNATSInput{
		OwnerID: input.UserId,
		Content: input.Content,
	}
	var output CommandOutput
	err := handler.c.Request("post.create", postInput, &output, 5*time.Second)
	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	response := internal.NewResponse(output.OK)
	err = c.JSON(200, response)
	return err
}

func (handler *PostHandler) UpdatePost(c echo.Context) error {
	postID := c.Param("postId")
	if postID == "" {
		response := internal.NewResponse(errors.New("no userId provided"))
		err := c.JSON(400, response)
		return err
	}

	var input UpdatePostInput
	if err := c.Bind(&input); err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(400, response)
		return err
	}

	postInput := UpdatePostNATSInput{
		PostId:  postID,
		OwnerId: input.OwnerId,
		Content: input.Content,
	}

	var output CommandOutput
	err := handler.c.Request("post.update", postInput, &output, 5*time.Second)
	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	response := internal.NewResponse(output.OK)
	err = c.JSON(200, response)
	return err
}
