import { PlusIcon } from "lucide-react";
import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";

export function ComposeAction(props: PropsWithChildren) {
  return (
    <LinkButton href="/new" variant="ghost" size="sm">
      <PlusIcon /> {props.children}
    </LinkButton>
  );
}
