package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, playwrightArgs []string) error {
	timestamp := time.Now().Format(time.RFC3339)
	e2eDir := filepath.Join("tests", "e2e-data")
	dataDir := filepath.Join(e2eDir, timestamp)
	relDir := "./" + dataDir

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(relDir, "data.db")
	dbURL := fmt.Sprintf("sqlite://./%s?_pragma=foreign_keys(1)", dbPath)

	backendBinary := "./backend.exe"

	log.Printf("Starting e2e test with data directory: %s api binary: %s", dataDir, backendBinary)

	backendCmd := exec.CommandContext(ctx, backendBinary)
	backendCmd.Dir = e2eDir
	backendCmd.Env = []string{
		fmt.Sprintf("DATABASE_URL=%s", dbURL),
		"LISTEN_ADDR=0.0.0.0:8001",
		"PROXY_FRONTEND_ADDRESS=http://localhost:3001",
		"PUBLIC_API_ADDRESS=http://localhost:8001",
		"PUBLIC_WEB_ADDRESS=http://localhost:3001",
	}
	backendCmd.Stdout = os.Stdout
	backendCmd.Stderr = os.Stderr
	if err := backendCmd.Start(); err != nil {
		return fmt.Errorf("failed to start backend: %w", err)
	}
	defer stopProcess("backend", backendCmd)

	frontendCmd := exec.CommandContext(ctx, "yarn", "start", "--port", "3001")
	frontendCmd.Dir = "web"
	frontendCmd.Env = append(os.Environ(),
		"PUBLIC_API_ADDRESS=http://localhost:8001",
		"PUBLIC_WEB_ADDRESS=http://localhost:3001",
	)
	frontendCmd.Stdout = os.Stdout
	frontendCmd.Stderr = os.Stderr
	if err := frontendCmd.Start(); err != nil {
		return fmt.Errorf("failed to start frontend: %w", err)
	}
	defer stopProcess("frontend", frontendCmd)

	log.Println("Waiting for services to be ready...")
	if err := waitForBackend(ctx, "http://localhost:8001/", 60*time.Second); err != nil {
		return fmt.Errorf("backend did not become ready: %w", err)
	}

	log.Println("Services ready, running Playwright tests...")
	playwrightCmdArgs := append([]string{"playwright", "test"}, playwrightArgs...)
	playwrightCmd := exec.CommandContext(ctx, "npx", playwrightCmdArgs...)
	playwrightCmd.Dir = "web"
	playwrightCmd.Stdout = os.Stdout
	playwrightCmd.Stderr = os.Stderr
	if err := playwrightCmd.Run(); err != nil {
		return fmt.Errorf("playwright tests failed: %w", err)
	}

	log.Println("E2E tests completed successfully!")
	return nil
}

func stopProcess(name string, cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	log.Println("Stopping", name, "...")

	_ = cmd.Process.Signal(syscall.SIGTERM)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			log.Println(name, "exited with:", err)
		}
	case <-time.After(10 * time.Second):
		log.Println(name, "did not stop in time; killing")
		_ = cmd.Process.Kill()
		<-done
	}
}

func waitForBackend(ctx context.Context, url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for backend at %s", url)
}
