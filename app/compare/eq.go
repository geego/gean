package compare

// Eqer can be used to determine if this value is equal to the other.
// The semantics of equals is that the two value are interchangeable
type Eqer interface {
	Eq(other interface{}) bool
}
