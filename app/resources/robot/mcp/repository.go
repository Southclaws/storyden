package mcp

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_mcp_server "github.com/Southclaws/storyden/internal/ent/robotmcpserver"
	ent_robot_mcp_tool "github.com/Southclaws/storyden/internal/ent/robotmcptool"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListServers(ctx context.Context) ([]Server, error) {
	rows, err := r.db.RobotMCPServer.Query().
		WithOauthRemoteConnection().
		WithTools(func(q *ent.RobotMCPToolQuery) {
			q.Order(ent_robot_mcp_tool.ByRemoteName(sql.OrderAsc()))
		}).
		Order(ent_robot_mcp_server.ByCreatedAt(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(rows, MapServer), nil
}

func (r *Repository) ListEnabledServers(ctx context.Context) ([]Server, error) {
	rows, err := r.db.RobotMCPServer.Query().
		Where(ent_robot_mcp_server.Enabled(true)).
		WithOauthRemoteConnection().
		WithTools(func(q *ent.RobotMCPToolQuery) {
			q.Order(ent_robot_mcp_tool.ByRemoteName(sql.OrderAsc()))
		}).
		Order(ent_robot_mcp_server.ByCreatedAt(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(rows, MapServer), nil
}

func (r *Repository) GetServer(ctx context.Context, id ServerID) (Server, error) {
	row, err := r.db.RobotMCPServer.Query().
		Where(ent_robot_mcp_server.IDEQ(xid.ID(id))).
		WithOauthRemoteConnection().
		WithTools(func(q *ent.RobotMCPToolQuery) {
			q.Order(ent_robot_mcp_tool.ByRemoteName(sql.OrderAsc()))
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return Server{}, fault.Wrap(err, fctx.With(ctx))
	}

	return MapServer(row), nil
}

func (r *Repository) CreateServer(ctx context.Context, in ServerCreate) (Server, error) {
	create := r.db.RobotMCPServer.Create().
		SetName(in.Name).
		SetSlug(in.Slug).
		SetDescription(in.Description).
		SetEndpointURL(in.EndpointURL).
		SetEnabled(in.Enabled).
		SetBearerToken(in.BearerToken).
		SetAddedBy(xid.ID(in.AddedBy))
	if in.OAuthRemoteConnectionID != nil {
		id := xid.ID(*in.OAuthRemoteConnectionID)
		create.SetOauthRemoteConnectionID(id)
	}
	row, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return Server{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return Server{}, fault.Wrap(err, fctx.With(ctx))
	}

	return r.GetServer(ctx, ServerID(row.ID))
}

func (r *Repository) ListServersByOAuthRemoteConnection(ctx context.Context, id xid.ID) ([]Server, error) {
	rows, err := r.db.RobotMCPServer.Query().
		Where(ent_robot_mcp_server.OauthRemoteConnectionIDEQ(id)).
		WithOauthRemoteConnection().
		WithTools(func(q *ent.RobotMCPToolQuery) {
			q.Order(ent_robot_mcp_tool.ByRemoteName(sql.OrderAsc()))
		}).
		Order(ent_robot_mcp_server.ByCreatedAt(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(rows, MapServer), nil
}

func (r *Repository) UpdateServer(ctx context.Context, id ServerID, in ServerUpdate) (Server, error) {
	update := r.db.RobotMCPServer.UpdateOneID(xid.ID(id))
	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.Description != nil {
		update.SetDescription(*in.Description)
	}
	if in.EndpointURL != nil {
		update.SetEndpointURL(*in.EndpointURL)
	}
	if in.Enabled != nil {
		update.SetEnabled(*in.Enabled)
	}
	if in.ClearBearerToken {
		update.ClearBearerToken()
	}
	if in.BearerToken != nil {
		update.SetBearerToken(*in.BearerToken)
	}

	if err := update.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return Server{}, fault.Wrap(err, fctx.With(ctx))
	}

	return r.GetServer(ctx, id)
}

func (r *Repository) DeleteServer(ctx context.Context, id ServerID) error {
	row, err := r.db.RobotMCPServer.Query().
		Where(ent_robot_mcp_server.IDEQ(xid.ID(id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	tx, err := r.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	defer func() { _ = tx.Rollback() }()

	if row.OauthRemoteConnectionID != nil {
		err = tx.OAuthRemoteConnection.DeleteOneID(*row.OauthRemoteConnectionID).Exec(ctx)
	} else {
		err = tx.RobotMCPServer.DeleteOneID(xid.ID(id)).Exec(ctx)
	}
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) ListTools(ctx context.Context) ([]Tool, error) {
	rows, err := r.db.RobotMCPTool.Query().
		WithServer().
		Order(ent_robot_mcp_tool.ByToolID(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(rows, func(row *ent.RobotMCPTool) Tool {
		slug := ""
		if row.Edges.Server != nil {
			slug = row.Edges.Server.Slug
		}
		return MapTool(row, slug)
	}), nil
}

func (r *Repository) UpsertTools(ctx context.Context, server Server, discovered []Tool) error {
	now := time.Now()
	seen := make([]string, 0, len(discovered))

	for _, tool := range discovered {
		seen = append(seen, tool.RemoteName)
		existing, err := r.db.RobotMCPTool.Query().
			Where(
				ent_robot_mcp_tool.ServerIDEQ(xid.ID(server.ID)),
				ent_robot_mcp_tool.RemoteNameEQ(tool.RemoteName),
			).
			Only(ctx)
		if ent.IsNotFound(err) {
			_, err = r.db.RobotMCPTool.Create().
				SetServerID(xid.ID(server.ID)).
				SetToolID(tool.ID).
				SetRemoteName(tool.RemoteName).
				SetCallableName(tool.CallableName).
				SetTitle(tool.Title).
				SetDescription(tool.Description).
				SetInputSchema(tool.InputSchema).
				SetOutputSchema(tool.OutputSchema).
				SetAnnotations(tool.Annotations).
				SetEnabled(true).
				SetLastSeenAt(now).
				Save(ctx)
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
			continue
		}
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err = r.db.RobotMCPTool.UpdateOne(existing).
			SetCallableName(tool.CallableName).
			SetTitle(tool.Title).
			SetDescription(tool.Description).
			SetInputSchema(tool.InputSchema).
			SetOutputSchema(tool.OutputSchema).
			SetAnnotations(tool.Annotations).
			SetEnabled(true).
			SetLastSeenAt(now).
			Save(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	update := r.db.RobotMCPTool.Update().
		Where(ent_robot_mcp_tool.ServerIDEQ(xid.ID(server.ID))).
		SetEnabled(false)
	if len(seen) > 0 {
		update.Where(ent_robot_mcp_tool.RemoteNameNotIn(seen...))
	}
	if _, err := update.Save(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) MarkRefreshSuccess(ctx context.Context, id ServerID, refreshedAt time.Time) error {
	return fault.Wrap(
		r.db.RobotMCPServer.UpdateOneID(xid.ID(id)).
			SetLastRefreshedAt(refreshedAt).
			ClearLastError().
			Exec(ctx),
		fctx.With(ctx),
	)
}

func (r *Repository) MarkRefreshError(ctx context.Context, id ServerID, message string) error {
	return fault.Wrap(
		r.db.RobotMCPServer.UpdateOneID(xid.ID(id)).
			SetLastError(message).
			Exec(ctx),
		fctx.With(ctx),
	)
}
