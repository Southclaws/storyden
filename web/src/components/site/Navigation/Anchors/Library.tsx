import { BookOpenIcon } from "@heroicons/react/24/outline";

import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const LibraryID = "library";
export const LibraryRoute = "/l";
export const LibraryLabel = "Library";
export const LibraryIcon = <BookOpenIcon />;

export function LibraryAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={LibraryID}
      route={LibraryRoute}
      label={LibraryLabel}
      icon={LibraryIcon}
      {...props}
    />
  );
}

export function LibraryMenuItem() {
  return (
    <MenuItem
      id={LibraryID}
      route={LibraryRoute}
      label={LibraryLabel}
      icon={LibraryIcon}
    />
  );
}
