package plugin

// NOTE: Fairly arbitrary values here, sane-ish defaults.
// TODO: Make these configurable somehow.

const (
	// max compressed plugin archive size accepted from uploads or URL fetches
	MaxArchiveSizeBytes int64 = 50 * 1024 * 1024 // 50 MiB
	// max manifest.json size read from a plugin archive
	MaxManifestSizeBytes int64 = 1 * 1024 * 1024 // 1 MiB

	// extraction guardrails for supervised runtime archive unpacking
	MaxArchiveExtractFileBytes  uint64 = 50 * 1024 * 1024  // 50 MiB per file
	MaxArchiveExtractTotalBytes uint64 = 200 * 1024 * 1024 // 200 MiB total
	MaxArchiveExtractFileCount  int    = 2048
)
