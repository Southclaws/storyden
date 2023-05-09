// Code generated by ent, DO NOT EDIT.

package react

import (
	"time"

	"github.com/rs/xid"
)

const (
	// Label holds the string label denoting the react type in the database.
	Label = "react"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldAccountID holds the string denoting the account_id field in the database.
	FieldAccountID = "account_id"
	// FieldPostID holds the string denoting the post_id field in the database.
	FieldPostID = "post_id"
	// FieldEmoji holds the string denoting the emoji field in the database.
	FieldEmoji = "emoji"
	// EdgeAccount holds the string denoting the account edge name in mutations.
	EdgeAccount = "account"
	// EdgePost holds the string denoting the post edge name in mutations.
	EdgePost = "Post"
	// Table holds the table name of the react in the database.
	Table = "reacts"
	// AccountTable is the table that holds the account relation/edge.
	AccountTable = "reacts"
	// AccountInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	AccountInverseTable = "accounts"
	// AccountColumn is the table column denoting the account relation/edge.
	AccountColumn = "account_id"
	// PostTable is the table that holds the Post relation/edge.
	PostTable = "reacts"
	// PostInverseTable is the table name for the Post entity.
	// It exists in this package in order to avoid circular dependency with the "post" package.
	PostInverseTable = "posts"
	// PostColumn is the table column denoting the Post relation/edge.
	PostColumn = "post_id"
)

// Columns holds all SQL columns for react fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldAccountID,
	FieldPostID,
	FieldEmoji,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() xid.ID
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)
