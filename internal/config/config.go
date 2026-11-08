// Package config contains all environment variable based configuration.
// THIS FILE IS GENERATED. DO NOT EDIT MANUALLY.
// To edit configuration variables, edit the config.yaml file and run codegen.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents environment variable configuration parameters

type Config struct {

	// -
	// General
	// -

	/*
	   Can be set to either:

	   - `debug`
	   - `info`
	   - `warn`
	   - `error`
	*/
	LogLevel slog.Level
	/*
	   Can be set to either:

	   - `(not set)` (default) somewhat human readable "logfmt" format logs for simple setups
	   - `dev` for developer-friendly logs, with colours and attributes on separate lines for readability
	   - `json` for machine-readable logs, mainly for log aggregators, etc.
	*/
	LogFormat string
	/*
	   Determines whether or not the backend service will also start the frontend Node.js process. When empty, it will not

	   When a path is provided, Storyden will execute `node <path>` to start the frontend process. This is used by the fullstack Docker image to start the frontend process in the same container as the backend.
	*/
	RunFrontend string
	/*
	   Used in conjunction with `RUN_FRONTEND`. This is the address that the frontend will be available at. This is used by the fullstack Docker image to proxy requests that don't match any `/api` or other routes to the frontend process.

	   In the default `fullstack` image, this is set to `http://localhost:3000` which is the default port for the Next.js process.
	*/
	FrontendProxy url.URL
	/*
	   Comma or newline separated list of CIDR ranges (or plain IPs) that are allowed as direct source addresses for SSR-origin API requests (`X-Storyden-SSR: 1`) when resolving `X-Forwarded-For` in `xff_trusted_proxies` mode.

	   This is useful when your frontend SSR process runs outside the backend host and therefore SSR API calls do not originate from loopback.

	   Storyden also derives trusted SSR source addresses from `PROXY_FRONTEND_ADDRESS` when possible:
	   - `localhost` trusts both `127.0.0.1/32` and `::1/128`
	   - literal IP hosts are trusted automatically
	*/
	SSRTrustedSourceCIDRs string

	// -
	// Development tools
	// -

	/*
	   Simulates slow requests.

	   This will add a random delay between zero and this value to all requests. This is useful for testing how the client handles slow responses.
	*/
	DevChaosSlowMode time.Duration
	/*
	   Simulates slow message delivery in the internal message queue.

	   This will add a random delay between zero and this value to all messages in the internal message queue. This is useful for testing how the client handles delayed message processing.
	*/
	DevChaosSlowModeQueue time.Duration
	/*
	   A value between 0 and 1 which simulates failed requests.

	   This will add a random failure to all requests. This is useful for testing how the client handles "internal server error" responses.
	*/
	DevChaosFailRate float64

	// -
	// Core configuration
	// -

	/*
	   The database URL to connect to. This can be a SQLite, PostgreSQL, or CockroachDB URL.

	   The accepted schemes for this URL are:
	   - `sqlite://` or `sqlite3://` for SQLite or Litestream.
	   - `postgres://` or `postgresql://` for PostgreSQL, CockroachDB and any other PostgreSQL-compatible database
	   - `libsql://` for Turso/libSQL.
	     - Remote Turso example: `libsql://your-db.turso.io?authToken=...`
	     - Local file-backed libSQL example: `libsql://./relative/path` or `libsql:///absolute/path`
	*/
	DatabaseURL string
	/*
	   The interface on which the API service will for HTTP requests.

	   Typically, in a containerised environment, this should be all interfaces (`0.0.0.0`.)
	*/
	ListenAddr string
	/*
	   The address at which the web frontend will be hosted.

	   This must be set to the public URL that users of the instance will access the frontend client. It is used to determine things such as cookie domain attributes, CORS policy, WebAuthn attributes and other necessary settings. The scheme may be used by some internal components to determine whether the instance is running in a secure context or not.

	   This is by default `http://localhost:3000` when running locally, or when deploying to production, `https://<your-domain>`.
	*/
	PublicWebAddress url.URL
	/*
	   The address at which the public API will be accessible.

	   This is also used for things such as cookies, CORS, etc.

	   Please note that both the public API address and public web address must share the same root domain name as Storyden cookies are configured to be issued under this assumption. It also makes a lot of cross-origin and cookie configurations easier to make secure.
	*/
	PublicAPIAddress url.URL
	/*
	   Additional browser origins permitted to make credentialed cross-origin requests.

	   The public web and API addresses are always trusted. Set this only when other first-party frontends (such as a separate marketing site or admin console) must call the API from the browser with cookies. Each entry is an origin (`scheme://host[:port]`), comma-separated. Arbitrary third-party origins must not be added here; programmatic clients should authenticate with bearer tokens, which do not depend on CORS credentials.
	*/
	CORSAllowedOrigins []string

	// -
	// Rate limiting
	// -

	/*
	   Maximum number of "units" allowed within the sliding window defined by `RATE_LIMIT_PERIOD`.

	   Most incoming requests will consume `1` unit for authenticated users, and `RATE_LIMIT_GUEST_COST`
	   units for unauthenticated (guest) visitors.

	   Certain endpoints are more expensive by default, such as those that trigger password resets, sending emails, etc.

	   You can also configure custom overrides in the System Settings screen. However, you cannot configure custom overrides via environment variables.
	*/
	RateLimit int
	/*
	   Sliding window duration used to enforce `RATE_LIMIT`.

	   On each request, Storyden considers the total number of units consumed in the last
	   `RATE_LIMIT_PERIOD` (not aligned to the hour/minute boundary).
	*/
	RateLimitPeriod time.Duration
	/*
	   Bucket size (granularity) used to approximate the sliding window.

	   Requests are counted into discrete time buckets of this size (e.g. 1 minute). When
	   checking the limit, Storyden sums all buckets whose timestamps fall within the last
	   `RATE_LIMIT_PERIOD` and discards older buckets.

	   Smaller buckets = more accurate sliding-window behaviour but more storage/CPU overhead.
	   Larger buckets = cheaper but "chunkier" enforcement.

	   Rule of thumb: set this to ~1/60 of `RATE_LIMIT_PERIOD` (e.g. 1m buckets for a 1h period).
	*/
	RateLimitBucket time.Duration
	/*
	   Cost multiplier applied to unauthenticated (guest) requests.

	   Example: if set to 5, each guest request consumes 5 units from the same `RATE_LIMIT`
	   budget, effectively allowing guests ~1/5th the throughput of authenticated users.
	*/
	RateLimitGuestCost int

	// -
	// Telemetry and monitoring
	// -

	/*
	   Either:
	   - `otlp` for any standard OpenTelemetry collector, or `file://...` JSONL output when `OTEL_EXPORTER_OTLP_ENDPOINT` is a file URL.
	   - `sentry` for Sentry (which is OpenTelemetry-compatible, however requires its own specific configuration.)
	   - `logger` for local logging to the console. This is only really useful for Storyden developers and is very noisy.
	*/
	OTELProvider string
	/*
	   The collector endpoint for sending OTEL data.

	   This also supports file output for local profiling when `OTEL_PROVIDER=otlp`
	   by setting a `file://` URL, for example: `file:///tmp/storyden-traces.jsonl`
	   or `file:./otel-logs.jsonl`.
	*/
	OTELEndpoint url.URL
	// When `OTEL_PROVIDER` is set to `sentry`, this is the DSN for the Sentry project.
	SentryDSN string

	// -
	// Web Push
	// -

	// Enable Web Push delivery and subscription registration. When false, all other Web Push configuration is ignored.
	WebPushEnabled bool
	// Maximum push subscriptions per account. Registering a new endpoint at the limit silently removes the oldest subscription. Updating an existing endpoint does not change its age. Must be positive when Web Push is enabled; ignored otherwise.
	WebPushMaxSubscriptions int
	// Unpadded base64url-encoded uncompressed P-256 public key. Required when Web Push is enabled. Keep the key pair stable across restarts; replacing it requires browsers to subscribe again.
	WebPushVAPIDPublicKey string
	// Unpadded base64url-encoded 32-byte P-256 private key matching the public key. Required when Web Push is enabled. Keep this secret.
	WebPushVAPIDPrivateKey string
	// Contact URI for the instance operator, such as mailto:admin@example.com or https://example.com/contact. Required when Web Push is enabled.
	WebPushVAPIDSubject string

	// -
	// Email
	// -

	/*
	   Either:

	   - unset (default) for no email sending. Email sending is not a requirement for a production deployment.
	   - `smtp` for SMTP-based email sending through any SMTP server.
	   - `sendgrid` for SendGrid based email sending.
	   - `mock` for logging emails to the console. Only useful for Storyden developers and testing.
	*/
	EmailProvider string
	/*
	   The hostname of the SMTP server. This is required for sending emails via SMTP.

	   This is typically something like `smtp.gmail.com`, `smtp.outlook.com`, or your organization's SMTP server.
	*/
	SMTPHost string
	/*
	   The port of the SMTP server. This is required for sending emails via SMTP.

	   Common values are:
	   - `587` for STARTTLS (recommended)
	   - `465` for implicit TLS/SSL
	   - `25` for unencrypted SMTP (not recommended for production)
	*/
	SMTPPort int
	// The username for SMTP authentication. This is typically your email address.
	SMTPUsername string
	// The password for SMTP authentication. For services like Gmail, this may be an "app password" rather than your regular password.
	SMTPPassword string
	/*
	   The name that will be used as the sender name for emails sent via SMTP.

	   This is typically the name of your community or organisation.
	*/
	SMTPFromName string
	/*
	   The email address that will be used as the sender address for emails sent via SMTP.

	   This is typically a no-reply address, such as `no-reply@<your-domain>`.
	*/
	SMTPFromAddress string
	/*
	   Whether to require TLS encryption for SMTP connections.

	   When enabled, Storyden uses implicit TLS on port `465` and STARTTLS on other ports.
	   Set to `false` only if your SMTP server does not support encryption (not recommended).
	*/
	SMTPUseTLS bool
	/*
	   The name that will be used as the sender name for emails sent via SendGrid.

	   This is typically the name of your community or organisation.
	*/
	SendGridFromName string
	/*
	   The email address that will be used as the sender address for emails sent via SendGrid.

	   This is typically a no-reply address, such as `no-reply@<your-domain>`.
	*/
	SendGridFromAddress string
	/*
	   The API key for the SendGrid account. This is required for sending emails via SendGrid.

	   This is typically a long string of characters that you can generate in the SendGrid dashboard.
	*/
	SendGridAPIKey string

	// -
	// Authentication
	// -

	/*
	   The secret key used to sign JWT tokens. This is used for authentication and should be kept secret.

	   This is typically a long string of characters that you can generate using a secure random generator such as `openssl rand -hex 12`.

	   The JWT secret is required if you enable any of the OAuth providers or enable email features. This is because JWTs are used to verify callbacks as well as verify password-reset and other tokens.
	*/
	JWTSecret []byte
	// Enable Google SSO authentication.
	GoogleEnabled bool
	// The client ID for the Google OAuth2 application.
	GoogleClientID string
	// The client secret for the Google OAuth2 application.
	GoogleClientSecret string
	// Enable GitHub SSO authentication.
	GitHubEnabled bool
	// The client ID for the GitHub OAuth2 application.
	GitHubClientID string
	// The client secret for the GitHub OAuth2 application.
	GitHubClientSecret string
	// Enable Discord SSO authentication.
	DiscordEnabled bool
	// The client ID for the Discord OAuth2 application.
	DiscordClientID string
	// The client secret for the Discord OAuth2 application.
	DiscordClientSecret string
	// Enable Keycloak OIDC authentication.
	KeycloakEnabled bool
	// The client ID for the Keycloak OAuth2 application.
	KeycloakClientID string
	// The client secret for the Keycloak OAuth2 application.
	KeycloakClientSecret string
	// The issuer/discovery URL for the Keycloak realm (e.g. https://auth.example.com/realms/YourRealm).
	KeycloakIssuerURL url.URL
	// The display name shown for the Keycloak login button.
	KeycloakDisplayName string
	// Enable Storyden's built-in OAuth2/OIDC authorization server endpoints.
	OAuthEnabled bool
	// Access token lifetime for Storyden OAuth tokens.
	OAuthAccessTokenTTL time.Duration
	// Refresh token lifetime for Storyden OAuth tokens.
	OAuthRefreshTokenTTL time.Duration
	// Device code lifetime for the OAuth Device Authorization Grant.
	OAuthDeviceCodeTTL time.Duration
	// Poll interval used by device flow responses.
	OAuthDevicePollEvery time.Duration
	// Frontend URL used by OAuth Device Authorization Grant users to approve or deny consent.
	OAuthDeviceAuthorisationConsentURL url.URL
	// Frontend URL used by OAuth Authorization Code Grant users to approve or deny consent.
	OAuthAuthorisationCodeConsentURL url.URL
	// Frontend URL to redirect unauthenticated users to when starting the OAuth Authorization Code flow without a session. Defaults to `{PUBLIC_WEB_ADDRESS}/login` when unset.
	OAuthAuthorisationLoginURL url.URL
	// Base64-encoded PEM private signing key used for OAuth2/OIDC JWT signing.
	OAuthSigningKeyBase64 string
	// Optional JWT key ID (kid) for OAuth signing keys.
	OAuthSigningKeyID string
	// Advertise and accept OAuth Client ID Metadata Documents (CIMD), letting clients identify themselves with an https URL client_id that resolves to a hosted metadata document instead of pre-registering.
	OAuthClientIDMetadataDocumentEnabled bool
	// Comma-separated list of permission scopes a CIMD client may request. Empty uses a conservative read-only default. Privileged scopes are still gated by OAUTH_CIMD_ALLOW_PRIVILEGED_SCOPES.
	OAuthCIMDAllowedScopes []string
	/*
	   Allow privileged/administrative permission scopes to be granted to CIMD clients.

	   Off by default; only enable for trusted, controlled deployments.
	*/
	OAuthCIMDAllowPrivilegedScopes bool
	// Relax CIMD metadata fetching restrictions (allow http, private hosts and skip TLS verification). For local development and testing only; never set in production.
	OAuthCIMDAllowInsecureFetch bool

	// -
	// SMS
	// -

	/*
	   Either:

	   - unset (default) for no SMS sending. SMS sending is not a requirement for a production deployment.
	   - `twilio` for Twilio based SMS sending.
	   - `mock` for logging SMS to the console. Only useful for Storyden developers and testing.
	*/
	SMSProvider string
	/*
	   The account SID for the Twilio account.

	   This is typically a long string of characters that you can view in the Twilio dashboard.
	*/
	TwilioAccountSID string
	// The phone number that will be used as the sender number for SMS sent via Twilio.
	TwilioPhoneNumber string
	/*
	   The auth token for the Twilio account. This is required for sending SMS via Twilio.

	   This is typically a long string of characters that you can generate in the Twilio dashboard.
	*/
	TwilioAuthToken string

	// -
	// Assets/file storage
	// -

	/*
	   Either:

	   - `local` for local file storage.
	   - `s3` for any Amazon S3-compatible storage, such as S3 itself (obviously...), Google Cloud Storage, Cloudflare R2, Minio, etc.
	*/
	AssetStorageType string
	// When `ASSET_STORAGE_TYPE` is set to `local`, this is the path to the directory where files will be stored.
	AssetStorageLocalPath string
	// When `ASSET_STORAGE_TYPE` is set to `s3`, this determines whether or not to use HTTPS for the S3 connection. You should always set this to `true` unless your S3-compatible storage provider is internally but not publicly accessible, such as in a Kubernetes cluster or running on the same host.
	S3Secure bool
	// The endpoint for the S3-compatible storage provider. This is typically the base URL of the provider, such as `https://s3.amazonaws.com` for AWS S3, or `https://storage.googleapis.com` for Google Cloud Storage, etc.
	S3Endpoint string
	// The bucket name for Storyden assets to be stored in.
	S3Bucket string
	/*
	   Most S3-compatible storage providers require a region to be specified. This is typically the region in which the bucket is located, such as `us-east-1` for AWS S3.

	   However, some providers do not use regions but S3-compatible clients still require this to be set. In most cases, the provider will give you a value for this, such as `auto` when using Cloudflare R2.
	*/
	S3Region string
	// The access key for the S3-compatible storage provider.
	S3AccessKey string
	// The secret key for the S3-compatible storage provider.
	S3SecretKey string

	// -
	// Image processing
	// -

	// Maximum simultaneously active image jobs per process.
	ImageWorkers int
	// Admission budget in bytes for image work, separate from application memory. This is not a process RSS limit.
	ImageMemoryBudget int64

	// -
	// Cache
	// -

	/*
	   When empty, caching will use an efficient in-memory store. This is usually fine for small to medium-sized deployments however it's worth keeping an eye on your deployment's machine memory usage.

	   When set to `redis`, Storyden will use Redis as a cache provider. This is recommended for larger deployments that receive a lot of traffic. The cache provider is also used for the rate limiter so that it can be shared across multiple instances of Storyden.

	   This is necessary for deploying replica instances of Storyden that are backed by the same persistence layers (database, asset storage, etc.)
	*/
	CacheProvider string
	/*
	   The Redis URL to connect to.

	   This is a full URL with `redis://` as the scheme. You can set the username and password using the URL format, for example: `redis://<username>:<password>@<host>:<port>`.
	*/
	RedisURL url.URL

	// -
	// Search features
	// -

	/*
	   Either:

	   - `database` for the default database-driven search.  This is not recommended for larger deployments as it does not scale well and has limited search quality.
	   - `bleve` for Bleve. This is a local full-text search engine that is fast and efficient for small to medium-sized deployments. This is best used when you are using local disk storage, such as SQLite and local asset storage.
	   - `redis` for Redisearch. This is a fast and efficient search provider that is recommended for larger deployments. This is recommended if your environment is ephemeral and you're already using external providers for database and asset storage.
	*/
	SearchProvider string
	/*
	   When using `SEARCH_PROVIDER` set to either `bleve` or `redis`, this is the number of items that will be indexed in a single batch.

	   Increasing this value will improve indexing performance, but will also increase memory usage during indexing.
	*/
	SearchIndexChunkSize int
	// The path to the directory where Bleve will store search indexes. Only used when `SEARCH_PROVIDER` is set to `bleve`.
	BlevePath string
	// The name of the Redis search index. Only used when `SEARCH_PROVIDER` is set to `redis`.
	RedisSearchIndexName string

	// -
	// Message queue
	// -

	/*
	   Either:

	   - Default (no value): in-memory Go channels. This is fast and efficient, but not persistent across restarts and will add a bit of memory usage to the process.
	   - `amqp`: RabbitMQ. This is a persistent message queue that is fast and reliable. It is recommended for larger deployments and is necessary for deploying replica instances of Storyden.
	*/
	QueueType string
	/*
	   The RabbitMQ URL to connect to.

	   This is a full URL with `amqp://` as the scheme. You can set the username and password using the URL format, for example: `amqp://<username>:<password>@<host>:<port>`.

	   The default value is `amqp://guest:guest@localhost:5672/` which is the default RabbitMQ URL.

	   Storyden does not currently support `amqps://` (secure) URLs, but this will be added soon.
	*/
	AmqpURL string
	/*
	   The maximum number of times a failed message will be retried before being moved to the dead letter queue.

	   Messages are retried with exponential backoff starting at 1 second and doubling each time up to a maximum of 1 minute between retries.
	*/
	QueueMaxRetries int
	// The initial interval to wait before the first retry attempt.
	QueueRetryInitialInterval time.Duration
	// The maximum interval to wait between retry attempts. The exponential backoff will not exceed this value.
	QueueRetryMaxInterval time.Duration

	// -
	// Plugins
	// -

	/*
	   The directory where plugin files will be extracted and stored. Each plugin gets its own subdirectory keyed by the plugin installation ID.

	   This directory should be persistent and writable by the Storyden process.
	*/
	PluginDataPath string
	/*
	   The directory where local Robot workspace instance files will be stored. Each workspace instance gets its own subdirectory keyed by the workspace instance ID.

	   This directory should be persistent and writable by the Storyden process when local Robot workspaces are used.
	*/
	RobotWorkspaceDataPath string
	/*
	   The plugin runtime provider. Different runtime providers offer different security guarantees. The simplest is `local` which just runs the plugin as a child process on the same machine as Storyden.

	   - `none`: disables plugins entirely. Plugin APIs return a permission error and instance capabilities will not include `plugins`.
	   - `local`: runs supervised plugins as local child processes on the same machine as Storyden.
	   - `sprites`: runs supervised plugins in isolated [Sprites](https://docs.sprites.dev/) runtimes.
	*/
	PluginRuntimeProvider string
	/*
	   Required when `PLUGIN_RUNTIME_PROVIDER` is set to `sprites`.

	   This is the API key Storyden uses to create and manage plugin runtimes on Sprites.
	*/
	SpritesAPIKey string
	/*
	   Maximum number of consecutive restart attempts for a supervised plugin before marking it as errored.

	   When a plugin crashes during startup (before the runtime crash threshold), the restart counter increments. After exceeding this limit, the plugin transitions to an error state and stops attempting restarts.
	*/
	PluginMaxRestartAttempts int
	/*
	   Maximum backoff duration between plugin restart attempts.

	   Restart delays use exponential backoff (1s, 2s, 4s, 8s, ...) capped at this maximum value.
	*/
	PluginMaxBackoffDuration time.Duration
	/*
	   Time threshold to distinguish between startup crashes and runtime crashes.

	   If a plugin runs successfully for longer than this duration before crashing, it's considered a "runtime crash" and the restart counter is reset. This allows plugins that crash after running for a while (hours/days) to restart without being penalized by the startup crash limit.
	*/
	PluginRuntimeCrashThreshold time.Duration
	/*
	   Backoff duration used specifically for runtime crashes (crashes that occur after the plugin has been running for longer than the runtime crash threshold).

	   This is typically shorter than the exponential backoff used for startup crashes, as runtime crashes are less likely to be configuration issues.
	*/
	PluginRuntimeCrashBackoff time.Duration

	// -
	// Artificial intelligence/language models
	// -

	/*
	   Enables the Model Context Protocol server, accessible via Streamable HTTP at `/mcp` on the public API origin. When enabled, the server publishes an MCP Server Card at `/mcp/server-card` and adds it to the web origin's `/.well-known/ai-catalog.json`.

	   This is used to integrate Storyden into agentic workflow engines and other language model tooling.

	   See [the documentation](https://storyden.org/docs/introduction/mcp) for more information.
	*/
	MCPEnabled bool
	/*
	   The provider for language model features.

	   `openai` is currently the only supported provider.
	*/
	LanguageModelProvider string
	// When `LANGUAGE_MODEL_PROVIDER` is set to `openai`, this is the API key for the OpenAI API.
	OpenAIKey string
	/*
	   Write a Markdown snapshot of each Robot model request after dynamic instructions, tools, and conversation history have been assembled. Includes the full system prompt, message parts, tool inputs and outputs, and tool declarations, without redaction or truncation. Dump metadata appears first, followed by tool declarations, system instructions, and messages.

	   Disabled by default. Snapshots use the configured object storage under `debug/model-requests/<session-id>/<timestamp>-<unique-id>.md`, so local storage writes beneath `ASSET_STORAGE_LOCAL_PATH` (normally `./data`) and S3 storage writes to the configured bucket. Each model call gets a separate snapshot, including multiple calls within a single conversation turn. Dumps capture the final ADK request before provider-specific conversion, not HTTP payloads. Callable-tool inventories and ADK request configuration are excluded. Section ordering follows the documented prompt-prefix layout; providers control the final model input.
	*/
	RobotModelDebugDump bool

	// -
	// Semdex
	// -

	/*
	   Either:
	   - `chromem` for an experimental local vector database. This is not recommended for use in large deployments as it's rather slow and memory-hungry.
	   - `pinecone` for Pinecone, a fully managed vector database.
	*/
	SemdexProvider string

	// -
	// Local Semdex
	// -

	// The path to the directory where Chromem will store vector indexes.
	SemdexLocalPath string

	// -
	// Pinecone Semdex
	// -

	// Your Pinecone API key. This is required for all Pinecone API requests.
	PineconeAPIKey string
	// The index name that Storyden will use in your Pinecone workspace.
	PineconeIndex string
	// This value is dependent on the underlying OpenAI configuration. Currently this is static and set to 3072 dimensions. In future, Storyden will provide more flexible configuration for language model providers.
	PineconeDimensions int32
	// Pinecone provides hosting on different cloud providers, see the Pinecone documentation for more information. The cloud provider you choose will be reflected in your Pinecone dashboard.
	PineconeCloud string
	// Same as above, but for the region. As with any third party providers, it's recommended to choose the region closest to both your Storyden deployment and your community members for best performance and experience.
	PineconeRegion string
}

