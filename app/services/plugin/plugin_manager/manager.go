package plugin_manager

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/wrun"
)

type Manager struct {
	session       *session.Provider
	pluginWriter  *plugin_writer.Writer
	pluginQuerier *plugin_reader.Reader
	run           wrun.Runner
}

func New(
	session *session.Provider,
	pluginWriter *plugin_writer.Writer,
	pluginQuerier *plugin_reader.Reader,
	run wrun.Runner,
) *Manager {
	return &Manager{
		session:       session,
		pluginWriter:  pluginWriter,
		pluginQuerier: pluginQuerier,
		run:           run,
	}
}

func (m *Manager) AddFromFile(ctx context.Context, r io.Reader) (*plugin.Available, error) {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pl, err := plugin.Binary(b).Validate(ctx, m.run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pa, err := m.pluginWriter.Add(ctx, acc.ID, pl)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pa, nil
}

func (m *Manager) AddFromURL(ctx context.Context, u url.URL) (*plugin.Available, error) {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
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

	pl, err := plugin.Binary(b).Validate(ctx, m.run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pa, err := m.pluginWriter.Add(ctx, acc.ID, pl)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pa, nil
}

func (m *Manager) Get(ctx context.Context, id plugin.ID) (*plugin.Record, error) {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return record, nil
}

func (m *Manager) List(ctx context.Context) ([]*plugin.Record, error) {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	records, err := m.pluginQuerier.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return records, nil
}

func (m *Manager) Delete(ctx context.Context, id plugin.ID) error {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.pluginWriter.Remove(ctx, id); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
