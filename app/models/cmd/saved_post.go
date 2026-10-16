package cmd

// SetPostSaved changes only the current user's private saved feedback.
type SetPostSaved struct {
	Number int
	Saved  bool
}
