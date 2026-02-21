package plugin_logger

import (
	"os"
	"path/filepath"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/config"
)

type Writer struct {
	dataPath string
}

func newWriter(cfg config.Config) *Writer {
	return &Writer{
		dataPath: cfg.PluginDataPath,
	}
}

type RotatingWriter struct {
	Rotator *lumberjack.Logger
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	return w.Rotator.Write(p)
}

func (w *Writer) NewWriter(dataPath string, pluginID plugin.InstallationID) (*RotatingWriter, error) {
	logPath := getOutputPath(w.dataPath, pluginID)
	logDir := filepath.Dir(logPath)

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create log directory"))
	}

	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
	}

	return &RotatingWriter{
		Rotator: writer,
	}, nil
}
