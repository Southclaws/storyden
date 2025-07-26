package event_participation

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_querier"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/event/participation/participant_querier"
	"github.com/Southclaws/storyden/app/resources/event/participation/participant_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/notification/notify"
)

var (
	ErrEventClosed                   = fault.New("event is closed", ftag.With(ftag.PermissionDenied))
	ErrCannotUpdateOtherParticipants = fault.New("only hosts can update other participants", ftag.With(ftag.PermissionDenied))
	ErrCannotUpdateOwnRole           = fault.New("cannot update own role", ftag.With(ftag.PermissionDenied))
)

type Manager struct {
	logger         *slog.Logger
	accountQuerier *account_querier.Querier
	eventQuerier   *event_querier.Querier
	querier        *participant_querier.Querier
	writer         *participant_writer.Writer
	notifier       *notify.Notifier
}

func New(
	logger *slog.Logger,
	accountQuerier *account_querier.Querier,
	eventQuerier *event_querier.Querier,
	querier *participant_querier.Querier,
	writer *participant_writer.Writer,
	notifier *notify.Notifier,
) Manager {
	return Manager{
		logger:         logger,
		accountQuerier: accountQuerier,
		eventQuerier:   eventQuerier,
		querier:        querier,
		writer:         writer,
		notifier:       notifier,
	}
}

type Change struct {
	AccountID account.AccountID
	Role      opt.Optional[participation.Role]
	Status    opt.Optional[participation.Status]
	Delete    bool
}

