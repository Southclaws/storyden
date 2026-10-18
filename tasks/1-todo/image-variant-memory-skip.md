# Image variant memory skips

During image delivery profiling, a source whose only pending preset was too
expensive for the memory budget became a terminal skipped background job. Later
requests for an admissible smaller preset merged intent but could not restart
the terminal job. The download resolver continued returning the original with
`no-cache`, causing repeated full-size downloads.

Handle admission failures per preset, preserving the ability to request smaller
variants. Terminal sources should not continue presenting a pending fallback.
Existing skipped jobs are not automatically retried when the operator raises
`IMAGE_MEMORY_BUDGET`.

Relevant code: `app/services/asset/asset_image/service.go`,
`app/resources/background_job/payload.go`, and
`app/services/asset/asset_download/downloader.go`.
