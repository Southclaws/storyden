// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/cluster"
	"github.com/Southclaws/storyden/internal/ent/item"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/rs/xid"
)

// ItemUpdate is the builder for updating Item entities.
type ItemUpdate struct {
	config
	hooks     []Hook
	mutation  *ItemMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ItemUpdate builder.
func (iu *ItemUpdate) Where(ps ...predicate.Item) *ItemUpdate {
	iu.mutation.Where(ps...)
	return iu
}

// SetUpdatedAt sets the "updated_at" field.
func (iu *ItemUpdate) SetUpdatedAt(t time.Time) *ItemUpdate {
	iu.mutation.SetUpdatedAt(t)
	return iu
}

// SetDeletedAt sets the "deleted_at" field.
func (iu *ItemUpdate) SetDeletedAt(t time.Time) *ItemUpdate {
	iu.mutation.SetDeletedAt(t)
	return iu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableDeletedAt(t *time.Time) *ItemUpdate {
	if t != nil {
		iu.SetDeletedAt(*t)
	}
	return iu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (iu *ItemUpdate) ClearDeletedAt() *ItemUpdate {
	iu.mutation.ClearDeletedAt()
	return iu
}

// SetName sets the "name" field.
func (iu *ItemUpdate) SetName(s string) *ItemUpdate {
	iu.mutation.SetName(s)
	return iu
}

// SetSlug sets the "slug" field.
func (iu *ItemUpdate) SetSlug(s string) *ItemUpdate {
	iu.mutation.SetSlug(s)
	return iu
}

// SetImageURL sets the "image_url" field.
func (iu *ItemUpdate) SetImageURL(s string) *ItemUpdate {
	iu.mutation.SetImageURL(s)
	return iu
}

// SetNillableImageURL sets the "image_url" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableImageURL(s *string) *ItemUpdate {
	if s != nil {
		iu.SetImageURL(*s)
	}
	return iu
}

// ClearImageURL clears the value of the "image_url" field.
func (iu *ItemUpdate) ClearImageURL() *ItemUpdate {
	iu.mutation.ClearImageURL()
	return iu
}

// SetURL sets the "url" field.
func (iu *ItemUpdate) SetURL(s string) *ItemUpdate {
	iu.mutation.SetURL(s)
	return iu
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableURL(s *string) *ItemUpdate {
	if s != nil {
		iu.SetURL(*s)
	}
	return iu
}

// ClearURL clears the value of the "url" field.
func (iu *ItemUpdate) ClearURL() *ItemUpdate {
	iu.mutation.ClearURL()
	return iu
}

// SetURLTitle sets the "url_title" field.
func (iu *ItemUpdate) SetURLTitle(s string) *ItemUpdate {
	iu.mutation.SetURLTitle(s)
	return iu
}

// SetNillableURLTitle sets the "url_title" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableURLTitle(s *string) *ItemUpdate {
	if s != nil {
		iu.SetURLTitle(*s)
	}
	return iu
}

// ClearURLTitle clears the value of the "url_title" field.
func (iu *ItemUpdate) ClearURLTitle() *ItemUpdate {
	iu.mutation.ClearURLTitle()
	return iu
}

// SetURLDescription sets the "url_description" field.
func (iu *ItemUpdate) SetURLDescription(s string) *ItemUpdate {
	iu.mutation.SetURLDescription(s)
	return iu
}

// SetNillableURLDescription sets the "url_description" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableURLDescription(s *string) *ItemUpdate {
	if s != nil {
		iu.SetURLDescription(*s)
	}
	return iu
}

// ClearURLDescription clears the value of the "url_description" field.
func (iu *ItemUpdate) ClearURLDescription() *ItemUpdate {
	iu.mutation.ClearURLDescription()
	return iu
}

// SetDescription sets the "description" field.
func (iu *ItemUpdate) SetDescription(s string) *ItemUpdate {
	iu.mutation.SetDescription(s)
	return iu
}

// SetContent sets the "content" field.
func (iu *ItemUpdate) SetContent(s string) *ItemUpdate {
	iu.mutation.SetContent(s)
	return iu
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableContent(s *string) *ItemUpdate {
	if s != nil {
		iu.SetContent(*s)
	}
	return iu
}

// ClearContent clears the value of the "content" field.
func (iu *ItemUpdate) ClearContent() *ItemUpdate {
	iu.mutation.ClearContent()
	return iu
}

// SetAccountID sets the "account_id" field.
func (iu *ItemUpdate) SetAccountID(x xid.ID) *ItemUpdate {
	iu.mutation.SetAccountID(x)
	return iu
}

// SetProperties sets the "properties" field.
func (iu *ItemUpdate) SetProperties(a any) *ItemUpdate {
	iu.mutation.SetProperties(a)
	return iu
}

// ClearProperties clears the value of the "properties" field.
func (iu *ItemUpdate) ClearProperties() *ItemUpdate {
	iu.mutation.ClearProperties()
	return iu
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (iu *ItemUpdate) SetOwnerID(id xid.ID) *ItemUpdate {
	iu.mutation.SetOwnerID(id)
	return iu
}

// SetOwner sets the "owner" edge to the Account entity.
func (iu *ItemUpdate) SetOwner(a *Account) *ItemUpdate {
	return iu.SetOwnerID(a.ID)
}

// AddClusterIDs adds the "clusters" edge to the Cluster entity by IDs.
func (iu *ItemUpdate) AddClusterIDs(ids ...xid.ID) *ItemUpdate {
	iu.mutation.AddClusterIDs(ids...)
	return iu
}

// AddClusters adds the "clusters" edges to the Cluster entity.
func (iu *ItemUpdate) AddClusters(c ...*Cluster) *ItemUpdate {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iu.AddClusterIDs(ids...)
}

// AddAssetIDs adds the "assets" edge to the Asset entity by IDs.
func (iu *ItemUpdate) AddAssetIDs(ids ...string) *ItemUpdate {
	iu.mutation.AddAssetIDs(ids...)
	return iu
}

// AddAssets adds the "assets" edges to the Asset entity.
func (iu *ItemUpdate) AddAssets(a ...*Asset) *ItemUpdate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return iu.AddAssetIDs(ids...)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (iu *ItemUpdate) AddTagIDs(ids ...xid.ID) *ItemUpdate {
	iu.mutation.AddTagIDs(ids...)
	return iu
}

// AddTags adds the "tags" edges to the Tag entity.
func (iu *ItemUpdate) AddTags(t ...*Tag) *ItemUpdate {
	ids := make([]xid.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return iu.AddTagIDs(ids...)
}

// Mutation returns the ItemMutation object of the builder.
func (iu *ItemUpdate) Mutation() *ItemMutation {
	return iu.mutation
}

// ClearOwner clears the "owner" edge to the Account entity.
func (iu *ItemUpdate) ClearOwner() *ItemUpdate {
	iu.mutation.ClearOwner()
	return iu
}

// ClearClusters clears all "clusters" edges to the Cluster entity.
func (iu *ItemUpdate) ClearClusters() *ItemUpdate {
	iu.mutation.ClearClusters()
	return iu
}

// RemoveClusterIDs removes the "clusters" edge to Cluster entities by IDs.
func (iu *ItemUpdate) RemoveClusterIDs(ids ...xid.ID) *ItemUpdate {
	iu.mutation.RemoveClusterIDs(ids...)
	return iu
}

// RemoveClusters removes "clusters" edges to Cluster entities.
func (iu *ItemUpdate) RemoveClusters(c ...*Cluster) *ItemUpdate {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iu.RemoveClusterIDs(ids...)
}

// ClearAssets clears all "assets" edges to the Asset entity.
func (iu *ItemUpdate) ClearAssets() *ItemUpdate {
	iu.mutation.ClearAssets()
	return iu
}

// RemoveAssetIDs removes the "assets" edge to Asset entities by IDs.
func (iu *ItemUpdate) RemoveAssetIDs(ids ...string) *ItemUpdate {
	iu.mutation.RemoveAssetIDs(ids...)
	return iu
}

// RemoveAssets removes "assets" edges to Asset entities.
func (iu *ItemUpdate) RemoveAssets(a ...*Asset) *ItemUpdate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return iu.RemoveAssetIDs(ids...)
}

// ClearTags clears all "tags" edges to the Tag entity.
func (iu *ItemUpdate) ClearTags() *ItemUpdate {
	iu.mutation.ClearTags()
	return iu
}

// RemoveTagIDs removes the "tags" edge to Tag entities by IDs.
func (iu *ItemUpdate) RemoveTagIDs(ids ...xid.ID) *ItemUpdate {
	iu.mutation.RemoveTagIDs(ids...)
	return iu
}

// RemoveTags removes "tags" edges to Tag entities.
func (iu *ItemUpdate) RemoveTags(t ...*Tag) *ItemUpdate {
	ids := make([]xid.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return iu.RemoveTagIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (iu *ItemUpdate) Save(ctx context.Context) (int, error) {
	iu.defaults()
	return withHooks(ctx, iu.sqlSave, iu.mutation, iu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (iu *ItemUpdate) SaveX(ctx context.Context) int {
	affected, err := iu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (iu *ItemUpdate) Exec(ctx context.Context) error {
	_, err := iu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iu *ItemUpdate) ExecX(ctx context.Context) {
	if err := iu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (iu *ItemUpdate) defaults() {
	if _, ok := iu.mutation.UpdatedAt(); !ok {
		v := item.UpdateDefaultUpdatedAt()
		iu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (iu *ItemUpdate) check() error {
	if _, ok := iu.mutation.OwnerID(); iu.mutation.OwnerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Item.owner"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (iu *ItemUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ItemUpdate {
	iu.modifiers = append(iu.modifiers, modifiers...)
	return iu
}

func (iu *ItemUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := iu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(item.Table, item.Columns, sqlgraph.NewFieldSpec(item.FieldID, field.TypeString))
	if ps := iu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iu.mutation.UpdatedAt(); ok {
		_spec.SetField(item.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := iu.mutation.DeletedAt(); ok {
		_spec.SetField(item.FieldDeletedAt, field.TypeTime, value)
	}
	if iu.mutation.DeletedAtCleared() {
		_spec.ClearField(item.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := iu.mutation.Name(); ok {
		_spec.SetField(item.FieldName, field.TypeString, value)
	}
	if value, ok := iu.mutation.Slug(); ok {
		_spec.SetField(item.FieldSlug, field.TypeString, value)
	}
	if value, ok := iu.mutation.ImageURL(); ok {
		_spec.SetField(item.FieldImageURL, field.TypeString, value)
	}
	if iu.mutation.ImageURLCleared() {
		_spec.ClearField(item.FieldImageURL, field.TypeString)
	}
	if value, ok := iu.mutation.URL(); ok {
		_spec.SetField(item.FieldURL, field.TypeString, value)
	}
	if iu.mutation.URLCleared() {
		_spec.ClearField(item.FieldURL, field.TypeString)
	}
	if value, ok := iu.mutation.URLTitle(); ok {
		_spec.SetField(item.FieldURLTitle, field.TypeString, value)
	}
	if iu.mutation.URLTitleCleared() {
		_spec.ClearField(item.FieldURLTitle, field.TypeString)
	}
	if value, ok := iu.mutation.URLDescription(); ok {
		_spec.SetField(item.FieldURLDescription, field.TypeString, value)
	}
	if iu.mutation.URLDescriptionCleared() {
		_spec.ClearField(item.FieldURLDescription, field.TypeString)
	}
	if value, ok := iu.mutation.Description(); ok {
		_spec.SetField(item.FieldDescription, field.TypeString, value)
	}
	if value, ok := iu.mutation.Content(); ok {
		_spec.SetField(item.FieldContent, field.TypeString, value)
	}
	if iu.mutation.ContentCleared() {
		_spec.ClearField(item.FieldContent, field.TypeString)
	}
	if value, ok := iu.mutation.Properties(); ok {
		_spec.SetField(item.FieldProperties, field.TypeJSON, value)
	}
	if iu.mutation.PropertiesCleared() {
		_spec.ClearField(item.FieldProperties, field.TypeJSON)
	}
	if iu.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   item.OwnerTable,
			Columns: []string{item.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   item.OwnerTable,
			Columns: []string{item.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iu.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.ClustersTable,
			Columns: item.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedClustersIDs(); len(nodes) > 0 && !iu.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.ClustersTable,
			Columns: item.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.ClustersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.ClustersTable,
			Columns: item.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iu.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.AssetsTable,
			Columns: item.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedAssetsIDs(); len(nodes) > 0 && !iu.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.AssetsTable,
			Columns: item.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.AssetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.AssetsTable,
			Columns: item.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iu.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.TagsTable,
			Columns: item.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedTagsIDs(); len(nodes) > 0 && !iu.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.TagsTable,
			Columns: item.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.TagsTable,
			Columns: item.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(iu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, iu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{item.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	iu.mutation.done = true
	return n, nil
}

// ItemUpdateOne is the builder for updating a single Item entity.
type ItemUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ItemMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (iuo *ItemUpdateOne) SetUpdatedAt(t time.Time) *ItemUpdateOne {
	iuo.mutation.SetUpdatedAt(t)
	return iuo
}

// SetDeletedAt sets the "deleted_at" field.
func (iuo *ItemUpdateOne) SetDeletedAt(t time.Time) *ItemUpdateOne {
	iuo.mutation.SetDeletedAt(t)
	return iuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableDeletedAt(t *time.Time) *ItemUpdateOne {
	if t != nil {
		iuo.SetDeletedAt(*t)
	}
	return iuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (iuo *ItemUpdateOne) ClearDeletedAt() *ItemUpdateOne {
	iuo.mutation.ClearDeletedAt()
	return iuo
}

// SetName sets the "name" field.
func (iuo *ItemUpdateOne) SetName(s string) *ItemUpdateOne {
	iuo.mutation.SetName(s)
	return iuo
}

// SetSlug sets the "slug" field.
func (iuo *ItemUpdateOne) SetSlug(s string) *ItemUpdateOne {
	iuo.mutation.SetSlug(s)
	return iuo
}

// SetImageURL sets the "image_url" field.
func (iuo *ItemUpdateOne) SetImageURL(s string) *ItemUpdateOne {
	iuo.mutation.SetImageURL(s)
	return iuo
}

// SetNillableImageURL sets the "image_url" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableImageURL(s *string) *ItemUpdateOne {
	if s != nil {
		iuo.SetImageURL(*s)
	}
	return iuo
}

// ClearImageURL clears the value of the "image_url" field.
func (iuo *ItemUpdateOne) ClearImageURL() *ItemUpdateOne {
	iuo.mutation.ClearImageURL()
	return iuo
}

// SetURL sets the "url" field.
func (iuo *ItemUpdateOne) SetURL(s string) *ItemUpdateOne {
	iuo.mutation.SetURL(s)
	return iuo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableURL(s *string) *ItemUpdateOne {
	if s != nil {
		iuo.SetURL(*s)
	}
	return iuo
}

// ClearURL clears the value of the "url" field.
func (iuo *ItemUpdateOne) ClearURL() *ItemUpdateOne {
	iuo.mutation.ClearURL()
	return iuo
}

// SetURLTitle sets the "url_title" field.
func (iuo *ItemUpdateOne) SetURLTitle(s string) *ItemUpdateOne {
	iuo.mutation.SetURLTitle(s)
	return iuo
}

// SetNillableURLTitle sets the "url_title" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableURLTitle(s *string) *ItemUpdateOne {
	if s != nil {
		iuo.SetURLTitle(*s)
	}
	return iuo
}

// ClearURLTitle clears the value of the "url_title" field.
func (iuo *ItemUpdateOne) ClearURLTitle() *ItemUpdateOne {
	iuo.mutation.ClearURLTitle()
	return iuo
}

// SetURLDescription sets the "url_description" field.
func (iuo *ItemUpdateOne) SetURLDescription(s string) *ItemUpdateOne {
	iuo.mutation.SetURLDescription(s)
	return iuo
}

// SetNillableURLDescription sets the "url_description" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableURLDescription(s *string) *ItemUpdateOne {
	if s != nil {
		iuo.SetURLDescription(*s)
	}
	return iuo
}

// ClearURLDescription clears the value of the "url_description" field.
func (iuo *ItemUpdateOne) ClearURLDescription() *ItemUpdateOne {
	iuo.mutation.ClearURLDescription()
	return iuo
}

// SetDescription sets the "description" field.
func (iuo *ItemUpdateOne) SetDescription(s string) *ItemUpdateOne {
	iuo.mutation.SetDescription(s)
	return iuo
}

// SetContent sets the "content" field.
func (iuo *ItemUpdateOne) SetContent(s string) *ItemUpdateOne {
	iuo.mutation.SetContent(s)
	return iuo
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableContent(s *string) *ItemUpdateOne {
	if s != nil {
		iuo.SetContent(*s)
	}
	return iuo
}

// ClearContent clears the value of the "content" field.
func (iuo *ItemUpdateOne) ClearContent() *ItemUpdateOne {
	iuo.mutation.ClearContent()
	return iuo
}

// SetAccountID sets the "account_id" field.
func (iuo *ItemUpdateOne) SetAccountID(x xid.ID) *ItemUpdateOne {
	iuo.mutation.SetAccountID(x)
	return iuo
}

// SetProperties sets the "properties" field.
func (iuo *ItemUpdateOne) SetProperties(a any) *ItemUpdateOne {
	iuo.mutation.SetProperties(a)
	return iuo
}

// ClearProperties clears the value of the "properties" field.
func (iuo *ItemUpdateOne) ClearProperties() *ItemUpdateOne {
	iuo.mutation.ClearProperties()
	return iuo
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (iuo *ItemUpdateOne) SetOwnerID(id xid.ID) *ItemUpdateOne {
	iuo.mutation.SetOwnerID(id)
	return iuo
}

// SetOwner sets the "owner" edge to the Account entity.
func (iuo *ItemUpdateOne) SetOwner(a *Account) *ItemUpdateOne {
	return iuo.SetOwnerID(a.ID)
}

// AddClusterIDs adds the "clusters" edge to the Cluster entity by IDs.
func (iuo *ItemUpdateOne) AddClusterIDs(ids ...xid.ID) *ItemUpdateOne {
	iuo.mutation.AddClusterIDs(ids...)
	return iuo
}

// AddClusters adds the "clusters" edges to the Cluster entity.
func (iuo *ItemUpdateOne) AddClusters(c ...*Cluster) *ItemUpdateOne {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iuo.AddClusterIDs(ids...)
}

// AddAssetIDs adds the "assets" edge to the Asset entity by IDs.
func (iuo *ItemUpdateOne) AddAssetIDs(ids ...string) *ItemUpdateOne {
	iuo.mutation.AddAssetIDs(ids...)
	return iuo
}

// AddAssets adds the "assets" edges to the Asset entity.
func (iuo *ItemUpdateOne) AddAssets(a ...*Asset) *ItemUpdateOne {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return iuo.AddAssetIDs(ids...)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (iuo *ItemUpdateOne) AddTagIDs(ids ...xid.ID) *ItemUpdateOne {
	iuo.mutation.AddTagIDs(ids...)
	return iuo
}

// AddTags adds the "tags" edges to the Tag entity.
func (iuo *ItemUpdateOne) AddTags(t ...*Tag) *ItemUpdateOne {
	ids := make([]xid.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return iuo.AddTagIDs(ids...)
}

// Mutation returns the ItemMutation object of the builder.
func (iuo *ItemUpdateOne) Mutation() *ItemMutation {
	return iuo.mutation
}

// ClearOwner clears the "owner" edge to the Account entity.
func (iuo *ItemUpdateOne) ClearOwner() *ItemUpdateOne {
	iuo.mutation.ClearOwner()
	return iuo
}

// ClearClusters clears all "clusters" edges to the Cluster entity.
func (iuo *ItemUpdateOne) ClearClusters() *ItemUpdateOne {
	iuo.mutation.ClearClusters()
	return iuo
}

// RemoveClusterIDs removes the "clusters" edge to Cluster entities by IDs.
func (iuo *ItemUpdateOne) RemoveClusterIDs(ids ...xid.ID) *ItemUpdateOne {
	iuo.mutation.RemoveClusterIDs(ids...)
	return iuo
}

// RemoveClusters removes "clusters" edges to Cluster entities.
func (iuo *ItemUpdateOne) RemoveClusters(c ...*Cluster) *ItemUpdateOne {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iuo.RemoveClusterIDs(ids...)
}

// ClearAssets clears all "assets" edges to the Asset entity.
func (iuo *ItemUpdateOne) ClearAssets() *ItemUpdateOne {
	iuo.mutation.ClearAssets()
	return iuo
}

// RemoveAssetIDs removes the "assets" edge to Asset entities by IDs.
func (iuo *ItemUpdateOne) RemoveAssetIDs(ids ...string) *ItemUpdateOne {
	iuo.mutation.RemoveAssetIDs(ids...)
	return iuo
}

// RemoveAssets removes "assets" edges to Asset entities.
func (iuo *ItemUpdateOne) RemoveAssets(a ...*Asset) *ItemUpdateOne {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return iuo.RemoveAssetIDs(ids...)
}

// ClearTags clears all "tags" edges to the Tag entity.
func (iuo *ItemUpdateOne) ClearTags() *ItemUpdateOne {
	iuo.mutation.ClearTags()
	return iuo
}

// RemoveTagIDs removes the "tags" edge to Tag entities by IDs.
func (iuo *ItemUpdateOne) RemoveTagIDs(ids ...xid.ID) *ItemUpdateOne {
	iuo.mutation.RemoveTagIDs(ids...)
	return iuo
}

// RemoveTags removes "tags" edges to Tag entities.
func (iuo *ItemUpdateOne) RemoveTags(t ...*Tag) *ItemUpdateOne {
	ids := make([]xid.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return iuo.RemoveTagIDs(ids...)
}

// Where appends a list predicates to the ItemUpdate builder.
func (iuo *ItemUpdateOne) Where(ps ...predicate.Item) *ItemUpdateOne {
	iuo.mutation.Where(ps...)
	return iuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (iuo *ItemUpdateOne) Select(field string, fields ...string) *ItemUpdateOne {
	iuo.fields = append([]string{field}, fields...)
	return iuo
}

// Save executes the query and returns the updated Item entity.
func (iuo *ItemUpdateOne) Save(ctx context.Context) (*Item, error) {
	iuo.defaults()
	return withHooks(ctx, iuo.sqlSave, iuo.mutation, iuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (iuo *ItemUpdateOne) SaveX(ctx context.Context) *Item {
	node, err := iuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (iuo *ItemUpdateOne) Exec(ctx context.Context) error {
	_, err := iuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iuo *ItemUpdateOne) ExecX(ctx context.Context) {
	if err := iuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (iuo *ItemUpdateOne) defaults() {
	if _, ok := iuo.mutation.UpdatedAt(); !ok {
		v := item.UpdateDefaultUpdatedAt()
		iuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (iuo *ItemUpdateOne) check() error {
	if _, ok := iuo.mutation.OwnerID(); iuo.mutation.OwnerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Item.owner"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (iuo *ItemUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ItemUpdateOne {
	iuo.modifiers = append(iuo.modifiers, modifiers...)
	return iuo
}

func (iuo *ItemUpdateOne) sqlSave(ctx context.Context) (_node *Item, err error) {
	if err := iuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(item.Table, item.Columns, sqlgraph.NewFieldSpec(item.FieldID, field.TypeString))
	id, ok := iuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Item.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := iuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, item.FieldID)
		for _, f := range fields {
			if !item.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != item.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := iuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iuo.mutation.UpdatedAt(); ok {
		_spec.SetField(item.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := iuo.mutation.DeletedAt(); ok {
		_spec.SetField(item.FieldDeletedAt, field.TypeTime, value)
	}
	if iuo.mutation.DeletedAtCleared() {
		_spec.ClearField(item.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := iuo.mutation.Name(); ok {
		_spec.SetField(item.FieldName, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Slug(); ok {
		_spec.SetField(item.FieldSlug, field.TypeString, value)
	}
	if value, ok := iuo.mutation.ImageURL(); ok {
		_spec.SetField(item.FieldImageURL, field.TypeString, value)
	}
	if iuo.mutation.ImageURLCleared() {
		_spec.ClearField(item.FieldImageURL, field.TypeString)
	}
	if value, ok := iuo.mutation.URL(); ok {
		_spec.SetField(item.FieldURL, field.TypeString, value)
	}
	if iuo.mutation.URLCleared() {
		_spec.ClearField(item.FieldURL, field.TypeString)
	}
	if value, ok := iuo.mutation.URLTitle(); ok {
		_spec.SetField(item.FieldURLTitle, field.TypeString, value)
	}
	if iuo.mutation.URLTitleCleared() {
		_spec.ClearField(item.FieldURLTitle, field.TypeString)
	}
	if value, ok := iuo.mutation.URLDescription(); ok {
		_spec.SetField(item.FieldURLDescription, field.TypeString, value)
	}
	if iuo.mutation.URLDescriptionCleared() {
		_spec.ClearField(item.FieldURLDescription, field.TypeString)
	}
	if value, ok := iuo.mutation.Description(); ok {
		_spec.SetField(item.FieldDescription, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Content(); ok {
		_spec.SetField(item.FieldContent, field.TypeString, value)
	}
	if iuo.mutation.ContentCleared() {
		_spec.ClearField(item.FieldContent, field.TypeString)
	}
	if value, ok := iuo.mutation.Properties(); ok {
		_spec.SetField(item.FieldProperties, field.TypeJSON, value)
	}
	if iuo.mutation.PropertiesCleared() {
		_spec.ClearField(item.FieldProperties, field.TypeJSON)
	}
	if iuo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   item.OwnerTable,
			Columns: []string{item.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   item.OwnerTable,
			Columns: []string{item.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iuo.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.ClustersTable,
			Columns: item.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedClustersIDs(); len(nodes) > 0 && !iuo.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.ClustersTable,
			Columns: item.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.ClustersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.ClustersTable,
			Columns: item.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iuo.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.AssetsTable,
			Columns: item.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedAssetsIDs(); len(nodes) > 0 && !iuo.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.AssetsTable,
			Columns: item.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.AssetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.AssetsTable,
			Columns: item.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iuo.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.TagsTable,
			Columns: item.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedTagsIDs(); len(nodes) > 0 && !iuo.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.TagsTable,
			Columns: item.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   item.TagsTable,
			Columns: item.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(iuo.modifiers...)
	_node = &Item{config: iuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, iuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{item.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	iuo.mutation.done = true
	return _node, nil
}
