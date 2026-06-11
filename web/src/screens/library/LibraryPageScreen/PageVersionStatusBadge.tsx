import { NodeVersionStatus } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import { CheckIcon } from "@/components/ui/icons/Check";
import { DraftIcon } from "@/components/ui/icons/Draft";

export function PageVersionStatusBadge({
  status,
}: {
  status: NodeVersionStatus;
}) {
  const isApplied = status === NodeVersionStatus.applied;
  const Icon = isApplied ? CheckIcon : DraftIcon;

  return (
    <Badge
      borderColor={isApplied ? "border.success" : "visibility.draft.border"}
      backgroundColor={isApplied ? "bg.success" : "visibility.draft.bg"}
      color={isApplied ? "fg.success" : "visibility.draft.fg"}
      css={{
        "& svg": {
          color: "current",
        },
      }}
    >
      <Icon />
      {status}
    </Badge>
  );
}
