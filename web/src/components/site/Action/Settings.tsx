import { Cog6ToothIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function SettingsAction() {
  return (
    <Link href="/settings" kind="ghost" size="sm" p="0">
      <Cog6ToothIcon width="1.25em" />
    </Link>
  );
}
