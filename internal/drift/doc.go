// Package drift evaluates the overall severity of configuration drift
// detected between two service configs.
//
// It maps raw diff entry counts into a named Level (none, low, moderate, high)
// so callers can make policy decisions or trigger alerts based on severity
// rather than raw numbers.
package drift
