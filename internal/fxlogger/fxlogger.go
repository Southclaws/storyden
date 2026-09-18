package fxlogger

import (
	"log/slog"
	"strings"
	"time"

	"go.uber.org/fx/fxevent"
)

// SlogLogger is an FX event logger that uses slog for structured logging.
// It tracks dependency construction timing by measuring the time between
// consecutive Provided events, giving insight into which constructors are slow.
//
// The logger emits:
// - "elapsed": time since application start when this dependency was constructed
// - "duration": approximate time taken by this constructor (time since last Provided event)
// - Constructor details and lifecycle hook execution timing

type SlogLogger struct {
	Logger          *slog.Logger
	startTimes      map[string]time.Time
	appStart        time.Time
	lastProvided    time.Time
	lastConstructor string
}

func New(logger *slog.Logger) *SlogLogger {
	now := time.Now()
	return &SlogLogger{
		Logger:       logger,
		startTimes:   make(map[string]time.Time),
		appStart:     now,
		lastProvided: now,
	}
}

func (l *SlogLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.startTimes[e.FunctionName] = time.Now()
		l.Logger.Debug("lifecycle hook executing",
			slog.String("hook", "OnStart"),
			slog.String("callee", e.CallerName),
			slog.String("function", e.FunctionName),
		)

	case *fxevent.OnStartExecuted:
		duration := time.Since(l.startTimes[e.FunctionName])
		delete(l.startTimes, e.FunctionName)

		if e.Err != nil {
			l.Logger.Error("lifecycle hook failed",
				slog.String("hook", "OnStart"),
				slog.String("callee", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.Duration("duration", duration),
				slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
				slog.String("error", e.Err.Error()),
			)
		} else {
			l.Logger.Debug("lifecycle hook completed",
				slog.String("hook", "OnStart"),
				slog.String("callee", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.Duration("duration", duration),
				slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
			)
		}

	case *fxevent.OnStopExecuting:
		l.Logger.Debug("lifecycle hook executing",
			slog.String("hook", "OnStop"),
			slog.String("callee", e.CallerName),
			slog.String("function", e.FunctionName),
		)

	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Error("lifecycle hook failed",
				slog.String("hook", "OnStop"),
				slog.String("callee", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.String("error", e.Err.Error()),
			)
		} else {
			l.Logger.Debug("lifecycle hook completed",
				slog.String("hook", "OnStop"),
				slog.String("callee", e.CallerName),
				slog.String("function", e.FunctionName),
			)
		}

	case *fxevent.Supplied:
		l.Logger.Debug("dependency supplied",
			slog.String("type", e.TypeName),
			slog.String("module", e.ModuleName),
		)

	case *fxevent.Provided:
		now := time.Now()
		constructorName := e.ConstructorName
		if e.ModuleName != "" {
			constructorName = e.ModuleName + "." + constructorName
		}

		elapsed := now.Sub(l.appStart)
		duration := now.Sub(l.lastProvided)

		if e.Err != nil {
			l.Logger.Error("constructor failed",
				slog.String("constructor", constructorName),
				slog.Duration("elapsed", elapsed),
				slog.Duration("duration", duration),
				slog.Float64("elapsed_ms", float64(elapsed.Microseconds())/1000.0),
				slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
				slog.String("error", e.Err.Error()),
			)
		} else {
			typeNames := strings.Join(e.OutputTypeNames, ", ")

			l.Logger.Info("dependency constructed",
				slog.String("constructor", constructorName),
				slog.String("types", typeNames),
				slog.Duration("elapsed", elapsed),
				slog.Duration("duration", duration),
				slog.Float64("elapsed_ms", float64(elapsed.Microseconds())/1000.0),
				slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
			)
		}

		l.lastProvided = now
		l.lastConstructor = constructorName

	case *fxevent.Invoking:
		l.startTimes[e.FunctionName] = time.Now()
		l.Logger.Debug("invoking function",
			slog.String("function", e.FunctionName),
			slog.String("module", e.ModuleName),
		)

	case *fxevent.Invoked:
		duration := time.Since(l.startTimes[e.FunctionName])
		delete(l.startTimes, e.FunctionName)

		if e.Err != nil {
			l.Logger.Error("function invocation failed",
				slog.String("function", e.FunctionName),
				slog.Duration("duration", duration),
				slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
				slog.String("error", e.Err.Error()),
			)
		} else {
			l.Logger.Debug("function invoked",
				slog.String("function", e.FunctionName),
				slog.Duration("duration", duration),
				slog.Float64("duration_ms", float64(duration.Microseconds())/1000.0),
			)
		}

	case *fxevent.Stopping:
		l.Logger.Info("application stopping",
			slog.String("signal", strings.ToUpper(e.Signal.String())),
		)

	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error("application stop failed",
				slog.String("error", e.Err.Error()),
			)
		}

	case *fxevent.RollingBack:
		l.Logger.Error("application startup failed, rolling back",
			slog.String("error", e.StartErr.Error()),
		)

	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error("rollback failed",
				slog.String("error", e.Err.Error()),
			)
		}

	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error("application start failed",
				slog.String("error", e.Err.Error()),
			)
		} else {
			totalTime := time.Since(l.appStart)
			l.Logger.Info("application started successfully",
				slog.Duration("total_startup_time", totalTime),
				slog.Float64("total_startup_ms", float64(totalTime.Microseconds())/1000.0),
			)
		}

	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error("logger initialization failed",
				slog.String("constructor", e.ConstructorName),
				slog.String("error", e.Err.Error()),
			)
		}
	}
}
