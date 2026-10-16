import React from "react"
import { render, screen, fireEvent, waitFor, act } from "@testing-library/react"
import { SavePostButton } from "./SavePostButton"
import { http } from "@fider/services"
import { useFider } from "@fider/hooks"

jest.mock("@lingui/react", () => ({
  Trans: ({ message, children }: { message?: string; children?: React.ReactNode }) => <>{message || children}</>,
}))

jest.mock("@fider/services", () => ({ http: { get: jest.fn(), put: jest.fn(), delete: jest.fn() } }))
jest.mock("@fider/hooks", () => ({ useFider: jest.fn() }))
const get = http.get as jest.Mock
const put = http.put as jest.Mock
const remove = http.delete as jest.Mock
const session = { session: { isAuthenticated: true, user: { id: 1 }, tenant: { id: 1 } }, isReadOnly: false }

beforeEach(() => {
  jest.clearAllMocks()
  ;(useFider as jest.Mock).mockReturnValue(session)
  get.mockResolvedValue({ ok: true, data: { saved: false } })
  put.mockResolvedValue({ ok: true, data: { saved: true } })
  remove.mockResolvedValue({ ok: true, data: { saved: false } })
})

test("anonymous users cannot save or fetch private state", () => {
  ;(useFider as jest.Mock).mockReturnValue({ ...session, session: { ...session.session, isAuthenticated: false } })
  render(<SavePostButton postNumber={7} />)
  expect(screen.queryByRole("button")).not.toBeInTheDocument()
  expect(get).not.toHaveBeenCalled()
})

test("loads persisted state, unsaves and saves independently", async () => {
  get.mockResolvedValue({ ok: true, data: { saved: true } })
  render(<SavePostButton postNumber={7} />)
  fireEvent.click(await screen.findByRole("button", { name: "Unsave" }))
  await waitFor(() => expect(remove).toHaveBeenCalledWith("/api/v1/posts/7/saved"))
  const save = await screen.findByRole("button", { name: "Save" })
  expect(save).toHaveAttribute("aria-pressed", "false")
  fireEvent.click(save)
  expect(await screen.findByRole("button", { name: "Unsave" })).toHaveAttribute("aria-pressed", "true")
  expect(put).toHaveBeenCalledWith("/api/v1/posts/7/saved")
})

test("keeps confirmed state and reports failed writes", async () => {
  put.mockResolvedValue({ ok: false })
  render(<SavePostButton postNumber={7} />)
  const button = screen.getByRole("button", { name: "Save" })
  await waitFor(() => expect(button).toBeEnabled())
  fireEvent.click(button)
  expect(await screen.findByRole("alert")).toBeInTheDocument()
  expect(button).toHaveAttribute("aria-pressed", "false")
  put.mockResolvedValue({ ok: true, data: { saved: true } })
  fireEvent.click(button)
  expect(await screen.findByRole("button", { name: "Unsave" })).toBeInTheDocument()
})

test("does not write when state cannot be loaded, and offers retry", async () => {
  get.mockRejectedValue(new Error("offline"))
  render(<SavePostButton postNumber={7} />)
  await screen.findByRole("alert")
  expect(screen.getByRole("button", { name: "Save" })).toBeDisabled()
  expect(put).not.toHaveBeenCalled()
  get.mockResolvedValue({ ok: true, data: { saved: true } })
  fireEvent.click(screen.getByRole("button", { name: "Retry" }))
  expect(await screen.findByRole("button", { name: "Unsave" })).toBeEnabled()
})

test("ignores a stale response when the selected post changes", async () => {
  let resolveOld!: (value: any) => void
  get.mockReturnValueOnce(
    new Promise((resolve) => {
      resolveOld = resolve
    })
  )
  const { rerender } = render(<SavePostButton postNumber={7} />)
  rerender(<SavePostButton postNumber={8} />)
  await waitFor(() => expect(screen.getByRole("button", { name: "Save" })).toBeEnabled())
  await act(async () => {
    resolveOld({ ok: true, data: { saved: true } })
  })
  await waitFor(() => expect(get).toHaveBeenLastCalledWith("/api/v1/posts/8/saved"))
  expect(screen.getByRole("button", { name: "Save" })).toHaveAttribute("aria-pressed", "false")
})

test("prevents duplicate requests while a write is pending", async () => {
  put.mockReturnValue(
    new Promise(() => {
      /* Remain pending to exercise duplicate clicks. */
    })
  )
  render(<SavePostButton postNumber={7} />)
  const button = screen.getByRole("button", { name: "Save" })
  await waitFor(() => expect(button).toBeEnabled())
  fireEvent.click(button)
  fireEvent.click(button)
  expect(button).toBeDisabled()
  expect(put).toHaveBeenCalledTimes(1)
})
