# Concurrent settings writes

Status: deferred

## Context

`SettingsRepository.Set` serializes writes within one repository instance, but
separate replicas can read the same settings snapshot and then overwrite each
other's changes when saving the entire JSON value.

This surfaced while adding conditional branding-reference publication for the
legacy Asset migration. The general settings write path is outside that task.

## Deferred investigation

- Reproduce writes from independent repository instances against the same row.
- Preserve concurrent changes to different settings fields.
- Keep cache refreshes ordered with successful writes.
- Test branding uploads racing with other settings updates and migration.
