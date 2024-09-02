// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/likepost"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/rs/xid"
)

// LikePost is the model entity for the LikePost schema.
type LikePost struct {
	config `json:"-"`
	// ID of the ent.
	ID xid.ID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// AccountID holds the value of the "account_id" field.
	AccountID xid.ID `json:"account_id,omitempty"`
	// PostID holds the value of the "post_id" field.
	PostID xid.ID `json:"post_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LikePostQuery when eager-loading is set.
	Edges        LikePostEdges `json:"edges"`
	selectValues sql.SelectValues
}

// LikePostEdges holds the relations/edges for other nodes in the graph.
type LikePostEdges struct {
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
func (e LikePostEdges) AccountOrErr() (*Account, error) {
	if e.Account != nil {
		return e.Account, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: account.Label}
	}
	return nil, &NotLoadedError{edge: "account"}
}

// PostOrErr returns the Post value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LikePostEdges) PostOrErr() (*Post, error) {
	if e.Post != nil {
		return e.Post, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: post.Label}
	}
	return nil, &NotLoadedError{edge: "Post"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LikePost) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case likepost.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case likepost.FieldID, likepost.FieldAccountID, likepost.FieldPostID:
			values[i] = new(xid.ID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LikePost fields.
func (lp *LikePost) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case likepost.FieldID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				lp.ID = *value
			}
		case likepost.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				lp.CreatedAt = value.Time
			}
		case likepost.FieldAccountID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field account_id", values[i])
			} else if value != nil {
				lp.AccountID = *value
			}
		case likepost.FieldPostID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field post_id", values[i])
			} else if value != nil {
				lp.PostID = *value
			}
		default:
			lp.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LikePost.
// This includes values selected through modifiers, order, etc.
func (lp *LikePost) Value(name string) (ent.Value, error) {
	return lp.selectValues.Get(name)
}

// QueryAccount queries the "account" edge of the LikePost entity.
func (lp *LikePost) QueryAccount() *AccountQuery {
	return NewLikePostClient(lp.config).QueryAccount(lp)
}

// QueryPost queries the "Post" edge of the LikePost entity.
func (lp *LikePost) QueryPost() *PostQuery {
	return NewLikePostClient(lp.config).QueryPost(lp)
}

// Update returns a builder for updating this LikePost.
// Note that you need to call LikePost.Unwrap() before calling this method if this LikePost
// was returned from a transaction, and the transaction was committed or rolled back.
func (lp *LikePost) Update() *LikePostUpdateOne {
	return NewLikePostClient(lp.config).UpdateOne(lp)
}

// Unwrap unwraps the LikePost entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (lp *LikePost) Unwrap() *LikePost {
	_tx, ok := lp.config.driver.(*txDriver)
	if !ok {
		panic("ent: LikePost is not a transactional entity")
	}
	lp.config.driver = _tx.drv
	return lp
}

// String implements the fmt.Stringer.
func (lp *LikePost) String() string {
	var builder strings.Builder
	builder.WriteString("LikePost(")
	builder.WriteString(fmt.Sprintf("id=%v, ", lp.ID))
	builder.WriteString("created_at=")
	builder.WriteString(lp.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("account_id=")
	builder.WriteString(fmt.Sprintf("%v", lp.AccountID))
	builder.WriteString(", ")
	builder.WriteString("post_id=")
	builder.WriteString(fmt.Sprintf("%v", lp.PostID))
	builder.WriteByte(')')
	return builder.String()
}

// LikePosts is a parsable slice of LikePost.
type LikePosts []*LikePost
