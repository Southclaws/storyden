module github.com/Southclaws/storyden

go 1.25.3

tool (
	entgo.io/ent/cmd/ent
	github.com/Southclaws/enumerator
	github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
)

require (
	ariga.io/atlas v0.37.0
	dario.cat/mergo v1.0.2
	entgo.io/ent v0.14.5
	github.com/Southclaws/dt v1.0.1
	github.com/alexedwards/argon2id v1.0.0
	github.com/forPelevin/gomoji v1.4.0
	github.com/getkin/kin-openapi v0.133.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.13.4
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.6.0
	github.com/samber/lo v1.52.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/dig v1.19.0
	go.uber.org/fx v1.24.0
	golang.org/x/oauth2 v0.32.0
)

require github.com/Southclaws/fault v0.8.2

require (
	github.com/JohannesKaufmann/html-to-markdown/v2 v2.4.0
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/Southclaws/lexorank v1.2.3
	github.com/Southclaws/opt v0.6.1
	github.com/Southclaws/swirl v1.0.1
	github.com/ThreeDotsLabs/watermill v1.5.1
	github.com/ThreeDotsLabs/watermill-amqp/v3 v3.0.2
	github.com/alitto/pond/v2 v2.5.0
	github.com/blevesearch/bleve/v2 v2.5.5
	github.com/bwmarrin/discordgo v0.29.0
	github.com/cixtor/readability v1.0.0
	github.com/coder/websocket v1.8.12
	github.com/coreos/go-oidc/v3 v3.16.0
	github.com/dave/jennifer v1.7.1
	github.com/dboslee/lru v0.0.1
	github.com/dgraph-io/ristretto/v2 v2.3.0
	github.com/disintegration/imaging v1.6.2
	github.com/dustinkirkland/golang-petname v0.0.0-20240428194347-eebcea082ee0
	github.com/gabriel-vasile/mimetype v1.4.10
	github.com/getsentry/sentry-go v0.35.3
	github.com/getsentry/sentry-go/otel v0.35.3
	github.com/glebarez/go-sqlite v1.22.0
	github.com/go-viper/mapstructure/v2 v2.5.0
	github.com/goccy/go-yaml v1.19.2
	github.com/golang-cz/devslog v0.0.15
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/go-github/v75 v75.0.0
	github.com/google/jsonschema-go v0.4.2
	github.com/iancoleman/strcase v0.3.0
	github.com/invopop/jsonschema v0.13.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/jmoiron/sqlx v1.4.0
	github.com/klippa-app/go-pdfium v1.17.2
	github.com/mark3labs/mcp-go v0.41.1
	github.com/matcornic/hermes/v2 v2.1.0
	github.com/mazznoer/colorgrad v0.10.0
	github.com/mazznoer/csscolorparser v0.1.6
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mileusna/useragent v1.3.5
	github.com/minimaxir/big-list-of-naughty-strings/naughtystrings v0.0.0-20210417190545-db33ec7b1d5d
	github.com/oapi-codegen/echo-middleware v1.0.2
	github.com/oapi-codegen/nullable v1.1.0
	github.com/oapi-codegen/runtime v1.1.2
	github.com/openai/openai-go v1.12.0
	github.com/pb33f/libopenapi v0.28.0
	github.com/philippgille/chromem-go v0.7.0
	github.com/pinecone-io/go-pinecone/v4 v4.1.4
	github.com/puzpuzpuz/xsync/v4 v4.2.0
	github.com/redis/rueidis v1.0.66
	github.com/rs/cors v1.11.1
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/sendgrid/sendgrid-go v3.16.1+incompatible
	github.com/shirou/gopsutil/v4 v4.25.9
	github.com/tursodatabase/libsql-client-go v0.0.0-20240902231107-85af5b9d094d
	github.com/twilio/twilio-go v1.28.3
	github.com/weaviate/weaviate v1.33.0
	github.com/weaviate/weaviate-go-client/v5 v5.5.0
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.38.0
	go.opentelemetry.io/otel/sdk v1.38.0
	go.opentelemetry.io/otel/trace v1.38.0
	golang.org/x/exp v0.0.0-20251219203646-944ab1f22d93
	golang.org/x/sync v0.19.0
	google.golang.org/api v0.252.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	cloud.google.com/go/auth v0.17.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	github.com/JohannesKaufmann/dom v0.2.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/RoaringBitmap/roaring/v2 v2.4.5 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/bits-and-blooms/bitset v1.22.0 // indirect
	github.com/blevesearch/bleve_index_api v1.2.11 // indirect
	github.com/blevesearch/geo v0.2.4 // indirect
	github.com/blevesearch/go-faiss v1.0.26 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/gtreap v0.1.1 // indirect
	github.com/blevesearch/mmap-go v1.0.4 // indirect
	github.com/blevesearch/scorch_segment_api/v2 v2.3.13 // indirect
	github.com/blevesearch/segment v0.9.1 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/upsidedown_store_api v1.0.2 // indirect
	github.com/blevesearch/vellum v1.1.0 // indirect
	github.com/blevesearch/zapx/v11 v11.4.2 // indirect
	github.com/blevesearch/zapx/v12 v12.4.2 // indirect
	github.com/blevesearch/zapx/v13 v13.4.2 // indirect
	github.com/blevesearch/zapx/v14 v14.4.2 // indirect
	github.com/blevesearch/zapx/v15 v15.4.2 // indirect
	github.com/blevesearch/zapx/v16 v16.2.7 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dprotaso/go-yit v0.0.0-20250909171706-0a81c39169bc // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/runtime v0.28.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-openapi/swag/cmdutils v0.25.1 // indirect
	github.com/go-openapi/swag/conv v0.25.1 // indirect
	github.com/go-openapi/swag/fileutils v0.25.1 // indirect
	github.com/go-openapi/swag/jsonname v0.25.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.1 // indirect
	github.com/go-openapi/swag/loading v0.25.1 // indirect
	github.com/go-openapi/swag/mangling v0.25.1 // indirect
	github.com/go-openapi/swag/netutils v0.25.1 // indirect
	github.com/go-openapi/swag/stringutils v0.25.1 // indirect
	github.com/go-openapi/swag/typeutils v0.25.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.1 // indirect
	github.com/go-openapi/validate v0.24.0 // indirect
	github.com/go-webauthn/x v0.1.25 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/go-tpm v0.9.5 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/jolestar/go-commons-pool/v2 v2.1.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lithammer/shortuuid/v3 v3.0.7 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/crc64nvme v1.0.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/oapi-codegen/oapi-codegen/v2 v2.5.0 // indirect
	github.com/oasdiff/yaml v0.0.0-20250309154309-f31be36b4037 // indirect
	github.com/oasdiff/yaml3 v0.0.0-20250309153720-d2182401db90 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pb33f/jsonpath v0.1.2 // indirect
	github.com/pb33f/ordered-map/v2 v2.3.0 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/sony/gobreaker v1.0.0 // indirect
	github.com/speakeasy-api/jsonpath v0.6.2 // indirect
	github.com/speakeasy-api/openapi-overlay v0.10.3 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/vanng822/css v1.0.1 // indirect
	github.com/vanng822/go-premailer v1.24.0 // indirect
	github.com/vmware-labs/yaml-jsonpath v0.3.2 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.9-0.20250401010720-46d686821e33 // indirect
	github.com/woodsbury/decimal128 v1.4.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	go.etcd.io/bbolt v1.4.3 // indirect
	go.mongodb.org/mongo-driver v1.17.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	go.yaml.in/yaml/v4 v4.0.0-rc.2 // indirect
	golang.org/x/image v0.26.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251002232023-7c0ddcbb5797 // indirect
	google.golang.org/grpc v1.75.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	modernc.org/libc v1.65.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.10.0 // indirect
	modernc.org/sqlite v1.37.0 // indirect
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.95
)

require (
	github.com/Southclaws/enumerator v1.4.1
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-openapi/inflect v0.21.3 // indirect
	github.com/go-openapi/jsonpointer v0.22.1 // indirect
	github.com/go-openapi/swag v0.25.1 // indirect
	github.com/go-webauthn/webauthn v0.14.0
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mailru/easyjson v0.9.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/zclconf/go-cty v1.17.0 // indirect
	go.uber.org/multierr v1.11.0
	golang.org/x/crypto v0.46.0
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.48.0
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0
	golang.org/x/time v0.13.0 // indirect
	google.golang.org/protobuf v1.36.10
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
