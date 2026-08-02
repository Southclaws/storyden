package download

import (
	"fmt"
	"io"
	"time"

	"charm.land/bubbles/v2/progress"
	"github.com/dustin/go-humanize"

	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

const progressRefreshInterval = 100 * time.Millisecond

type downloadProgress struct {
	out        io.Writer
	bar        progress.Model
	enabled    bool
	lastRender time.Time
}

func newDownloadProgress(out io.Writer) *downloadProgress {
	barWidth := max(10, output.TerminalWidth(out, 80)-42)

	return &downloadProgress{
		out:     out,
		bar:     progress.New(progress.WithDefaultBlend(), progress.WithWidth(barWidth)),
		enabled: output.IsTerminal(out),
	}
}

func (p *downloadProgress) Update(downloaded, total int64) {
	if !p.enabled {
		return
	}

	now := time.Now()
	if downloaded > 0 && downloaded < total && now.Sub(p.lastRender) < progressRefreshInterval {
		return
	}
	p.lastRender = now

	if total <= 0 {
		fmt.Fprintf(p.out, "\r\x1b[2KDownloading plugin package: %s", humanize.Bytes(uint64(max(downloaded, 0))))
		return
	}

	percent := min(float64(downloaded)/float64(total), 1)
	fmt.Fprintf(
		p.out,
		"\r\x1b[2KDownloading %s / %s %s",
		humanize.Bytes(uint64(max(downloaded, 0))),
		humanize.Bytes(uint64(total)),
		p.bar.ViewAs(percent),
	)
}

func (p *downloadProgress) Finish(downloaded int64) {
	if !p.enabled {
		return
	}
	p.Update(downloaded, downloaded)
	fmt.Fprintln(p.out)
}

func (p *downloadProgress) Clear() {
	if p.enabled {
		fmt.Fprint(p.out, "\r\x1b[2K")
	}
}
