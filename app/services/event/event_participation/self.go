package event_participation

import (
	"context"
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/event/participation/participant_writer"
)

func (m *Manager) updateSelf(ctx context.Context, acc *account.Account, mk event_ref.QueryKey, evt *event.Event, change Change) ([]notificationTarget, error) {
	// Cannot change own role.
	if change.Role.Ok() {
		return nil, fault.Wrap(ErrCannotUpdateOwnRole, fctx.With(ctx))
	}

	status, ok := change.Status.Get()
	if !ok {
		switch evt.Policy {
		case participation.PolicyClosed:
			return nil, fault.Wrap(ErrEventClosed, fctx.With(ctx))

		case participation.PolicyInviteOnly:
			status = participation.StatusRequested

		case participation.PolicyOpen:
			status = participation.StatusAttending
		}
	}

	// Grab all hosts
	hosts := dt.Filter(evt.Participants, func(p *participation.EventParticipant) bool { return p.Role == participation.RoleHost })

	// Hosts other than the currently authenticated member, if they are a host.
	otherHosts := dt.Filter(hosts, func(p *participation.EventParticipant) bool { return p.Account.ID != acc.ID })

	// Is the requesting account a participant?
	selfParticipation, isAttending := lo.Find(evt.Participants, func(p *participation.EventParticipant) bool { return p.Account.ID == acc.ID })

	// Is the requesting account a host of this event?
	isHost := selfParticipation.Role == participation.RoleHost

	logger := m.logger.With(
		slog.String("event_id", evt.ID.String()),
		slog.Bool("is_host", isHost),
		slog.Bool("is_attending", isAttending),
		slog.Int("participant_count", len(evt.Participants)),
		slog.Int("host_count", len(hosts)),
		slog.Int("other_host_count", len(otherHosts)),
	)

	logger.Info("updating self event participation")

	var notifications []notificationTarget
	if isAttending {
		err := m.writer.Update(ctx, mk, acc.ID, participant_writer.WithStatus(status))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		notifications = append(notifications, dt.Map(otherHosts, func(host *participation.EventParticipant) notificationTarget {
			return notificationTarget{
				Event:     notification.EventMemberAttendingEvent,
				AccountID: host.Account.ID,
			}
		})...)
	} else {
		err := m.writer.Add(ctx, mk, acc.ID, participant_writer.WithStatus(status))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return notifications, nil
}
