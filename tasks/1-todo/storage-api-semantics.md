# Consistent storage API semantics

Status: deferred

## Context

`internal/infrastructure/object.Storer` exposes one interface for local files and
S3-compatible object storage, but callers still encounter provider differences.

The legacy avatar import exposed several differences in `List`:

- Local listing reads a directory; S3 listing matches a key prefix.
- `avatar` and `avatar/` can produce different results in S3. The local adapter
  now trims the trailing slash because `os.DirFS` rejects it.
- Local listing returns entry names; S3 listing returns full keys.
- Local listing includes directory entries; non-recursive S3 listing can return
  grouped prefixes. Recursive listing changes which objects are returned.

The migration currently handles these differences. Define a clearer shared
contract so future consumers do not need provider-specific knowledge.

## Deferred investigation

- Review existing consumers before choosing the API and changing behavior.
- Define object-key and path normalization, including trailing slashes.
- Make directory listing versus prefix matching, recursion, and returned entry
  types explicit. Decide whether results contain names or full storage keys.
- Standardize missing-prefix and missing-object behavior, while distinguishing
  absence from permissions, network, and other storage failures.
- Add shared contract tests for the local and S3 adapters, including nested keys,
  similarly named prefixes, empty results, and trailing slashes.
- Validate S3 changes against a real S3-compatible test instance and Cloudflare
  R2 where applicable, rather than relying only on mocks or local-store tests.

Keep the current infrastructure behavior until the replacement contract and
provider tests are ready.
