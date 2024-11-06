import { PropsWithChildren } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { InfoIcon } from "@/components/ui/icons/Info";
import * as Popover from "@/components/ui/popover";
import { LStack } from "@/styled-system/jsx";

type Props = {
  title: string;
};

export function InfoTip({ title, children }: PropsWithChildren<Props>) {
  return (
    <Popover.Root>
      <Popover.Trigger asChild>
        <IconButton size="xs" variant="ghost" borderRadius="full">
          <InfoIcon />
        </IconButton>
      </Popover.Trigger>
      <Popover.Positioner>
        <Popover.Content>
          <Popover.Arrow>
            <Popover.ArrowTip />
          </Popover.Arrow>
          <LStack gap="1">
            <Popover.Title>{title}</Popover.Title>
            <Popover.Description>{children}</Popover.Description>
          </LStack>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  );
}
