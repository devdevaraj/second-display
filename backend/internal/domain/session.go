package domain

// SessionStatus represents the current state of a virtual display session.
type SessionStatus string

const (
	StatusCreated SessionStatus = "created"
	StatusRunning SessionStatus = "running"
	StatusStopped SessionStatus = "stopped"
	StatusError   SessionStatus = "error"
)

// Session represents a virtual display and its streaming configuration.
type Session struct {
	ID         string        `json:"id"`
	OutputName string        `json:"outputName"`
	Resolution string        `json:"resolution"`
	FPS        int           `json:"fps"`
	Bitrate    int           `json:"bitrate"`
	Status     SessionStatus `json:"status"`
	StreamURL  string        `json:"streamUrl,omitempty"`
}

// SessionRepository defines how sessions are persisted/retrieved.
type SessionRepository interface {
	Save(session *Session) error
	Get(id string) (*Session, error)
	Delete(id string) error
	List() ([]*Session, error)
}
