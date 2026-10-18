# Pubsub session context boundary

Status: deferred

## Context

`internal/infrastructure/pubsub/middleware.go` imports the account and RBAC
resources and the authentication session service to propagate session context.
This crosses the repository's infrastructure boundary and couples the bus to
OAuth, avatars and Assets through the session validator.

Replacing `job_dispatcher.Dispatcher` with `*background_job.Runner` in the Asset
uploader and downloader exposes this import cycle:

```text
background_job → pubsub → session → oauth → avatar → asset_download → background_job
```

## Current decision

Keep `job_dispatcher.Dispatcher` and its Fx adapter until the pubsub/session
boundary is addressed separately. The interface currently breaks this cycle;
it is not required for multiple production implementations.

## Deferred investigation

- Separate application session serialization from the infrastructure bus using
  the existing repository conventions.
- Preserve account identity, guest permissions, access-key security schemes,
  malformed-metadata rejection and message cancellation behavior.
- Remove the Dispatcher interface and inject `*background_job.Runner` directly
  once the dependency cycle is resolved.
- Verify real Fx startup, session propagation, immediate dispatch and recovery.
