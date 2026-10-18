package imagery

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"github.com/gen2brain/webp"
	"github.com/stretchr/testify/require"
)

type memoryVariant struct {
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Mode   ResizeMode `json:"mode"`
}

type memoryWorkload struct {
	Path     string          `json:"path"`
	Variants []memoryVariant `json:"variants"`
	Results  string          `json:"results"`
}

type memoryMeasurement struct {
	Breakdown    MemoryEstimate `json:"breakdown"`
	Estimate     int64          `json:"estimate"`
	PeakHeap     uint64         `json:"peak_heap"`
	PeakRuntime  uint64         `json:"peak_runtime"`
	Allocated    uint64         `json:"allocated"`
	EncodedBytes int64          `json:"encoded_bytes"`
	OutputBytes  int64          `json:"output_bytes"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
}

var memoryScenarios = []struct {
	name     string
	variants []memoryVariant
}{
	{name: "eager", variants: []memoryVariant{
		{64, 64, ThumbnailMode},
		{128, 128, ThumbnailMode},
		{256, 256, ThumbnailMode},
		{320, 320, FitMode},
	}},
	{name: "medium", variants: []memoryVariant{{960, 960, FitMode}}},
	{name: "large", variants: []memoryVariant{{1920, 1920, FitMode}}},
	{name: "all", variants: []memoryVariant{
		{64, 64, ThumbnailMode},
		{128, 128, ThumbnailMode},
		{256, 256, ThumbnailMode},
		{320, 320, FitMode},
		{960, 960, FitMode},
		{1920, 1920, FitMode},
	}},
	{name: "branding", variants: []memoryVariant{
		{32, 32, FillMode},
		{120, 120, FillMode},
		{152, 152, FillMode},
		{167, 167, FillMode},
		{180, 180, FillMode},
		{512, 512, FillMode},
		{1200, 630, FillMode},
	}},
}

func BenchmarkEstimate(b *testing.B) {
	if webp.Dynamic() == nil {
		b.Fatal("use -tags nodynamic to measure the portable Go encoder")
	}

	paths := []string{
		"testdata/avatar.png",
		"testdata/mountains.png",
		"testdata/screenshot.png",
		"testdata/screenshot.webp",
	}
	paths = append(paths, makeMemoryCorpus(b)...)

	if directory := os.Getenv("STORYDEN_IMAGE_CORPUS"); directory != "" {
		entries, err := os.ReadDir(directory)
		require.NoError(b, err)

		for _, entry := range entries {
			if !entry.IsDir() {
				paths = append(paths, filepath.Join(directory, entry.Name()))
			}
		}
	}

	for _, path := range paths {
		b.Run(filepath.Base(path), func(b *testing.B) {
			for _, workload := range memoryScenarios {
				b.Run(workload.name, func(b *testing.B) {
					request := memoryWorkload{Path: path, Variants: workload.variants}
					benchmarkMemory(b, request)
				})
			}
		})
	}
}

func benchmarkMemory(b *testing.B, workload memoryWorkload) {
	workload.Results = filepath.Join(b.TempDir(), "memory.json")
	request, err := json.Marshal(workload)
	require.NoError(b, err)

	for range b.N {
		ctx, cancel := context.WithTimeout(b.Context(), 2*time.Minute)
		command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestMeasureImageMemory$")
		command.Env = append(os.Environ(), "STORYDEN_IMAGE_MEMORY_WORKLOAD="+string(request))
		output, err := command.CombinedOutput()
		cancel()
		require.NoError(b, err, "%s", output)

		data, err := os.ReadFile(workload.Results)
		require.NoError(b, err)

		var result memoryMeasurement
		require.NoError(b, json.Unmarshal(data, &result))
		require.GreaterOrEqual(b, result.Estimate, int64(result.PeakRuntime))

		b.ReportMetric(float64(result.Estimate), "estimate-B")
		b.ReportMetric(float64(result.Breakdown.Source), "source-B")
		b.ReportMetric(float64(result.Breakdown.Output), "output-B")
		b.ReportMetric(float64(result.Breakdown.Encoded), "encoded-B")
		b.ReportMetric(float64(result.Breakdown.Resize), "resize-B")
		b.ReportMetric(float64(result.Breakdown.Contributions), "contributions-B")
		b.ReportMetric(float64(result.Breakdown.Codec), "codec-B")
		b.ReportMetric(float64(result.Breakdown.Safety), "safety-B")
		b.ReportMetric(float64(result.PeakHeap), "peak-heap-B")
		b.ReportMetric(float64(result.PeakRuntime), "peak-runtime-B")
		b.ReportMetric(float64(result.Allocated), "allocated-B")
		b.ReportMetric(float64(result.EncodedBytes+result.OutputBytes), "temp-B")
		b.ReportMetric(float64(result.Estimate)/float64(result.PeakRuntime), "estimate/peak")
	}
}

func TestMeasureImageMemory(t *testing.T) {
	request := os.Getenv("STORYDEN_IMAGE_MEMORY_WORKLOAD")
	if request == "" {
		t.Skip("subprocess helper for BenchmarkEstimate")
	}
	require.Error(t, webp.Dynamic(), "use -tags nodynamic")

	var workload memoryWorkload
	require.NoError(t, json.Unmarshal([]byte(request), &workload))

	file, err := os.Open(workload.Path)
	require.NoError(t, err)
	defer file.Close()

	debug.FreeOSMemory()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	stop := make(chan struct{})
	measurements := make(chan memoryMeasurement, 1)
	go sampleImageMemory(before, stop, measurements)

	source, err := Open(file, 32<<20, 12000000)
	require.NoError(t, err)
	defer source.Close()

	var estimate MemoryEstimate
	for _, variant := range workload.Variants {
		candidate, err := Estimate(source, variant.Width, variant.Height, variant.Mode)
		require.NoError(t, err)
		if candidate.Total() > estimate.Total() {
			estimate = candidate
		}
	}

	decoded, err := source.Decode()
	require.NoError(t, err)

	var outputBytes int64
	for _, variant := range workload.Variants {
		outputBytes = max(outputBytes, measureImageOutput(t, decoded, variant))
	}
	runtime.KeepAlive(decoded)

	close(stop)
	result := <-measurements
	result.Estimate = estimate.Total()
	result.Breakdown = estimate
	result.EncodedBytes = source.Size
	result.OutputBytes = outputBytes
	result.Width = source.Config.Width
	result.Height = source.Config.Height

	data, err := json.Marshal(result)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(workload.Results, data, 0600))
}

func sampleImageMemory(before runtime.MemStats, stop <-chan struct{}, results chan<- memoryMeasurement) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	var result memoryMeasurement
	sample := func() {
		var current runtime.MemStats
		runtime.ReadMemStats(&current)
		if current.HeapAlloc > before.HeapAlloc {
			result.PeakHeap = max(result.PeakHeap, current.HeapAlloc-before.HeapAlloc)
		}
		used := current.Sys - current.HeapReleased
		baseline := before.Sys - before.HeapReleased
		if used > baseline {
			result.PeakRuntime = max(result.PeakRuntime, used-baseline)
		}
		result.Allocated = current.TotalAlloc - before.TotalAlloc
	}

	for {
		select {
		case <-stop:
			sample()
			results <- result
			return
		case <-ticker.C:
			sample()
		}
	}
}

func measureImageOutput(t *testing.T, source image.Image, variant memoryVariant) int64 {
	t.Helper()

	var resized *image.NRGBA
	if variant.Mode == FillMode {
		resized = Fill(source, variant.Width, variant.Height)
	} else {
		resized = Resize(source, variant.Width, variant.Height, variant.Mode == ThumbnailMode)
	}
	output, err := os.CreateTemp(t.TempDir(), "variant-*")
	require.NoError(t, err)
	defer output.Close()

	require.NoError(t, Encode(output, resized, variant.Mode == FillMode))
	info, err := output.Stat()
	require.NoError(t, err)
	require.NoError(t, output.Close())
	require.NoError(t, os.Remove(output.Name()))
	return info.Size()
}

func makeMemoryCorpus(b *testing.B) []string {
	b.Helper()
	directory := b.TempDir()
	var paths []string

	for _, fixture := range []struct {
		name   string
		width  int
		height int
		alpha  bool
	}{
		{name: "landscape.jpg", width: 3072, height: 1536},
		{name: "portrait.jpg", width: 1536, height: 3072},
		{name: "maximum.jpg", width: 4000, height: 3000},
		{name: "transparent.png", width: 3024, height: 1964, alpha: true},
		{name: "maximum-transparent.png", width: 4000, height: 3000, alpha: true},
		{name: "panorama.jpg", width: 7680, height: 1024},
		{name: "thin.png", width: 1, height: 1000000, alpha: true},
		{name: "small.png", width: 80, height: 40, alpha: true},
		{name: "square.jpg", width: 1024, height: 1024},
	} {
		pixels := image.NewNRGBA(image.Rect(0, 0, fixture.width, fixture.height))
		for y := 0; y < fixture.height; y++ {
			for x := 0; x < fixture.width; x++ {
				alpha := uint8(255)
				if fixture.alpha {
					alpha = uint8(x + y)
				}
				pixels.SetNRGBA(x, y, color.NRGBA{
					R: uint8(x*17 + y*7),
					G: uint8(x + y*11),
					B: uint8(y * 5),
					A: alpha,
				})
			}
		}

		var encoded bytes.Buffer
		if fixture.alpha {
			require.NoError(b, png.Encode(&encoded, pixels))
		} else {
			require.NoError(b, jpeg.Encode(&encoded, pixels, &jpeg.Options{Quality: 90}))
		}

		path := filepath.Join(directory, fixture.name)
		require.NoError(b, os.WriteFile(path, encoded.Bytes(), 0600))
		paths = append(paths, path)

		if fixture.name == "landscape.jpg" {
			path := filepath.Join(directory, "oriented.jpg")
			require.NoError(b, os.WriteFile(path, memoryOrientedJPEG(encoded.Bytes()), 0600))
			paths = append(paths, path)
		}
	}
	return paths
}

func memoryOrientedJPEG(data []byte) []byte {
	exif := []byte{
		'E', 'x', 'i', 'f', 0, 0,
		'I', 'I', 42, 0, 8, 0, 0, 0,
		1, 0,
		0x12, 1, 3, 0, 1, 0, 0, 0, 6, 0, 0, 0,
		0, 0, 0, 0,
	}

	oriented := append([]byte{}, data[:2]...)
	oriented = append(oriented, 0xff, 0xe1, 0, byte(len(exif)+2))
	oriented = append(oriented, exif...)
	oriented = append(oriented, data[2:]...)
	return oriented
}
