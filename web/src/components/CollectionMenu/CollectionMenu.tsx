import {
  Menu,
  MenuButton,
  MenuDivider,
  MenuGroup,
  MenuItem,
  MenuList,
} from "@chakra-ui/react";
import { MinusIcon, PlusIcon } from "@heroicons/react/24/solid";

import { Bookmark, BookmarkSolid } from "../Action/Action";

import { Props, useCollectionMenu } from "./useCollectionMenu";

export function CollectionMenu(props: Props) {
  const { collections, isAlreadySaved, onSelect } = useCollectionMenu(props);

  if (!collections) return null;

  return (
    <Menu
      preventOverflow={true}
      modifiers={[
        {
          name: "preventOverflow",
          options: {
            altAxis: true,
            padding: { bottom: 82 },
          },
        },
      ]}
    >
      <MenuButton
        title="Add to collections"
        as={isAlreadySaved ? BookmarkSolid : Bookmark}
      />
      <MenuList title="Add to collections">
        <MenuGroup title="Add to collections">
          <MenuDivider />
          {collections.map((c) => (
            <MenuItem
              key={c.id}
              icon={
                c.hasPost ? (
                  <MinusIcon width="1.4em" />
                ) : (
                  <PlusIcon width="1.4em" />
                )
              }
              onClick={onSelect(c)}
            >
              {c.name}
            </MenuItem>
          ))}
        </MenuGroup>
      </MenuList>
    </Menu>
  );
}
