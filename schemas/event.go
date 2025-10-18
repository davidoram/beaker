package schemas

type Event interface {
	// Returns the NATS subject that this event will be published to.
	Subject() string
}
