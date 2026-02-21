package plugin_logger

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/config"
)

type Reader struct {
	dataPath string
}

func newReader(cfg config.Config) *Reader {
	return &Reader{
		dataPath: cfg.PluginDataPath,
	}
}

type LogStream struct {
	Lines <-chan string
	Done  <-chan struct{}
}

func (r *Reader) StreamPluginLogs(ctx context.Context, pluginID plugin.InstallationID) (*LogStream, error) {
	logDir := getPluginLogDirectory(r.dataPath, pluginID)

	entries, err := os.ReadDir(logDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to read plugin log directory"))
		}

		entries = nil
	}

	currentLogPath := getOutputPath(r.dataPath, pluginID)

	var logFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".log" {
			fullPath := filepath.Join(logDir, entry.Name())
			if fullPath != currentLogPath {
				logFiles = append(logFiles, fullPath)
			}
		}
	}

	sort.Strings(logFiles)

	lines := make(chan string, 100)
	done := make(chan struct{})

	go func() {
		defer func() {
			close(lines)
			close(done)
		}()

		for _, logFile := range logFiles {
			file, err := os.Open(logFile)
			if err != nil {
				return
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					file.Close()
					return
				case lines <- scanner.Text():
				}
			}
			file.Close()
		}

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		var currentPos int64

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stat, err := os.Stat(currentLogPath)
				if err != nil {
					if os.IsNotExist(err) {
						currentPos = 0
						continue
					}
					return
				}

				fileSize := stat.Size()
				if fileSize < currentPos {
					currentPos = 0
				}

				if fileSize == currentPos {
					continue
				}

				file, err := os.Open(currentLogPath)
				if err != nil {
					return
				}

				scanner := bufio.NewScanner(io.NewSectionReader(file, currentPos, fileSize-currentPos))
				for scanner.Scan() {
					select {
					case <-ctx.Done():
						file.Close()
						return
					case lines <- scanner.Text():
					}
				}
				file.Close()

				currentPos = fileSize
			}
		}
	}()

	return &LogStream{
		Lines: lines,
		Done:  done,
	}, nil
}

type multiFileReader struct {
	files   []string
	current int
	reader  io.ReadCloser
}

func (m *multiFileReader) Read(p []byte) (n int, err error) {
	for {
		if m.reader == nil {
			m.current++
			if m.current >= len(m.files) {
				return 0, io.EOF
			}

			file, err := os.Open(m.files[m.current])
			if err != nil {
				return 0, err
			}
			m.reader = file
		}

		n, err = m.reader.Read(p)
		if err == io.EOF {
			m.reader.Close()
			m.reader = nil
			if m.current+1 >= len(m.files) {
				return n, io.EOF
			}
			if n > 0 {
				return n, nil
			}
			continue
		}
		return n, err
	}
}

func (m *multiFileReader) Close() error {
	if m.reader != nil {
		return m.reader.Close()
	}
	return nil
}