func (m *Manager) Update(ctx context.Context, mk event_ref.QueryKey, change Change) (*event.Event, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.accountQuerier.GetByID(ctx, session)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err := m.eventQuerier.Get(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	notifications, err := func() ([]notificationTarget, error) {
		if change.AccountID == session {
			return m.updateSelf(ctx, &acc.Account, mk, evt, change)
		} else {
			return m.updateGuest(ctx, acc, mk, evt, change)
		}
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err = m.eventQuerier.Get(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	for _, notify := range notifications {
		m.notifier.Send(ctx, notify.AccountID, opt.New(acc.ID), notify.Event, datagraph.NewRef(evt))
	}

	return evt, nil
}

type notificationTarget struct {
	Event     notification.Event
	AccountID account.AccountID
}

// 	switch evt.Policy {
// 	case participation.PolicyClosed:

// 	case participation.PolicyInviteOnly:

// 	case participation.PolicyOpen:
// 	}

// 	// Grab all hosts
// 	hosts := dt.Filter(evt.Participants, func(p *participation.EventParticipant) bool { return p.Role == participation.RoleHost })

// 	// Hosts other than the currently authenticated member, if they are a host.
// 	otherHosts := dt.Filter(hosts, func(p *participation.EventParticipant) bool { return p.Account.ID != session })

// 	// Is the requesting account a participant?
// 	selfParticipation, isAttending := lo.Find(evt.Participants, func(p *participation.EventParticipant) bool { return p.Account.ID == session })

// 	// Is the requesting account a host of this event?
// 	isHost := selfParticipation.Role == participation.RoleHost

// 	// Is the requesting account making changes to other participants?
// 	_, isAffectingOtherParticipants := lo.Find(updates, func(update Spec) bool { return update.AccountID != session })

// 	logger := m.logger.With(
// 		slog.String("event", mark.Key(eventID).String()),
// 		slog.Bool("is_host", isHost),
// 		slog.Bool("is_affecting_other_participants", isAffectingOtherParticipants),
// 		slog.Bool("is_attending", isAttending),
// 		slog.Int("participant_count", len(evt.Participants)),
// 		slog.Int("update_count", len(updates)),
// 	)

// 	logger.Info("updating event participation")

// 	err = acc.Roles.Permissions().Authorise(ctx, func() error {
// 		if isHost {
// 			return nil
// 		}

// 		if isAffectingOtherParticipants {
// 			return ErrCannotUpdateOtherParticipants
// 		}

// 		return nil
// 	}, rbac.PermissionManageEvents)
// 	if err != nil {
// 		return nil, fault.Wrap(err, fctx.With(ctx))
// 	}

// 	partMap := lo.KeyBy(evt.Participants, func(p *participation.EventParticipant) account.AccountID { return p.Account.ID })

// 	defaultStatus, err := func() (participation.Status, error) {
// 		switch evt.Policy {
// 		case participation.PolicyClosed:
// 			if isHost {
// 				return participation.StatusAttending, nil
// 			} else {
// 				return participation.Status{}, ErrEventClosed
// 			}

// 		case participation.PolicyInviteOnly:
// 			if isHost {
// 				return participation.StatusInvited, nil
// 			} else {
// 				return participation.StatusRequested, nil
// 			}
// 		}

// 		return participation.StatusAttending, nil
// 	}()
// 	if err != nil {
// 		return nil, fault.Wrap(err, fctx.With(ctx))
// 	}

// 	type notificationTarget struct {
// 		Event     notification.Event
// 		AccountID account.AccountID
// 	}

// 	type mutation struct {
// 		opts   []participant_writer.Option
// 		delete bool
// 		errors error
// 	}

// 	notifications := []notificationTarget{}

// 	optMap := dt.Reduce(updates, func(updateMap map[account.AccountID]mutation, spec Spec) map[account.AccountID]mutation {
// 		accountID := spec.AccountID
// 		isSelf := accountID == session

// 		existing, exists := partMap[accountID]
// 		mut := mutation{
// 			opts: []participant_writer.Option{},
// 		}

// 		// TODO: Not sure if we need this.
// 		status := spec.Status.Or(defaultStatus)
// 		if !exists {
// 			mut.opts = append(mut.opts, participant_writer.WithStatus(status), participant_writer.WithRole(participation.RoleAttendee))
// 		}

// 		// Status rules:
// 		// - Hosts can update any status
// 		// - Non-hosts can only update their own status
// 		// - Non-hosts can only select a certain set of statuses
// 		switch status {
// 		case participation.StatusDeclined:
// 			// Members can only decline their own attendance. Hosts cannot set
// 			// another guest's status to declined.

// 			if !isSelf {
// 				mut.errors = ErrCannotUpdateOtherParticipants
// 			} else if !exists {
// 				mut.errors = fault.New("cannot decline non-existent participation", ftag.With(ftag.InvalidArgument))
// 			} else {
// 				mut.opts = append(mut.opts, participant_writer.WithStatus(participation.StatusDeclined))

// 				if existing.Status != status {
// 					notifications = append(notifications, dt.Map(hosts, func(host *participation.EventParticipant) notificationTarget {
// 						return notificationTarget{
// 							Event:     notification.EventMemberDeclinedEvent,
// 							AccountID: host.Account.ID,
// 						}
// 					})...)
// 				}
// 			}

// 		case participation.StatusAttending:
// 			// Only open events allow any member to set their status to
// 			// attending. If the event is not open, change status to requested.

// 			if isSelf {
// 				if evt.Policy == participation.PolicyOpen {
// 					mut.opts = append(mut.opts, participant_writer.WithStatus(participation.StatusAttending))

// 					if existing.Status != status {
// 						notifications = append(notifications, dt.Map(hosts, func(host *participation.EventParticipant) notificationTarget {
// 							return notificationTarget{
// 								Event:     notification.EventMemberAttendingEvent,
// 								AccountID: host.Account.ID,
// 							}
// 						})...)
// 					}
// 				} else {
// 					mut.errors = fault.New("cannot set status to attending for non-open event", ftag.With(ftag.PermissionDenied))
// 				}
// 			} else {
// 			}

// 		case participation.StatusRequested:
// 			// Only invite-only events can have requested status.
// 			// If the event is not invite-only, change status to attending.

// 			if isSelf {
// 			} else {
// 			}

// 		case participation.StatusInvited:
// 			// If the event is invite only, any host can set a participant to
// 			// invited. If the event is open any guest can invite any non-guest.

// 			if isSelf {
// 			} else {
// 			}
// 		}

// 		// Role rules:
// 		// - Hosts cannot remove themselves
// 		// - Non-hosts cannot update roles
// 		if role, ok := spec.Role.Get(); ok {
// 			if !isHost {
// 				mut.errors = ErrCannotUpdateOwnRole
// 			} else {
// 				if isSelf && role != participation.RoleHost {
// 					mut.errors = fault.New("cannot remove self as host", ftag.With(ftag.InvalidArgument))
// 				} else {
// 					mut.opts = append(mut.opts, participant_writer.WithRole(role))

// 					notifications = append(notifications, dt.Map(otherHosts, func(host *participation.EventParticipant) notificationTarget {
// 						return notificationTarget{
// 							Event:     notification.EventEventHostAdded,
// 							AccountID: host.Account.ID,
// 						}
// 					})...)
// 				}
// 			}
// 		}

// 		updateMap[accountID] = mut

// 		return updateMap
// 	}, map[account.AccountID]mutation{})

// 	errors := dt.Map(lo.Values(optMap), func(update mutation) error { return update.errors })
// 	err = multierr.Combine(errors...)
// 	if err != nil {
// 		return nil, fault.Wrap(err, fctx.With(ctx))
// 	}

// 	for accountID, mut := range optMap {
// 		if mut.delete {
// 			err := m.writer.Remove(ctx, eventID, accountID)
// 			if err != nil {
// 				return nil, fault.Wrap(err, fctx.With(ctx))
// 			}
// 		} else if len(mut.opts) > 0 {
// 			if _, exists := partMap[accountID]; exists {
// 				err := m.writer.Update(ctx, eventID, accountID, mut.opts...)
// 				if err != nil {
// 					return nil, fault.Wrap(err, fctx.With(ctx))
// 				}
// 			} else {
// 				err := m.writer.Add(ctx, eventID, accountID, mut.opts...)
// 				if err != nil {
// 					return nil, fault.Wrap(err, fctx.With(ctx))
// 				}
// 			}
// 		}
// 	}

// 	for _, notify := range notifications {
// 		m.notifier.Send(ctx, notify.AccountID, notify.Event, datagraph.NewRef(evt))
// 	}

// 	return nil, nil
// }
