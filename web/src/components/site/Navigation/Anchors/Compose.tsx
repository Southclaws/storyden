import { PlusCircleIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Link } from "src/theme/components/Link";

export function ComposeAction(props: PropsWithChildren) {
  return (
    <Link href="/new" variant="ghost" size="sm">
      <PlusCircleIcon /> {props.children}
    </Link>
  );
}
