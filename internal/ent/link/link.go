// Code generated by ent, DO NOT EDIT.

package link

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/rs/xid"
)

const (
	// Label holds the string label denoting the link type in the database.
	Label = "link"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldSlug holds the string denoting the slug field in the database.
	FieldSlug = "slug"
	// FieldDomain holds the string denoting the domain field in the database.
	FieldDomain = "domain"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldPrimaryAssetID holds the string denoting the primary_asset_id field in the database.
	FieldPrimaryAssetID = "primary_asset_id"
	// FieldFaviconAssetID holds the string denoting the favicon_asset_id field in the database.
	FieldFaviconAssetID = "favicon_asset_id"
	// EdgePosts holds the string denoting the posts edge name in mutations.
	EdgePosts = "posts"
	// EdgePostContentReferences holds the string denoting the post_content_references edge name in mutations.
	EdgePostContentReferences = "post_content_references"
	// EdgeNodes holds the string denoting the nodes edge name in mutations.
	EdgeNodes = "nodes"
	// EdgeNodeContentReferences holds the string denoting the node_content_references edge name in mutations.
	EdgeNodeContentReferences = "node_content_references"
	// EdgePrimaryImage holds the string denoting the primary_image edge name in mutations.
	EdgePrimaryImage = "primary_image"
	// EdgeFaviconImage holds the string denoting the favicon_image edge name in mutations.
	EdgeFaviconImage = "favicon_image"
	// EdgeAssets holds the string denoting the assets edge name in mutations.
	EdgeAssets = "assets"
	// Table holds the table name of the link in the database.
	Table = "links"
	// PostsTable is the table that holds the posts relation/edge.
	PostsTable = "posts"
	// PostsInverseTable is the table name for the Post entity.
	// It exists in this package in order to avoid circular dependency with the "post" package.
	PostsInverseTable = "posts"
	// PostsColumn is the table column denoting the posts relation/edge.
	PostsColumn = "link_id"
	// PostContentReferencesTable is the table that holds the post_content_references relation/edge. The primary key declared below.
	PostContentReferencesTable = "link_post_content_references"
	// PostContentReferencesInverseTable is the table name for the Post entity.
	// It exists in this package in order to avoid circular dependency with the "post" package.
	PostContentReferencesInverseTable = "posts"
	// NodesTable is the table that holds the nodes relation/edge.
	NodesTable = "nodes"
	// NodesInverseTable is the table name for the Node entity.
	// It exists in this package in order to avoid circular dependency with the "node" package.
	NodesInverseTable = "nodes"
	// NodesColumn is the table column denoting the nodes relation/edge.
	NodesColumn = "link_id"
	// NodeContentReferencesTable is the table that holds the node_content_references relation/edge. The primary key declared below.
	NodeContentReferencesTable = "link_node_content_references"
	// NodeContentReferencesInverseTable is the table name for the Node entity.
	// It exists in this package in order to avoid circular dependency with the "node" package.
	NodeContentReferencesInverseTable = "nodes"
	// PrimaryImageTable is the table that holds the primary_image relation/edge.
	PrimaryImageTable = "links"
	// PrimaryImageInverseTable is the table name for the Asset entity.
	// It exists in this package in order to avoid circular dependency with the "asset" package.
	PrimaryImageInverseTable = "assets"
	// PrimaryImageColumn is the table column denoting the primary_image relation/edge.
	PrimaryImageColumn = "primary_asset_id"
	// FaviconImageTable is the table that holds the favicon_image relation/edge.
	FaviconImageTable = "links"
	// FaviconImageInverseTable is the table name for the Asset entity.
	// It exists in this package in order to avoid circular dependency with the "asset" package.
	FaviconImageInverseTable = "assets"
	// FaviconImageColumn is the table column denoting the favicon_image relation/edge.
	FaviconImageColumn = "favicon_asset_id"
	// AssetsTable is the table that holds the assets relation/edge. The primary key declared below.
	AssetsTable = "link_assets"
	// AssetsInverseTable is the table name for the Asset entity.
	// It exists in this package in order to avoid circular dependency with the "asset" package.
	AssetsInverseTable = "assets"
)

