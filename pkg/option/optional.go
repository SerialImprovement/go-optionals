package option

import (
	"encoding/json"
	"reflect"
	"sync"
)

// Option is an optional value.
type Option[T any] struct {
	present bool
	value   T

	l sync.Mutex
}

// Some represents an optional with a value.
func Some[T any](value T) Option[T] {
	var opt = Option[T]{value: value, present: true}

	// handle pointer types values
	if val := reflect.ValueOf(value); val.Kind() == reflect.Ptr && val.IsNil() {
		opt.present = false
	}

	return opt
}

// None represents an optional with no value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// Present is true when a value is present.
func (o Option[T]) Present() bool {
	return o.present
}

// Get a value or none.
//
// returns a some & none channel.
// When the optional is Some, then the value is sent on that channel.
// If the optional is None, then an empty struct is sent on none channel.
func (o Option[T]) Get() (some chan T, none chan struct{}) {
	some = make(chan T)
	none = make(chan struct{})

	go func() {
		defer close(none)
		defer close(some)

		if o.present {
			o.l.Lock()
			some <- o.value
			o.l.Unlock()
		} else {
			none <- struct{}{}
		}
	}()

	return some, none
}

// MustGet either returns the value or panics.
func (o Option[T]) MustGet() T {
	if !o.Present() {
		panic("value not present")
	}

	return o.value
}

// If calls fn only if a value is present.
func (o Option[T]) If(fn func(v T)) {
	if o.present {
		fn(o.value)
	}
}

// OrElse returns the value if it is present or returns def.
func (o Option[T]) OrElse(def T) T {
	if o.present {
		return o.value
	}

	return def
}

// JSON handling

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	var value T

	if string(data) == "null" {
		o.present = false
		return nil
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.present = true
	o.value = value

	return nil
}
