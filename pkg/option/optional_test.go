package option

import (
	"encoding/json"
	"testing"
)

func TestOption_Present(t *testing.T) {
	t.Run("test valid", func(t *testing.T) {
		var opt = Some("test")

		if !opt.Present() {
			t.Errorf("should be valid")
		}
	})

	t.Run("test none", func(t *testing.T) {
		var opt = None[string]()

		if opt.Present() {
			t.Error("none should not be present")
		}
	})
}

func TestJsonEncoding(t *testing.T) {
	t.Run("json encoding empty", func(t *testing.T) {
		var opt = None[string]()

		data, err := json.Marshal(struct {
			Opt Option[string] `json:"opt"`
		}{opt})
		if err != nil {
			t.Errorf("err encoding json: %s", err)
		}

		if string(data) != `{"opt":null}` {
			t.Error("none was not encoded properly")
		}
	})

	t.Run("json encoding with value", func(t *testing.T) {
		var opt = Some("i am here")

		data, err := json.Marshal(struct {
			Opt Option[string] `json:"opt"`
		}{opt})
		if err != nil {
			t.Errorf("err encoding json: %s", err)
		}

		if string(data) != `{"opt":"i am here"}` {
			t.Error("value was not encoded properly")
		}
	})

	t.Run("json decoding with value", func(t *testing.T) {
		var jsonB = []byte(`{"opt": "i am here"}`)

		var optHolder = struct {
			Opt Option[string] `json:"opt"`
		}{}

		err := json.Unmarshal(jsonB, &optHolder)
		if err != nil {
			t.Errorf("err decoding json: %s", err)
		}

		if !optHolder.Opt.Present() {
			t.Errorf("should be present")
		}

		if optHolder.Opt.MustGet() != "i am here" {
			t.Error("value was not decoded correctly")
		}
	})

	t.Run("json decoding none", func(t *testing.T) {
		var jsonB = []byte(`{"opt": null}`)

		var optNull = struct {
			Opt Option[testStr] `json:"opt"`
		}{}

		err := json.Unmarshal(jsonB, &optNull)
		if err != nil {
			t.Errorf("err encoding json: %s", err)
		}

		if optNull.Opt.Present() {
			t.Error("none should not be present")
		}
	})
}

func TestOption_Get(t *testing.T) {
	t.Run("gets a value", func(t *testing.T) {
		var opt = Some("blah blah")

		some, none := opt.Get()

		select {
		case v := <-some:
			if v != "blah blah" {
				t.Error("optional value was incorrect")
			}
		case <-none:
			t.Error("was not expecting none")
		}
	})

	t.Run("gets a none", func(t *testing.T) {
		var opt = None[string]()

		some, none := opt.Get()

		select {
		case <-some:
			t.Error("was not expecting some")
		case <-none:
			// we got a none :)
		}
	})
}

func TestOption_If(t *testing.T) {
	t.Run("gets a value", func(t *testing.T) {
		var opt = Some(123456)

		opt.If(func(v int) {
			if v != 123456 {
				t.Errorf("did not find expected value")
			}
		})
	})

	t.Run("gets a value", func(t *testing.T) {
		var opt = None[int]()

		opt.If(func(v int) {
			t.Errorf("did not expect execution")
		})
	})
}

func TestOption_OrElse(t *testing.T) {
	t.Run("existing value", func(t *testing.T) {
		var opt = Some("hello jimmy")

		value := opt.OrElse("hello scarlet")

		if value != "hello jimmy" {
			t.Errorf("default should not have been used")
		}
	})

	t.Run("missing value", func(t *testing.T) {
		var opt = None[string]()

		value := opt.OrElse("hello scarlet")

		if value != "hello scarlet" {
			t.Errorf("default should have been used")
		}
	})
}

type testStr struct {
	Label string `json:"label"`
}
