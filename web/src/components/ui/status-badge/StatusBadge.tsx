import { Badge, type BadgeProps } from "@/components/ui/badge";

export type StatusBadgeTone =
  "danger" | "info" | "neutral" | "success" | "warning";

export type StatusBadgeProps = BadgeProps & {
  tone: StatusBadgeTone;
};

const STATUS_BADGE_STYLE = {
  danger: {
    borderColor: "status.danger.border",
    backgroundColor: "status.danger.surface",
    color: "status.danger.content",
  },
  info: {
    borderColor: "status.info.border",
    backgroundColor: "status.info.surface",
    color: "status.info.content",
  },
  neutral: {
    borderColor: "border.default",
    backgroundColor: "background.controlDisabled",
    color: "text.muted",
  },
  success: {
    borderColor: "status.success.border",
    backgroundColor: "status.success.surface",
    color: "status.success.content",
  },
  warning: {
    borderColor: "status.warning.border",
    backgroundColor: "status.warning.surface",
    color: "status.warning.content",
  },
} as const;

export function StatusBadge({ tone, ...props }: StatusBadgeProps) {
  return <Badge {...props} {...STATUS_BADGE_STYLE[tone]} />;
}
