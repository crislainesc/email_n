package campaign

type Status uint8

const (
	Created Status = iota
	Pending
	Started
	Done
)

func (s Status) String() string {
	switch s {
	case Created:
		return "Created"
	case Pending:
		return "Pending"
	case Started:
		return "Started"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}
}
