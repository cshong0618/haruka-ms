package post

import (
	"context"
	"errors"
	"github.com/cshong0618/haruka/post/pkg/domain/post"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type EventSource struct {
	mongoCollection *mongo.Collection
}

func (e *EventSource) CreatePost(ctx context.Context, ownerID string, content string) error {
	postID := primitive.NewObjectID()
	command := Command{
		ID:      primitive.NewObjectID(),
		PostID:  postID,
		Command: CreatePost,
		CreatePostCommand: &CreatePostCommand{
			PostID:    postID,
			Content:   content,
			Owner:     ownerID,
			CreatedOn: time.Now(),
		},
	}

	return e.insertCommand(ctx, command)
}

func (e *EventSource) UpdatePost(ctx context.Context, postID string, content string) error {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	command := Command{
		ID:      primitive.NewObjectID(),
		PostID:  postObjectID,
		Command: UpdatePost,
		UpdatePostCommand: &UpdatePostCommand{
			Content:   content,
			CreatedOn: time.Now(),
		},
	}

	return e.insertCommand(ctx, command)
}

func (e *EventSource) UpdatePostLike(ctx context.Context, postID string, likerID string) error {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	command := Command{
		ID:      primitive.NewObjectID(),
		PostID:  postObjectID,
		Command: UpdatePostLike,
		UpdatePostLikeCommand: &UpdatePostLikeCommand{
			Owner:     likerID,
			CreatedOn: time.Now(),
		},
	}

	return e.insertCommand(ctx, command)
}

func (e *EventSource) CreateComment(ctx context.Context, postID string, ownerID string, content string) error {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	command := Command{
		ID:      primitive.NewObjectID(),
		PostID:  postObjectID,
		Command: CreateComment,
		CreateCommentCommand: &CreateCommentCommand{
			ID:        primitive.NewObjectID(),
			Owner:     ownerID,
			Content:   content,
			CreatedOn: time.Now(),
		},
	}

	return e.insertCommand(ctx, command)
}

func (e *EventSource) UpdateCommentLike(ctx context.Context, postID string, commentID string, likerID string) error {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	command := Command{
		ID:      primitive.NewObjectID(),
		PostID:  postObjectID,
		Command: UpdateCommentLike,
		UpdateCommentLikeCommand: &UpdateCommentLikeCommand{
			CommentID: commentObjectID,
			Owner:     likerID,
		},
	}

	return e.insertCommand(ctx, command)
}

func (e *EventSource) PostExists(ctx context.Context, postID string) (bool, error) {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, err
	}

	findQuery := bson.M{
		"postId": postObjectID,
	}

	findOptions := options.Find().
		SetLimit(1).
		SetProjection(bson.M{
			"postId": 1,
		})

	cur, err := e.mongoCollection.Find(ctx, findQuery, findOptions)
	if err != nil {
		return false, err
	}
	defer cur.Close(ctx)

	ok := cur.Next(ctx)
	if !ok {
		return false, nil
	}

	return cur.Current != nil, nil
}

func (e *EventSource) PostOwner(ctx context.Context, postID string) (string, error) {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return "", err
	}

	findQuery := bson.M {
		"postId": postObjectID,
		"command": CreatePost,
	}

	var command Command
	err = e.mongoCollection.FindOne(ctx, findQuery).Decode(&command)
	if err != nil {
		return "", err
	}

	createCommentCommand := command.CreatePostCommand
	if createCommentCommand == nil {
		return "", errors.New("broken command")
	}

	return createCommentCommand.Owner, nil
}

func (e *EventSource) GetPostByID(ctx context.Context, postID string) (post.Post, error) {
	query := bson.M{
		"postId": postID,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{
		"_id": 1,
	})

	cursor, err := e.mongoCollection.Find(ctx, query, findOptions)
	if err != nil {
		return post.Post{}, err
	}

	var commands []Command
	err = cursor.All(ctx, &commands)
	if err != nil {
		return post.Post{}, err
	}

	return BuildPostFromCommands(commands)
}

func (e *EventSource) GetPostsByOwner(ctx context.Context, ownerID string) ([]post.Post, error) {
	panic("implement me")
}

func (e *EventSource) insertCommand(ctx context.Context, command Command) error {
	_, err := e.mongoCollection.InsertOne(ctx, command)
	return err
}

func NewEventSource(
	mongoClient *mongo.Client,
	database string,
	collection string) *EventSource {
	db := mongoClient.Database(database)
	if db == nil {
		panic("failed to create user mongo repository. db is nil")
	}

	dbCollection := db.Collection(collection)
	if dbCollection == nil {
		panic("failed to create user mongo repository. collection is nil")
	}

	return &EventSource{
		mongoCollection: dbCollection,
	}
}
