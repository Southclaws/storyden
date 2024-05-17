import { QueueListIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function QueueAction() {
  return (
    <Link href="/queue" variant="ghost" size="sm" p="0">
      <QueueListIcon width="1.5em" />
    </Link>
  );
}
