package plugin_runner

// import (
// 	"context"
// 	"log/slog"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"golang.org/x/sync/errgroup"
// )

// func Test_wazeroRunner_RunOnce(t *testing.T) {
// 	r := require.New(t)
// 	a := assert.New(t)

// 	ctx := context.Background()

// 	runner := newWazeroRunner(ctx, slog.Default())

// 	f, err := os.ReadFile("./testdata/test1.wasm")
// 	r.NoError(err)

// 	want := []byte(`{"name":"test_plugin","version":"1.0"}`)

// 	b, err := runner.RunOnce(ctx, f, nil)
// 	r.NoError(err)
// 	a.Equal(want, b)
// }

// func Test_wazeroRunner_Session(t *testing.T) {
// 	r := require.New(t)
// 	// a := assert.New(t)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	errs := make(chan error, 1)

// 	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	}))

// 	runner := newWazeroRunner(ctx, logger)

// 	f, err := os.ReadFile("./testdata/test2.wasm")
// 	r.NoError(err)

// 	eg := errgroup.Group{}

// 	s := runner.NewSession(ctx, f)

// 	go func() {
// 		err := s.Start(ctx)
// 		errs <- err
// 	}()

// 	resch := make(chan any, 2)
// 	eg.Go(func() error {
// 		res, err := s.Send(ctx, "test", map[string]any{
// 			"message": "long",
// 			"wait":    5000,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		resch <- res

// 		return nil
// 	})
// 	eg.Go(func() error {
// 		// Force this RPC to be sent after the one above.
// 		time.Sleep(time.Millisecond * 500)

// 		res, err := s.Send(ctx, "test", map[string]any{
// 			"message": "short",
// 			"wait":    100,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		resch <- res

// 		return nil
// 	})

// 	err = eg.Wait()
// 	r.NoError(err)

// 	cancel()

// 	err = <-errs
// 	r.ErrorAs(err, &context.Canceled)

// 	res1 := <-resch
// 	r.Equal(`handled method "test": short for 100`, res1)

// 	res2 := <-resch
// 	r.Equal(`handled method "test": long for 250`, res2)
// }
