import { CrownIcon } from "src/components/graphics/CrownIcon";
import { Link } from "src/theme/components/Link";

export function AdminAction() {
  return (
    <Link href="/admin" variant="ghost" size="sm" p="0">
      <CrownIcon width="1.5em" />
    </Link>
  );
}
