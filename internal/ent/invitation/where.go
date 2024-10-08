// Code generated by ent, DO NOT EDIT.

package invitation

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// ID filters vertices based on their ID field.
func ID(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldDeletedAt, v))
}

// Message applies equality check predicate on the "message" field. It's identical to MessageEQ.
func Message(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldMessage, v))
}

// CreatorAccountID applies equality check predicate on the "creator_account_id" field. It's identical to CreatorAccountIDEQ.
func CreatorAccountID(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldCreatorAccountID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Invitation {
	return predicate.Invitation(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Invitation {
	return predicate.Invitation(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Invitation {
	return predicate.Invitation(sql.FieldNotNull(FieldDeletedAt))
}

// MessageEQ applies the EQ predicate on the "message" field.
func MessageEQ(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldMessage, v))
}

// MessageNEQ applies the NEQ predicate on the "message" field.
func MessageNEQ(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldNEQ(FieldMessage, v))
}

// MessageIn applies the In predicate on the "message" field.
func MessageIn(vs ...string) predicate.Invitation {
	return predicate.Invitation(sql.FieldIn(FieldMessage, vs...))
}

// MessageNotIn applies the NotIn predicate on the "message" field.
func MessageNotIn(vs ...string) predicate.Invitation {
	return predicate.Invitation(sql.FieldNotIn(FieldMessage, vs...))
}

// MessageGT applies the GT predicate on the "message" field.
func MessageGT(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldGT(FieldMessage, v))
}

// MessageGTE applies the GTE predicate on the "message" field.
func MessageGTE(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldGTE(FieldMessage, v))
}

// MessageLT applies the LT predicate on the "message" field.
func MessageLT(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldLT(FieldMessage, v))
}

// MessageLTE applies the LTE predicate on the "message" field.
func MessageLTE(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldLTE(FieldMessage, v))
}

// MessageContains applies the Contains predicate on the "message" field.
func MessageContains(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldContains(FieldMessage, v))
}

// MessageHasPrefix applies the HasPrefix predicate on the "message" field.
func MessageHasPrefix(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldHasPrefix(FieldMessage, v))
}

// MessageHasSuffix applies the HasSuffix predicate on the "message" field.
func MessageHasSuffix(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldHasSuffix(FieldMessage, v))
}

// MessageIsNil applies the IsNil predicate on the "message" field.
func MessageIsNil() predicate.Invitation {
	return predicate.Invitation(sql.FieldIsNull(FieldMessage))
}

// MessageNotNil applies the NotNil predicate on the "message" field.
func MessageNotNil() predicate.Invitation {
	return predicate.Invitation(sql.FieldNotNull(FieldMessage))
}

// MessageEqualFold applies the EqualFold predicate on the "message" field.
func MessageEqualFold(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldEqualFold(FieldMessage, v))
}

// MessageContainsFold applies the ContainsFold predicate on the "message" field.
func MessageContainsFold(v string) predicate.Invitation {
	return predicate.Invitation(sql.FieldContainsFold(FieldMessage, v))
}

// CreatorAccountIDEQ applies the EQ predicate on the "creator_account_id" field.
func CreatorAccountIDEQ(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldEQ(FieldCreatorAccountID, v))
}

// CreatorAccountIDNEQ applies the NEQ predicate on the "creator_account_id" field.
func CreatorAccountIDNEQ(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldNEQ(FieldCreatorAccountID, v))
}

// CreatorAccountIDIn applies the In predicate on the "creator_account_id" field.
func CreatorAccountIDIn(vs ...xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldIn(FieldCreatorAccountID, vs...))
}

// CreatorAccountIDNotIn applies the NotIn predicate on the "creator_account_id" field.
func CreatorAccountIDNotIn(vs ...xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldNotIn(FieldCreatorAccountID, vs...))
}

// CreatorAccountIDGT applies the GT predicate on the "creator_account_id" field.
func CreatorAccountIDGT(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldGT(FieldCreatorAccountID, v))
}

// CreatorAccountIDGTE applies the GTE predicate on the "creator_account_id" field.
func CreatorAccountIDGTE(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldGTE(FieldCreatorAccountID, v))
}

// CreatorAccountIDLT applies the LT predicate on the "creator_account_id" field.
func CreatorAccountIDLT(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldLT(FieldCreatorAccountID, v))
}

// CreatorAccountIDLTE applies the LTE predicate on the "creator_account_id" field.
func CreatorAccountIDLTE(v xid.ID) predicate.Invitation {
	return predicate.Invitation(sql.FieldLTE(FieldCreatorAccountID, v))
}

// CreatorAccountIDContains applies the Contains predicate on the "creator_account_id" field.
func CreatorAccountIDContains(v xid.ID) predicate.Invitation {
	vc := v.String()
	return predicate.Invitation(sql.FieldContains(FieldCreatorAccountID, vc))
}

// CreatorAccountIDHasPrefix applies the HasPrefix predicate on the "creator_account_id" field.
func CreatorAccountIDHasPrefix(v xid.ID) predicate.Invitation {
	vc := v.String()
	return predicate.Invitation(sql.FieldHasPrefix(FieldCreatorAccountID, vc))
}

// CreatorAccountIDHasSuffix applies the HasSuffix predicate on the "creator_account_id" field.
func CreatorAccountIDHasSuffix(v xid.ID) predicate.Invitation {
	vc := v.String()
	return predicate.Invitation(sql.FieldHasSuffix(FieldCreatorAccountID, vc))
}

// CreatorAccountIDEqualFold applies the EqualFold predicate on the "creator_account_id" field.
func CreatorAccountIDEqualFold(v xid.ID) predicate.Invitation {
	vc := v.String()
	return predicate.Invitation(sql.FieldEqualFold(FieldCreatorAccountID, vc))
}

// CreatorAccountIDContainsFold applies the ContainsFold predicate on the "creator_account_id" field.
func CreatorAccountIDContainsFold(v xid.ID) predicate.Invitation {
	vc := v.String()
	return predicate.Invitation(sql.FieldContainsFold(FieldCreatorAccountID, vc))
}

// HasCreator applies the HasEdge predicate on the "creator" edge.
func HasCreator() predicate.Invitation {
	return predicate.Invitation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CreatorTable, CreatorColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCreatorWith applies the HasEdge predicate on the "creator" edge with a given conditions (other predicates).
func HasCreatorWith(preds ...predicate.Account) predicate.Invitation {
	return predicate.Invitation(func(s *sql.Selector) {
		step := newCreatorStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasInvited applies the HasEdge predicate on the "invited" edge.
func HasInvited() predicate.Invitation {
	return predicate.Invitation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, InvitedTable, InvitedColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasInvitedWith applies the HasEdge predicate on the "invited" edge with a given conditions (other predicates).
func HasInvitedWith(preds ...predicate.Account) predicate.Invitation {
	return predicate.Invitation(func(s *sql.Selector) {
		step := newInvitedStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Invitation) predicate.Invitation {
	return predicate.Invitation(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Invitation) predicate.Invitation {
	return predicate.Invitation(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Invitation) predicate.Invitation {
	return predicate.Invitation(sql.NotPredicates(p))
}
