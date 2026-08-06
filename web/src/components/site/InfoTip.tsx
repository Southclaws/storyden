import { PropsWithChildren } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { InfoIcon } from "@/components/ui/icons/Info";
import { Description, Popover, Title } from "@/components/ui/popover";
import { LStack } from "@/styled-system/jsx";

type Props = {
  title: string;
};

export function InfoTip({ title, children }: PropsWithChildren<Props>) {
  return (
    <Popover
      trigger={
        <IconButton aria-label={title} variant="ghost" borderRadius="full">
          <InfoIcon />
        </IconButton>
      }
    >
      <LStack gap="1">
        <Title>{title}</Title>
        <Description>{children}</Description>
      </LStack>
    </Popover>
  );
}
