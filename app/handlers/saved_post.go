package handlers

import (
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"net/http"
)

func SavedPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		if !c.IsAuthenticated() {
			return c.Unauthorized()
		}
		q := &query.GetSavedPosts{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}
		return c.Page(http.StatusOK, web.Props{Page: "Saved/Saved.page", Title: "Saved", Data: web.Map{"posts": q.Result}})
	}
}
