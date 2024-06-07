import { LinkIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function LinksAction() {
  return (
    <LinkButton href="/l" variant="ghost" size="sm">
      <LinkIcon />
    </LinkButton>
  );
}
