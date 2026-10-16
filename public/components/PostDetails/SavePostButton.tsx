import React, { useEffect, useRef, useState } from "react"
import { Trans } from "@lingui/react/macro"
import { useFider } from "@fider/hooks"
import { http } from "@fider/services"
import "@fider/pages/ShowPost/components/ActionButton.scss"

export const SavePostButton = ({ postNumber }: { postNumber: number }) => {
  const fider = useFider()
  const authenticated = fider.session.isAuthenticated
  const userID = authenticated ? fider.session.user.id : undefined
  const tenantID = fider.session.tenant.id
  const [saved, setSaved] = useState(false)
  const [ready, setReady] = useState(false)
  const [busy, setBusy] = useState(false)
  const [failed, setFailed] = useState(false)
  const [retry, setRetry] = useState(0)
  const pending = useRef(false)
  const generation = useRef(0)
  const url = `/api/v1/posts/${postNumber}/saved`

  useEffect(() => {
    const current = ++generation.current
    setReady(false)
    setFailed(false)
    setSaved(false)
    setBusy(false)
    pending.current = false
    if (authenticated) {
      http
        .get<{ saved: boolean }>(url)
        .then((result) => {
          if (generation.current !== current) return
          if (result.ok) {
            setSaved(result.data.saved)
            setReady(true)
          } else {
            setFailed(true)
          }
        })
        .catch(() => {
          if (generation.current === current) setFailed(true)
        })
    }
    return () => {
      generation.current++
    }
  }, [url, authenticated, userID, tenantID, retry])

  if (!authenticated) return null

  const toggle = async () => {
    if (pending.current || !ready) return
    pending.current = true
    const current = generation.current
    setBusy(true)
    setFailed(false)
    try {
      const result = saved ? await http.delete<{ saved: boolean }>(url) : await http.put<{ saved: boolean }>(url)
      if (generation.current !== current) return
      if (result.ok) setSaved(result.data.saved)
      else setFailed(true)
    } catch {
      if (generation.current === current) setFailed(true)
    } finally {
      if (generation.current === current) {
        pending.current = false
        setBusy(false)
      }
    }
  }

  return (
    <div>
      <button
        type="button"
        className="c-action-button"
        aria-pressed={saved}
        aria-busy={busy || (!ready && !failed)}
        disabled={!ready || busy || fider.isReadOnly}
        onClick={toggle}
      >
        {saved ? <Trans id="saved.unsave">Unsave</Trans> : <Trans id="saved.save">Save</Trans>}
      </button>
      {failed && (
        <div role="alert" className="text-sm">
          <Trans id="saved.error">Could not load or update saved feedback. Please try again.</Trans>{" "}
          {!ready && (
            <button type="button" className="c-action-button" onClick={() => setRetry(retry + 1)}>
              <Trans id="saved.retry">Retry</Trans>
            </button>
          )}
        </div>
      )}
    </div>
  )
}
