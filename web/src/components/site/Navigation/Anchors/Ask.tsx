import { AskIcon } from "@/components/ui/icons/Ask";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const AskID = "ask";
export const AskRoute = "/ask";
export const AskLabel = "Ask";

type Props = AnchorProps & LinkButtonStyleProps;

export function AskAnchor(props: Props) {
  return (
    <Anchor
      id={AskID}
      route={AskRoute}
      label={AskLabel}
      icon={<AskIcon />}
      {...props}
    />
  );
}

export function AskMenuItem() {
  return (
    <MenuItem id={AskID} route={AskRoute} label={AskLabel} icon={<AskIcon />} />
  );
}
