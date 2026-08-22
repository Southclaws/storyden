import {
  TrailActionRunStatus,
  TrailRunStatus,
  TrailStatus,
} from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";

type StatusPresentation = {
  label: string;
  borderColor:
    | "border.default"
    | "status.danger.border"
    | "status.info.border"
    | "status.success.border"
    | "status.warning.border";
  backgroundColor:
    | "background.controlDisabled"
    | "status.danger.surface"
    | "status.info.surface"
    | "status.success.surface"
    | "status.warning.surface";
  color:
    | "status.danger.content"
    | "status.info.content"
    | "status.success.content"
    | "status.warning.content"
    | "text.muted";
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
  return <StatusBadge presentation={TRAIL_STATUS[status]} />;
}

export function TrailRunStatusBadge({ status }: { status: TrailRunStatus }) {
  return <StatusBadge presentation={TRAIL_RUN_STATUS[status]} />;
}

export function TrailActionRunStatusBadge({
  status,
}: {
  status: TrailActionRunStatus;
}) {
  return <StatusBadge presentation={TRAIL_ACTION_RUN_STATUS[status]} />;
}

function StatusBadge({ presentation }: { presentation: StatusPresentation }) {
  return (
    <Badge
      size="sm"
      borderColor={presentation.borderColor}
      backgroundColor={presentation.backgroundColor}
      color={presentation.color}
    >
      {presentation.label}
    </Badge>
  );
}

function success(label: string): StatusPresentation {
  return {
    label,
    borderColor: "status.success.border",
    backgroundColor: "status.success.surface",
    color: "status.success.content",
  };
}

function info(label: string): StatusPresentation {
  return {
    label,
    borderColor: "status.info.border",
    backgroundColor: "status.info.surface",
    color: "status.info.content",
  };
}

function warning(label: string): StatusPresentation {
  return {
    label,
    borderColor: "status.warning.border",
    backgroundColor: "status.warning.surface",
    color: "status.warning.content",
  };
}

function danger(label: string): StatusPresentation {
  return {
    label,
    borderColor: "status.danger.border",
    backgroundColor: "status.danger.surface",
    color: "status.danger.content",
  };
}

function neutral(label: string): StatusPresentation {
  return {
    label,
    borderColor: "border.default",
    backgroundColor: "background.controlDisabled",
    color: "text.muted",
  };
}