// Columns holds all SQL columns for link fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldURL,
	FieldSlug,
	FieldDomain,
	FieldTitle,
	FieldDescription,
	FieldPrimaryAssetID,
	FieldFaviconAssetID,
}

var (
	// PostContentReferencesPrimaryKey and PostContentReferencesColumn2 are the table columns denoting the
	// primary key for the post_content_references relation (M2M).
	PostContentReferencesPrimaryKey = []string{"link_id", "post_id"}
	// NodeContentReferencesPrimaryKey and NodeContentReferencesColumn2 are the table columns denoting the
	// primary key for the node_content_references relation (M2M).
	NodeContentReferencesPrimaryKey = []string{"link_id", "node_id"}
	// AssetsPrimaryKey and AssetsColumn2 are the table columns denoting the
	// primary key for the assets relation (M2M).
	AssetsPrimaryKey = []string{"link_id", "asset_id"}
)

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

// OrderOption defines the ordering options for the Link queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByURL orders the results by the url field.
func ByURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURL, opts...).ToFunc()
}

// BySlug orders the results by the slug field.
func BySlug(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSlug, opts...).ToFunc()
}

// ByDomain orders the results by the domain field.
func ByDomain(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDomain, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByPrimaryAssetID orders the results by the primary_asset_id field.
func ByPrimaryAssetID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPrimaryAssetID, opts...).ToFunc()
}

// ByFaviconAssetID orders the results by the favicon_asset_id field.
func ByFaviconAssetID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFaviconAssetID, opts...).ToFunc()
}

// ByPostsCount orders the results by posts count.
func ByPostsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPostsStep(), opts...)
	}
}

// ByPosts orders the results by posts terms.
func ByPosts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPostsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByPostContentReferencesCount orders the results by post_content_references count.
func ByPostContentReferencesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPostContentReferencesStep(), opts...)
	}
}

// ByPostContentReferences orders the results by post_content_references terms.
func ByPostContentReferences(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPostContentReferencesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByNodesCount orders the results by nodes count.
func ByNodesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newNodesStep(), opts...)
	}
}

// ByNodes orders the results by nodes terms.
func ByNodes(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNodesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByNodeContentReferencesCount orders the results by node_content_references count.
func ByNodeContentReferencesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newNodeContentReferencesStep(), opts...)
	}
}

// ByNodeContentReferences orders the results by node_content_references terms.
func ByNodeContentReferences(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNodeContentReferencesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByPrimaryImageField orders the results by primary_image field.
func ByPrimaryImageField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPrimaryImageStep(), sql.OrderByField(field, opts...))
	}
}

// ByFaviconImageField orders the results by favicon_image field.
func ByFaviconImageField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFaviconImageStep(), sql.OrderByField(field, opts...))
	}
}

// ByAssetsCount orders the results by assets count.
func ByAssetsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAssetsStep(), opts...)
	}
}

// ByAssets orders the results by assets terms.
func ByAssets(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAssetsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newPostsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PostsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, PostsTable, PostsColumn),
	)
}
func newPostContentReferencesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PostContentReferencesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, PostContentReferencesTable, PostContentReferencesPrimaryKey...),
	)
}
func newNodesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NodesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, NodesTable, NodesColumn),
	)
}
func newNodeContentReferencesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NodeContentReferencesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, NodeContentReferencesTable, NodeContentReferencesPrimaryKey...),
	)
}
func newPrimaryImageStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PrimaryImageInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, PrimaryImageTable, PrimaryImageColumn),
	)
}
func newFaviconImageStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FaviconImageInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, FaviconImageTable, FaviconImageColumn),
	)
}
func newAssetsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AssetsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, AssetsTable, AssetsPrimaryKey...),
	)
}
