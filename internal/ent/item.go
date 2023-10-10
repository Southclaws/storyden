// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/item"
	"github.com/rs/xid"
)

// Item is the model entity for the Item schema.
type Item struct {
	config `json:"-"`
	// ID of the ent.
	ID xid.ID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// ImageURL holds the value of the "image_url" field.
	ImageURL *string `json:"image_url,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// AccountID holds the value of the "account_id" field.
	AccountID xid.ID `json:"account_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ItemQuery when eager-loading is set.
	Edges ItemEdges `json:"edges"`
}

// ItemEdges holds the relations/edges for other nodes in the graph.
type ItemEdges struct {
	// Owner holds the value of the owner edge.
	Owner *Account `json:"owner,omitempty"`
	// Clusters holds the value of the clusters edge.
	Clusters []*Cluster `json:"clusters,omitempty"`
	// Assets holds the value of the assets edge.
	Assets []*Asset `json:"assets,omitempty"`
	// Tags holds the value of the tags edge.
	Tags []*Tag `json:"tags,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ItemEdges) OwnerOrErr() (*Account, error) {
	if e.loadedTypes[0] {
		if e.Owner == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: account.Label}
		}
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// ClustersOrErr returns the Clusters value or an error if the edge
// was not loaded in eager-loading.
func (e ItemEdges) ClustersOrErr() ([]*Cluster, error) {
	if e.loadedTypes[1] {
		return e.Clusters, nil
	}
	return nil, &NotLoadedError{edge: "clusters"}
}

// AssetsOrErr returns the Assets value or an error if the edge
// was not loaded in eager-loading.
func (e ItemEdges) AssetsOrErr() ([]*Asset, error) {
	if e.loadedTypes[2] {
		return e.Assets, nil
	}
	return nil, &NotLoadedError{edge: "assets"}
}

// TagsOrErr returns the Tags value or an error if the edge
// was not loaded in eager-loading.
func (e ItemEdges) TagsOrErr() ([]*Tag, error) {
	if e.loadedTypes[3] {
		return e.Tags, nil
	}
	return nil, &NotLoadedError{edge: "tags"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Item) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case item.FieldName, item.FieldSlug, item.FieldImageURL, item.FieldDescription:
			values[i] = new(sql.NullString)
		case item.FieldCreatedAt, item.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case item.FieldID, item.FieldAccountID:
			values[i] = new(xid.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Item", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Item fields.
func (i *Item) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case item.FieldID:
			if value, ok := values[j].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value != nil {
				i.ID = *value
			}
		case item.FieldCreatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[j])
			} else if value.Valid {
				i.CreatedAt = value.Time
			}
		case item.FieldUpdatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[j])
			} else if value.Valid {
				i.UpdatedAt = value.Time
			}
		case item.FieldName:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[j])
			} else if value.Valid {
				i.Name = value.String
			}
		case item.FieldSlug:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field slug", values[j])
			} else if value.Valid {
				i.Slug = value.String
			}
		case item.FieldImageURL:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field image_url", values[j])
			} else if value.Valid {
				i.ImageURL = new(string)
				*i.ImageURL = value.String
			}
		case item.FieldDescription:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[j])
			} else if value.Valid {
				i.Description = value.String
			}
		case item.FieldAccountID:
			if value, ok := values[j].(*xid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field account_id", values[j])
			} else if value != nil {
				i.AccountID = *value
			}
		}
	}
	return nil
}

// QueryOwner queries the "owner" edge of the Item entity.
func (i *Item) QueryOwner() *AccountQuery {
	return NewItemClient(i.config).QueryOwner(i)
}

// QueryClusters queries the "clusters" edge of the Item entity.
func (i *Item) QueryClusters() *ClusterQuery {
	return NewItemClient(i.config).QueryClusters(i)
}

// QueryAssets queries the "assets" edge of the Item entity.
func (i *Item) QueryAssets() *AssetQuery {
	return NewItemClient(i.config).QueryAssets(i)
}

// QueryTags queries the "tags" edge of the Item entity.
func (i *Item) QueryTags() *TagQuery {
	return NewItemClient(i.config).QueryTags(i)
}

// Update returns a builder for updating this Item.
// Note that you need to call Item.Unwrap() before calling this method if this Item
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Item) Update() *ItemUpdateOne {
	return NewItemClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Item entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Item) Unwrap() *Item {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Item is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Item) String() string {
	var builder strings.Builder
	builder.WriteString("Item(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("created_at=")
	builder.WriteString(i.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(i.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(i.Name)
	builder.WriteString(", ")
	builder.WriteString("slug=")
	builder.WriteString(i.Slug)
	builder.WriteString(", ")
	if v := i.ImageURL; v != nil {
		builder.WriteString("image_url=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(i.Description)
	builder.WriteString(", ")
	builder.WriteString("account_id=")
	builder.WriteString(fmt.Sprintf("%v", i.AccountID))
	builder.WriteByte(')')
	return builder.String()
}

// Items is a parsable slice of Item.
type Items []*Item
