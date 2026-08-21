package store

var bucketNames = [][]byte{
	[]byte("journal_entries"),
	[]byte("review_records"),
	[]byte("resources"),
	[]byte("archive_records"),
	[]byte("domain_events"),
	[]byte("developer_profiles"),
}

const (
	EntriesBucket   = "journal_entries"
	ReviewsBucket   = "review_records"
	ResourcesBucket = "resources"
	ArchivesBucket  = "archive_records"
	EventsBucket    = "domain_events"
	ProfilesBucket  = "developer_profiles"
)

func bucketFor(kind string) string {
	switch kind {
	case EntriesBucket:
		return EntriesBucket
	case ReviewsBucket:
		return ReviewsBucket
	case ResourcesBucket:
		return ResourcesBucket
	case ArchivesBucket:
		return ArchivesBucket
	case EventsBucket:
		return EventsBucket
	case ProfilesBucket:
		return ProfilesBucket
	default:
		return ""
	}
}
