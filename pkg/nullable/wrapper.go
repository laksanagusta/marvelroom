package nullable

import (
	"reflect"
	"time"
)

// Nullable represents a nullable value with additional methods
type Nullable[T any] struct {
	Val      T
	IsExists bool
}

// NewNullable creates a new Nullable value
func NewNullable[T any](val T) Nullable[T] {
	zero := *new(T)
	if reflect.DeepEqual(val, zero) {
		return Nullable[T]{IsExists: false}
	}
	return Nullable[T]{Val: val, IsExists: true}
}

// NewNullablePtr creates a new Nullable value from a pointer
func NewNullablePtr[T any](val *T) Nullable[T] {
	if val == nil {
		return Nullable[T]{IsExists: false}
	}
	return Nullable[T]{Val: *val, IsExists: true}
}

// GetOrDefault returns the value if exists, otherwise returns the zero value
func (n Nullable[T]) GetOrDefault() T {
	if n.IsExists {
		return n.Val
	}
	var zero T
	return zero
}

// GetOrDefaultWithDefault returns the value if exists, otherwise returns the provided default
func (n Nullable[T]) GetOrDefaultWithDefault(defaultVal T) T {
	if n.IsExists {
		return n.Val
	}
	return defaultVal
}

// String represents a nullable string
type String = Nullable[string]

// NewString creates a new nullable string
func NewString(val string) String {
	return NewNullable(val)
}

// NewStringPtr creates a new nullable string from a pointer
func NewStringPtr(val *string) String {
	if val == nil {
		return String{IsExists: false}
	}
	return String{Val: *val, IsExists: true}
}

// Time represents a nullable time
type Time = Nullable[time.Time]

// NewTime creates a new nullable time
func NewTime(val time.Time) Time {
	return NewNullable(val)
}

// NewTimePtr creates a new nullable time from a pointer
func NewTimePtr(val *time.Time) Time {
	if val == nil {
		return Time{IsExists: false}
	}
	return Time{Val: *val, IsExists: true}
}

// Int represents a nullable int
type Int = Nullable[int]

// NewInt creates a new nullable int
func NewInt(val int) Int {
	return NewNullable(val)
}

// NewIntPtr creates a new nullable int from a pointer
func NewIntPtr(val *int) Int {
	if val == nil {
		return Int{IsExists: false}
	}
	return Int{Val: *val, IsExists: true}
}

// Float64 represents a nullable float64
type Float64 = Nullable[float64]

// NewFloat64 creates a new nullable float64
func NewFloat64(val float64) Float64 {
	return NewNullable(val)
}

// NewFloat64Ptr creates a new nullable float64 from a pointer
func NewFloat64Ptr(val *float64) Float64 {
	if val == nil {
		return Float64{IsExists: false}
	}
	return Float64{Val: *val, IsExists: true}
}
