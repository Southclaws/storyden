import { Portal } from "@ark-ui/react";

import { WarningIcon } from "@/components/ui/icons/Warning";
import * as Popover from "@/components/ui/popover";
import { deriveError } from "@/utils/error";

import { IconButton } from "./icon-button";

type Props = {
  error?: unknown;
};

export function ErrorTooltip({ error }: Props) {
  if (!error) {
    return null;
  }

  const errorMessage = deriveError(error);

  return (
    <Popover.Root>
      <Popover.Trigger asChild>
        <IconButton
          type="button"
          size="xs"
          variant="ghost"
          aria-label="Show error details"
        >
          <WarningIcon color="fg.error" />
        </IconButton>
      </Popover.Trigger>
      <Portal>
        <Popover.Positioner>
          <Popover.Content>
            <Popover.Arrow>
              <Popover.ArrowTip />
            </Popover.Arrow>
            {errorMessage}
          </Popover.Content>
        </Popover.Positioner>
      </Portal>
    </Popover.Root>
  );
}
