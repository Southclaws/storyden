import {
  Box,
  Checkbox,
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
  const {
    collections,
    isAlreadySaved,
    onSelect,
    multiSelect,
    onKeyDown,
    onKeyUp,
    isOpen,
    onOpen,
    onClose,
  } = useCollectionMenu(props);

  if (!collections) return null;

  return (
    <Box onKeyDown={onKeyDown} onKeyUp={onKeyUp} tabIndex={1}>
      <Menu
        isOpen={isOpen}
        onOpen={onOpen}
        onClose={onClose}
        closeOnSelect={!multiSelect}
        preventOverflow={true}
        modifiers={[
          {
            name: "preventOverflow",
            options: {
              altAxis: true,
              offset: { bottom: 82 },
              padding: { bottom: 82 },
            },
          },
        ]}
      >
        <MenuButton
          title="Add to collections"
          as={isAlreadySaved ? BookmarkSolid : Bookmark}
        />
        <MenuList>
          <MenuGroup title="Add to collections">
            <MenuDivider />
            {collections.map((c) => (
              <MenuItem
                key={c.id}
                icon={
                  multiSelect ? (
                    <Checkbox isChecked={c.hasPost} width="1.4em" />
                  ) : c.hasPost ? (
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
    </Box>
  );
}
