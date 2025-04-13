package frontend

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
)

type NextjsProcess struct {
	logger *slog.Logger
}

func (p *NextjsProcess) Run(ctx context.Context, path string) {
	cmd := exec.Command("node", path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	p.logger.Info("storyden frontend server starting", slog.String("path", path))

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return
}
