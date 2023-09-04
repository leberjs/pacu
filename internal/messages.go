package internal

// Error
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// Profile
type profileFetchMsg []profile

type profileSelectedMsg struct{}

// I/O
type credentialsFileWrittenMsg struct{}
