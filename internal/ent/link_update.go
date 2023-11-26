// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/cluster"
	"github.com/Southclaws/storyden/internal/ent/item"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// LinkUpdate is the builder for updating Link entities.
type LinkUpdate struct {
	config
	hooks     []Hook
	mutation  *LinkMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the LinkUpdate builder.
func (lu *LinkUpdate) Where(ps ...predicate.Link) *LinkUpdate {
	lu.mutation.Where(ps...)
	return lu
}

// SetDomain sets the "domain" field.
func (lu *LinkUpdate) SetDomain(s string) *LinkUpdate {
	lu.mutation.SetDomain(s)
	return lu
}

// SetTitle sets the "title" field.
func (lu *LinkUpdate) SetTitle(s string) *LinkUpdate {
	lu.mutation.SetTitle(s)
	return lu
}

// SetDescription sets the "description" field.
func (lu *LinkUpdate) SetDescription(s string) *LinkUpdate {
	lu.mutation.SetDescription(s)
	return lu
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (lu *LinkUpdate) AddPostIDs(ids ...xid.ID) *LinkUpdate {
	lu.mutation.AddPostIDs(ids...)
	return lu
}

// AddPosts adds the "posts" edges to the Post entity.
func (lu *LinkUpdate) AddPosts(p ...*Post) *LinkUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.AddPostIDs(ids...)
}

// AddClusterIDs adds the "clusters" edge to the Cluster entity by IDs.
func (lu *LinkUpdate) AddClusterIDs(ids ...xid.ID) *LinkUpdate {
	lu.mutation.AddClusterIDs(ids...)
	return lu
}

// AddClusters adds the "clusters" edges to the Cluster entity.
func (lu *LinkUpdate) AddClusters(c ...*Cluster) *LinkUpdate {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return lu.AddClusterIDs(ids...)
}

// AddItemIDs adds the "items" edge to the Item entity by IDs.
func (lu *LinkUpdate) AddItemIDs(ids ...xid.ID) *LinkUpdate {
	lu.mutation.AddItemIDs(ids...)
	return lu
}

// AddItems adds the "items" edges to the Item entity.
func (lu *LinkUpdate) AddItems(i ...*Item) *LinkUpdate {
	ids := make([]xid.ID, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return lu.AddItemIDs(ids...)
}

// AddAssetIDs adds the "assets" edge to the Asset entity by IDs.
func (lu *LinkUpdate) AddAssetIDs(ids ...string) *LinkUpdate {
	lu.mutation.AddAssetIDs(ids...)
	return lu
}

// AddAssets adds the "assets" edges to the Asset entity.
func (lu *LinkUpdate) AddAssets(a ...*Asset) *LinkUpdate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return lu.AddAssetIDs(ids...)
}

// Mutation returns the LinkMutation object of the builder.
func (lu *LinkUpdate) Mutation() *LinkMutation {
	return lu.mutation
}

// ClearPosts clears all "posts" edges to the Post entity.
func (lu *LinkUpdate) ClearPosts() *LinkUpdate {
	lu.mutation.ClearPosts()
	return lu
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (lu *LinkUpdate) RemovePostIDs(ids ...xid.ID) *LinkUpdate {
	lu.mutation.RemovePostIDs(ids...)
	return lu
}

// RemovePosts removes "posts" edges to Post entities.
func (lu *LinkUpdate) RemovePosts(p ...*Post) *LinkUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.RemovePostIDs(ids...)
}

// ClearClusters clears all "clusters" edges to the Cluster entity.
func (lu *LinkUpdate) ClearClusters() *LinkUpdate {
	lu.mutation.ClearClusters()
	return lu
}

// RemoveClusterIDs removes the "clusters" edge to Cluster entities by IDs.
func (lu *LinkUpdate) RemoveClusterIDs(ids ...xid.ID) *LinkUpdate {
	lu.mutation.RemoveClusterIDs(ids...)
	return lu
}

// RemoveClusters removes "clusters" edges to Cluster entities.
func (lu *LinkUpdate) RemoveClusters(c ...*Cluster) *LinkUpdate {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return lu.RemoveClusterIDs(ids...)
}

// ClearItems clears all "items" edges to the Item entity.
func (lu *LinkUpdate) ClearItems() *LinkUpdate {
	lu.mutation.ClearItems()
	return lu
}

// RemoveItemIDs removes the "items" edge to Item entities by IDs.
func (lu *LinkUpdate) RemoveItemIDs(ids ...xid.ID) *LinkUpdate {
	lu.mutation.RemoveItemIDs(ids...)
	return lu
}

// RemoveItems removes "items" edges to Item entities.
func (lu *LinkUpdate) RemoveItems(i ...*Item) *LinkUpdate {
	ids := make([]xid.ID, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return lu.RemoveItemIDs(ids...)
}

// ClearAssets clears all "assets" edges to the Asset entity.
func (lu *LinkUpdate) ClearAssets() *LinkUpdate {
	lu.mutation.ClearAssets()
	return lu
}

// RemoveAssetIDs removes the "assets" edge to Asset entities by IDs.
func (lu *LinkUpdate) RemoveAssetIDs(ids ...string) *LinkUpdate {
	lu.mutation.RemoveAssetIDs(ids...)
	return lu
}

// RemoveAssets removes "assets" edges to Asset entities.
func (lu *LinkUpdate) RemoveAssets(a ...*Asset) *LinkUpdate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return lu.RemoveAssetIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lu *LinkUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, lu.sqlSave, lu.mutation, lu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LinkUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LinkUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LinkUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (lu *LinkUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *LinkUpdate {
	lu.modifiers = append(lu.modifiers, modifiers...)
	return lu
}

func (lu *LinkUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(link.Table, link.Columns, sqlgraph.NewFieldSpec(link.FieldID, field.TypeString))
	if ps := lu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lu.mutation.Domain(); ok {
		_spec.SetField(link.FieldDomain, field.TypeString, value)
	}
	if value, ok := lu.mutation.Title(); ok {
		_spec.SetField(link.FieldTitle, field.TypeString, value)
	}
	if value, ok := lu.mutation.Description(); ok {
		_spec.SetField(link.FieldDescription, field.TypeString, value)
	}
	if lu.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.PostsTable,
			Columns: link.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.RemovedPostsIDs(); len(nodes) > 0 && !lu.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.PostsTable,
			Columns: link.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.PostsTable,
			Columns: link.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lu.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ClustersTable,
			Columns: link.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.RemovedClustersIDs(); len(nodes) > 0 && !lu.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ClustersTable,
			Columns: link.ClustersPrimaryKey,
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
	if nodes := lu.mutation.ClustersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ClustersTable,
			Columns: link.ClustersPrimaryKey,
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
	if lu.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ItemsTable,
			Columns: link.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.RemovedItemsIDs(); len(nodes) > 0 && !lu.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ItemsTable,
			Columns: link.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.ItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ItemsTable,
			Columns: link.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lu.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.AssetsTable,
			Columns: link.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.RemovedAssetsIDs(); len(nodes) > 0 && !lu.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.AssetsTable,
			Columns: link.AssetsPrimaryKey,
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
	if nodes := lu.mutation.AssetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.AssetsTable,
			Columns: link.AssetsPrimaryKey,
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
	_spec.AddModifiers(lu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, lu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{link.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	lu.mutation.done = true
	return n, nil
}

// LinkUpdateOne is the builder for updating a single Link entity.
type LinkUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *LinkMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetDomain sets the "domain" field.
func (luo *LinkUpdateOne) SetDomain(s string) *LinkUpdateOne {
	luo.mutation.SetDomain(s)
	return luo
}

// SetTitle sets the "title" field.
func (luo *LinkUpdateOne) SetTitle(s string) *LinkUpdateOne {
	luo.mutation.SetTitle(s)
	return luo
}

// SetDescription sets the "description" field.
func (luo *LinkUpdateOne) SetDescription(s string) *LinkUpdateOne {
	luo.mutation.SetDescription(s)
	return luo
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (luo *LinkUpdateOne) AddPostIDs(ids ...xid.ID) *LinkUpdateOne {
	luo.mutation.AddPostIDs(ids...)
	return luo
}

// AddPosts adds the "posts" edges to the Post entity.
func (luo *LinkUpdateOne) AddPosts(p ...*Post) *LinkUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.AddPostIDs(ids...)
}

// AddClusterIDs adds the "clusters" edge to the Cluster entity by IDs.
func (luo *LinkUpdateOne) AddClusterIDs(ids ...xid.ID) *LinkUpdateOne {
	luo.mutation.AddClusterIDs(ids...)
	return luo
}

// AddClusters adds the "clusters" edges to the Cluster entity.
func (luo *LinkUpdateOne) AddClusters(c ...*Cluster) *LinkUpdateOne {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return luo.AddClusterIDs(ids...)
}

// AddItemIDs adds the "items" edge to the Item entity by IDs.
func (luo *LinkUpdateOne) AddItemIDs(ids ...xid.ID) *LinkUpdateOne {
	luo.mutation.AddItemIDs(ids...)
	return luo
}

// AddItems adds the "items" edges to the Item entity.
func (luo *LinkUpdateOne) AddItems(i ...*Item) *LinkUpdateOne {
	ids := make([]xid.ID, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return luo.AddItemIDs(ids...)
}

// AddAssetIDs adds the "assets" edge to the Asset entity by IDs.
func (luo *LinkUpdateOne) AddAssetIDs(ids ...string) *LinkUpdateOne {
	luo.mutation.AddAssetIDs(ids...)
	return luo
}

// AddAssets adds the "assets" edges to the Asset entity.
func (luo *LinkUpdateOne) AddAssets(a ...*Asset) *LinkUpdateOne {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return luo.AddAssetIDs(ids...)
}

// Mutation returns the LinkMutation object of the builder.
func (luo *LinkUpdateOne) Mutation() *LinkMutation {
	return luo.mutation
}

// ClearPosts clears all "posts" edges to the Post entity.
func (luo *LinkUpdateOne) ClearPosts() *LinkUpdateOne {
	luo.mutation.ClearPosts()
	return luo
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (luo *LinkUpdateOne) RemovePostIDs(ids ...xid.ID) *LinkUpdateOne {
	luo.mutation.RemovePostIDs(ids...)
	return luo
}

// RemovePosts removes "posts" edges to Post entities.
func (luo *LinkUpdateOne) RemovePosts(p ...*Post) *LinkUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.RemovePostIDs(ids...)
}

// ClearClusters clears all "clusters" edges to the Cluster entity.
func (luo *LinkUpdateOne) ClearClusters() *LinkUpdateOne {
	luo.mutation.ClearClusters()
	return luo
}

// RemoveClusterIDs removes the "clusters" edge to Cluster entities by IDs.
func (luo *LinkUpdateOne) RemoveClusterIDs(ids ...xid.ID) *LinkUpdateOne {
	luo.mutation.RemoveClusterIDs(ids...)
	return luo
}

// RemoveClusters removes "clusters" edges to Cluster entities.
func (luo *LinkUpdateOne) RemoveClusters(c ...*Cluster) *LinkUpdateOne {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return luo.RemoveClusterIDs(ids...)
}

// ClearItems clears all "items" edges to the Item entity.
func (luo *LinkUpdateOne) ClearItems() *LinkUpdateOne {
	luo.mutation.ClearItems()
	return luo
}

// RemoveItemIDs removes the "items" edge to Item entities by IDs.
func (luo *LinkUpdateOne) RemoveItemIDs(ids ...xid.ID) *LinkUpdateOne {
	luo.mutation.RemoveItemIDs(ids...)
	return luo
}

// RemoveItems removes "items" edges to Item entities.
func (luo *LinkUpdateOne) RemoveItems(i ...*Item) *LinkUpdateOne {
	ids := make([]xid.ID, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return luo.RemoveItemIDs(ids...)
}

// ClearAssets clears all "assets" edges to the Asset entity.
func (luo *LinkUpdateOne) ClearAssets() *LinkUpdateOne {
	luo.mutation.ClearAssets()
	return luo
}

// RemoveAssetIDs removes the "assets" edge to Asset entities by IDs.
func (luo *LinkUpdateOne) RemoveAssetIDs(ids ...string) *LinkUpdateOne {
	luo.mutation.RemoveAssetIDs(ids...)
	return luo
}

// RemoveAssets removes "assets" edges to Asset entities.
func (luo *LinkUpdateOne) RemoveAssets(a ...*Asset) *LinkUpdateOne {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return luo.RemoveAssetIDs(ids...)
}

// Where appends a list predicates to the LinkUpdate builder.
func (luo *LinkUpdateOne) Where(ps ...predicate.Link) *LinkUpdateOne {
	luo.mutation.Where(ps...)
	return luo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (luo *LinkUpdateOne) Select(field string, fields ...string) *LinkUpdateOne {
	luo.fields = append([]string{field}, fields...)
	return luo
}

// Save executes the query and returns the updated Link entity.
func (luo *LinkUpdateOne) Save(ctx context.Context) (*Link, error) {
	return withHooks(ctx, luo.sqlSave, luo.mutation, luo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LinkUpdateOne) SaveX(ctx context.Context) *Link {
	node, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (luo *LinkUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LinkUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (luo *LinkUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *LinkUpdateOne {
	luo.modifiers = append(luo.modifiers, modifiers...)
	return luo
}

func (luo *LinkUpdateOne) sqlSave(ctx context.Context) (_node *Link, err error) {
	_spec := sqlgraph.NewUpdateSpec(link.Table, link.Columns, sqlgraph.NewFieldSpec(link.FieldID, field.TypeString))
	id, ok := luo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Link.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := luo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, link.FieldID)
		for _, f := range fields {
			if !link.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != link.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := luo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := luo.mutation.Domain(); ok {
		_spec.SetField(link.FieldDomain, field.TypeString, value)
	}
	if value, ok := luo.mutation.Title(); ok {
		_spec.SetField(link.FieldTitle, field.TypeString, value)
	}
	if value, ok := luo.mutation.Description(); ok {
		_spec.SetField(link.FieldDescription, field.TypeString, value)
	}
	if luo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.PostsTable,
			Columns: link.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.RemovedPostsIDs(); len(nodes) > 0 && !luo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.PostsTable,
			Columns: link.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.PostsTable,
			Columns: link.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if luo.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ClustersTable,
			Columns: link.ClustersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.RemovedClustersIDs(); len(nodes) > 0 && !luo.mutation.ClustersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ClustersTable,
			Columns: link.ClustersPrimaryKey,
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
	if nodes := luo.mutation.ClustersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ClustersTable,
			Columns: link.ClustersPrimaryKey,
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
	if luo.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ItemsTable,
			Columns: link.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.RemovedItemsIDs(); len(nodes) > 0 && !luo.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ItemsTable,
			Columns: link.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.ItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.ItemsTable,
			Columns: link.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if luo.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.AssetsTable,
			Columns: link.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.RemovedAssetsIDs(); len(nodes) > 0 && !luo.mutation.AssetsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.AssetsTable,
			Columns: link.AssetsPrimaryKey,
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
	if nodes := luo.mutation.AssetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   link.AssetsTable,
			Columns: link.AssetsPrimaryKey,
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
	_spec.AddModifiers(luo.modifiers...)
	_node = &Link{config: luo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, luo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{link.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	luo.mutation.done = true
	return _node, nil
}
