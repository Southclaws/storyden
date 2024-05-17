import { ArrowLeftIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { LinkButton, LinkProps } from "@/components/ui/link-button";

export function BackAction({
  children,
  ...props
}: PropsWithChildren<LinkProps>) {
  return (
    <LinkButton {...props}>
      {children ?? <ArrowLeftIcon width="1.4em" />}
    </LinkButton>
  );
}
