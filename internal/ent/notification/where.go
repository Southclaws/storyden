// Code generated by ent, DO NOT EDIT.

package notification

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// ID filters vertices based on their ID field.
func ID(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id xid.ID) predicate.Notification {
	return predicate.Notification(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldCreatedAt, v))
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldTitle, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldDescription, v))
}

// Link applies equality check predicate on the "link" field. It's identical to LinkEQ.
func Link(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldLink, v))
}

// Read applies equality check predicate on the "read" field. It's identical to ReadEQ.
func Read(v bool) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldRead, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Notification {
	return predicate.Notification(sql.FieldLTE(FieldCreatedAt, v))
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldTitle, v))
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.Notification {
	return predicate.Notification(sql.FieldNEQ(FieldTitle, v))
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.Notification {
	return predicate.Notification(sql.FieldIn(FieldTitle, vs...))
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.Notification {
	return predicate.Notification(sql.FieldNotIn(FieldTitle, vs...))
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.Notification {
	return predicate.Notification(sql.FieldGT(FieldTitle, v))
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.Notification {
	return predicate.Notification(sql.FieldGTE(FieldTitle, v))
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.Notification {
	return predicate.Notification(sql.FieldLT(FieldTitle, v))
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.Notification {
	return predicate.Notification(sql.FieldLTE(FieldTitle, v))
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.Notification {
	return predicate.Notification(sql.FieldContains(FieldTitle, v))
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.Notification {
	return predicate.Notification(sql.FieldHasPrefix(FieldTitle, v))
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.Notification {
	return predicate.Notification(sql.FieldHasSuffix(FieldTitle, v))
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEqualFold(FieldTitle, v))
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.Notification {
	return predicate.Notification(sql.FieldContainsFold(FieldTitle, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Notification {
	return predicate.Notification(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Notification {
	return predicate.Notification(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Notification {
	return predicate.Notification(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Notification {
	return predicate.Notification(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Notification {
	return predicate.Notification(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Notification {
	return predicate.Notification(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Notification {
	return predicate.Notification(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Notification {
	return predicate.Notification(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Notification {
	return predicate.Notification(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Notification {
	return predicate.Notification(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Notification {
	return predicate.Notification(sql.FieldContainsFold(FieldDescription, v))
}

// LinkEQ applies the EQ predicate on the "link" field.
func LinkEQ(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldLink, v))
}

// LinkNEQ applies the NEQ predicate on the "link" field.
func LinkNEQ(v string) predicate.Notification {
	return predicate.Notification(sql.FieldNEQ(FieldLink, v))
}

// LinkIn applies the In predicate on the "link" field.
func LinkIn(vs ...string) predicate.Notification {
	return predicate.Notification(sql.FieldIn(FieldLink, vs...))
}

// LinkNotIn applies the NotIn predicate on the "link" field.
func LinkNotIn(vs ...string) predicate.Notification {
	return predicate.Notification(sql.FieldNotIn(FieldLink, vs...))
}

// LinkGT applies the GT predicate on the "link" field.
func LinkGT(v string) predicate.Notification {
	return predicate.Notification(sql.FieldGT(FieldLink, v))
}

// LinkGTE applies the GTE predicate on the "link" field.
func LinkGTE(v string) predicate.Notification {
	return predicate.Notification(sql.FieldGTE(FieldLink, v))
}

// LinkLT applies the LT predicate on the "link" field.
func LinkLT(v string) predicate.Notification {
	return predicate.Notification(sql.FieldLT(FieldLink, v))
}

// LinkLTE applies the LTE predicate on the "link" field.
func LinkLTE(v string) predicate.Notification {
	return predicate.Notification(sql.FieldLTE(FieldLink, v))
}

// LinkContains applies the Contains predicate on the "link" field.
func LinkContains(v string) predicate.Notification {
	return predicate.Notification(sql.FieldContains(FieldLink, v))
}

// LinkHasPrefix applies the HasPrefix predicate on the "link" field.
func LinkHasPrefix(v string) predicate.Notification {
	return predicate.Notification(sql.FieldHasPrefix(FieldLink, v))
}

// LinkHasSuffix applies the HasSuffix predicate on the "link" field.
func LinkHasSuffix(v string) predicate.Notification {
	return predicate.Notification(sql.FieldHasSuffix(FieldLink, v))
}

// LinkEqualFold applies the EqualFold predicate on the "link" field.
func LinkEqualFold(v string) predicate.Notification {
	return predicate.Notification(sql.FieldEqualFold(FieldLink, v))
}

// LinkContainsFold applies the ContainsFold predicate on the "link" field.
func LinkContainsFold(v string) predicate.Notification {
	return predicate.Notification(sql.FieldContainsFold(FieldLink, v))
}

// ReadEQ applies the EQ predicate on the "read" field.
func ReadEQ(v bool) predicate.Notification {
	return predicate.Notification(sql.FieldEQ(FieldRead, v))
}

// ReadNEQ applies the NEQ predicate on the "read" field.
func ReadNEQ(v bool) predicate.Notification {
	return predicate.Notification(sql.FieldNEQ(FieldRead, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Notification) predicate.Notification {
	return predicate.Notification(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Notification) predicate.Notification {
	return predicate.Notification(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Notification) predicate.Notification {
	return predicate.Notification(sql.NotPredicates(p))
}
