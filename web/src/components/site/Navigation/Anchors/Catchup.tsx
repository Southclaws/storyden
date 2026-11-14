import { CatchupIcon } from "@/components/ui/icons/Catchup";
import type { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, type AnchorProps, MenuItem } from "./Anchor";

export const CatchupID = "catchup";
export const CatchupRoute = "/catchup";
export const CatchupLabel = "Catch Up";

export function CatchupAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={CatchupID}
      route={CatchupRoute}
      label={CatchupLabel}
      icon={<CatchupIcon />}
      {...props}
    />
  );
}

export function CatchupMenuItem() {
  return (
    <MenuItem
      id={CatchupID}
      route={CatchupRoute}
      label={CatchupLabel}
      icon={<CatchupIcon />}
    />
  );
}
