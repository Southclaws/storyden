import { PlusCircleIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";

export function ComposeAction(props: PropsWithChildren) {
  return (
    <LinkButton href="/new" variant="ghost" size="sm">
      <PlusCircleIcon /> {props.children}
    </LinkButton>
  );
}
