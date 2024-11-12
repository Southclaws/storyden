import { SearchIcon } from "@/components/ui/icons/Search";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const SearchID = "search";
export const SearchRoute = "/search";
export const SearchLabel = "Search";

type Props = AnchorProps & LinkButtonStyleProps;

export function SearchAnchor(props: Props) {
  return (
    <Anchor
      id={SearchID}
      route={SearchRoute}
      label={SearchLabel}
      icon={<SearchIcon />}
      {...props}
    />
  );
}

export function SearchMenuItem() {
  return (
    <MenuItem
      id={SearchID}
      route={SearchRoute}
      label={SearchLabel}
      icon={<SearchIcon />}
    />
  );
}
