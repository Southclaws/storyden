package plugin_manager

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	lib_plugin "github.com/Southclaws/storyden/lib/plugin"
)

type Manager struct {
	pluginWriter  *plugin_writer.Writer
	pluginQuerier *plugin_reader.Reader
	runner        plugin_runner.Runner
}

func New(
	pluginWriter *plugin_writer.Writer,
	pluginQuerier *plugin_reader.Reader,
	runner plugin_runner.Runner,
) *Manager {
	return &Manager{
		pluginWriter:  pluginWriter,
		pluginQuerier: pluginQuerier,
		runner:        runner,
	}
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

func (m *Manager) addFromBuffer(ctx context.Context, b []byte) (*plugin.Available, error) {
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Validating here is probably not necessary any more, Load can return
	// the session + manifest and perform validation.
	pl, err := plugin.Binary(b).Validate(ctx, func(b []byte) (*lib_plugin.Manifest, error) {
		mb, err := m.runner.Validate(ctx, b)
		if err != nil {
			return nil, err
		}
		return mb, nil
	})
	if err != nil {
		details := getValidationDetails(err)

		return nil, fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid", "The provided plugin file is invalid: "+details),
		)
	}

	_, err = m.runner.Load(ctx, pl.Binary)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
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

func (m *Manager) Get(ctx context.Context, id plugin.ID) (*plugin.Record, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := m.runner.GetSession(ctx, record.Manifest.ID)
	if err != nil {
		// TODO: Better distinguish between error paths.
		// particularly around the state on the record (desired?) and runtime.
		record.State = plugin.ActiveStateError
		record.StatusMessage = "Plugin is not running"
	} else {
		hydrateSession(record, *sess)
	}

	return record, nil
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
	sessionMap := lo.KeyBy(sessions, func(s *plugin_runner.PluginSession) lib_plugin.ID {
		return s.ID()
	})

	// match up sessions to records
	for _, record := range records {
		if sess, ok := sessionMap[record.Manifest.ID]; ok {
			hydrateSession(record, *sess)
		} else {
			// TODO: If desired state is running, but session not here, mark as error?
			record.StatusMessage = "Plugin is not running"
		}
	}

	return records, nil
}

func hydrateSession(record *plugin.Record, sess plugin_runner.PluginSession) {
	record.StartedAt = sess.GetStartedAt().OrZero()
}

func (m *Manager) Delete(ctx context.Context, id plugin.ID) error {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err = m.runner.Unload(ctx, rec.Manifest.ID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.pluginWriter.Remove(ctx, id); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
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
