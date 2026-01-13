// Package config contains all environment variable based configuration.
// THIS FILE IS GENERATED. DO NOT EDIT MANUALLY.
// To edit configuration variables, edit the config.yaml file and run codegen.
package config

import (
	"log/slog"
	"net/url"
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
	LogLevel slog.Level `default:"info" envconfig:"LOG_LEVEL"`
	/*
	   Can be set to either:

	   - `(not set)` (default) somewhat human readable "logfmt" format logs for simple setups
	   - `dev` for developer-friendly logs, with colours and attributes on separate lines for readability
	   - `json` for machine-readable logs, mainly for log aggregators, etc.
	*/
	LogFormat string `envconfig:"LOG_FORMAT"`
	/*
	   Determines whether or not the backend service will also start the frontend Node.js process. When empty, it will not

	   When a path is provided, Storyden will execute `node <path>` to start the frontend process. This is used by the fullstack Docker image to start the frontend process in the same container as the backend.
	*/
	RunFrontend string `default:"" envconfig:"RUN_FRONTEND"`
	/*
	   Used in conjunction with `RUN_FRONTEND`. This is the address that the frontend will be available at. This is used by the fullstack Docker image to proxy requests that don't match any `/api` or other routes to the frontend process.

	   In the default `fullstack` image, this is set to `http://localhost:3000` which is the default port for the Next.js process.
	*/
	FrontendProxy url.URL `default:"" envconfig:"PROXY_FRONTEND_ADDRESS"`

	// -
	// Development tools
	// -

	/*
	   Simulates slow requests.

	   This will add a random delay between zero and this value to all requests. This is useful for testing how the client handles slow responses.
	*/
	DevChaosSlowMode time.Duration `envconfig:"DEV_CHAOS_SLOW_MODE"`
	/*
	   Simulates slow message delivery in the internal message queue.

	   This will add a random delay between zero and this value to all messages in the internal message queue. This is useful for testing how the client handles delayed message processing.
	*/
	DevChaosSlowModeQueue time.Duration `envconfig:"DEV_CHAOS_SLOW_MODE_QUEUE"`
	/*
	   A value between 0 and 1 which simulates failed requests.

	   This will add a random failure to all requests. This is useful for testing how the client handles "internal server error" responses.
	*/
	DevChaosFailRate float64 `envconfig:"DEV_CHAOS_FAIL_RATE"`

	// -
	// Core configuration
	// -

	/*
	   The database URL to connect to. This can be a SQLite, PostgreSQL, or CockroachDB URL.

	   The accepted schemes for this URL are:
	   - `sqlite://` or `sqlite3://` for SQLite or Litestream.
	   - `postgres://` or `postgresql://` for PostgreSQL, CockroachDB and any other PostgreSQL-compatible database
	   - `libsql://` for Turso remote SQLite. **Note:** This is currently experimental, only remote Turso databases are supported.
	*/
	DatabaseURL string `default:"sqlite://data/data.db?_pragma=foreign_keys(1)" envconfig:"DATABASE_URL"`
	/*
	   The interface on which the API service will for HTTP requests.

	   Typically, in a containerised environment, this should be all interfaces (`0.0.0.0`.)
	*/
	ListenAddr string `default:"0.0.0.0:8000" envconfig:"LISTEN_ADDR"`
	/*
	   The address at which the web frontend will be hosted.

	   This must be set to the public URL that users of the instance will access the frontend client. It is used to determine things such as cookie domain attributes, CORS policy, WebAuthn attributes and other necessary settings. The scheme may be used by some internal components to determine whether the instance is running in a secure context or not.

	   This is by default `http://localhost:3000` when running locally, or when deploying to production, `https://<your-domain>`.
	*/
	PublicWebAddress url.URL `default:"http://localhost:3000" envconfig:"PUBLIC_WEB_ADDRESS"`
	/*
	   The address at which the public API will be accessible.

	   This is also used for things such as cookies, CORS, etc.

	   Please note that both the public API address and public web address must share the same root domain name as Storyden cookies are configured to be issued under this assumption. It also makes a lot of cross-origin and cookie configurations easier to make secure.
	*/
	PublicAPIAddress url.URL `default:"http://localhost:8000" envconfig:"PUBLIC_API_ADDRESS"`

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
	RateLimit int `default:"5000" envconfig:"RATE_LIMIT"`
	/*
	   Sliding window duration used to enforce `RATE_LIMIT`.

	   On each request, Storyden considers the total number of units consumed in the last
	   `RATE_LIMIT_PERIOD` (not aligned to the hour/minute boundary).
	*/
	RateLimitPeriod time.Duration `default:"1h" envconfig:"RATE_LIMIT_PERIOD"`
	/*
	   Bucket size (granularity) used to approximate the sliding window.

	   Requests are counted into discrete time buckets of this size (e.g. 1 minute). When
	   checking the limit, Storyden sums all buckets whose timestamps fall within the last
	   `RATE_LIMIT_PERIOD` and discards older buckets.

	   Smaller buckets = more accurate sliding-window behaviour but more storage/CPU overhead.
	   Larger buckets = cheaper but "chunkier" enforcement.

	   Rule of thumb: set this to ~1/60 of `RATE_LIMIT_PERIOD` (e.g. 1m buckets for a 1h period).
	*/
	RateLimitBucket time.Duration `default:"1m" envconfig:"RATE_LIMIT_BUCKET"`
	/*
	   Cost multiplier applied to unauthenticated (guest) requests.

	   Example: if set to 5, each guest request consumes 5 units from the same `RATE_LIMIT`
	   budget, effectively allowing guests ~1/5th the throughput of authenticated users.
	*/
	RateLimitGuestCost int `default:"1" envconfig:"RATE_LIMIT_GUEST_COST"`

	// -
	// Telemetry and monitoring
	// -

	/*
	   Either:
	   - `otlp` for any standard OpenTelemetry collector.
	   - `sentry` for Sentry (which is OpenTelemetry-compatible, however requires its own specific configuration.)
	   - `logger` for local logging to the console. This is only really useful for Storyden developers and is very noisy.
	*/
	OTELProvider string `default:"" envconfig:"OTEL_PROVIDER"`
	// The collector endpoint for sending OTEL data.
	OTELEndpoint url.URL `default:"" envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	// When `OTEL_PROVIDER` is set to `sentry`, this is the DSN for the Sentry project.
	SentryDSN string `default:"" envconfig:"SENTRY_DSN"`

	// -
	// Email
	// -

	/*
	   Either:

	   - unset (default) for no email sending. Email sending is not a requirement for a production deployment.
	   - `sendgrid` for SendGrid based email sending.
	   - `mock` for logging emails to the console. Only useful for Storyden developers and testing.
	*/
	EmailProvider string `envconfig:"EMAIL_PROVIDER"`
	/*
	   The name that will be used as the sender name for emails sent via SendGrid.

	   This is typically the name of your community or organisation.
	*/
	SendGridFromName string `envconfig:"SENDGRID_FROM_NAME"`
	/*
	   The email address that will be used as the sender address for emails sent via SendGrid.

	   This is typically a no-reply address, such as `no-reply@<your-domain>`.
	*/
	SendGridFromAddress string `envconfig:"SENDGRID_FROM_ADDRESS"`
	/*
	   The API key for the SendGrid account. This is required for sending emails via SendGrid.

	   This is typically a long string of characters that you can generate in the SendGrid dashboard.
	*/
	SendGridAPIKey string `envconfig:"SENDGRID_API_KEY"`

	// -
	// Authentication
	// -

	/*
	   The secret key used to sign JWT tokens. This is used for authentication and should be kept secret.

	   This is typically a long string of characters that you can generate using a secure random generator such as `openssl rand -hex 12`.

	   The JWT secret is required if you enable any of the OAuth providers or enable email features. This is because JWTs are used to verify callbacks as well as verify password-reset and other tokens.
	*/
	JWTSecret []byte `envconfig:"JWT_SECRET"`
	// Enable Google SSO authentication.
	GoogleEnabled bool `envconfig:"OAUTH_GOOGLE_ENABLED"`
	// The client ID for the Google OAuth2 application.
	GoogleClientID string `envconfig:"OAUTH_GOOGLE_CLIENT_ID"`
	// The client secret for the Google OAuth2 application.
	GoogleClientSecret string `envconfig:"OAUTH_GOOGLE_CLIENT_SECRET"`
	// Enable GitHub SSO authentication.
	GitHubEnabled bool `envconfig:"OAUTH_GITHUB_ENABLED"`
	// The client ID for the GitHub OAuth2 application.
	GitHubClientID string `envconfig:"OAUTH_GITHUB_CLIENT_ID"`
	// The client secret for the GitHub OAuth2 application.
	GitHubClientSecret string `envconfig:"OAUTH_GITHUB_CLIENT_SECRET"`
	// Enable Discord SSO authentication.
	DiscordEnabled bool `envconfig:"OAUTH_DISCORD_ENABLED"`
	// The client ID for the Discord OAuth2 application.
	DiscordClientID string `envconfig:"OAUTH_DISCORD_CLIENT_ID"`
	// The client secret for the Discord OAuth2 application.
	DiscordClientSecret string `envconfig:"OAUTH_DISCORD_CLIENT_SECRET"`
	// Enable Keycloak OIDC authentication.
	KeycloakEnabled bool `envconfig:"OAUTH_KEYCLOAK_ENABLED"`
	// The client ID for the Keycloak OAuth2 application.
	KeycloakClientID string `envconfig:"OAUTH_KEYCLOAK_CLIENT_ID"`
	// The client secret for the Keycloak OAuth2 application.
	KeycloakClientSecret string `envconfig:"OAUTH_KEYCLOAK_CLIENT_SECRET"`
	// The issuer/discovery URL for the Keycloak realm (e.g. https://auth.example.com/realms/YourRealm).
	KeycloakIssuerURL url.URL `envconfig:"OAUTH_KEYCLOAK_ISSUER_URL"`

	// -
	// SMS
	// -

	/*
	   Either:

	   - unset (default) for no SMS sending. SMS sending is not a requirement for a production deployment.
	   - `twilio` for Twilio based SMS sending.
	   - `mock` for logging SMS to the console. Only useful for Storyden developers and testing.
	*/
	SMSProvider string `envconfig:"SMS_PROVIDER"`
	/*
	   The account SID for the Twilio account.

	   This is typically a long string of characters that you can view in the Twilio dashboard.
	*/
	TwilioAccountSID string `envconfig:"TWILIO_ACCOUNT_SID"`
	// The phone number that will be used as the sender number for SMS sent via Twilio.
	TwilioPhoneNumber string `envconfig:"TWILIO_PHONE_NUMBER"`
	/*
	   The auth token for the Twilio account. This is required for sending SMS via Twilio.

	   This is typically a long string of characters that you can generate in the Twilio dashboard.
	*/
	TwilioAuthToken string `envconfig:"TWILIO_AUTH_TOKEN"`

	// -
	// Assets/file storage
	// -

	/*
	   Either:

	   - `local` for local file storage.
	   - `s3` for any Amazon S3-compatible storage, such as S3 itself (obviously...), Google Cloud Storage, Cloudflare R2, Minio, etc.
	*/
	AssetStorageType string `envconfig:"ASSET_STORAGE_TYPE"`
	// When `ASSET_STORAGE_TYPE` is set to `local`, this is the path to the directory where files will be stored.
	AssetStorageLocalPath string `envconfig:"ASSET_STORAGE_LOCAL_PATH"`
	// When `ASSET_STORAGE_TYPE` is set to `s3`, this determines whether or not to use HTTPS for the S3 connection. You should always set this to `true` unless your S3-compatible storage provider is internally but not publicly accessible, such as in a Kubernetes cluster or running on the same host.
	S3Secure bool `default:"true" envconfig:"S3_SECURE"`
	// The endpoint for the S3-compatible storage provider. This is typically the base URL of the provider, such as `https://s3.amazonaws.com` for AWS S3, or `https://storage.googleapis.com` for Google Cloud Storage, etc.
	S3Endpoint string `envconfig:"S3_ENDPOINT"`
	// The bucket name for Storyden assets to be stored in.
	S3Bucket string `envconfig:"S3_BUCKET"`
	/*
	   Most S3-compatible storage providers require a region to be specified. This is typically the region in which the bucket is located, such as `us-east-1` for AWS S3.

	   However, some providers do not use regions but S3-compatible clients still require this to be set. In most cases, the provider will give you a value for this, such as `auto` when using Cloudflare R2.
	*/
	S3Region string `envconfig:"S3_REGION"`
	// The access key for the S3-compatible storage provider.
	S3AccessKey string `envconfig:"S3_ACCESS_KEY"`
	// The secret key for the S3-compatible storage provider.
	S3SecretKey string `envconfig:"S3_SECRET_KEY"`

	// -
	// Cache
	// -

	/*
	   When empty, caching will use an efficient in-memory store. This is usually fine for small to medium-sized deployments however it's worth keeping an eye on your deployment's machine memory usage.

	   When set to `redis`, Storyden will use Redis as a cache provider. This is recommended for larger deployments that receive a lot of traffic. The cache provider is also used for the rate limiter so that it can be shared across multiple instances of Storyden.

	   This is necessary for deploying replica instances of Storyden that are backed by the same persistence layers (database, asset storage, etc.)
	*/
	CacheProvider string `default:"" envconfig:"CACHE_PROVIDER"`
	/*
	   The Redis URL to connect to.

	   This is a full URL with `redis://` as the scheme. You can set the username and password using the URL format, for example: `redis://<username>:<password>@<host>:<port>`.
	*/
	RedisURL url.URL `default:"" envconfig:"REDIS_URL"`

	// -
	// Search features
	// -

	/*
	   Either:

	   - `database` for the default database-driven search.  This is not recommended for larger deployments as it does not scale well and has limited search quality.
	   - `bleve` for Bleve. This is a local full-text search engine that is fast and efficient for small to medium-sized deployments. This is best used when you are using local disk storage, such as SQLite and local asset storage.
	   - `redis` for Redisearch. This is a fast and efficient search provider that is recommended for larger deployments. This is recommended if your environment is ephemeral and you're already using external providers for database and asset storage.
	*/
	SearchProvider string `default:"database" envconfig:"SEARCH_PROVIDER"`
	/*
	   When using `SEARCH_PROVIDER` set to either `bleve` or `redis`, this is the number of items that will be indexed in a single batch.

	   Increasing this value will improve indexing performance, but will also increase memory usage during indexing.
	*/
	SearchIndexChunkSize int `default:"1000" envconfig:"SEARCH_INDEX_CHUNK_SIZE"`
	// The path to the directory where Bleve will store search indexes. Only used when `SEARCH_PROVIDER` is set to `bleve`.
	BlevePath string `default:"data/bleve" envconfig:"BLEVE_PATH"`
	// The name of the Redis search index. Only used when `SEARCH_PROVIDER` is set to `redis`.
	RedisSearchIndexName string `default:"storyden" envconfig:"REDIS_SEARCH_INDEX_NAME"`

	// -
	// Message queue
	// -

	/*
	   Either:

	   - Default (no value): in-memory Go channels. This is fast and efficient, but not persistent across restarts and will add a bit of memory usage to the process.
	   - `amqp`: RabbitMQ. This is a persistent message queue that is fast and reliable. It is recommended for larger deployments and is necessary for deploying replica instances of Storyden.
	*/
	QueueType string `default:"internal" envconfig:"QUEUE_TYPE"`
	/*
	   The RabbitMQ URL to connect to.

	   This is a full URL with `amqp://` as the scheme. You can set the username and password using the URL format, for example: `amqp://<username>:<password>@<host>:<port>`.

	   The default value is `amqp://guest:guest@localhost:5672/` which is the default RabbitMQ URL.

	   Storyden does not currently support `amqps://` (secure) URLs, but this will be added soon.
	*/
	AmqpURL string `default:"amqp://guest:guest@localhost:5672/" envconfig:"AMQP_URL"`
	/*
	   The maximum number of times a failed message will be retried before being moved to the dead letter queue.

	   Messages are retried with exponential backoff starting at 1 second and doubling each time up to a maximum of 1 minute between retries.
	*/
	QueueMaxRetries int `default:"5" envconfig:"QUEUE_MAX_RETRIES"`
	// The initial interval to wait before the first retry attempt.
	QueueRetryInitialInterval time.Duration `default:"1s" envconfig:"QUEUE_RETRY_INITIAL_INTERVAL"`
	// The maximum interval to wait between retry attempts. The exponential backoff will not exceed this value.
	QueueRetryMaxInterval time.Duration `default:"1m" envconfig:"QUEUE_RETRY_MAX_INTERVAL"`

	// -
	// Artificial intelligence/language models
	// -

	/*
	   Enables the Model Context Provider server, accessible via SSE at `/mcp`.

	   This is used to integrate Storyden into agentic workflow engines and other language model tooling.

	   See [the documentation](https://storyden.org/docs/introduction/mcp) for more information.
	*/
	MCPEnabled bool `default:"false" envconfig:"MCP_ENABLED"`
	/*
	   The provider for language model features.

	   `openai` is currently the only supported provider.
	*/
	LanguageModelProvider string `envconfig:"LANGUAGE_MODEL_PROVIDER"`
	// When `LANGUAGE_MODEL_PROVIDER` is set to `openai`, this is the API key for the OpenAI API.
	OpenAIKey string `envconfig:"OPENAI_API_KEY"`
	/*
	   The Asker feature provides a conversational interface for exploring the community's content across library pages, threads, links, profiles, etc. It is separate from the language model provider as some providers support different features.

	   This can be set to either:

	   - `openai` for OpenAI
	   - `perplexity` for Perplexity AI - note that Perplexity does not currently support all the features necessary to be a `LANGUAGE_MODEL_PROVIDER` so it is only available as an `ASKER_PROVIDER`.
	*/
	AskerProvider string `default:"" envconfig:"ASKER_PROVIDER"`
	// If `ASKER_PROVIDER` is set to `perplexity`, this is the API key for the Perplexity API.
	PerplexityAPIKey string `envconfig:"PERPLEXITY_API_KEY"`

	// -
	// Semdex
	// -

	/*
	   Either:
	   - `chromem` for an experimental local vector database. This is not recommended for use in large deployments as it's rather slow and memory-hungry.
	   - `weaviate` for Weaviate, a self-hostable or managed vector database.
	   - `pinecone` for Pinecone, a fully managed vector database.
	*/
	SemdexProvider string `default:"" envconfig:"SEMDEX_PROVIDER"`

	// -
	// Local Semdex
	// -

	// The path to the directory where Chromem will store vector indexes.
	SemdexLocalPath string `default:"data/semdex" envconfig:"SEMDEX_LOCAL_PATH"`

	// -
	// Weaviate Semdex
	// -

	// The Weaviate API URL. This can be set to a self-hosted instance of Weaviate or the Weaviate Cloud API.
	WeaviateURL string `envconfig:"WEAVIATE_URL"`
	// For self-hosted Weaviate where authentication is enabled, or when using Weaviate Cloud.
	WeaviateToken string `envconfig:"WEAVIATE_API_TOKEN"`
	/*
	   The class name for Weaviate. This value actually controls which model is used for embeddings and other configuration. In future, there will be a more flexible configuration for Weaviate.

	   Value values are:

	   - `text2vec-transformers`: requires that the Weaviate instance is using the `text2vec-transformers` module. This is for calculating embeddings locally using a GPU (or very slowly, using a CPU.) This is only available when self-hosting Weaviate.
	   - `text2vec-openai`: requires that the Weaviate instance is using the `text2vec-openai` module. This uses OpenAI's API to calculate embeddings. This works on both self-hosted and Weaviate Cloud instances.
	*/
	WeaviateClassName string `envconfig:"WEAVIATE_CLASS_NAME"`

	// -
	// Pinecone Semdex
	// -

	// Your Pinecone API key. This is required for all Pinecone API requests.
	PineconeAPIKey string `envconfig:"PINECONE_API_KEY"`
	// The index name that Storyden will use in your Pinecone workspace.
	PineconeIndex string `envconfig:"PINECONE_INDEX"`
	// This value is dependent on the underlying OpenAI configuration. Currently this is static and set to 3072 dimensions. In future, Storyden will provide more flexible configuration for language model providers.
	PineconeDimensions int32 `envconfig:"PINECONE_DIMENSIONS"`
	// Pinecone provides hosting on different cloud providers, see the Pinecone documentation for more information. The cloud provider you choose will be reflected in your Pinecone dashboard.
	PineconeCloud string `envconfig:"PINECONE_CLOUD"`
	// Same as above, but for the region. As with any third party providers, it's recommended to choose the region closest to both your Storyden deployment and your community members for best performance and experience.
	PineconeRegion string `envconfig:"PINECONE_REGION"`
}
