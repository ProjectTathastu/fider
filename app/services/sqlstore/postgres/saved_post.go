package postgres

import (
	"context"
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func savedUserMatches(tenant *entity.Tenant, user *entity.User) bool {
	return tenant != nil && user != nil && user.Tenant != nil && user.Tenant.ID == tenant.ID
}

func setPostSaved(ctx context.Context, c *cmd.SetPostSaved) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if !savedUserMatches(tenant, user) {
			return app.ErrNotFound
		}
		post := &query.GetPostByNumber{Number: c.Number}
		if err := getPostByNumber(ctx, post); err != nil {
			return err
		}
		if c.Saved {
			_, err := trx.Execute(`INSERT INTO saved_posts (tenant_id, user_id, post_id)
    VALUES ($1, $2, $3) ON CONFLICT (tenant_id, user_id, post_id) DO NOTHING`, tenant.ID, user.ID, post.Result.ID)
			return err
		}
		_, err := trx.Execute(`DELETE FROM saved_posts WHERE tenant_id = $1 AND user_id = $2 AND post_id = $3`, tenant.ID, user.ID, post.Result.ID)
		return err
	})
}

func getPostSaved(ctx context.Context, q *query.GetPostSaved) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if !savedUserMatches(tenant, user) {
			return app.ErrNotFound
		}
		post := &query.GetPostByNumber{Number: q.Number}
		if err := getPostByNumber(ctx, post); err != nil {
			return err
		}
		var err error
		q.Result, err = trx.Exists(`SELECT 1 FROM saved_posts WHERE tenant_id = $1 AND user_id = $2 AND post_id = $3`, tenant.ID, user.ID, post.Result.ID)
		return err
	})
}

func getSavedPosts(ctx context.Context, q *query.GetSavedPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if !savedUserMatches(tenant, user) {
			return app.ErrNotFound
		}
		// Use the same visibility rules as opening the original post, including moderation.
		sql := buildSinglePostQuery(user, `p.tenant_id = $1 AND EXISTS (
   SELECT 1 FROM saved_posts s WHERE s.tenant_id = $1 AND s.user_id = $2 AND s.post_id = p.id)`)
		var posts []*dbEntities.Post
		if err := trx.Select(&posts, sql+" ORDER BY p.id DESC", tenant.ID, user.ID); err != nil {
			return err
		}
		q.Result = make([]*entity.Post, len(posts))
		for i, post := range posts {
			q.Result[i] = post.ToModel(ctx)
		}
		return nil
	})
}
