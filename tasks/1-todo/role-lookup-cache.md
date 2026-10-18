# Role lookup cache

Status: deferred

## Context

`account_repo.GetRefByID` refreshes role membership and role definitions even
when the account itself comes from cache. `role_repo.GetMany` consults the shared
role-list cache, but its SQL fallback does not populate that cache. Repeated
account-cache hits can therefore continue issuing role queries while the
role-list cache is empty.

This surfaced during avatar query-count testing in the real Fx integration
harness: reloading an account issued both an `AccountQuery` and a `RoleQuery`.
Cached account reads must not be described as free of database queries without
accounting for role hydration.

## Current decision

Keep this separate from the avatar Asset work. The proposed change to make
`GetMany` use the cached `List` path was reverted.

## Deferred investigation

- Measure cold and warm query counts for account reads and avatar requests.
- Compare populating the shared role-list cache with caching selected roles.
- Preserve current role assignment, update, deletion, default-role and permission
  behaviour, including cache failures and invalidation.
- Add query-count regression coverage using the existing integration harness.
