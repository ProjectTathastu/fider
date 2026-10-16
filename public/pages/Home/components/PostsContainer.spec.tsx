import React from "react"
import { act, fireEvent, render, screen } from "@testing-library/react"
import { actions, Fider } from "@fider/services"
import { PostsContainer } from "./PostsContainer"

jest.mock("@fider/services/actions", () => ({ searchPosts: jest.fn() }))

jest.mock("./ListPosts", () => ({
  ListPosts: ({ posts }: { posts?: Array<{ title: string }> }) => {
    const React = require("react")
    return React.createElement(
      "div",
      { "data-testid": "feedback-list" },
      posts?.map((post) => React.createElement("div", { key: post.title }, post.title))
    )
  },
}))

jest.mock("./PostFilter", () => ({ PostFilter: () => null }))
jest.mock("./PostsSort", () => ({ PostsSort: () => null }))

const feedback = [{ title: "Matching feedback" }, { title: "Other feedback" }]

const renderPostsContainer = () => render(<PostsContainer posts={feedback as any} tags={[]} countPerStatus={{}} />)

const finishSearch = async () => {
  await act(async () => {
    jest.advanceTimersByTime(500)
    await Promise.resolve()
  })
}

describe("feedback search clear button", () => {
  beforeEach(() => {
    jest.useFakeTimers()
    Fider.initialize({ settings: { baseURL: "" }, page: "", contextID: "", props: {} })
    jest.mocked(actions.searchPosts).mockImplementation(async ({ query }) => ({
      ok: true,
      data: (query ? feedback.filter((post) => post.title.toLowerCase().includes(query)) : feedback) as any,
    }))
  })

  afterEach(() => {
    jest.clearAllMocks()
    jest.useRealTimers()
  })

  test("appears for a query, filters feedback, then clears by pointer and restores focus and all results", async () => {
    renderPostsContainer()
    const input = screen.getByPlaceholderText("Search")

    expect(screen.queryByRole("button", { name: "Clear search" })).not.toBeInTheDocument()
    expect(screen.getByText("Matching feedback")).toBeInTheDocument()
    expect(screen.getByText("Other feedback")).toBeInTheDocument()

    fireEvent.change(input, { target: { value: "matching" } })
    const clearButton = screen.getByRole("button", { name: "Clear search" })
    expect(clearButton).toBeInTheDocument()
    await finishSearch()

    expect(actions.searchPosts).toHaveBeenCalledWith(expect.objectContaining({ query: "matching" }))
    expect(screen.getByText("Matching feedback")).toBeInTheDocument()
    expect(screen.queryByText("Other feedback")).not.toBeInTheDocument()

    fireEvent.click(clearButton)
    await finishSearch()

    expect(input).toHaveValue("")
    expect(screen.queryByRole("button", { name: "Clear search" })).not.toBeInTheDocument()
    expect(screen.getByText("Matching feedback")).toBeInTheDocument()
    expect(screen.getByText("Other feedback")).toBeInTheDocument()
    expect(document.activeElement).toBe(input)
    expect(actions.searchPosts).toHaveBeenLastCalledWith(expect.objectContaining({ query: "" }))
  })

  test("clears search when the focused button is activated with the keyboard", async () => {
    renderPostsContainer()
    const input = screen.getByPlaceholderText("Search")
    fireEvent.change(input, { target: { value: "matching" } })
    await finishSearch()

    const clearButton = screen.getByRole("button", { name: "Clear search" })
    clearButton.focus()
    expect(document.activeElement).toBe(clearButton)
    fireEvent.keyDown(clearButton, { key: "Enter" })
    // Browsers synthesize a click when Enter activates a native button; jsdom requires that activation explicitly.
    fireEvent.click(clearButton)
    await finishSearch()

    expect(input).toHaveValue("")
    expect(screen.getByText("Other feedback")).toBeInTheDocument()
    expect(document.activeElement).toBe(input)
  })
})
