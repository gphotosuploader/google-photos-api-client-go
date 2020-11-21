package resumable

// MockedSessionStorer mocks a service to store resumable upload data.
type MockedSessionStorer struct {
	GetFn    func(f string) []byte
	SetFn    func(f string, u []byte)
	DeleteFn func(f string)
}

// Get invokes the mock implementation and marks the function as invoked.
func (s MockedSessionStorer) Get(f string) []byte {
	return s.GetFn(f)
}

// Set invokes the mock implementation and marks the function as invoked.
func (s MockedSessionStorer) Set(f string, u []byte) {
	s.SetFn(f, u)
}

// Delete invokes the mock implementation and marks the function as invoked.
func (s MockedSessionStorer) Delete(f string) {
	s.DeleteFn(f)
}
