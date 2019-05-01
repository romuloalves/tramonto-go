package entities

// Callback represents an instance of a callback function
type Callback interface {
	Invoke(data string)
}
