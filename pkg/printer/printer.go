package printer

// Printer prints the entities (for now ns, engines and KV secrets).
type Printer interface {
	// Out prints out the entities.
	Out(secrets interface{}) error
}
