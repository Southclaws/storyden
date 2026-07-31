package cligen

import (
	"time"
)

type Asset struct {
	Filename string  `json:"filename"`
	Height   float64 `json:"height"`
	ID       string  `json:"id"`
	MimeType string  `json:"mime_type"`
	// The API path of the asset (GET /assets/...).
	Path  string  `json:"path"`
	Width float64 `json:"width"`
}

// The result of one item in a batch operation (node delete/move/visibility).
type BatchResult struct {
	Error *string `json:"error,omitempty"`
	ID    string  `json:"id"`
	// Per-command context; e.g. the node's name, or "name → visibility" for a visibility change.
	Name *string `json:"name,omitempty"`
	Ok   bool    `json:"ok"`
}

type CategoryReference struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// A discriminated search result item; `ref`'s exact shape depends on `kind`.
type DatagraphItem struct {
	Kind string `json:"kind"`
	// The referenced resource; shape depends on kind (Thread for thread/post, Node for node, ProfileReference-like for profile, and so on).
	Ref map[string]interface{} `json:"ref"`
}

type InfoMotd struct {
	Content *string `json:"content,omitempty"`
}

type Info struct {
	AccentColour       *string  `json:"accent_colour,omitempty"`
	APIAddress         string   `json:"api_address"`
	AuthenticationMode *string  `json:"authentication_mode,omitempty"`
	Capabilities       []string `json:"capabilities,omitempty"`
	// Rich HTML overview content for the instance.
	Content          *string   `json:"content,omitempty"`
	Description      *string   `json:"description,omitempty"`
	Motd             *InfoMotd `json:"motd,omitempty"`
	RegistrationMode *string   `json:"registration_mode,omitempty"`
	Title            string    `json:"title"`
	WebAddress       string    `json:"web_address"`
}

type InstanceInfo struct {
	BaseURL  string  `json:"base_url"`
	Context  *string `json:"context,omitempty"`
	Endpoint string  `json:"endpoint"`
	Info     Info    `json:"info"`
}

type LinkReference struct {
	Description  *string `json:"description,omitempty"`
	Domain       string  `json:"domain"`
	FaviconImage *Asset  `json:"favicon_image,omitempty"`
	PrimaryImage *Asset  `json:"primary_image,omitempty"`
	Slug         string  `json:"slug"`
	Title        *string `json:"title,omitempty"`
	URL          string  `json:"url"`
}

// Arbitrary metadata for a resource.
type Metadata = interface{}

// A minimal reference to an account.
type ProfileReference struct {
	Handle string   `json:"handle"`
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles,omitempty"`
}

type TagReference struct {
	Colour *string `json:"colour,omitempty"`
	Name   string  `json:"name"`
}

type Visibility string

const (
	VisibilityDraft     Visibility = "draft"
	VisibilityReview    Visibility = "review"
	VisibilityPublished Visibility = "published"
	VisibilityUnlisted  Visibility = "unlisted"
)

var VisibilityValues = []Visibility{
	VisibilityDraft,
	VisibilityReview,
	VisibilityPublished,
	VisibilityUnlisted,
}

// A node without its children (as returned by search results).
type Node struct {
	Assets []Asset `json:"assets"`
	// Rich HTML content.
	Content          *string          `json:"content,omitempty"`
	CreatedAt        *time.Time       `json:"createdAt,omitempty"`
	CurrentVersionID *string          `json:"current_version_id,omitempty"`
	DeletedAt        *time.Time       `json:"deletedAt,omitempty"`
	Description      string           `json:"description"`
	HideChildTree    bool             `json:"hide_child_tree"`
	ID               string           `json:"id"`
	Link             *LinkReference   `json:"link,omitempty"`
	Name             string           `json:"name"`
	Owner            ProfileReference `json:"owner"`
	PrimaryImage     *Asset           `json:"primary_image,omitempty"`
	Slug             string           `json:"slug"`
	Tags             []TagReference   `json:"tags"`
	UpdatedAt        *time.Time       `json:"updatedAt,omitempty"`
	Visibility       Visibility       `json:"visibility"`
}

