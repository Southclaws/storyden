// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/invitation"
	"github.com/rs/xid"
)

// Invitation is the model entity for the Invitation schema.
type Invitation struct {
	config `json:"-"`
	// ID of the ent.
	ID xid.ID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// Message holds the value of the "message" field.
	Message *string `json:"message,omitempty"`
	// CreatorAccountID holds the value of the "creator_account_id" field.
	CreatorAccountID xid.ID `json:"creator_account_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the InvitationQuery when eager-loading is set.
	Edges        InvitationEdges `json:"edges"`
	selectValues sql.SelectValues
}

// InvitationEdges holds the relations/edges for other nodes in the graph.
type InvitationEdges struct {
	// Creator holds the value of the creator edge.
	Creator *Account `json:"creator,omitempty"`
	// Invited holds the value of the invited edge.
	Invited []*Account `json:"invited,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// CreatorOrErr returns the Creator value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e InvitationEdges) CreatorOrErr() (*Account, error) {
	if e.Creator != nil {
		return e.Creator, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: account.Label}
	}
	return nil, &NotLoadedError{edge: "creator"}
}

// InvitedOrErr returns the Invited value or an error if the edge
// was not loaded in eager-loading.
func (e InvitationEdges) InvitedOrErr() ([]*Account, error) {
	if e.loadedTypes[1] {
		return e.Invited, nil
	}
	return nil, &NotLoadedError{edge: "invited"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Invitation) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case invitation.FieldMessage:
			values[i] = new(sql.NullString)
		case invitation.FieldCreatedAt, invitation.FieldUpdatedAt, invitation.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		case invitation.FieldID, invitation.FieldCreatorAccountID:
			values[i] = new(xid.ID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Invitation fields.
func (i *Invitation) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case invitation.FieldID:
			if value, ok := values[j].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value != nil {
				i.ID = *value
			}
		case invitation.FieldCreatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[j])
			} else if value.Valid {
				i.CreatedAt = value.Time
			}
		case invitation.FieldUpdatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[j])
			} else if value.Valid {
				i.UpdatedAt = value.Time
			}
		case invitation.FieldDeletedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[j])
			} else if value.Valid {
				i.DeletedAt = new(time.Time)
				*i.DeletedAt = value.Time
			}
		case invitation.FieldMessage:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[j])
			} else if value.Valid {
				i.Message = new(string)
				*i.Message = value.String
			}
		case invitation.FieldCreatorAccountID:
			if value, ok := values[j].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field creator_account_id", values[j])
			} else if value != nil {
				i.CreatorAccountID = *value
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Invitation.
// This includes values selected through modifiers, order, etc.
func (i *Invitation) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QueryCreator queries the "creator" edge of the Invitation entity.
func (i *Invitation) QueryCreator() *AccountQuery {
	return NewInvitationClient(i.config).QueryCreator(i)
}

// QueryInvited queries the "invited" edge of the Invitation entity.
func (i *Invitation) QueryInvited() *AccountQuery {
	return NewInvitationClient(i.config).QueryInvited(i)
}

// Update returns a builder for updating this Invitation.
// Note that you need to call Invitation.Unwrap() before calling this method if this Invitation
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Invitation) Update() *InvitationUpdateOne {
	return NewInvitationClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Invitation entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Invitation) Unwrap() *Invitation {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Invitation is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Invitation) String() string {
	var builder strings.Builder
	builder.WriteString("Invitation(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("created_at=")
	builder.WriteString(i.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(i.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := i.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := i.Message; v != nil {
		builder.WriteString("message=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("creator_account_id=")
	builder.WriteString(fmt.Sprintf("%v", i.CreatorAccountID))
	builder.WriteByte(')')
	return builder.String()
}

// Invitations is a parsable slice of Invitation.
type Invitations []*Invitation
