package apiv1_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSavedHandlers_RejectAnonymous(t *testing.T) {
	for _, handler := range []web.HandlerFunc{apiv1.SetPostSaved(true), apiv1.SetPostSaved(false), apiv1.GetPostSaved(), apiv1.SavedPosts()} {
		code, _ := mock.NewServer().OnTenant(mock.DemoTenant).AddHeader("Accept", "application/json").Execute(handler)
		if code != http.StatusUnauthorized {
			t.Fatalf("got %d, want 401", code)
		}
	}
}

func TestSavedHandlers_WriteAndFailures(t *testing.T) {
	for _, saved := range []bool{true, false} {
		for _, cause := range []error{nil, app.ErrNotFound, fmt.Errorf("storage unavailable")} {
			var received *cmd.SetPostSaved
			bus.AddHandler(func(ctx context.Context, c *cmd.SetPostSaved) error { received = c; return cause })
			code, _ := mock.NewServer().OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).AddParam("number", 42).AddHeader("Accept", "application/json").Execute(apiv1.SetPostSaved(saved))
			want := http.StatusOK
			if cause == app.ErrNotFound {
				want = http.StatusNotFound
			} else if cause != nil {
				want = http.StatusInternalServerError
			}
			if code != want || received == nil || received.Number != 42 || received.Saved != saved {
				t.Fatalf("code %d, command %+v", code, received)
			}
		}
	}
	code, _ := mock.NewServer().AsUser(mock.JonSnow).AddParam("number", "invalid").AddHeader("Accept", "application/json").Execute(apiv1.SetPostSaved(true))
	if code != http.StatusNotFound {
		t.Fatalf("invalid number: %d", code)
	}
}

func TestSavedHandlers_ReadAndFailures(t *testing.T) {
	for _, cause := range []error{nil, fmt.Errorf("storage unavailable")} {
		bus.AddHandler(func(ctx context.Context, q *query.GetPostSaved) error { q.Result = true; return cause })
		bus.AddHandler(func(ctx context.Context, q *query.GetSavedPosts) error { return cause })
		for _, handler := range []web.HandlerFunc{apiv1.GetPostSaved(), apiv1.SavedPosts()} {
			code, _ := mock.NewServer().OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).AddParam("number", 42).AddHeader("Accept", "application/json").Execute(handler)
			want := http.StatusOK
			if cause != nil {
				want = http.StatusInternalServerError
			}
			if code != want {
				t.Fatalf("got %d want %d", code, want)
			}
		}
	}
}
