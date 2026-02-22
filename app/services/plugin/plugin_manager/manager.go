package plugin_manager

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Manager struct {
	pluginWriter  *plugin_writer.Writer
	pluginQuerier *plugin_reader.Reader
	runner        plugin_runner.Host
	bus           *pubsub.Bus
	logger        *slog.Logger
}

func New(
	lc fx.Lifecycle,
	pluginWriter *plugin_writer.Writer,
	pluginQuerier *plugin_reader.Reader,
	runner plugin_runner.Host,
	bus *pubsub.Bus,
	logger *slog.Logger,
) *Manager {
	m := &Manager{
		pluginWriter:  pluginWriter,
		pluginQuerier: pluginQuerier,
		runner:        runner,
		bus:           bus,
		logger:        logger,
	}

	lc.Append(fx.Hook{
		OnStart: m.onStart,
	})

	return m
}

func (m *Manager) AddFromFile(ctx context.Context, r io.Reader) (*plugin.Available, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return m.addFromBuffer(ctx, b)
}

func (m *Manager) AddFromURL(ctx context.Context, u url.URL) (*plugin.Available, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fault.Newf("failed to fetch plugin from URL: %s, status code: %d", u.String(), resp.StatusCode)
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return m.addFromBuffer(ctx, b)
}

func (m *Manager) AddExternal(ctx context.Context, manifest rpc.Manifest) (*plugin.Record, string, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	rec, token, err := m.pluginWriter.AddExternal(ctx, acc.ID, manifest)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.runner.Load(ctx, *rec); err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := m.runner.GetSession(ctx, rec.InstallationID)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	if err := sess.SetActiveState(ctx, plugin.ActiveStateActive); err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	return rec, token, nil
}

func (m *Manager) addFromBuffer(ctx context.Context, b []byte) (*plugin.Available, error) {
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pl, err := plugin.Binary(b).Validate(ctx)
	if err != nil {
		details := getValidationDetails(err)

		return nil, fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid", "The provided plugin file is invalid: "+details),
		)
	}

	pa, err := m.pluginWriter.Add(ctx, acc.ID, pl)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to add plugin", "An error occurred while adding the plugin."),
		)
	}

	return pa, nil
}

func (m *Manager) Get(ctx context.Context, id plugin.InstallationID) (*plugin.Record, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := m.runner.GetSession(ctx, record.InstallationID)
	if err != nil {
		record.StatusMessage = "Session not found"
		record.State = plugin.ActiveStateInactive
	} else {
		hydrateSession(record, sess)
	}

	return record, nil
}

func (m *Manager) GetSession(ctx context.Context, id plugin.InstallationID) (plugin_runner.Session, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return m.runner.GetSession(ctx, id)
}

func (m *Manager) List(ctx context.Context) ([]*plugin.Record, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	records, err := m.pluginQuerier.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sessions, err := m.runner.GetSessions(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	sessionMap := lo.KeyBy(sessions, func(s plugin_runner.Session) plugin.InstallationID {
		return s.ID()
	})

	// match up sessions to records
	for _, record := range records {
		if sess, ok := sessionMap[record.InstallationID]; ok {
			hydrateSession(record, sess)
		} else {
			record.StatusMessage = "Session not found"
			record.State = plugin.ActiveStateInactive
		}
	}

	return records, nil
}

func hydrateSession(record *plugin.Record, sess plugin_runner.Session) {
	record.ReportedState = sess.GetReportedState()
	record.StartedAt = sess.GetStartedAt().OrZero()
	record.StatusMessage = sess.GetErrorMessage()
	record.Details = sess.GetErrorDetails()
}

func (m *Manager) Delete(ctx context.Context, id plugin.InstallationID) error {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err = m.runner.Unload(ctx, rec.InstallationID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.pluginWriter.Remove(ctx, id); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (m *Manager) onStart(ctx context.Context) error {
	plugins, err := m.pluginQuerier.List(ctx)
	if err != nil {
		m.logger.Error("failed to list plugins on startup", slog.Any("error", err))
		return nil
	}

	bootCtx := context.WithoutCancel(ctx)

	for _, rec := range plugins {
		if rec.State != plugin.ActiveStateActive {
			continue
		}

		rec := rec
		logger := m.logger.With(
			slog.String("plugin_id", rec.InstallationID.String()),
			slog.String("plugin_name", rec.Manifest.Name),
		)

		go func() {
			logger.Info("starting active plugin")

			if err := m.activatePlugin(bootCtx, rec); err != nil {
				logger.Warn("failed to start plugin on startup", slog.Any("error", err))
				return
			}

			logger.Info("successfully started plugin")
		}()
	}

	return nil
}

func getValidationDetails(err error) string {
	type unwrapper interface {
		Unwrap() []error
	}

	if uw, ok := err.(unwrapper); ok {
		errs := uw.Unwrap()

		s := dt.Map(errs, func(e error) string {
			return e.Error()
		})
		return strings.Join(s, "; ")
	}

	return err.Error()
}
