package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestSavedPage_RequiresUserAndPropagatesStorageFailure(t *testing.T) {
	code, _ := mock.NewServer().OnTenant(mock.DemoTenant).AddHeader("Accept", "application/json").Execute(handlers.SavedPosts())
	if code != http.StatusUnauthorized {
		t.Fatalf("anonymous status %d", code)
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetSavedPosts) error { return fmt.Errorf("database unavailable") })
	code, _ = mock.NewServer().OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).AddHeader("Accept", "application/json").Execute(handlers.SavedPosts())
	if code != http.StatusInternalServerError {
		t.Fatalf("storage failure status %d", code)
	}
}

func TestSavedPage_UsesPrivatePosts(t *testing.T) {
	bus.AddHandler(func(ctx context.Context, q *query.GetSavedPosts) error {
		q.Result = []*entity.Post{{ID: 7, Number: 42, Title: "My saved post", Slug: "my-saved-post"}}
		return nil
	})
	code, props := mock.NewServer().OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).ExecuteAsPage(handlers.SavedPosts())
	if code != http.StatusOK || props.Page != "Saved/Saved.page" {
		t.Fatalf("status %d, page %s", code, props.Page)
	}
	posts, ok := props.Data["posts"].([]interface{})
	if !ok || len(posts) != 1 {
		t.Fatalf("saved props: %+v", props.Data)
	}
}
