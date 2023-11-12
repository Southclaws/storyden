import { BellIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function NotificationsAction() {
  return (
    <Link href="/new" kind="ghost" size="sm">
      <BellIcon />
    </Link>
  );
}
