package viewdata

// MaxWordLen is maximum a word can go in editor view
const MaxWordLen int = 15

// WordViewData stores content for word editor view
type WordViewData struct {
	// mistakes done during race
	Mistyped int

	// target word to type next
	TargetWord string
}
