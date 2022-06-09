package api

// MainAPI exports the internal User related functions.
type MainAPI struct {
	// validator auth.Validator
}

// NewMainAPI creates a new api for user funcs.
func NewMainAPI() *MainAPI { // validator auth.Validator
	return &MainAPI{
		// validator: validator,
	}
}
