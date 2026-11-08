package config

import (
	"errors"
	"testing"
)

func TestFixture(t *testing.T) {
	values := map[string]string{}
	lookup := func(key string) (string, bool) {
		value, present := values[key]
		return value, present
	}

	if _, err := load(lookup); !errors.Is(err, errInvalidConfiguration) {
		t.Fatal("missing required token accepted", err)
	}

	values["TOKEN"] = ""
	if _, err := load(lookup); !errors.Is(err, errInvalidConfiguration) {
		t.Fatal("empty required token accepted", err)
	}

	values["TOKEN"] = "present"
	c, err := load(lookup)
	if err != nil {
		t.Fatal(err)
	}

	if c.Token != "present" || c.Limit != 5 || c.RequiredDefault != "fallback" {
		t.Fatal("incorrect defaults", c)
	}

	if string(c.Payload) != "raw" || len(c.Origins) != 2 || c.Origins[1] != "two" {
		t.Fatal("incorrect collection defaults", c)
	}

	values["LIMIT"] = "0"
	c, err = load(lookup)
	if err != nil || c.Limit != 5 {
		t.Fatal("zero default not applied", c, err)
	}

	values["LIMIT"] = "11"
	if _, err := load(lookup); !errors.Is(err, errInvalidConfiguration) {
		t.Fatal("maximum not enforced", err)
	}

	values["LIMIT"] = "-1"
	if _, err := load(lookup); !errors.Is(err, errInvalidConfiguration) {
		t.Fatal("minimum not enforced", err)
	}

	values["LIMIT"] = "5"
	values["REQUIRED_DEFAULT"] = ""
	if _, err := load(lookup); !errors.Is(err, errInvalidConfiguration) {
		t.Fatal("explicit empty required value accepted", err)
	}
}
