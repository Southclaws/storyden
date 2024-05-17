import { BookOpenIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function KnowledgebaseAction() {
  return (
    <Link href="/directory" variant="ghost" size="sm" p="0">
      <BookOpenIcon width="1.5em" />
    </Link>
  );
}
