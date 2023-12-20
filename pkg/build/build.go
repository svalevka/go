package build

// revision is the reference of the Git HEAD at build time, or an empty string
// if not set.
var revision string = ""

// GetRevision returns the optionally truncated reference of the Git HEAD at
// build time. If the truncation length is longer than the revision length, the
// whole revision is returned. If the revision is not set, an empty string is
// returned.
func GetRevision(n int) string {
	if n < 0 || n > len(revision) {
		return revision
	}

	return revision[:n]
}
