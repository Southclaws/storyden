package frontend

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var fallbackTimeout = 5 * time.Second

type NextjsProcess struct {
	logger *slog.Logger
	ready  chan struct{}
	once   sync.Once
}

func (p *NextjsProcess) Ready() <-chan struct{} {
	return p.ready
}

func (p *NextjsProcess) Run(ctx context.Context, path string) {
	p.ready = make(chan struct{})

	cmd := exec.CommandContext(ctx, "node", path)
	cmd.Stderr = os.Stderr

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	p.logger.Info("storyden frontend server starting", slog.String("path", path))

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go p.waitForReady(stdoutPipe)
	go p.fallbackWait()

	err = cmd.Wait()
	if err != nil {
		p.logger.Error("frontend process exited with error", slog.Any("error", err))
	}
}

func (p *NextjsProcess) waitForReady(stdout io.Reader) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		os.Stdout.WriteString(line + "\n")

		// NOTE: This is a Next.js specific log line we're looking for. THIS
		// MIGHT CHANGE! And if it does, this will break. To account for this,
		// there is a fallback after 5 seconds to simply mark the app as ready
		// and allow requests through. If this does happen, a warning is logged.
		// It would be nice to have a better startup signal but... another day.
		if strings.Contains(line, "Ready in") {
			p.once.Do(func() {
				close(p.ready)
				p.logger.Info("frontend server is ready")
			})
		}
	}
	if err := scanner.Err(); err != nil {
		p.logger.Error("error reading frontend output", slog.Any("error", err))
	}
}

func (p *NextjsProcess) fallbackWait() {
	// Fallback: If we don't see the "Ready in" log line in 30 seconds, just
	// assume the frontend is ready anyway. This is to prevent the entire app
	// from being unusable if the log line changes or something else goes wrong.
	// We still log an error in that case, but we don't want to block the entire
	// app from working.
	select {
	case <-p.ready:

	case <-time.After(fallbackTimeout):
		p.once.Do(func() {
			close(p.ready)
			p.logger.Warn("timeout waiting for frontend to be ready, assuming it's ready anyway, if you see this message please open an issue! https://github.com/Southclaws/storyden/issues/new")
		})
	}
}
