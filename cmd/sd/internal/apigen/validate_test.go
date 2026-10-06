package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestBindingsRejectContractDrift(t *testing.T) {
	api, err := openapi3.NewLoader().LoadFromData([]byte(`{"openapi":"3.0.3","info":{"title":"fixture","version":"1"},"paths":{"/things/{thing_id}":{"parameters":[{"in":"path","name":"thing_id","required":true,"schema":{"type":"string"}}],"get":{"operationId":"ThingGet","parameters":[{"in":"query","name":"tags","schema":{"type":"array","items":{"type":"string"}}}],"responses":{"200":{"description":"OK"}}}}}}`))
	require.NoError(t, err)
	valid := command{ID: "apiThingGet", Arguments: []parameter{{Name: "THING_ID"}}, Flags: []parameter{{Name: "tags", TrackChanged: true, Repeatable: true}}}
	require.NoError(t, validateBindings([]command{valid}, api))
	for _, tc := range []struct {
		name    string
		command command
		error   string
	}{
		{"missing operation", command{ID: "apiMissing"}, "no matching"},
		{"missing path", command{ID: "apiThingGet"}, "missing path:thing_id"},
		{"lost presence", command{ID: "apiThingGet", Arguments: valid.Arguments, Flags: []parameter{{Name: "tags", Repeatable: true}}}, "trackChanged"},
		{"lost repetition", command{ID: "apiThingGet", Arguments: valid.Arguments, Flags: []parameter{{Name: "tags", TrackChanged: true}}}, "repeatability"},
		{"unexpected parameter", command{ID: "apiThingGet", Arguments: valid.Arguments, Flags: []parameter{{Name: "tags", TrackChanged: true, Repeatable: true}, {Name: "typo", TrackChanged: true}}}, "unknown API parameters"},
	} {
		t.Run(tc.name, func(t *testing.T) { require.ErrorContains(t, validateBindings([]command{tc.command}, api), tc.error) })
	}
}

func TestSharedFlagsResolveBeforeBinding(t *testing.T) {
	commands := []command{{Commands: []command{{Flags: []parameter{{Ref: "#/components/flags/api-output"}}}}}}
	require.NoError(t, resolveFlags(commands, map[string]parameter{"api-output": {Name: "output"}}))
	require.Equal(t, "output", commands[0].Commands[0].Flags[0].Name)
	require.ErrorContains(t, resolveFlags([]command{{Flags: []parameter{{Ref: "#/components/flags/missing"}}}}, nil), "unknown flag reference")
}
