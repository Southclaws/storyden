import { ReadStateUnreadIcon } from "@/components/ui/icons/ReadState";
import type { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, type AnchorProps } from "./Anchor";

export const CatchupID = "catchup";
export const CatchupRoute = "/catchup";
export const CatchupLabel = "Catch Up";

export function CatchupAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={CatchupID}
      route={CatchupRoute}
      label={CatchupLabel}
      icon={<ReadStateUnreadIcon />}
      {...props}
    />
  );
}
