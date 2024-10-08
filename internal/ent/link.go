// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/rs/xid"
)

// Link is the model entity for the Link schema.
type Link struct {
	config `json:"-"`
	// ID of the ent.
	ID xid.ID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// Domain holds the value of the "domain" field.
	Domain string `json:"domain,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// PrimaryAssetID holds the value of the "primary_asset_id" field.
	PrimaryAssetID *xid.ID `json:"primary_asset_id,omitempty"`
	// FaviconAssetID holds the value of the "favicon_asset_id" field.
	FaviconAssetID *xid.ID `json:"favicon_asset_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LinkQuery when eager-loading is set.
	Edges        LinkEdges `json:"edges"`
	selectValues sql.SelectValues
}

// LinkEdges holds the relations/edges for other nodes in the graph.
type LinkEdges struct {
	// Link aggregation posts that have shared this link.
	Posts []*Post `json:"posts,omitempty"`
	// Posts that reference this link in their content.
	PostContentReferences []*Post `json:"post_content_references,omitempty"`
	// Nodes holds the value of the nodes edge.
	Nodes []*Node `json:"nodes,omitempty"`
	// NodeContentReferences holds the value of the node_content_references edge.
	NodeContentReferences []*Node `json:"node_content_references,omitempty"`
	// PrimaryImage holds the value of the primary_image edge.
	PrimaryImage *Asset `json:"primary_image,omitempty"`
	// FaviconImage holds the value of the favicon_image edge.
	FaviconImage *Asset `json:"favicon_image,omitempty"`
	// Assets holds the value of the assets edge.
	Assets []*Asset `json:"assets,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
}

// PostsOrErr returns the Posts value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) PostsOrErr() ([]*Post, error) {
	if e.loadedTypes[0] {
		return e.Posts, nil
	}
	return nil, &NotLoadedError{edge: "posts"}
}

// PostContentReferencesOrErr returns the PostContentReferences value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) PostContentReferencesOrErr() ([]*Post, error) {
	if e.loadedTypes[1] {
		return e.PostContentReferences, nil
	}
	return nil, &NotLoadedError{edge: "post_content_references"}
}

// NodesOrErr returns the Nodes value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) NodesOrErr() ([]*Node, error) {
	if e.loadedTypes[2] {
		return e.Nodes, nil
	}
	return nil, &NotLoadedError{edge: "nodes"}
}

// NodeContentReferencesOrErr returns the NodeContentReferences value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) NodeContentReferencesOrErr() ([]*Node, error) {
	if e.loadedTypes[3] {
		return e.NodeContentReferences, nil
	}
	return nil, &NotLoadedError{edge: "node_content_references"}
}

// PrimaryImageOrErr returns the PrimaryImage value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LinkEdges) PrimaryImageOrErr() (*Asset, error) {
	if e.PrimaryImage != nil {
		return e.PrimaryImage, nil
	} else if e.loadedTypes[4] {
		return nil, &NotFoundError{label: asset.Label}
	}
	return nil, &NotLoadedError{edge: "primary_image"}
}

// FaviconImageOrErr returns the FaviconImage value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LinkEdges) FaviconImageOrErr() (*Asset, error) {
	if e.FaviconImage != nil {
		return e.FaviconImage, nil
	} else if e.loadedTypes[5] {
		return nil, &NotFoundError{label: asset.Label}
	}
	return nil, &NotLoadedError{edge: "favicon_image"}
}

// AssetsOrErr returns the Assets value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) AssetsOrErr() ([]*Asset, error) {
	if e.loadedTypes[6] {
		return e.Assets, nil
	}
	return nil, &NotLoadedError{edge: "assets"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Link) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case link.FieldPrimaryAssetID, link.FieldFaviconAssetID:
			values[i] = &sql.NullScanner{S: new(xid.ID)}
		case link.FieldURL, link.FieldSlug, link.FieldDomain, link.FieldTitle, link.FieldDescription:
			values[i] = new(sql.NullString)
		case link.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case link.FieldID:
			values[i] = new(xid.ID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Link fields.
func (l *Link) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case link.FieldID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				l.ID = *value
			}
		case link.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				l.CreatedAt = value.Time
			}
		case link.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				l.URL = value.String
			}
		case link.FieldSlug:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field slug", values[i])
			} else if value.Valid {
				l.Slug = value.String
			}
		case link.FieldDomain:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field domain", values[i])
			} else if value.Valid {
				l.Domain = value.String
			}
		case link.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				l.Title = value.String
			}
		case link.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				l.Description = value.String
			}
		case link.FieldPrimaryAssetID:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field primary_asset_id", values[i])
			} else if value.Valid {
				l.PrimaryAssetID = new(xid.ID)
				*l.PrimaryAssetID = *value.S.(*xid.ID)
			}
		case link.FieldFaviconAssetID:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field favicon_asset_id", values[i])
			} else if value.Valid {
				l.FaviconAssetID = new(xid.ID)
				*l.FaviconAssetID = *value.S.(*xid.ID)
			}
		default:
			l.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Link.
// This includes values selected through modifiers, order, etc.
func (l *Link) Value(name string) (ent.Value, error) {
	return l.selectValues.Get(name)
}

// QueryPosts queries the "posts" edge of the Link entity.
func (l *Link) QueryPosts() *PostQuery {
	return NewLinkClient(l.config).QueryPosts(l)
}

// QueryPostContentReferences queries the "post_content_references" edge of the Link entity.
func (l *Link) QueryPostContentReferences() *PostQuery {
	return NewLinkClient(l.config).QueryPostContentReferences(l)
}

// QueryNodes queries the "nodes" edge of the Link entity.
func (l *Link) QueryNodes() *NodeQuery {
	return NewLinkClient(l.config).QueryNodes(l)
}

// QueryNodeContentReferences queries the "node_content_references" edge of the Link entity.
func (l *Link) QueryNodeContentReferences() *NodeQuery {
	return NewLinkClient(l.config).QueryNodeContentReferences(l)
}

// QueryPrimaryImage queries the "primary_image" edge of the Link entity.
func (l *Link) QueryPrimaryImage() *AssetQuery {
	return NewLinkClient(l.config).QueryPrimaryImage(l)
}

// QueryFaviconImage queries the "favicon_image" edge of the Link entity.
func (l *Link) QueryFaviconImage() *AssetQuery {
	return NewLinkClient(l.config).QueryFaviconImage(l)
}

// QueryAssets queries the "assets" edge of the Link entity.
func (l *Link) QueryAssets() *AssetQuery {
	return NewLinkClient(l.config).QueryAssets(l)
}

// Update returns a builder for updating this Link.
// Note that you need to call Link.Unwrap() before calling this method if this Link
// was returned from a transaction, and the transaction was committed or rolled back.
func (l *Link) Update() *LinkUpdateOne {
	return NewLinkClient(l.config).UpdateOne(l)
}

// Unwrap unwraps the Link entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (l *Link) Unwrap() *Link {
	_tx, ok := l.config.driver.(*txDriver)
	if !ok {
		panic("ent: Link is not a transactional entity")
	}
	l.config.driver = _tx.drv
	return l
}

// String implements the fmt.Stringer.
func (l *Link) String() string {
	var builder strings.Builder
	builder.WriteString("Link(")
	builder.WriteString(fmt.Sprintf("id=%v, ", l.ID))
	builder.WriteString("created_at=")
	builder.WriteString(l.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("url=")
	builder.WriteString(l.URL)
	builder.WriteString(", ")
	builder.WriteString("slug=")
	builder.WriteString(l.Slug)
	builder.WriteString(", ")
	builder.WriteString("domain=")
	builder.WriteString(l.Domain)
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(l.Title)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(l.Description)
	builder.WriteString(", ")
	if v := l.PrimaryAssetID; v != nil {
		builder.WriteString("primary_asset_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := l.FaviconAssetID; v != nil {
		builder.WriteString("favicon_asset_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Links is a parsable slice of Link.
type Links []*Link
