package post

type CreatePostInput struct {
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type CreatePostNATSInput struct {
	OwnerID string `json:"ownerId"`
	Content string `json:"content"`
}

type CommandOutput struct {
	OK    bool `json:"ok"`
	Error *struct {
		Reason string `json:"reason"`
	} `json:"error"`
}

type UpdatePostInput struct {
	OwnerId string `json:"userId"`
	Content string `json:"content"`
}

type UpdatePostNATSInput struct {
	PostId  string `json:"postId"`
	OwnerId string `json:"ownerId"`
	Content string `json:"content"`
}
