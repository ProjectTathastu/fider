package apiv1

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

func SetPostSaved(saved bool) web.HandlerFunc {
	return func(c *web.Context) error {
		if !c.IsAuthenticated() {
			return c.Unauthorized()
		}
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}
		if err := bus.Dispatch(c, &cmd.SetPostSaved{Number: number, Saved: saved}); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{"saved": saved})
	}
}

func GetPostSaved() web.HandlerFunc {
	return func(c *web.Context) error {
		if !c.IsAuthenticated() {
			return c.Unauthorized()
		}
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}
		q := &query.GetPostSaved{Number: number}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{"saved": q.Result})
	}
}

func SavedPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		if !c.IsAuthenticated() {
			return c.Unauthorized()
		}
		q := &query.GetSavedPosts{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}
		return c.Ok(q.Result)
	}
}
