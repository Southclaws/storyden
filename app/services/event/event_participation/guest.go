package event_participation

import (
	"context"
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/event/participation/participant_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

var ErrNotHost = fault.New("not a host", ftag.With(ftag.PermissionDenied))

func (m *Manager) updateGuest(ctx context.Context, acc *account.AccountWithEdges, mk event_ref.QueryKey, evt *event.Event, change Change) ([]notificationTarget, error) {
	// Grab all hosts
	hosts := dt.Filter(evt.Participants, func(p *participation.EventParticipant) bool { return p.Role == participation.RoleHost })

	// Hosts other than the currently authenticated member, if they are a host.
	otherHosts := dt.Filter(hosts, func(p *participation.EventParticipant) bool { return p.Account.ID != acc.ID })

	// Is the requesting account a participant?
	selfParticipation, isSelfAttending := lo.Find(evt.Participants, func(p *participation.EventParticipant) bool { return p.Account.ID == acc.ID })

	// Is the requesting account a host of this event?
	isSelfHost := selfParticipation.Role == participation.RoleHost

	// Is the target account a participant?
	targetParticipation, isTargetAttending := lo.Find(evt.Participants, func(p *participation.EventParticipant) bool { return p.Account.ID == change.AccountID })

	// Is the target account a host of this event?
	isTargetHost := targetParticipation.Role == participation.RoleHost

	logger := m.logger.With(
		slog.String("event_id", evt.ID.String()),
		slog.Bool("is_self_host", isSelfHost),
		slog.Bool("is_self_attending", isSelfAttending),
		slog.Bool("is_target_host", isTargetHost),
		slog.Bool("is_target_attending", isTargetAttending),
		slog.Int("participant_count", len(evt.Participants)),
		slog.Int("host_count", len(hosts)),
		slog.Int("other_host_count", len(otherHosts)),
	)

	logger.Info("updating other member event participation")

	err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if !isSelfHost {
			return ErrNotHost
		}

		return nil
	}, rbac.PermissionManageEvents)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var notifications []notificationTarget

	switch {
	case change.Delete:
		// Perform deletion for target member
		err = m.writer.Remove(ctx, mk, change.AccountID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		notifications = append(notifications, notificationTarget{
			AccountID: change.AccountID,
			Event:     notification.EventAttendeeRemoved,
		})

	case isTargetAttending:
		// Perform updates to already attending target member
		opts := []participant_writer.Option{}
		if role, ok := change.Role.Get(); ok {
			opts = append(opts, participant_writer.WithRole(role))
		}

		if status, ok := change.Status.Get(); ok {
			opts = append(opts, participant_writer.WithStatus(status))
		}

		if len(opts) > 0 {
			err = m.writer.Update(ctx, mk, change.AccountID, opts...)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			// Branch on event policy
			// open:
			// - self attending -> notify hosts
			// - self declined -> notify hosts
			// - guest invited -> notify guest
			// closed:
			// - self request -> notify hosts
			// - self declined -> notify hosts
			// - host accepted -> notify guest
		}

	case !isTargetAttending:
		// Perform creation for new target member
		opts := []participant_writer.Option{}
		if role, ok := change.Role.Get(); ok {
			opts = append(opts, participant_writer.WithRole(role))
		}

		if status, ok := change.Status.Get(); ok {
			opts = append(opts, participant_writer.WithStatus(status))
		}

		err = m.writer.Add(ctx, mk, change.AccountID, opts...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

	}

	return notifications, nil
}
