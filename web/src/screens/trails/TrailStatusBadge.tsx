import {
  TrailActionRunStatus,
  TrailRunStatus,
  TrailStatus,
} from "@/api/openapi-schema";
import {
  StatusBadge,
  type StatusBadgeTone,
} from "@/components/ui/status-badge";

type StatusPresentation = {
  label: string;
  tone: StatusBadgeTone;
};

const TRAIL_STATUS: Record<TrailStatus, StatusPresentation> = {
  active: success("Active"),
  paused: warning("Paused"),
  finished: neutral("Finished"),
  archived: neutral("Archived"),
};

const TRAIL_RUN_STATUS: Record<TrailRunStatus, StatusPresentation> = {
  queued: info("Queued"),
  running: info("Running"),
  completed: success("Completed"),
  attention_required: warning("Needs attention"),
  cancelled: neutral("Cancelled"),
  skipped: neutral("Skipped"),
};

const TRAIL_ACTION_RUN_STATUS: Record<
  TrailActionRunStatus,
  StatusPresentation
> = {
  queued: info("Queued"),
  running: info("Running"),
  completed: success("Completed"),
  blocked: warning("Blocked"),
  failed: danger("Failed"),
  cancelled: neutral("Cancelled"),
};

export function TrailStatusBadge({ status }: { status: TrailStatus }) {
  return <MappedStatusBadge presentation={TRAIL_STATUS[status]} />;
}

export function TrailRunStatusBadge({ status }: { status: TrailRunStatus }) {
  return <MappedStatusBadge presentation={TRAIL_RUN_STATUS[status]} />;
}

export function TrailActionRunStatusBadge({
  status,
}: {
  status: TrailActionRunStatus;
}) {
  return <MappedStatusBadge presentation={TRAIL_ACTION_RUN_STATUS[status]} />;
}

function MappedStatusBadge({
  presentation,
}: {
  presentation: StatusPresentation;
}) {
  return (
    <StatusBadge size="sm" tone={presentation.tone}>
      {presentation.label}
    </StatusBadge>
  );
}

function success(label: string): StatusPresentation {
  return {
    label,
    tone: "success",
  };
}

function info(label: string): StatusPresentation {
  return {
    label,
    tone: "info",
  };
}

function warning(label: string): StatusPresentation {
  return {
    label,
    tone: "warning",
  };
}

function danger(label: string): StatusPresentation {
  return {
    label,
    tone: "danger",
  };
}

function neutral(label: string): StatusPresentation {
  return {
    label,
    tone: "neutral",
  };
}
