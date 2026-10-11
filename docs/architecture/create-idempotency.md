# Idempotent REST creates

`POST /api/threads`, `POST /api/threads/{thread_mark}/replies`, and
`POST /api/nodes` accept an optional `Idempotency-Key` header. The key must be
1–128 printable ASCII characters without spaces. A client should generate a
fresh random key for each intended create and keep it across network retries.

Claims are stored in the application database, so requests handled by separate
backend replicas share a claim. The scope is the authenticated account, the
operation, and a SHA-256 digest of the key. The key itself is not stored. The
fingerprint covers the parsed request body and, for replies, the thread mark.
Identical requests on a completed claim receive the original `200` JSON body.
The current request must still pass authentication and the applicable create
permission, and the created resource must still be visible to the caller.

Conflicting claims return HTTP `409` with the following
`metadata.idempotency_reason` values:

| Reason | Meaning | Client action |
| --- | --- | --- |
| `key_mismatch` | The key was used with different input. | Use a new key only for a new intended create. |
| `in_progress` | Another request has claimed the key and has not finished. | Retry the same request later. |
| `uncertain` | The earlier request failed after claiming or was interrupted, and may have written data. | Read the collection to reconcile before deciding whether to create with a new key. |

Only successful JSON responses are replayed. An error after a claim is never
automatically re-executed, because the create services can write a resource
before later moderation, property mutation, cache, or event steps fail. The
receipt therefore provides duplicate suppression across replicas, not atomic
exactly-once execution of those services. Successful receipts remain replayable
for 24 hours; after expiry the same key may create another resource. Unfinished
receipts remain uncertain after 24 hours and are not taken over automatically.
Abandoned unfinished receipts are retained for up to 30 days before cleanup.
Cleanup is opportunistic during API traffic, so storage collection may occur
later on an idle instance.

The request fingerprint input is capped at 1 MiB and the stored successful
response at 2 MiB. Only the response body is stored; cookies and authorization
headers are never stored or replayed. Calls without `Idempotency-Key` follow
the existing create behavior.

`TrailRunNow` and `RobotSessionCreate` do not accept this contract. Trail runs
start independent actions, and Robot session creation streams execution. Their
current service boundaries cannot safely replay a completed result or reconcile
an interrupted operation from a receipt alone.
