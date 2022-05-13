package sne

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

// SoNoEr is an optional value.
type SoNoEr[T any] struct {
	present bool
	value   T
	err     error

	l sync.Mutex
}

// Present is true when a value is present.
func (o SoNoEr[T]) Present() bool {
	return o.present
}

// Some represents an optional with a value.
func Some[T any](value T) SoNoEr[T] {
	var opt = SoNoEr[T]{value: value, present: true}

	// handle pointer types values
	if val := reflect.ValueOf(value); val.Kind() == reflect.Ptr && val.IsNil() {
		opt.present = false
	}

	return opt
}

// None represents an optional with no value.
func None[T any]() SoNoEr[T] {
	return SoNoEr[T]{}
}

func Err[T any](err error) SoNoEr[T] {
	n := None[T]()
	n.err = err
	return n
}

// Get a some, none or err.
// When the optional is Some, then the value is sent on that channel.
// If there is an error then it is sent to err channel.
// If the optional is None, then an empty struct is sent on none channel.
func (o SoNoEr[T]) Get() (some chan T, none chan struct{}, err chan error) {
	some = make(chan T)
	none = make(chan struct{})
	err = make(chan error)

	go func() {
		defer close(none)
		defer close(some)
		defer close(err)

		if o.present {
			o.l.Lock()
			some <- o.value
			o.l.Unlock()
		} else if o.err != nil {
			err <- o.err
		} else {
			none <- struct{}{}
		}
	}()

	return some, none, err
}

// MustGet either returns the value or panics.
func (o SoNoEr[T]) MustGet() T {
	if !o.Present() {
		panic("value not present")
	}

	return o.value
}

// If calls fn only if a value is present.
func (o SoNoEr[T]) If(fn func(v T)) {
	if o.present {
		fn(o.value)
	}
}

// OrElse returns the value if it is present or returns def.
func (o SoNoEr[T]) OrElse(def T) T {
	if o.present {
		return o.value
	}

	return def
}

// JSON handling

func (o SoNoEr[T]) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	} else if o.err != nil {
		// we do not allow json encoding of errs because decoding
		// them is ambiguous. We cannot tell if a string is an error or a regular
		// string value, so it is better practice not to allow encoding.
		return nil, errors.New("err cannot encode errors to JSON")
	}
	return json.Marshal(nil)
}

func (o *SoNoEr[T]) UnmarshalJSON(data []byte) error {
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
