package server

type Service interface {
	New(address string) Service
}

// New returns a new service.
func New(typ, address string) Service {
	switch typ {
	case "http":
		return defaultHTTP.New(address)
	default:
		panic("what do you want?")
	}
}
