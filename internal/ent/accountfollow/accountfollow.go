// Code generated by ent, DO NOT EDIT.

package accountfollow

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/rs/xid"
)

const (
	// Label holds the string label denoting the accountfollow type in the database.
	Label = "account_follow"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldFollowerAccountID holds the string denoting the follower_account_id field in the database.
	FieldFollowerAccountID = "follower_account_id"
	// FieldFollowingAccountID holds the string denoting the following_account_id field in the database.
	FieldFollowingAccountID = "following_account_id"
	// EdgeFollower holds the string denoting the follower edge name in mutations.
	EdgeFollower = "follower"
	// EdgeFollowing holds the string denoting the following edge name in mutations.
	EdgeFollowing = "following"
	// Table holds the table name of the accountfollow in the database.
	Table = "account_follows"
	// FollowerTable is the table that holds the follower relation/edge.
	FollowerTable = "account_follows"
	// FollowerInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	FollowerInverseTable = "accounts"
	// FollowerColumn is the table column denoting the follower relation/edge.
	FollowerColumn = "follower_account_id"
	// FollowingTable is the table that holds the following relation/edge.
	FollowingTable = "account_follows"
	// FollowingInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	FollowingInverseTable = "accounts"
	// FollowingColumn is the table column denoting the following relation/edge.
	FollowingColumn = "following_account_id"
)

// Columns holds all SQL columns for accountfollow fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldFollowerAccountID,
	FieldFollowingAccountID,
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

// OrderOption defines the ordering options for the AccountFollow queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByFollowerAccountID orders the results by the follower_account_id field.
func ByFollowerAccountID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFollowerAccountID, opts...).ToFunc()
}

// ByFollowingAccountID orders the results by the following_account_id field.
func ByFollowingAccountID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFollowingAccountID, opts...).ToFunc()
}

// ByFollowerField orders the results by follower field.
func ByFollowerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFollowerStep(), sql.OrderByField(field, opts...))
	}
}

// ByFollowingField orders the results by following field.
func ByFollowingField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFollowingStep(), sql.OrderByField(field, opts...))
	}
}
func newFollowerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FollowerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, FollowerTable, FollowerColumn),
	)
}
func newFollowingStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FollowingInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, FollowingTable, FollowingColumn),
	)
}
