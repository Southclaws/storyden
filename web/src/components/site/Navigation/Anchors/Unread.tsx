import { ReadStateUnreadIcon } from "@/components/ui/icons/ReadState";
import type { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, type AnchorProps } from "./Anchor";

export const UnreadID = "unread";
export const UnreadRoute = "/unread";
export const UnreadLabel = "Unread";

export function UnreadAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={UnreadID}
      route={UnreadRoute}
      label={UnreadLabel}
      icon={<ReadStateUnreadIcon />}
      {...props}
    />
  );
}
