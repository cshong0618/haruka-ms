package post

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreatePostCommand struct {
	PostID    primitive.ObjectID `bson:"postId"`
	Content   string             `bson:"content"`
	Owner     string             `bson:"owner"`
	CreatedOn time.Time          `bson:"createdOn"`
}

type UpdatePostCommand struct {
	Content   string    `bson:"content"`
	CreatedOn time.Time `bson:"createdOn"`
}

type UpdatePostLikeCommand struct {
	Owner     string    `bson:"owner"`
	CreatedOn time.Time `bson:"createdOn"`
}

type CreateCommentCommand struct {
	ID        primitive.ObjectID `bson:"_id"`
	Owner     string             `bson:"owner"`
	Content   string             `bson:"content"`
	CreatedOn time.Time          `bson:"createdOn"`
}

type UpdateCommentLikeCommand struct {
	CommentID primitive.ObjectID `bson:"commentId"`
	Owner     string             `bson:"owner"`
}

type CommandType string

const (
	CreatePost        = CommandType("CREATE_POST")
	UpdatePost        = CommandType("UPDATE_POST")
	UpdatePostLike    = CommandType("UPDATE_POST_TYPE")
	CreateComment     = CommandType("CREATE_COMMENT")
	UpdateCommentLike = CommandType("UPDATE_COMMENT_LIKE")
)

// TODO: think of a better way
type Command struct {
	ID                       primitive.ObjectID        `bson:"_id"`
	PostID                   primitive.ObjectID        `bson:"postId"`
	Command                  CommandType               `bson:"command"`
	CreatePostCommand        *CreatePostCommand        `bson:"createPost"`
	UpdatePostCommand        *UpdatePostCommand        `bson:"updatePost"`
	UpdatePostLikeCommand    *UpdatePostLikeCommand    `bson:"updatePostLike"`
	CreateCommentCommand     *CreateCommentCommand     `bson:"createComment"`
	UpdateCommentLikeCommand *UpdateCommentLikeCommand `bson:"updateCommentLike"`
}
