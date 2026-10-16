import React from "react"
import { render, screen } from "@testing-library/react"
import SavedPage from "./Saved.page"
import { Post } from "@fider/models"

jest.mock("@lingui/react", () => ({
  Trans: ({ message, children }: { message?: string; children?: React.ReactNode }) => <>{message || children}</>,
}))

jest.mock("@fider/components", () => ({
  Header: () => <header />,
  Button: ({ href, children }: { href: string; children: React.ReactNode }) => <a href={href}>{children}</a>,
}))

test("private empty state explains saving and links home", () => {
  render(<SavedPage posts={[]} />)
  expect(screen.getByText(/Only you can see/)).toBeInTheDocument()
  expect(screen.getByText(/No saved feedback yet/)).toBeInTheDocument()
  expect(screen.getByRole("link", { name: "Browse feedback" })).toHaveAttribute("href", "/")
})

test("saved links open the original post", () => {
  render(<SavedPage posts={[{ id: 4, number: 7, slug: "original-feedback", title: "Original feedback" } as Post]} />)
  expect(screen.getByRole("link", { name: "Original feedback" })).toHaveAttribute("href", "/posts/7/original-feedback")
  expect(screen.queryByText(/No saved feedback yet/)).not.toBeInTheDocument()
})
