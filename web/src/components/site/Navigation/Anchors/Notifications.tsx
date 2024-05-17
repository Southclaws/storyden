import { BellIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function NotificationsAction() {
  return (
    <LinkButton href="/notifications" variant="ghost" size="sm">
      <BellIcon />
    </LinkButton>
  );
}
