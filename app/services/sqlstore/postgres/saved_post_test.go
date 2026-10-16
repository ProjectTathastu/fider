package postgres_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestSavedPosts_PersistenceAndIsolation(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()
	post := &cmd.AddNewPost{Title: "Save me", Description: "Private bookmark"}
	Expect(bus.Dispatch(jonSnowCtx, post)).IsNil()
	save := &cmd.SetPostSaved{Number: post.Result.Number, Saved: true}
	Expect(bus.Dispatch(aryaStarkCtx, save, save)).IsNil()
	state := &query.GetPostSaved{Number: post.Result.Number}
	Expect(bus.Dispatch(aryaStarkCtx, state)).IsNil()
	Expect(state.Result).IsTrue()
	list := &query.GetSavedPosts{}
	// A fresh context for the same identity reads the persisted row.
	Expect(bus.Dispatch(withUser(demoTenantCtx, aryaStark), list)).IsNil()
	Expect(list.Result).HasLen(1)
	Expect(list.Result[0].ID).Equals(post.Result.ID)
	Expect(list.Result[0].VotesCount).Equals(0)
	sub := &query.UserSubscribedTo{PostID: post.Result.ID}
	Expect(bus.Dispatch(aryaStarkCtx, sub)).IsNil()
	Expect(sub.Result).IsFalse()
	var count int
	Expect(trx.Scalar(&count, "SELECT COUNT(*) FROM saved_posts WHERE post_id=$1", post.Result.ID)).IsNil()
	Expect(count).Equals(1)
	Expect(bus.Dispatch(sansaStarkCtx, list)).IsNil()
	Expect(list.Result).HasLen(0)
	Expect(bus.Dispatch(tonyStarkCtx, list)).IsNil()
	Expect(list.Result).HasLen(0)
	// Same post number in another tenant resolves only to that tenant's post.
	other := &cmd.AddNewPost{Title: "Other tenant"}
	Expect(bus.Dispatch(tonyStarkCtx, other)).IsNil()
	Expect(bus.Dispatch(tonyStarkCtx, &cmd.SetPostSaved{Number: other.Result.Number, Saved: true})).IsNil()
	Expect(bus.Dispatch(tonyStarkCtx, list)).IsNil()
	Expect(list.Result).HasLen(1)
	Expect(list.Result[0].ID).Equals(other.Result.ID)
	Expect(bus.Dispatch(aryaStarkCtx, &cmd.SetPostSaved{Number: post.Result.Number}, &cmd.SetPostSaved{Number: post.Result.Number})).IsNil()
	Expect(bus.Dispatch(aryaStarkCtx, state, list)).IsNil()
	Expect(state.Result).IsFalse()
	Expect(list.Result).HasLen(0)
	Expect(bus.Dispatch(tonyStarkCtx, list)).IsNil()
	Expect(list.Result).HasLen(1)
}

func TestSavedPosts_AuthorizationAndVisibility(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()
	post := &cmd.AddNewPost{Title: "Visibility"}
	Expect(bus.Dispatch(jonSnowCtx, post)).IsNil()
	save := &cmd.SetPostSaved{Number: post.Result.Number, Saved: true}
	list := &query.GetSavedPosts{}
	for _, ctx := range []context.Context{demoTenantCtx, withTenant(aryaStarkCtx, avengersTenant)} {
		Expect(bus.Dispatch(ctx, save)).IsNotNil()
		Expect(bus.Dispatch(ctx, list)).IsNotNil()
		Expect(bus.Dispatch(ctx, &query.GetPostSaved{Number: post.Result.Number})).IsNotNil()
	}
	Expect(bus.Dispatch(aryaStarkCtx, save)).IsNil()
	_, err := trx.Execute("UPDATE posts SET is_approved=false WHERE id=$1", post.Result.ID)
	Expect(err).IsNil()
	Expect(bus.Dispatch(aryaStarkCtx, list)).IsNil()
	Expect(list.Result).HasLen(0)
	Expect(errors.Cause(bus.Dispatch(aryaStarkCtx, save))).Equals(app.ErrNotFound)
	// Author and staff retain the same access as the original detail view.
	Expect(bus.Dispatch(jonSnowCtx, save, list)).IsNil()
	Expect(list.Result).HasLen(1)
	_, err = trx.Execute("UPDATE posts SET status=$2 WHERE id=$1", post.Result.ID, enum.PostDeleted)
	Expect(err).IsNil()
	Expect(bus.Dispatch(jonSnowCtx, list)).IsNil()
	Expect(list.Result).HasLen(0)
	Expect(errors.Cause(bus.Dispatch(jonSnowCtx, save))).Equals(app.ErrNotFound)
	Expect(errors.Cause(bus.Dispatch(aryaStarkCtx, &cmd.SetPostSaved{Number: 999999, Saved: true}))).Equals(app.ErrNotFound)
}
