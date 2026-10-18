# Image memory corpus

Run the portable encoder without CGO or a system libwebp:

```sh
CGO_ENABLED=0 go test -tags nodynamic ./internal/infrastructure/imagery \
  -run '^$' -bench '^BenchmarkEstimate$' -benchtime=1x -count=3
```

The benchmark refuses the dynamic encoder. CI runs it on Linux, macOS and Windows.
Each workload runs in a fresh copy of the test binary, decodes once, then writes
variants sequentially to temporary files, matching the worker's image operations.
There are four workloads: eager thumbnails and S, M alone, L alone, and all six
presets together. Measurements include opening the encoded source, decoding,
orientation, resizing, and encoding. Fixture generation happens in the parent
process and is excluded.

The committed corpus contains existing Storyden images:

| Fixture | Dimensions | Copied from |
| --- | --- | --- |
| `avatar.png` | 513 × 513 | `app/services/branding/icon/embed/default_icon.png` |
| `mountains.png` | 1200 × 600 | `app/services/branding/banner/storyden-mountains-1200x600.png` |
| `screenshot.png` | 1469 × 961 | `home/public/2026_app_screenshot.png` |
| `screenshot.webp` | 1024 × 670 | `home/public/2026_app_screenshot_1024.webp` |

Generated fixtures add landscape, portrait, square, undersized, panoramic,
transparent, EXIF-rotated, and extreme-aspect-ratio images. JPEG and transparent
PNG fixtures also exercise the 12-million-pixel source limit. Together these give
56 fixture/workload combinations. Generated patterns are repeatable, not a
substitute for photographs with different compression characteristics.

To include real uploads, put copies of eligible static JPEG, PNG or WebP files in
a directory and set `STORYDEN_IMAGE_CORPUS` to its absolute path. Filenames do not
need extensions. Every regular file in that directory is measured; invalid or
oversized sources fail the benchmark. No application database or server is used.

## Reading the results

- `estimate-B`: the largest reservation returned by `Estimate` for the workload.
- `peak-heap-B`: the largest sampled increase in `runtime.MemStats.HeapAlloc`.
- `peak-runtime-B`: the largest sampled increase in `Sys - HeapReleased`.
- `allocated-B`: total allocated bytes, including memory subsequently collected.
- `temp-B`: encoded source bytes plus the largest encoded output.
- `estimate/peak`: estimate divided by the sampled runtime memory increase.

The benchmark fails a workload whose observed runtime increase exceeds its
estimate. This makes codec or runtime changes that invalidate the calibration
visible in CI, without asserting exact allocation counts.

The child starts after `debug.FreeOSMemory`. Memory is sampled every millisecond,
with a final sample after encoding. Sampling can miss short peaks and adds some
overhead. These are observed memory increases, not exact allocation bounds or
process RSS. `temp-B` excludes the source object already in storage, filesystem
block overhead and temporary directory entries. Timing includes subprocess
startup and sampling; use the existing image benchmarks for throughput.

The benchmark inherits GC settings. Record Go version, OS, architecture, `GOGC`
and `GOMEMLIMIT` when comparing results. The measurements below used Go 1.27.2,
arm64, the default GC settings, and the Go encoder. They do not establish a hard
bound for other runtime settings, future codecs, or every accepted input.

## Initial calibration, 18 October 2026

On macOS, the three real uploads rejected during frontend profiling measured:

| Source | Dimensions | Previous estimate for L | Observed runtime increase for L | Process peak RSS |
| --- | --- | --- | --- | --- |
| Common Voice JPEG, 1.68 MB | 3072 × 1536 | 265.6 MiB | 158.8 MiB | 169.8 MiB |
| Untitled PNG, 1.23 MB | 3024 × 1964 | 308.5 MiB | 208.8 MiB | 220.2 MiB |
| Screenshot PNG, 1.23 MB | 2902 × 1620 | 269.4 MiB | 174.7 MiB | 185.7 MiB |

RSS was measured separately with macOS `/usr/bin/time -l` around an individual
child test process. It includes the executable and runtime baseline, unlike the
benchmark's memory increases. These real uploads are local data and are not
redistributed in the committed corpus.

The previous estimate was not uniformly conservative. The generated transparent
3024 × 1964 PNG reached about 371 MiB for all variants against a 306 MiB estimate.
The transparent 4000 × 3000 PNG reached about 595 MiB against a 460 MiB estimate.
The 1 × 1,000,000 PNG reached about 99 MiB for eager variants against a 64 MiB
estimate. A medium WebP workload also exceeded its estimate.

`Estimate` now includes Catmull-Rom contribution tables and a 50% allowance for
encoder backing-buffer growth and unreclaimed intermediate allocations. The
default admission budget is 768 MiB, with one active worker per process. This
permits the corpus's ordinary large sources while retaining a rejection boundary
for expensive extreme dimensions. The original square output allowance remains
conservative; reducing it needs an estimate that distinguishes fitted output
from fixed-size branding output.

After calibration, all 56 combinations passed on macOS and Linux arm64. Three
macOS repetitions gave 168 observations. The smallest estimate/peak ratio was
1.163 on macOS and 1.164 on Linux; the largest estimate was 711.3 MiB. Windows
amd64 was cross-compiled locally; its measurements run in CI.

## Portable memory controls

Go allocation exhaustion can terminate the process; `recover` is not an OOM
handling mechanism. `GOMEMLIMIT` and `debug.SetMemoryLimit` provide a process-wide
soft limit. They influence collection and memory release, rather than providing
a per-job allocation quota or a recoverable allocation error. A telemetry check
can delay admitting another job, but cannot interrupt an image codec before its
next allocation. See the [Go GC guide](https://go.dev/doc/gc-guide#Memory_limit)
and [runtime/debug](https://pkg.go.dev/runtime/debug#SetMemoryLimit).

Experimental arenas do not redirect allocations made inside these codecs or
`x/image/draw`, and are not an allocation quota. There is no arena or memory
watchdog in this implementation.

Encoded input and output already use temporary files. The current decoder and
resizer require pixel buffers in memory. The Go WebP encoder constructs its
output before calling the writer, so passing it a file does not make encoding
streaming. Spilling these internal buffers would require a different codec API
or resizer, beyond this calibration.

`IMAGE_MEMORY_BUDGET` remains admission control, separate from normal application
memory. Leave room for the API, database, caches, runtime, and other processes.
Replicas multiply image capacity; a larger cap is not protection against an
operator setting it above available memory.
