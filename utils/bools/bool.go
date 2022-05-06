package bools

// True is true pointer
var True = func() *bool {
	b := true
	return &b
}()

// False is false pointer
var False = func() *bool {
	b := false
	return &b
}()

// IsTrue returns true when p is not nil and with true value
func IsTrue(p *bool) bool {
	return p != nil && *p
}
