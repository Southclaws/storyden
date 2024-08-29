// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AccountsColumns holds the columns for the "accounts" table.
	AccountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "handle", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString},
		{Name: "bio", Type: field.TypeString, Nullable: true},
		{Name: "admin", Type: field.TypeBool, Default: false},
		{Name: "links", Type: field.TypeJSON, Nullable: true},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
	}
	// AccountsTable holds the schema information for the "accounts" table.
	AccountsTable = &schema.Table{
		Name:       "accounts",
		Columns:    AccountsColumns,
		PrimaryKey: []*schema.Column{AccountsColumns[0]},
	}
	// AssetsColumns holds the columns for the "assets" table.
	AssetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "filename", Type: field.TypeString},
		{Name: "url", Type: field.TypeString},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
		{Name: "account_id", Type: field.TypeString, Size: 20},
	}
	// AssetsTable holds the schema information for the "assets" table.
	AssetsTable = &schema.Table{
		Name:       "assets",
		Columns:    AssetsColumns,
		PrimaryKey: []*schema.Column{AssetsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "assets_accounts_assets",
				Columns:    []*schema.Column{AssetsColumns[6]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "asset_filename",
				Unique:  false,
				Columns: []*schema.Column{AssetsColumns[3]},
			},
		},
	}
	// AuthenticationsColumns holds the columns for the "authentications" table.
	AuthenticationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "service", Type: field.TypeString},
		{Name: "identifier", Type: field.TypeString},
		{Name: "token", Type: field.TypeString},
		{Name: "name", Type: field.TypeString, Nullable: true},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
		{Name: "account_authentication", Type: field.TypeString, Size: 20},
	}
	// AuthenticationsTable holds the schema information for the "authentications" table.
	AuthenticationsTable = &schema.Table{
		Name:       "authentications",
		Columns:    AuthenticationsColumns,
		PrimaryKey: []*schema.Column{AuthenticationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "authentications_accounts_authentication",
				Columns:    []*schema.Column{AuthenticationsColumns[7]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "authentication_service_identifier",
				Unique:  true,
				Columns: []*schema.Column{AuthenticationsColumns[2], AuthenticationsColumns[3]},
			},
		},
	}
	// CategoriesColumns holds the columns for the "categories" table.
	CategoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "slug", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Default: "(No description)"},
		{Name: "colour", Type: field.TypeString, Default: "#8577ce"},
		{Name: "sort", Type: field.TypeInt, Default: -1},
		{Name: "admin", Type: field.TypeBool, Default: false},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
	}
	// CategoriesTable holds the schema information for the "categories" table.
	CategoriesTable = &schema.Table{
		Name:       "categories",
		Columns:    CategoriesColumns,
		PrimaryKey: []*schema.Column{CategoriesColumns[0]},
	}
	// CollectionsColumns holds the columns for the "collections" table.
	CollectionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "visibility", Type: field.TypeEnum, Enums: []string{"draft", "unlisted", "review", "published"}, Default: "draft"},
		{Name: "account_collections", Type: field.TypeString, Nullable: true, Size: 20},
	}
	// CollectionsTable holds the schema information for the "collections" table.
	CollectionsTable = &schema.Table{
		Name:       "collections",
		Columns:    CollectionsColumns,
		PrimaryKey: []*schema.Column{CollectionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "collections_accounts_collections",
				Columns:    []*schema.Column{CollectionsColumns[6]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// CollectionNodesColumns holds the columns for the "collection_nodes" table.
	CollectionNodesColumns = []*schema.Column{
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "membership_type", Type: field.TypeString, Default: "normal"},
		{Name: "collection_id", Type: field.TypeString, Size: 20},
		{Name: "node_id", Type: field.TypeString, Size: 20},
	}
	// CollectionNodesTable holds the schema information for the "collection_nodes" table.
	CollectionNodesTable = &schema.Table{
		Name:       "collection_nodes",
		Columns:    CollectionNodesColumns,
		PrimaryKey: []*schema.Column{CollectionNodesColumns[2], CollectionNodesColumns[3]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "collection_nodes_collections_collection",
				Columns:    []*schema.Column{CollectionNodesColumns[2]},
				RefColumns: []*schema.Column{CollectionsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "collection_nodes_nodes_node",
				Columns:    []*schema.Column{CollectionNodesColumns[3]},
				RefColumns: []*schema.Column{NodesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "unique_collection_node",
				Unique:  true,
				Columns: []*schema.Column{CollectionNodesColumns[2], CollectionNodesColumns[3]},
			},
		},
	}
	// CollectionPostsColumns holds the columns for the "collection_posts" table.
	CollectionPostsColumns = []*schema.Column{
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "membership_type", Type: field.TypeString, Default: "normal"},
		{Name: "collection_id", Type: field.TypeString, Size: 20},
		{Name: "post_id", Type: field.TypeString, Size: 20},
	}
	// CollectionPostsTable holds the schema information for the "collection_posts" table.
	CollectionPostsTable = &schema.Table{
		Name:       "collection_posts",
		Columns:    CollectionPostsColumns,
		PrimaryKey: []*schema.Column{CollectionPostsColumns[2], CollectionPostsColumns[3]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "collection_posts_collections_collection",
				Columns:    []*schema.Column{CollectionPostsColumns[2]},
				RefColumns: []*schema.Column{CollectionsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "collection_posts_posts_post",
				Columns:    []*schema.Column{CollectionPostsColumns[3]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "unique_collection_post",
				Unique:  true,
				Columns: []*schema.Column{CollectionPostsColumns[2], CollectionPostsColumns[3]},
			},
		},
	}
	// EmailsColumns holds the columns for the "emails" table.
	EmailsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "email_address", Type: field.TypeString, Unique: true, Size: 254},
		{Name: "verification_code", Type: field.TypeString, Size: 6},
		{Name: "verified", Type: field.TypeBool, Default: "false"},
		{Name: "account_id", Type: field.TypeString, Nullable: true, Size: 20},
		{Name: "authentication_record_id", Type: field.TypeString, Nullable: true, Size: 20},
	}
	// EmailsTable holds the schema information for the "emails" table.
	EmailsTable = &schema.Table{
		Name:       "emails",
		Columns:    EmailsColumns,
		PrimaryKey: []*schema.Column{EmailsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "emails_accounts_emails",
				Columns:    []*schema.Column{EmailsColumns[5]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "emails_authentications_email_address",
				Columns:    []*schema.Column{EmailsColumns[6]},
				RefColumns: []*schema.Column{AuthenticationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// LinksColumns holds the columns for the "links" table.
	LinksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "url", Type: field.TypeString, Unique: true},
		{Name: "slug", Type: field.TypeString, Unique: true},
		{Name: "domain", Type: field.TypeString},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "primary_asset_id", Type: field.TypeString, Nullable: true, Size: 20},
		{Name: "favicon_asset_id", Type: field.TypeString, Nullable: true, Size: 20},
	}
	// LinksTable holds the schema information for the "links" table.
	LinksTable = &schema.Table{
		Name:       "links",
		Columns:    LinksColumns,
		PrimaryKey: []*schema.Column{LinksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "links_assets_primary_image",
				Columns:    []*schema.Column{LinksColumns[7]},
				RefColumns: []*schema.Column{AssetsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "links_assets_favicon_image",
				Columns:    []*schema.Column{LinksColumns[8]},
				RefColumns: []*schema.Column{AssetsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// NodesColumns holds the columns for the "nodes" table.
	NodesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "name", Type: field.TypeString},
		{Name: "slug", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "content", Type: field.TypeString, Nullable: true},
		{Name: "visibility", Type: field.TypeEnum, Enums: []string{"draft", "unlisted", "review", "published"}, Default: "draft"},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
		{Name: "account_id", Type: field.TypeString, Size: 20},
		{Name: "link_id", Type: field.TypeString, Nullable: true, Size: 20},
		{Name: "parent_node_id", Type: field.TypeString, Nullable: true, Size: 20},
	}
	// NodesTable holds the schema information for the "nodes" table.
	NodesTable = &schema.Table{
		Name:       "nodes",
		Columns:    NodesColumns,
		PrimaryKey: []*schema.Column{NodesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "nodes_accounts_nodes",
				Columns:    []*schema.Column{NodesColumns[10]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "nodes_links_nodes",
				Columns:    []*schema.Column{NodesColumns[11]},
				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "nodes_nodes_nodes",
				Columns:    []*schema.Column{NodesColumns[12]},
				RefColumns: []*schema.Column{NodesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "node_slug",
				Unique:  false,
				Columns: []*schema.Column{NodesColumns[5]},
			},
		},
	}
	// NotificationsColumns holds the columns for the "notifications" table.
	NotificationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "link", Type: field.TypeString},
		{Name: "read", Type: field.TypeBool},
	}
	// NotificationsTable holds the schema information for the "notifications" table.
	NotificationsTable = &schema.Table{
		Name:       "notifications",
		Columns:    NotificationsColumns,
		PrimaryKey: []*schema.Column{NotificationsColumns[0]},
	}
	// PostsColumns holds the columns for the "posts" table.
	PostsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "first", Type: field.TypeBool},
		{Name: "title", Type: field.TypeString, Nullable: true},
		{Name: "slug", Type: field.TypeString, Nullable: true},
		{Name: "pinned", Type: field.TypeBool, Default: false},
		{Name: "body", Type: field.TypeString},
		{Name: "short", Type: field.TypeString},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
		{Name: "visibility", Type: field.TypeEnum, Enums: []string{"draft", "unlisted", "review", "published"}, Default: "draft"},
		{Name: "account_posts", Type: field.TypeString, Size: 20},
		{Name: "category_id", Type: field.TypeString, Nullable: true, Size: 20},
		{Name: "link_id", Type: field.TypeString, Nullable: true, Size: 20},
		{Name: "root_post_id", Type: field.TypeString, Nullable: true, Size: 20},
		{Name: "reply_to_post_id", Type: field.TypeString, Nullable: true, Size: 20},
	}
	// PostsTable holds the schema information for the "posts" table.
	PostsTable = &schema.Table{
		Name:       "posts",
		Columns:    PostsColumns,
		PrimaryKey: []*schema.Column{PostsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "posts_accounts_posts",
				Columns:    []*schema.Column{PostsColumns[12]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "posts_categories_posts",
				Columns:    []*schema.Column{PostsColumns[13]},
				RefColumns: []*schema.Column{CategoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "posts_links_posts",
				Columns:    []*schema.Column{PostsColumns[14]},
				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "posts_posts_posts",
				Columns:    []*schema.Column{PostsColumns[15]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "posts_posts_replies",
				Columns:    []*schema.Column{PostsColumns[16]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ReactsColumns holds the columns for the "reacts" table.
	ReactsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "emoji", Type: field.TypeString},
		{Name: "account_id", Type: field.TypeString, Size: 20},
		{Name: "post_id", Type: field.TypeString, Size: 20},
	}
	// ReactsTable holds the schema information for the "reacts" table.
	ReactsTable = &schema.Table{
		Name:       "reacts",
		Columns:    ReactsColumns,
		PrimaryKey: []*schema.Column{ReactsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "reacts_accounts_reacts",
				Columns:    []*schema.Column{ReactsColumns[3]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "reacts_posts_reacts",
				Columns:    []*schema.Column{ReactsColumns[4]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// RolesColumns holds the columns for the "roles" table.
	RolesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// RolesTable holds the schema information for the "roles" table.
	RolesTable = &schema.Table{
		Name:       "roles",
		Columns:    RolesColumns,
		PrimaryKey: []*schema.Column{RolesColumns[0]},
	}
	// SettingsColumns holds the columns for the "settings" table.
	SettingsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "value", Type: field.TypeString},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// SettingsTable holds the schema information for the "settings" table.
	SettingsTable = &schema.Table{
		Name:       "settings",
		Columns:    SettingsColumns,
		PrimaryKey: []*schema.Column{SettingsColumns[0]},
	}
	// TagsColumns holds the columns for the "tags" table.
	TagsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Size: 20},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// TagsTable holds the schema information for the "tags" table.
	TagsTable = &schema.Table{
		Name:       "tags",
		Columns:    TagsColumns,
		PrimaryKey: []*schema.Column{TagsColumns[0]},
	}
	// AccountTagsColumns holds the columns for the "account_tags" table.
	AccountTagsColumns = []*schema.Column{
		{Name: "account_id", Type: field.TypeString, Size: 20},
		{Name: "tag_id", Type: field.TypeString, Size: 20},
	}
	// AccountTagsTable holds the schema information for the "account_tags" table.
	AccountTagsTable = &schema.Table{
		Name:       "account_tags",
		Columns:    AccountTagsColumns,
		PrimaryKey: []*schema.Column{AccountTagsColumns[0], AccountTagsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "account_tags_account_id",
				Columns:    []*schema.Column{AccountTagsColumns[0]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "account_tags_tag_id",
				Columns:    []*schema.Column{AccountTagsColumns[1]},
				RefColumns: []*schema.Column{TagsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// LinkPostContentReferencesColumns holds the columns for the "link_post_content_references" table.
	LinkPostContentReferencesColumns = []*schema.Column{
		{Name: "link_id", Type: field.TypeString, Size: 20},
		{Name: "post_id", Type: field.TypeString, Size: 20},
	}
	// LinkPostContentReferencesTable holds the schema information for the "link_post_content_references" table.
	LinkPostContentReferencesTable = &schema.Table{
		Name:       "link_post_content_references",
		Columns:    LinkPostContentReferencesColumns,
		PrimaryKey: []*schema.Column{LinkPostContentReferencesColumns[0], LinkPostContentReferencesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "link_post_content_references_link_id",
				Columns:    []*schema.Column{LinkPostContentReferencesColumns[0]},
				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "link_post_content_references_post_id",
				Columns:    []*schema.Column{LinkPostContentReferencesColumns[1]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// LinkNodeContentReferencesColumns holds the columns for the "link_node_content_references" table.
	LinkNodeContentReferencesColumns = []*schema.Column{
		{Name: "link_id", Type: field.TypeString, Size: 20},
		{Name: "node_id", Type: field.TypeString, Size: 20},
	}
	// LinkNodeContentReferencesTable holds the schema information for the "link_node_content_references" table.
	LinkNodeContentReferencesTable = &schema.Table{
		Name:       "link_node_content_references",
		Columns:    LinkNodeContentReferencesColumns,
		PrimaryKey: []*schema.Column{LinkNodeContentReferencesColumns[0], LinkNodeContentReferencesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "link_node_content_references_link_id",
				Columns:    []*schema.Column{LinkNodeContentReferencesColumns[0]},
				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "link_node_content_references_node_id",
				Columns:    []*schema.Column{LinkNodeContentReferencesColumns[1]},
				RefColumns: []*schema.Column{NodesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// LinkAssetsColumns holds the columns for the "link_assets" table.
	LinkAssetsColumns = []*schema.Column{
		{Name: "link_id", Type: field.TypeString, Size: 20},
		{Name: "asset_id", Type: field.TypeString, Size: 20},
	}
	// LinkAssetsTable holds the schema information for the "link_assets" table.
	LinkAssetsTable = &schema.Table{
		Name:       "link_assets",
		Columns:    LinkAssetsColumns,
		PrimaryKey: []*schema.Column{LinkAssetsColumns[0], LinkAssetsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "link_assets_link_id",
				Columns:    []*schema.Column{LinkAssetsColumns[0]},
				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "link_assets_asset_id",
				Columns:    []*schema.Column{LinkAssetsColumns[1]},
				RefColumns: []*schema.Column{AssetsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// NodeAssetsColumns holds the columns for the "node_assets" table.
	NodeAssetsColumns = []*schema.Column{
		{Name: "node_id", Type: field.TypeString, Size: 20},
		{Name: "asset_id", Type: field.TypeString, Size: 20},
	}
	// NodeAssetsTable holds the schema information for the "node_assets" table.
	NodeAssetsTable = &schema.Table{
		Name:       "node_assets",
		Columns:    NodeAssetsColumns,
		PrimaryKey: []*schema.Column{NodeAssetsColumns[0], NodeAssetsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "node_assets_node_id",
				Columns:    []*schema.Column{NodeAssetsColumns[0]},
				RefColumns: []*schema.Column{NodesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "node_assets_asset_id",
				Columns:    []*schema.Column{NodeAssetsColumns[1]},
				RefColumns: []*schema.Column{AssetsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// PostAssetsColumns holds the columns for the "post_assets" table.
	PostAssetsColumns = []*schema.Column{
		{Name: "post_id", Type: field.TypeString, Size: 20},
		{Name: "asset_id", Type: field.TypeString, Size: 20},
	}
	// PostAssetsTable holds the schema information for the "post_assets" table.
	PostAssetsTable = &schema.Table{
		Name:       "post_assets",
		Columns:    PostAssetsColumns,
		PrimaryKey: []*schema.Column{PostAssetsColumns[0], PostAssetsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "post_assets_post_id",
				Columns:    []*schema.Column{PostAssetsColumns[0]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "post_assets_asset_id",
				Columns:    []*schema.Column{PostAssetsColumns[1]},
				RefColumns: []*schema.Column{AssetsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// RoleAccountsColumns holds the columns for the "role_accounts" table.
	RoleAccountsColumns = []*schema.Column{
		{Name: "role_id", Type: field.TypeString, Size: 20},
		{Name: "account_id", Type: field.TypeString, Size: 20},
	}
	// RoleAccountsTable holds the schema information for the "role_accounts" table.
	RoleAccountsTable = &schema.Table{
		Name:       "role_accounts",
		Columns:    RoleAccountsColumns,
		PrimaryKey: []*schema.Column{RoleAccountsColumns[0], RoleAccountsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "role_accounts_role_id",
				Columns:    []*schema.Column{RoleAccountsColumns[0]},
				RefColumns: []*schema.Column{RolesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "role_accounts_account_id",
				Columns:    []*schema.Column{RoleAccountsColumns[1]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TagPostsColumns holds the columns for the "tag_posts" table.
	TagPostsColumns = []*schema.Column{
		{Name: "tag_id", Type: field.TypeString, Size: 20},
		{Name: "post_id", Type: field.TypeString, Size: 20},
	}
	// TagPostsTable holds the schema information for the "tag_posts" table.
	TagPostsTable = &schema.Table{
		Name:       "tag_posts",
		Columns:    TagPostsColumns,
		PrimaryKey: []*schema.Column{TagPostsColumns[0], TagPostsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tag_posts_tag_id",
				Columns:    []*schema.Column{TagPostsColumns[0]},
				RefColumns: []*schema.Column{TagsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "tag_posts_post_id",
				Columns:    []*schema.Column{TagPostsColumns[1]},
				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TagNodesColumns holds the columns for the "tag_nodes" table.
	TagNodesColumns = []*schema.Column{
		{Name: "tag_id", Type: field.TypeString, Size: 20},
		{Name: "node_id", Type: field.TypeString, Size: 20},
	}
	// TagNodesTable holds the schema information for the "tag_nodes" table.
	TagNodesTable = &schema.Table{
		Name:       "tag_nodes",
		Columns:    TagNodesColumns,
		PrimaryKey: []*schema.Column{TagNodesColumns[0], TagNodesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tag_nodes_tag_id",
				Columns:    []*schema.Column{TagNodesColumns[0]},
				RefColumns: []*schema.Column{TagsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "tag_nodes_node_id",
				Columns:    []*schema.Column{TagNodesColumns[1]},
				RefColumns: []*schema.Column{NodesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AccountsTable,
		AssetsTable,
		AuthenticationsTable,
		CategoriesTable,
		CollectionsTable,
		CollectionNodesTable,
		CollectionPostsTable,
		EmailsTable,
		LinksTable,
		NodesTable,
		NotificationsTable,
		PostsTable,
		ReactsTable,
		RolesTable,
		SettingsTable,
		TagsTable,
		AccountTagsTable,
		LinkPostContentReferencesTable,
		LinkNodeContentReferencesTable,
		LinkAssetsTable,
		NodeAssetsTable,
		PostAssetsTable,
		RoleAccountsTable,
		TagPostsTable,
		TagNodesTable,
	}
)

func init() {
	AssetsTable.ForeignKeys[0].RefTable = AccountsTable
	AuthenticationsTable.ForeignKeys[0].RefTable = AccountsTable
	CollectionsTable.ForeignKeys[0].RefTable = AccountsTable
	CollectionNodesTable.ForeignKeys[0].RefTable = CollectionsTable
	CollectionNodesTable.ForeignKeys[1].RefTable = NodesTable
	CollectionPostsTable.ForeignKeys[0].RefTable = CollectionsTable
	CollectionPostsTable.ForeignKeys[1].RefTable = PostsTable
	EmailsTable.ForeignKeys[0].RefTable = AccountsTable
	EmailsTable.ForeignKeys[1].RefTable = AuthenticationsTable
	LinksTable.ForeignKeys[0].RefTable = AssetsTable
	LinksTable.ForeignKeys[1].RefTable = AssetsTable
	NodesTable.ForeignKeys[0].RefTable = AccountsTable
	NodesTable.ForeignKeys[1].RefTable = LinksTable
	NodesTable.ForeignKeys[2].RefTable = NodesTable
	PostsTable.ForeignKeys[0].RefTable = AccountsTable
	PostsTable.ForeignKeys[1].RefTable = CategoriesTable
	PostsTable.ForeignKeys[2].RefTable = LinksTable
	PostsTable.ForeignKeys[3].RefTable = PostsTable
	PostsTable.ForeignKeys[4].RefTable = PostsTable
	ReactsTable.ForeignKeys[0].RefTable = AccountsTable
	ReactsTable.ForeignKeys[1].RefTable = PostsTable
	AccountTagsTable.ForeignKeys[0].RefTable = AccountsTable
	AccountTagsTable.ForeignKeys[1].RefTable = TagsTable
	AccountTagsTable.Annotation = &entsql.Annotation{}
	LinkPostContentReferencesTable.ForeignKeys[0].RefTable = LinksTable
	LinkPostContentReferencesTable.ForeignKeys[1].RefTable = PostsTable
	LinkNodeContentReferencesTable.ForeignKeys[0].RefTable = LinksTable
	LinkNodeContentReferencesTable.ForeignKeys[1].RefTable = NodesTable
	LinkAssetsTable.ForeignKeys[0].RefTable = LinksTable
	LinkAssetsTable.ForeignKeys[1].RefTable = AssetsTable
	NodeAssetsTable.ForeignKeys[0].RefTable = NodesTable
	NodeAssetsTable.ForeignKeys[1].RefTable = AssetsTable
	PostAssetsTable.ForeignKeys[0].RefTable = PostsTable
	PostAssetsTable.ForeignKeys[1].RefTable = AssetsTable
	RoleAccountsTable.ForeignKeys[0].RefTable = RolesTable
	RoleAccountsTable.ForeignKeys[1].RefTable = AccountsTable
	TagPostsTable.ForeignKeys[0].RefTable = TagsTable
	TagPostsTable.ForeignKeys[1].RefTable = PostsTable
	TagNodesTable.ForeignKeys[0].RefTable = TagsTable
	TagNodesTable.ForeignKeys[1].RefTable = NodesTable
}
