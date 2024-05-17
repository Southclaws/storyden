import { LinkIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function LinksAction() {
  return (
    <Link href="/l" variant="ghost" size="sm">
      <LinkIcon />
    </Link>
  );
}
