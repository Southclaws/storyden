import { Cog6ToothIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function SettingsAction() {
  return (
    <LinkButton href="/settings" variant="ghost" size="sm" p="0">
      <Cog6ToothIcon width="1.25em" />
    </LinkButton>
  );
}
