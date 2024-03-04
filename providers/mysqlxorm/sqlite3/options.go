package sqlite3

type Options struct {
	JournalMode JournalMode
}

type OptionFunc func(options *Options)

type JournalMode string

const (
	DELETE   JournalMode = "delete"
	TRUNCATE JournalMode = "truncate"
	PERSIST  JournalMode = "persist"
	MEMORY   JournalMode = "memory"
	OFF      JournalMode = "off"
	WAL      JournalMode = "wal"
)

func WithJournalMode(mode JournalMode) OptionFunc {
	return func(o *Options) {
		o.JournalMode = mode
	}
}
