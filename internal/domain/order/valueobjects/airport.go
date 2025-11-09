package valueobjects

// Airport is a value object representing an airport.
type Airport struct {
	Code string // IATA code like "SFO"
	Name string
}

func NewAirport(code, name string) Airport { return Airport{Code: code, Name: name} }
