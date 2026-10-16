package query

import "github.com/getfider/fider/app/models/entity"

type GetSavedPosts struct {
	Result []*entity.Post
}

type GetPostSaved struct {
	Number int
	Result bool
}
