// Code generated by entc, DO NOT EDIT.

package react

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the react type in the database.
	Label = "react"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEmoji holds the string denoting the emoji field in the database.
	FieldEmoji = "emoji"
	// FieldCreatedAt holds the string denoting the createdat field in the database.
	FieldCreatedAt = "created_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgePost holds the string denoting the post edge name in mutations.
	EdgePost = "Post"
	// Table holds the table name of the react in the database.
	Table = "reacts"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "reacts"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "react_user"
	// PostTable is the table that holds the Post relation/edge.
	PostTable = "reacts"
	// PostInverseTable is the table name for the Post entity.
	// It exists in this package in order to avoid circular dependency with the "post" package.
	PostInverseTable = "posts"
	// PostColumn is the table column denoting the Post relation/edge.
	PostColumn = "react_post"
)

// Columns holds all SQL columns for react fields.
var Columns = []string{
	FieldID,
	FieldEmoji,
	FieldCreatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "reacts"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"post_reacts",
	"react_user",
	"react_post",
	"user_reacts",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "createdAt" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
