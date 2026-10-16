import React from "react"
import { Trans } from "@lingui/react/macro"
import { Header, Button } from "@fider/components"
import { Post } from "@fider/models"

export default function SavedPage({ posts }: { posts: Post[] }) {
  return (
    <>
      <Header />
      <main id="p-saved" className="page container">
        <h1 className="text-2xl text-bold mb-3">
          <Trans id="saved.title">Saved</Trans>
        </h1>
        <p className="text-muted mb-4">
          <Trans id="saved.private">Only you can see your saved feedback on this site.</Trans>
        </p>
        {posts.length === 0 ? (
          <div>
            <p className="mb-4">
              <Trans id="saved.empty">No saved feedback yet. Open a post and select Save to find it here later.</Trans>
            </p>
            <Button href="/">
              <Trans id="saved.browse">Browse feedback</Trans>
            </Button>
          </div>
        ) : (
          <ul className="p-0" style={{ listStyle: "none" }}>
            {posts.map((post) => (
              <li key={post.id} className="py-3 border-b">
                <a className="text-link text-break" href={`/posts/${post.number}/${post.slug}`}>
                  {post.title}
                </a>
              </li>
            ))}
          </ul>
        )}
      </main>
    </>
  )
}
