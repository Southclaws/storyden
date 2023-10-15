// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/react"
	"github.com/rs/xid"
)

// React is the model entity for the React schema.
type React struct {
	config `json:"-"`
	// ID of the ent.
	ID xid.ID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// AccountID holds the value of the "account_id" field.
	AccountID xid.ID `json:"account_id,omitempty"`
	// PostID holds the value of the "post_id" field.
	PostID xid.ID `json:"post_id,omitempty"`
	// Emoji holds the value of the "emoji" field.
	Emoji string `json:"emoji,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ReactQuery when eager-loading is set.
	Edges        ReactEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ReactEdges holds the relations/edges for other nodes in the graph.
type ReactEdges struct {
	// Account holds the value of the account edge.
	Account *Account `json:"account,omitempty"`
	// Post holds the value of the Post edge.
	Post *Post `json:"Post,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// AccountOrErr returns the Account value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ReactEdges) AccountOrErr() (*Account, error) {
	if e.loadedTypes[0] {
		if e.Account == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: account.Label}
		}
		return e.Account, nil
	}
	return nil, &NotLoadedError{edge: "account"}
}

// PostOrErr returns the Post value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ReactEdges) PostOrErr() (*Post, error) {
	if e.loadedTypes[1] {
		if e.Post == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: post.Label}
		}
		return e.Post, nil
	}
	return nil, &NotLoadedError{edge: "Post"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*React) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case react.FieldEmoji:
			values[i] = new(sql.NullString)
		case react.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case react.FieldID, react.FieldAccountID, react.FieldPostID:
			values[i] = new(xid.ID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the React fields.
func (r *React) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case react.FieldID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				r.ID = *value
			}
		case react.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				r.CreatedAt = value.Time
			}
		case react.FieldAccountID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field account_id", values[i])
			} else if value != nil {
				r.AccountID = *value
			}
		case react.FieldPostID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field post_id", values[i])
			} else if value != nil {
				r.PostID = *value
			}
		case react.FieldEmoji:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field emoji", values[i])
			} else if value.Valid {
				r.Emoji = value.String
			}
		default:
			r.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the React.
// This includes values selected through modifiers, order, etc.
func (r *React) Value(name string) (ent.Value, error) {
	return r.selectValues.Get(name)
}

// QueryAccount queries the "account" edge of the React entity.
func (r *React) QueryAccount() *AccountQuery {
	return NewReactClient(r.config).QueryAccount(r)
}

// QueryPost queries the "Post" edge of the React entity.
func (r *React) QueryPost() *PostQuery {
	return NewReactClient(r.config).QueryPost(r)
}

// Update returns a builder for updating this React.
// Note that you need to call React.Unwrap() before calling this method if this React
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *React) Update() *ReactUpdateOne {
	return NewReactClient(r.config).UpdateOne(r)
}

// Unwrap unwraps the React entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *React) Unwrap() *React {
	_tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: React is not a transactional entity")
	}
	r.config.driver = _tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *React) String() string {
	var builder strings.Builder
	builder.WriteString("React(")
	builder.WriteString(fmt.Sprintf("id=%v, ", r.ID))
	builder.WriteString("created_at=")
	builder.WriteString(r.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("account_id=")
	builder.WriteString(fmt.Sprintf("%v", r.AccountID))
	builder.WriteString(", ")
	builder.WriteString("post_id=")
	builder.WriteString(fmt.Sprintf("%v", r.PostID))
	builder.WriteString(", ")
	builder.WriteString("emoji=")
	builder.WriteString(r.Emoji)
	builder.WriteByte(')')
	return builder.String()
}

// Reacts is a parsable slice of React.
type Reacts []*React