const DefaultImageWorkers int = 1
const DefaultImageMemoryBudget int64 = 268435456

var errInvalidConfiguration = errors.New("invalid configuration")

func Load() (Config, error) {
	return load(os.LookupEnv)
}

func load(lookup func(string) (string, bool)) (Config, error) {
	var c Config

	{
		value, present := lookup("LOG_LEVEL")
		if !present {
			value = "info"
		}

		if err := c.LogLevel.UnmarshalText([]byte(value)); err != nil {
			return c, fmt.Errorf("%w: LOG_LEVEL must be a log level", errInvalidConfiguration)
		}
	}

	{
		value, present := lookup("LOG_FORMAT")
		if present {
			c.LogFormat = value
		}
	}

	{
		value, present := lookup("RUN_FRONTEND")
		if !present {
			value = ""
		}

		c.RunFrontend = value
	}

	{
		value, present := lookup("PROXY_FRONTEND_ADDRESS")
		if !present {
			value = ""
		}

		parsed, err := url.Parse(value)
		if err != nil {
			return c, fmt.Errorf("%w: PROXY_FRONTEND_ADDRESS has an invalid net/url.URL value", errInvalidConfiguration)
		}

		c.FrontendProxy = *parsed
	}

	{
		value, present := lookup("SSR_TRUSTED_SOURCE_CIDRS")
		if !present {
			value = ""
		}

		c.SSRTrustedSourceCIDRs = value
	}

	{
		value, present := lookup("DEV_CHAOS_SLOW_MODE")
		if present {
			parsed, err := time.ParseDuration(value)
			if err != nil {
				return c, fmt.Errorf("%w: DEV_CHAOS_SLOW_MODE has an invalid time.Duration value", errInvalidConfiguration)
			}

			c.DevChaosSlowMode = parsed
		}
	}

	{
		value, present := lookup("DEV_CHAOS_SLOW_MODE_QUEUE")
		if present {
			parsed, err := time.ParseDuration(value)
			if err != nil {
				return c, fmt.Errorf("%w: DEV_CHAOS_SLOW_MODE_QUEUE has an invalid time.Duration value", errInvalidConfiguration)
			}

			c.DevChaosSlowModeQueue = parsed
		}
	}

	{
		value, present := lookup("DEV_CHAOS_FAIL_RATE")
		if present {
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return c, fmt.Errorf("%w: DEV_CHAOS_FAIL_RATE has an invalid float64 value", errInvalidConfiguration)
			}

			c.DevChaosFailRate = parsed
		}
	}

	{
		value, present := lookup("DATABASE_URL")
		if !present {
			value = "sqlite://data/data.db?_pragma=foreign_keys(1)"
		}

		c.DatabaseURL = value
	}

	{
		value, present := lookup("LISTEN_ADDR")
		if !present {
			value = "0.0.0.0:8000"
		}

		c.ListenAddr = value
	}

	{
		value, present := lookup("PUBLIC_WEB_ADDRESS")
		if !present {
			value = "http://localhost:3000"
		}

		parsed, err := url.Parse(value)
		if err != nil {
			return c, fmt.Errorf("%w: PUBLIC_WEB_ADDRESS has an invalid net/url.URL value", errInvalidConfiguration)
		}

		c.PublicWebAddress = *parsed
	}

	{
		value, present := lookup("PUBLIC_API_ADDRESS")
		if !present {
			value = "http://localhost:8000"
		}

		parsed, err := url.Parse(value)
		if err != nil {
			return c, fmt.Errorf("%w: PUBLIC_API_ADDRESS has an invalid net/url.URL value", errInvalidConfiguration)
		}

		c.PublicAPIAddress = *parsed
	}

	{
		value, present := lookup("CORS_ALLOWED_ORIGINS")
		if present {
			c.CORSAllowedOrigins = []string{}
			if strings.TrimSpace(value) != "" {
				c.CORSAllowedOrigins = strings.Split(value, ",")
			}
		}
	}

	{
		value, present := lookup("RATE_LIMIT")
		if !present {
			value = "5000"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: RATE_LIMIT has an invalid int value", errInvalidConfiguration)
		}

		c.RateLimit = int(parsed)
	}

	{
		value, present := lookup("RATE_LIMIT_PERIOD")
		if !present {
			value = "1h"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: RATE_LIMIT_PERIOD has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.RateLimitPeriod = parsed
	}

	{
		value, present := lookup("RATE_LIMIT_BUCKET")
		if !present {
			value = "1m"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: RATE_LIMIT_BUCKET has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.RateLimitBucket = parsed
	}

	{
		value, present := lookup("RATE_LIMIT_GUEST_COST")
		if !present {
			value = "1"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: RATE_LIMIT_GUEST_COST has an invalid int value", errInvalidConfiguration)
		}

		c.RateLimitGuestCost = int(parsed)
	}

	{
		value, present := lookup("OTEL_PROVIDER")
		if !present {
			value = ""
		}

		c.OTELProvider = value
	}

	{
		value, present := lookup("OTEL_EXPORTER_OTLP_ENDPOINT")
		if !present {
			value = ""
		}

		parsed, err := url.Parse(value)
		if err != nil {
			return c, fmt.Errorf("%w: OTEL_EXPORTER_OTLP_ENDPOINT has an invalid net/url.URL value", errInvalidConfiguration)
		}

		c.OTELEndpoint = *parsed
	}

	{
		value, present := lookup("SENTRY_DSN")
		if !present {
			value = ""
		}

		c.SentryDSN = value
	}

	{
		value, present := lookup("WEB_PUSH_ENABLED")
		if !present {
			value = "false"
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("%w: WEB_PUSH_ENABLED has an invalid bool value", errInvalidConfiguration)
		}

		c.WebPushEnabled = parsed
	}

	{
		value, present := lookup("WEB_PUSH_MAX_SUBSCRIPTIONS")
		if !present {
			value = "5"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: WEB_PUSH_MAX_SUBSCRIPTIONS has an invalid int value", errInvalidConfiguration)
		}

		c.WebPushMaxSubscriptions = int(parsed)
	}

	{
		value, present := lookup("WEB_PUSH_VAPID_PUBLIC_KEY")
		if present {
			c.WebPushVAPIDPublicKey = value
		}
	}

	{
		value, present := lookup("WEB_PUSH_VAPID_PRIVATE_KEY")
		if present {
			c.WebPushVAPIDPrivateKey = value
		}
	}

	{
		value, present := lookup("WEB_PUSH_VAPID_SUBJECT")
		if present {
			c.WebPushVAPIDSubject = value
		}
	}

	{
		value, present := lookup("EMAIL_PROVIDER")
		if present {
			c.EmailProvider = value
		}
	}

	{
		value, present := lookup("SMTP_HOST")
		if present {
			c.SMTPHost = value
		}
	}

	{
		value, present := lookup("SMTP_PORT")
		if present {
			parsed, err := strconv.ParseInt(value, 0, 0)
			if err != nil {
				return c, fmt.Errorf("%w: SMTP_PORT has an invalid int value", errInvalidConfiguration)
			}

			c.SMTPPort = int(parsed)
		}
	}

	{
		value, present := lookup("SMTP_USERNAME")
		if present {
			c.SMTPUsername = value
		}
	}

	{
		value, present := lookup("SMTP_PASSWORD")
		if present {
			c.SMTPPassword = value
		}
	}

	{
		value, present := lookup("SMTP_FROM_NAME")
		if present {
			c.SMTPFromName = value
		}
	}

	{
		value, present := lookup("SMTP_FROM_ADDRESS")
		if present {
			c.SMTPFromAddress = value
		}
	}

	{
		value, present := lookup("SMTP_USE_TLS")
		if !present {
			value = "true"
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("%w: SMTP_USE_TLS has an invalid bool value", errInvalidConfiguration)
		}

		c.SMTPUseTLS = parsed
	}

	{
		value, present := lookup("SENDGRID_FROM_NAME")
		if present {
			c.SendGridFromName = value
		}
	}

	{
		value, present := lookup("SENDGRID_FROM_ADDRESS")
		if present {
			c.SendGridFromAddress = value
		}
	}

	{
		value, present := lookup("SENDGRID_API_KEY")
		if present {
			c.SendGridAPIKey = value
		}
	}

	{
		value, present := lookup("JWT_SECRET")
		if present {
			c.JWTSecret = []byte(value)
		}
	}

	{
		value, present := lookup("OAUTH_GOOGLE_ENABLED")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_GOOGLE_ENABLED has an invalid bool value", errInvalidConfiguration)
			}

			c.GoogleEnabled = parsed
		}
	}

	{
		value, present := lookup("OAUTH_GOOGLE_CLIENT_ID")
		if present {
			c.GoogleClientID = value
		}
	}

	{
		value, present := lookup("OAUTH_GOOGLE_CLIENT_SECRET")
		if present {
			c.GoogleClientSecret = value
		}
	}

	{
		value, present := lookup("OAUTH_GITHUB_ENABLED")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_GITHUB_ENABLED has an invalid bool value", errInvalidConfiguration)
			}

			c.GitHubEnabled = parsed
		}
	}

	{
		value, present := lookup("OAUTH_GITHUB_CLIENT_ID")
		if present {
			c.GitHubClientID = value
		}
	}

	{
		value, present := lookup("OAUTH_GITHUB_CLIENT_SECRET")
		if present {
			c.GitHubClientSecret = value
		}
	}

	{
		value, present := lookup("OAUTH_DISCORD_ENABLED")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_DISCORD_ENABLED has an invalid bool value", errInvalidConfiguration)
			}

			c.DiscordEnabled = parsed
		}
	}

	{
		value, present := lookup("OAUTH_DISCORD_CLIENT_ID")
		if present {
			c.DiscordClientID = value
		}
	}

	{
		value, present := lookup("OAUTH_DISCORD_CLIENT_SECRET")
		if present {
			c.DiscordClientSecret = value
		}
	}

	{
		value, present := lookup("OAUTH_KEYCLOAK_ENABLED")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_KEYCLOAK_ENABLED has an invalid bool value", errInvalidConfiguration)
			}

			c.KeycloakEnabled = parsed
		}
	}

	{
		value, present := lookup("OAUTH_KEYCLOAK_CLIENT_ID")
		if present {
			c.KeycloakClientID = value
		}
	}

	{
		value, present := lookup("OAUTH_KEYCLOAK_CLIENT_SECRET")
		if present {
			c.KeycloakClientSecret = value
		}
	}

	{
		value, present := lookup("OAUTH_KEYCLOAK_ISSUER_URL")
		if present {
			parsed, err := url.Parse(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_KEYCLOAK_ISSUER_URL has an invalid net/url.URL value", errInvalidConfiguration)
			}

			c.KeycloakIssuerURL = *parsed
		}
	}

	{
		value, present := lookup("OAUTH_KEYCLOAK_DISPLAY_NAME")
		if !present {
			value = "Keycloak"
		}

		c.KeycloakDisplayName = value
	}

	{
		value, present := lookup("OAUTH_ENABLED")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_ENABLED has an invalid bool value", errInvalidConfiguration)
			}

			c.OAuthEnabled = parsed
		}
	}

	{
		value, present := lookup("OAUTH_ACCESS_TOKEN_TTL")
		if !present {
			value = "15m"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: OAUTH_ACCESS_TOKEN_TTL has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.OAuthAccessTokenTTL = parsed
	}

	{
		value, present := lookup("OAUTH_REFRESH_TOKEN_TTL")
		if !present {
			value = "720h"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: OAUTH_REFRESH_TOKEN_TTL has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.OAuthRefreshTokenTTL = parsed
	}

	{
		value, present := lookup("OAUTH_DEVICE_CODE_TTL")
		if !present {
			value = "10m"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: OAUTH_DEVICE_CODE_TTL has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.OAuthDeviceCodeTTL = parsed
	}

	{
		value, present := lookup("OAUTH_DEVICE_POLL_EVERY")
		if !present {
			value = "5s"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: OAUTH_DEVICE_POLL_EVERY has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.OAuthDevicePollEvery = parsed
	}

	{
		value, present := lookup("OAUTH_DEVICE_AUTHORISATION_CONSENT_URL")
		if present {
			parsed, err := url.Parse(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_DEVICE_AUTHORISATION_CONSENT_URL has an invalid net/url.URL value", errInvalidConfiguration)
			}

			c.OAuthDeviceAuthorisationConsentURL = *parsed
		}
	}

	{
		value, present := lookup("OAUTH_AUTHORISATION_CODE_CONSENT_URL")
		if present {
			parsed, err := url.Parse(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_AUTHORISATION_CODE_CONSENT_URL has an invalid net/url.URL value", errInvalidConfiguration)
			}

			c.OAuthAuthorisationCodeConsentURL = *parsed
		}
	}

	{
		value, present := lookup("OAUTH_AUTHORISATION_LOGIN_URL")
		if present {
			parsed, err := url.Parse(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_AUTHORISATION_LOGIN_URL has an invalid net/url.URL value", errInvalidConfiguration)
			}

			c.OAuthAuthorisationLoginURL = *parsed
		}
	}

	{
		value, present := lookup("OAUTH_SIGNING_KEY_BASE64")
		if present {
			c.OAuthSigningKeyBase64 = value
		}
	}

	{
		value, present := lookup("OAUTH_SIGNING_KEY_ID")
		if present {
			c.OAuthSigningKeyID = value
		}
	}

	{
		value, present := lookup("OAUTH_CIMD_ENABLED")
		if !present {
			value = "false"
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("%w: OAUTH_CIMD_ENABLED has an invalid bool value", errInvalidConfiguration)
		}

		c.OAuthClientIDMetadataDocumentEnabled = parsed
	}

	{
		value, present := lookup("OAUTH_CIMD_ALLOWED_SCOPES")
		if present {
			c.OAuthCIMDAllowedScopes = []string{}
			if strings.TrimSpace(value) != "" {
				c.OAuthCIMDAllowedScopes = strings.Split(value, ",")
			}
		}
	}

	{
		value, present := lookup("OAUTH_CIMD_ALLOW_PRIVILEGED_SCOPES")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_CIMD_ALLOW_PRIVILEGED_SCOPES has an invalid bool value", errInvalidConfiguration)
			}

			c.OAuthCIMDAllowPrivilegedScopes = parsed
		}
	}

	{
		value, present := lookup("OAUTH_CIMD_ALLOW_INSECURE_FETCH")
		if present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return c, fmt.Errorf("%w: OAUTH_CIMD_ALLOW_INSECURE_FETCH has an invalid bool value", errInvalidConfiguration)
			}

			c.OAuthCIMDAllowInsecureFetch = parsed
		}
	}

	{
		value, present := lookup("SMS_PROVIDER")
		if present {
			c.SMSProvider = value
		}
	}

	{
		value, present := lookup("TWILIO_ACCOUNT_SID")
		if present {
			c.TwilioAccountSID = value
		}
	}

	{
		value, present := lookup("TWILIO_PHONE_NUMBER")
		if present {
			c.TwilioPhoneNumber = value
		}
	}

	{
		value, present := lookup("TWILIO_AUTH_TOKEN")
		if present {
			c.TwilioAuthToken = value
		}
	}

	{
		value, present := lookup("ASSET_STORAGE_TYPE")
		if present {
			c.AssetStorageType = value
		}
	}

	{
		value, present := lookup("ASSET_STORAGE_LOCAL_PATH")
		if present {
			c.AssetStorageLocalPath = value
		}
	}

	{
		value, present := lookup("S3_SECURE")
		if !present {
			value = "true"
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("%w: S3_SECURE has an invalid bool value", errInvalidConfiguration)
		}

		c.S3Secure = parsed
	}

	{
		value, present := lookup("S3_ENDPOINT")
		if present {
			c.S3Endpoint = value
		}
	}

	{
		value, present := lookup("S3_BUCKET")
		if present {
			c.S3Bucket = value
		}
	}

	{
		value, present := lookup("S3_REGION")
		if present {
			c.S3Region = value
		}
	}

	{
		value, present := lookup("S3_ACCESS_KEY")
		if present {
			c.S3AccessKey = value
		}
	}

	{
		value, present := lookup("S3_SECRET_KEY")
		if present {
			c.S3SecretKey = value
		}
	}

	{
		value, present := lookup("IMAGE_WORKERS")
		if !present {
			value = "1"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: IMAGE_WORKERS has an invalid int value", errInvalidConfiguration)
		}

		c.ImageWorkers = int(parsed)
	}

	{
		value, present := lookup("IMAGE_MEMORY_BUDGET")
		if !present {
			value = "268435456"
		}

		parsed, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return c, fmt.Errorf("%w: IMAGE_MEMORY_BUDGET has an invalid int64 value", errInvalidConfiguration)
		}

		c.ImageMemoryBudget = parsed
	}

	{
		value, present := lookup("CACHE_PROVIDER")
		if !present {
			value = ""
		}

		c.CacheProvider = value
	}

	{
		value, present := lookup("REDIS_URL")
		if !present {
			value = ""
		}

		parsed, err := url.Parse(value)
		if err != nil {
			return c, fmt.Errorf("%w: REDIS_URL has an invalid net/url.URL value", errInvalidConfiguration)
		}

		c.RedisURL = *parsed
	}

	{
		value, present := lookup("SEARCH_PROVIDER")
		if !present {
			value = "database"
		}

		c.SearchProvider = value
	}

	{
		value, present := lookup("SEARCH_INDEX_CHUNK_SIZE")
		if !present {
			value = "1000"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: SEARCH_INDEX_CHUNK_SIZE has an invalid int value", errInvalidConfiguration)
		}

		c.SearchIndexChunkSize = int(parsed)
	}

	{
		value, present := lookup("BLEVE_PATH")
		if !present {
			value = "data/bleve"
		}

		c.BlevePath = value
	}

	{
		value, present := lookup("REDIS_SEARCH_INDEX_NAME")
		if !present {
			value = "storyden"
		}

		c.RedisSearchIndexName = value
	}

	{
		value, present := lookup("QUEUE_TYPE")
		if !present {
			value = "internal"
		}

		c.QueueType = value
	}

	{
		value, present := lookup("AMQP_URL")
		if !present {
			value = "amqp://guest:guest@localhost:5672/"
		}

		c.AmqpURL = value
	}

	{
		value, present := lookup("QUEUE_MAX_RETRIES")
		if !present {
			value = "5"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: QUEUE_MAX_RETRIES has an invalid int value", errInvalidConfiguration)
		}

		c.QueueMaxRetries = int(parsed)
	}

	{
		value, present := lookup("QUEUE_RETRY_INITIAL_INTERVAL")
		if !present {
			value = "1s"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: QUEUE_RETRY_INITIAL_INTERVAL has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.QueueRetryInitialInterval = parsed
	}

	{
		value, present := lookup("QUEUE_RETRY_MAX_INTERVAL")
		if !present {
			value = "1m"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: QUEUE_RETRY_MAX_INTERVAL has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.QueueRetryMaxInterval = parsed
	}

	{
		value, present := lookup("PLUGIN_DATA_PATH")
		if !present {
			value = "./data/plugins"
		}

		c.PluginDataPath = value
	}

	{
		value, present := lookup("ROBOT_WORKSPACE_DATA_PATH")
		if !present {
			value = "./data/robot-workspaces"
		}

		c.RobotWorkspaceDataPath = value
	}

	{
		value, present := lookup("PLUGIN_RUNTIME_PROVIDER")
		if !present {
			value = "local"
		}

		c.PluginRuntimeProvider = value
	}

	{
		value, present := lookup("SPRITES_API_KEY")
		if present {
			c.SpritesAPIKey = value
		}
	}

	{
		value, present := lookup("PLUGIN_MAX_RESTART_ATTEMPTS")
		if !present {
			value = "3"
		}

		parsed, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return c, fmt.Errorf("%w: PLUGIN_MAX_RESTART_ATTEMPTS has an invalid int value", errInvalidConfiguration)
		}

		c.PluginMaxRestartAttempts = int(parsed)
	}

	{
		value, present := lookup("PLUGIN_MAX_BACKOFF_DURATION")
		if !present {
			value = "60s"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: PLUGIN_MAX_BACKOFF_DURATION has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.PluginMaxBackoffDuration = parsed
	}

	{
		value, present := lookup("PLUGIN_RUNTIME_CRASH_THRESHOLD")
		if !present {
			value = "30s"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: PLUGIN_RUNTIME_CRASH_THRESHOLD has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.PluginRuntimeCrashThreshold = parsed
	}

	{
		value, present := lookup("PLUGIN_RUNTIME_CRASH_BACKOFF")
		if !present {
			value = "5s"
		}

		parsed, err := time.ParseDuration(value)
		if err != nil {
			return c, fmt.Errorf("%w: PLUGIN_RUNTIME_CRASH_BACKOFF has an invalid time.Duration value", errInvalidConfiguration)
		}

		c.PluginRuntimeCrashBackoff = parsed
	}

	{
		value, present := lookup("MCP_ENABLED")
		if !present {
			value = "false"
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("%w: MCP_ENABLED has an invalid bool value", errInvalidConfiguration)
		}

		c.MCPEnabled = parsed
	}

	{
		value, present := lookup("LANGUAGE_MODEL_PROVIDER")
		if present {
			c.LanguageModelProvider = value
		}
	}

	{
		value, present := lookup("OPENAI_API_KEY")
		if present {
			c.OpenAIKey = value
		}
	}

	{
		value, present := lookup("ROBOT_MODEL_DEBUG_DUMP")
		if !present {
			value = "false"
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("%w: ROBOT_MODEL_DEBUG_DUMP has an invalid bool value", errInvalidConfiguration)
		}

		c.RobotModelDebugDump = parsed
	}

	{
		value, present := lookup("SEMDEX_PROVIDER")
		if !present {
			value = ""
		}

		c.SemdexProvider = value
	}

	{
		value, present := lookup("SEMDEX_LOCAL_PATH")
		if !present {
			value = "data/semdex"
		}

		c.SemdexLocalPath = value
	}

	{
		value, present := lookup("PINECONE_API_KEY")
		if present {
			c.PineconeAPIKey = value
		}
	}

	{
		value, present := lookup("PINECONE_INDEX")
		if present {
			c.PineconeIndex = value
		}
	}

	{
		value, present := lookup("PINECONE_DIMENSIONS")
		if present {
			parsed, err := strconv.ParseInt(value, 0, 32)
			if err != nil {
				return c, fmt.Errorf("%w: PINECONE_DIMENSIONS has an invalid int32 value", errInvalidConfiguration)
			}

			c.PineconeDimensions = int32(parsed)
		}
	}

	{
		value, present := lookup("PINECONE_CLOUD")
		if present {
			c.PineconeCloud = value
		}
	}

	{
		value, present := lookup("PINECONE_REGION")
		if present {
			c.PineconeRegion = value
		}
	}

	return ApplyDefaults(c)
}

func ApplyDefaults(c Config) (Config, error) {
	if c.ImageWorkers == 0 {
		c.ImageWorkers = DefaultImageWorkers
	}

	if c.ImageMemoryBudget == 0 {
		c.ImageMemoryBudget = DefaultImageMemoryBudget
	}

	return c, c.Validate()
}

func (c Config) Validate() error {
	if math.IsNaN(c.DevChaosFailRate) || math.IsInf(c.DevChaosFailRate, 0) {
		return fmt.Errorf("%w: DEV_CHAOS_FAIL_RATE must be finite", errInvalidConfiguration)
	}

	if c.DevChaosFailRate < 0.0 {
		return fmt.Errorf("%w: DEV_CHAOS_FAIL_RATE must be at least 0", errInvalidConfiguration)
	}

	if c.DevChaosFailRate > 1.0 {
		return fmt.Errorf("%w: DEV_CHAOS_FAIL_RATE must be at most 1", errInvalidConfiguration)
	}

	if c.ImageWorkers < 0.0 {
		return fmt.Errorf("%w: IMAGE_WORKERS must be at least 0", errInvalidConfiguration)
	}

	if c.ImageMemoryBudget < 0.0 {
		return fmt.Errorf("%w: IMAGE_MEMORY_BUDGET must be at least 0", errInvalidConfiguration)
	}

	return nil
}
