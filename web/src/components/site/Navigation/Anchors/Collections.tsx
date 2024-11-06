import { CollectionIcon } from "@/components/ui/icons/Collection";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const CollectionsID = "collections";
export const CollectionsRoute = "/c";
export const CollectionsLabel = "Collections";

export function CollectionsAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={CollectionsID}
      route={CollectionsRoute}
      label={CollectionsLabel}
      icon={<CollectionIcon />}
      {...props}
    />
  );
}

export function CollectionsMenuItem() {
  return (
    <MenuItem
      id={CollectionsID}
      route={CollectionsRoute}
      label={CollectionsLabel}
      icon={<CollectionIcon />}
    />
  );
}
