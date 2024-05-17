import { QueueListIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function QueueAction() {
  return (
    <LinkButton href="/queue" variant="ghost" size="sm" p="0">
      <QueueListIcon width="1.5em" />
    </LinkButton>
  );
}
