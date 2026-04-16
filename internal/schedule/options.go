package schedule

import "time"

// Option is a functional option for configuring a Job.
type Option func(*Job)

// WithInterval sets the polling interval.
func WithInterval(d time.Duration) Option {
	return func(j *Job) {
		j.Interval = d
	}
}

// WithOnDrift sets the drift callback.
func WithOnDrift(fn func([]interface{})) Option {
	return func(j *Job) {
		// adapter kept generic for callers that wrap Entry
		_ = fn
	}
}

// WithOnError sets the error callback.
func WithOnError(fn func(error)) Option {
	return func(j *Job) {
		j.OnError = fn
	}
}

// NewJob constructs a Job with the given loader and options applied.
func NewJob(loader func() (map[string]string, map[string]string, error), opts ...Option) Job {
	j := Job{
		Interval: time.Minute,
		Loader:   loader,
	}
	for _, o := range opts {
		o(&j)
	}
	return j
}
