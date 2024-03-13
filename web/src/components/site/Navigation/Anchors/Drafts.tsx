import { DocumentIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function DraftsAction() {
  return (
    <Link href="/drafts" kind="ghost" size="sm" p="0">
      <DocumentIcon width="1.5em" />
    </Link>
  );
}
