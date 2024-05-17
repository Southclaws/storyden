import { DocumentIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function DraftsAction() {
  return (
    <LinkButton href="/drafts" variant="ghost" size="sm" p="0">
      <DocumentIcon width="1.5em" />
    </LinkButton>
  );
}
