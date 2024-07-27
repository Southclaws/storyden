package frontend

import (
	"context"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

type NextjsProcess struct {
	l *zap.Logger
}

func (p *NextjsProcess) Run(ctx context.Context, path string) {
	cmd := exec.Command("node", path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	p.l.Info("storyden frontend server starting",
		zap.String("path", path),
	)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return
}
