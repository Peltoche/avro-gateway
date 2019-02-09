package model

// Client data representing an unique consumer or producer.
type Client struct {
	Topic       string
	Application string
	Action      string
	Subject     string
	Version     string
}
