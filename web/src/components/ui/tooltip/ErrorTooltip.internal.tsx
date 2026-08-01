import { Portal } from "@ark-ui/react";

import { IconButton } from "@/components/ui/icon-button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { deriveError } from "@/utils/error";

import * as Tooltip from "./Tooltip.internal";

type ErrorTooltipProps = {
  error?: unknown;
};

export function ErrorTooltip({ error }: ErrorTooltipProps) {
  if (!error) {
    return null;
  }

  return (
    <Tooltip.Root openDelay={100}>
      <Tooltip.Trigger asChild>
        <IconButton
          type="button"
          variant="ghost"
          aria-label="Show error details"
        >
          <WarningIcon color="fg.error" />
        </IconButton>
      </Tooltip.Trigger>
      <Portal>
        <Tooltip.Positioner>
          <Tooltip.Content>
            <Tooltip.Arrow>
              <Tooltip.ArrowTip />
            </Tooltip.Arrow>
            {deriveError(error)}
          </Tooltip.Content>
        </Tooltip.Positioner>
      </Portal>
    </Tooltip.Root>
  );
}
