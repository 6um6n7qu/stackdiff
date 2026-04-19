// Package throttle provides a rolling-window event rate limiter for drift
// notifications. It prevents alert storms when many config keys change at
// once by capping the number of drift events forwarded within a configurable
// time window.
package throttle
