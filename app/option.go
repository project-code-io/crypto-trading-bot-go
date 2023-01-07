package app

// Option allows for overriding of the internals of the application. These
// options are typically only meant for internal testing.
type Option func(a *App)

// WithIDGenerator overrides the internal ID generator of the app. Use this
// method for testing or for use with an exchange that does not support the
// id format in use by the app.
func WithIDGenerator(gen IDGenerator) Option {
	return func(a *App) {
		a.idGenerator = gen
	}
}
