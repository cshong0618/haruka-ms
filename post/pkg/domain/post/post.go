package post

import "time"

type Post struct {
	ID        string
	Content   string
	Owner     string
	UserLikes []string
	Comments  []Comment
	CreatedOn time.Time
}

type Comment struct {
	ID        string
	Content   string
	Owner     string
	UserLikes []string
	CreatedOn time.Time
}
