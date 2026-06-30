package outbound

// Adapter is an interface that all outbound services should implement to keep the application resilient.
// It mandates a Ping check before any task execution to ensure the external connection is alive and healthy.
type Adapter interface {
	// Ping checks if the external service is reachable and functional.
	// It should be called before performing any task to avoid unexpected breaks.
	Ping() bool
}
