import { BookOpenIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function KnowledgebaseAction() {
  return (
    <LinkButton href="/directory" variant="ghost" size="sm" p="0">
      <BookOpenIcon width="1.5em" />
    </LinkButton>
  );
}
