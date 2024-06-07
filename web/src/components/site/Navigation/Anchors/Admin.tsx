import { CrownIcon } from "src/components/graphics/CrownIcon";

import { LinkButton } from "@/components/ui/link-button";

export function AdminAction() {
  return (
    <LinkButton href="/admin" variant="ghost" size="sm" p="0">
      <CrownIcon width="1.5em" />
    </LinkButton>
  );
}