type Property struct {
	Fid   string `json:"fid"`
	Name  string `json:"name"`
	Sort  string `json:"sort"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PropertyList = []Property

type PropertySchema struct {
	Fid  string `json:"fid"`
	Name string `json:"name"`
	Sort string `json:"sort"`
	Type string `json:"type"`
}

type PropertySchemaList = []PropertySchema

// The full node shape, including children, properties, and child property schema.
type NodeWithChildren struct {
	Assets              []Asset            `json:"assets"`
	ChildPropertySchema PropertySchemaList `json:"child_property_schema"`
	Children            []NodeWithChildren `json:"children"`
	// Rich HTML content.
	Content          *string          `json:"content,omitempty"`
	CreatedAt        *time.Time       `json:"createdAt,omitempty"`
	CurrentVersionID *string          `json:"current_version_id,omitempty"`
	DeletedAt        *time.Time       `json:"deletedAt,omitempty"`
	Description      string           `json:"description"`
	HideChildTree    bool             `json:"hide_child_tree"`
	ID               string           `json:"id"`
	Link             *LinkReference   `json:"link,omitempty"`
	Name             string           `json:"name"`
	Owner            ProfileReference `json:"owner"`
	PrimaryImage     *Asset           `json:"primary_image,omitempty"`
	Properties       PropertyList     `json:"properties"`
	Slug             string           `json:"slug"`
	Tags             []TagReference   `json:"tags"`
	UpdatedAt        *time.Time       `json:"updatedAt,omitempty"`
	Visibility       Visibility       `json:"visibility"`
}

// Composed into every paginated list result. These fields are present for a single page; with --all, a command aggregates every page into one bare array and omits them (there is no single "current page").
type PaginatedResult struct {
	CurrentPage *int `json:"current_page,omitempty"`
	NextPage    *int `json:"next_page,omitempty"`
	PageSize    *int `json:"page_size,omitempty"`
	Results     *int `json:"results,omitempty"`
	TotalPages  *int `json:"total_pages,omitempty"`
}

type NodeListResult struct {
	CurrentPage *int               `json:"current_page,omitempty"`
	NextPage    *int               `json:"next_page,omitempty"`
	Nodes       []NodeWithChildren `json:"nodes"`
	PageSize    *int               `json:"page_size,omitempty"`
	Results     *int               `json:"results,omitempty"`
	TotalPages  *int               `json:"total_pages,omitempty"`
}

// How this plugin is run — supervised (mode only) or external (mode + a static bearer token).
type PluginConnection struct {
	Mode string `json:"mode"`
	// Present only when mode is external.
	Token *string `json:"token,omitempty"`
}

// The plugin's manifest.yaml, parsed.
type PluginManifest = interface{}

type Plugin struct {
	AddedAt time.Time `json:"added_at"`
	// How this plugin is run — supervised (mode only) or external (mode + a static bearer token).
	Connection  PluginConnection `json:"connection"`
	Description *string          `json:"description,omitempty"`
	ID          string           `json:"id"`
	Manifest    PluginManifest   `json:"manifest"`
	Name        string           `json:"name"`
	// e.g. running, stopped, error.
	Status  string  `json:"status"`
	Version *string `json:"version,omitempty"`
}

type ReplyStatus struct {
	// If requested by an authenticated account, how many of those replies are theirs.
	Replied int `json:"replied"`
	// The total number of replies to the thread.
	Replies int `json:"replies"`
}

type SearchNodeListResult struct {
	CurrentPage *int   `json:"current_page,omitempty"`
	NextPage    *int   `json:"next_page,omitempty"`
	Nodes       []Node `json:"nodes"`
	PageSize    *int   `json:"page_size,omitempty"`
	Results     *int   `json:"results,omitempty"`
	TotalPages  *int   `json:"total_pages,omitempty"`
}

type SearchResult struct {
	CurrentPage *int            `json:"current_page,omitempty"`
	Items       []DatagraphItem `json:"items"`
	NextPage    *int            `json:"next_page,omitempty"`
	PageSize    *int            `json:"page_size,omitempty"`
	Results     *int            `json:"results,omitempty"`
	TotalPages  *int            `json:"total_pages,omitempty"`
}

type Thread struct {
	Assets []Asset          `json:"assets,omitempty"`
	Author ProfileReference `json:"author"`
	// Rich HTML content.
	Body        string             `json:"body"`
	Category    *CategoryReference `json:"category,omitempty"`
	CreatedAt   *time.Time         `json:"createdAt,omitempty"`
	Description *string            `json:"description,omitempty"`
	ID          string             `json:"id"`
	Link        *LinkReference     `json:"link,omitempty"`
	Mark        string             `json:"mark"`
	// Pin rank among other pinned threads; 0 means not pinned.
	Pinned      *int           `json:"pinned,omitempty"`
	ReplyStatus *ReplyStatus   `json:"reply_status,omitempty"`
	Tags        []TagReference `json:"tags,omitempty"`
	Title       string         `json:"title"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
	Visibility  Visibility     `json:"visibility"`
}

type ThreadListResult struct {
	CurrentPage *int     `json:"current_page,omitempty"`
	NextPage    *int     `json:"next_page,omitempty"`
	PageSize    *int     `json:"page_size,omitempty"`
	Results     *int     `json:"results,omitempty"`
	Threads     []Thread `json:"threads"`
	TotalPages  *int     `json:"total_pages,omitempty"`
}
