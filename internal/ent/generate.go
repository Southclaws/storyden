package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/execquery --feature sql/upsert --feature sql/modifier --feature sql/upsert ./schema
//go:generate go run -mod=mod github.com/a8m/enter ./schema
