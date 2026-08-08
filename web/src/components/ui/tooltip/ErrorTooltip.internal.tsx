import { IconButton } from "@/components/ui/icon-button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { deriveError } from "@/utils/error";

import { Tooltip } from "./Tooltip";

type ErrorTooltipProps = {
  error?: unknown;
};

export function ErrorTooltip({ error }: ErrorTooltipProps) {
  if (!error) {
    return null;
  }

  return (
    <Tooltip content={deriveError(error)} openDelay={100}>
      <IconButton type="button" variant="ghost" aria-label="Show error details">
        <WarningIcon color="status.danger.content" />
      </IconButton>
    </Tooltip>
  );
}
