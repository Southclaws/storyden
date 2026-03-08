package rpc

import (
	"encoding/json"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/go-viper/mapstructure/v2"
)

func (m *Manifest) ToMap() map[string]any {
	var v map[string]any
	b, err := json.Marshal(m)
	if err != nil {
		return map[string]any{}
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return map[string]any{}
	}
	return v
}

func (r *DatagraphRef) ToDomain() datagraph.Ref {
	k, _ := datagraph.NewKind(r.Kind)
	return datagraph.Ref{
		Kind: k,
		ID:   r.ID,
	}
}

func DatagraphRefToRPC(r datagraph.Ref) DatagraphRef {
	return DatagraphRef{
		Kind: r.Kind.String(),
		ID:   r.ID,
	}
}

func SerialiseSettings(s settings.Settings) map[string]any {
	var v map[string]any
	mapstructure.Decode(s, &v)
	return v
}
