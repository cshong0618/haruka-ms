package post

import (
	"errors"
	"github.com/cshong0618/haruka/post/pkg/domain/post"
	"sort"
)

func BuildPostFromCommands(commands []Command) (post.Post, error) {
	domainPost := post.Post{}
	comments := map[string]post.Comment{}

	return buildPostFromCommands(domainPost, comments, commands)
}

func ContinuePostFromCommands(domainPost post.Post, commands []Command) (post.Post, error) {
	comments := map[string]post.Comment{}
	for _, comment := range domainPost.Comments {
		comments[comment.ID] = comment
	}

	return buildPostFromCommands(domainPost, comments, commands)
}

func buildPostFromCommands(domainPost post.Post, comments map[string]post.Comment, commands []Command) (post.Post, error) {
	for _, command := range commands {
		var err error
		switch command.Command {
		case CreatePost:
			err = handleCreatePostCommand(&domainPost, *command.CreatePostCommand)
		case UpdatePost:
			err = handleUpdatePostCommand(&domainPost, *command.UpdatePostCommand)
		case UpdatePostLike:
			err = handleUpdatePostLikeCommand(&domainPost, *command.UpdatePostLikeCommand)
		case CreateComment:
			err = handleCreateCommentCommand(comments, *command.CreateCommentCommand)
		case UpdateCommentLike:
			err = handleUpdateCommentLikeCommand(comments, *command.UpdateCommentLikeCommand)
		default:
			err = errors.New("unknown command")
		}

		if err != nil {
			break
		}
	}

	// create domain comments from comment map
	domainComments := make([]post.Comment, len(comments))
	i := 0
	for _, v := range comments {
		domainComments[i] = v
		i++
	}

	// sort domain comments
	sort.SliceStable(domainComments, func(i, j int) bool {
		return domainComments[i].CreatedOn.Before(domainComments[j].CreatedOn)
	})

	domainPost.Comments = domainComments

	return domainPost, nil
}

func handleCreatePostCommand(domainPost *post.Post, command CreatePostCommand) error {
	if command.PostID.IsZero() {
		return errors.New("invalid postId")
	}

	domainPost.ID = command.PostID.Hex()
	domainPost.Content = command.Content
	domainPost.Owner = command.Owner
	domainPost.CreatedOn = command.CreatedOn

	return nil
}

func handleUpdatePostCommand(domainPost *post.Post, command UpdatePostCommand) error {
	domainPost.Content = command.Content
	return nil
}

func handleUpdatePostLikeCommand(domainPost *post.Post, command UpdatePostLikeCommand) error {
	if domainPost.UserLikes == nil {
		domainPost.UserLikes = []string{}
	}

	// check if owner already liked the comment
	for _, v := range domainPost.UserLikes {
		if v == command.Owner {
			return nil
		}
	}

	domainPost.UserLikes = append(domainPost.UserLikes, command.Owner)
	return nil
}

func handleCreateCommentCommand(comments map[string]post.Comment, command CreateCommentCommand) error {
	if comments == nil {
		return errors.New("null comment container")
	}
	if command.ID.IsZero() {
		return errors.New("invalid comment Id")
	}

	commentID := command.ID.Hex()
	if _, ok := comments[commentID]; ok {
		return errors.New("comment already created")
	}

	comment := post.Comment{
		ID:        commentID,
		Content:   command.Content,
		Owner:     command.Owner,
		UserLikes: []string{},
		CreatedOn: command.CreatedOn,
	}

	comments[commentID] = comment
	return nil
}

func handleUpdateCommentLikeCommand(comments map[string]post.Comment, command UpdateCommentLikeCommand) error {
	if comments == nil {
		return errors.New("null comment container")
	}

	if command.CommentID.IsZero() {
		return errors.New("invalid comment Id")
	}

	commentID := command.CommentID.Hex()
	comment, ok := comments[commentID]
	if !ok {
		return errors.New("comment not found")
	}

	// check if owner already liked the comment
	for _, v := range comment.UserLikes {
		if v == command.Owner {
			return nil
		}
	}

	comment.UserLikes = append(comment.UserLikes, command.Owner)
	return nil
}