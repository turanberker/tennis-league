package outbox

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusPublished Status = "PUBLISHED"
	StatusProcessed Status = "PROCESSED"
	StatusFailed    Status = "FAILED"
)

type PersistEntity struct {
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
}

type EventToPublish struct {
	Id        string
	EventType string
	Payload   []byte
}
