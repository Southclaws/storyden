import { PlusCircleIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function ComposeAction() {
  return (
    <Link href="/new" kind="ghost" size="sm">
      <PlusCircleIcon />
    </Link>
  );
}
