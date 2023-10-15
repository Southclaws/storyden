// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/rs/xid"
)

// Asset is the model entity for the Asset schema.
type Asset struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Mimetype holds the value of the "mimetype" field.
	Mimetype string `json:"mimetype,omitempty"`
	// Width holds the value of the "width" field.
	Width int `json:"width,omitempty"`
	// Height holds the value of the "height" field.
	Height int `json:"height,omitempty"`
	// AccountID holds the value of the "account_id" field.
	AccountID xid.ID `json:"account_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AssetQuery when eager-loading is set.
	Edges        AssetEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AssetEdges holds the relations/edges for other nodes in the graph.
type AssetEdges struct {
	// Posts holds the value of the posts edge.
	Posts []*Post `json:"posts,omitempty"`
	// Clusters holds the value of the clusters edge.
	Clusters []*Cluster `json:"clusters,omitempty"`
	// Items holds the value of the items edge.
	Items []*Item `json:"items,omitempty"`
	// Owner holds the value of the owner edge.
	Owner *Account `json:"owner,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// PostsOrErr returns the Posts value or an error if the edge
// was not loaded in eager-loading.
func (e AssetEdges) PostsOrErr() ([]*Post, error) {
	if e.loadedTypes[0] {
		return e.Posts, nil
	}
	return nil, &NotLoadedError{edge: "posts"}
}

// ClustersOrErr returns the Clusters value or an error if the edge
// was not loaded in eager-loading.
func (e AssetEdges) ClustersOrErr() ([]*Cluster, error) {
	if e.loadedTypes[1] {
		return e.Clusters, nil
	}
	return nil, &NotLoadedError{edge: "clusters"}
}

// ItemsOrErr returns the Items value or an error if the edge
// was not loaded in eager-loading.
func (e AssetEdges) ItemsOrErr() ([]*Item, error) {
	if e.loadedTypes[2] {
		return e.Items, nil
	}
	return nil, &NotLoadedError{edge: "items"}
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssetEdges) OwnerOrErr() (*Account, error) {
	if e.loadedTypes[3] {
		if e.Owner == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: account.Label}
		}
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Asset) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case asset.FieldWidth, asset.FieldHeight:
			values[i] = new(sql.NullInt64)
		case asset.FieldID, asset.FieldURL, asset.FieldMimetype:
			values[i] = new(sql.NullString)
		case asset.FieldCreatedAt, asset.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case asset.FieldAccountID:
			values[i] = new(xid.ID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Asset fields.
func (a *Asset) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case asset.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				a.ID = value.String
			}
		case asset.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				a.CreatedAt = value.Time
			}
		case asset.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				a.UpdatedAt = value.Time
			}
		case asset.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				a.URL = value.String
			}
		case asset.FieldMimetype:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mimetype", values[i])
			} else if value.Valid {
				a.Mimetype = value.String
			}
		case asset.FieldWidth:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field width", values[i])
			} else if value.Valid {
				a.Width = int(value.Int64)
			}
		case asset.FieldHeight:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field height", values[i])
			} else if value.Valid {
				a.Height = int(value.Int64)
			}
		case asset.FieldAccountID:
			if value, ok := values[i].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field account_id", values[i])
			} else if value != nil {
				a.AccountID = *value
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Asset.
// This includes values selected through modifiers, order, etc.
func (a *Asset) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// QueryPosts queries the "posts" edge of the Asset entity.
func (a *Asset) QueryPosts() *PostQuery {
	return NewAssetClient(a.config).QueryPosts(a)
}

// QueryClusters queries the "clusters" edge of the Asset entity.
func (a *Asset) QueryClusters() *ClusterQuery {
	return NewAssetClient(a.config).QueryClusters(a)
}

// QueryItems queries the "items" edge of the Asset entity.
func (a *Asset) QueryItems() *ItemQuery {
	return NewAssetClient(a.config).QueryItems(a)
}

// QueryOwner queries the "owner" edge of the Asset entity.
func (a *Asset) QueryOwner() *AccountQuery {
	return NewAssetClient(a.config).QueryOwner(a)
}

// Update returns a builder for updating this Asset.
// Note that you need to call Asset.Unwrap() before calling this method if this Asset
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Asset) Update() *AssetUpdateOne {
	return NewAssetClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Asset entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Asset) Unwrap() *Asset {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Asset is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Asset) String() string {
	var builder strings.Builder
	builder.WriteString("Asset(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("created_at=")
	builder.WriteString(a.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(a.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("url=")
	builder.WriteString(a.URL)
	builder.WriteString(", ")
	builder.WriteString("mimetype=")
	builder.WriteString(a.Mimetype)
	builder.WriteString(", ")
	builder.WriteString("width=")
	builder.WriteString(fmt.Sprintf("%v", a.Width))
	builder.WriteString(", ")
	builder.WriteString("height=")
	builder.WriteString(fmt.Sprintf("%v", a.Height))
	builder.WriteString(", ")
	builder.WriteString("account_id=")
	builder.WriteString(fmt.Sprintf("%v", a.AccountID))
	builder.WriteByte(')')
	return builder.String()
}

// Assets is a parsable slice of Asset.
type Assets []*Asset
