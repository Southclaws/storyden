import { HomeIcon } from "@heroicons/react/24/outline";

import { Link } from "src/theme/components/Link";

export function HomeAction() {
  return (
    <Link href="/" variant="ghost" size="sm" p="0">
      <HomeIcon />
    </Link>
  );
}
