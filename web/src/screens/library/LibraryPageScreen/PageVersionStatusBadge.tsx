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
      borderColor={
        isApplied ? "status.success.border" : "visibility.draft.border"
      }
      backgroundColor={
        isApplied ? "status.success.surface" : "visibility.draft.surface"
      }
      color={isApplied ? "status.success.content" : "visibility.draft.content"}
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
