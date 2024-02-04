import { Portal } from "@ark-ui/react";
import { MinusIcon, PlusIcon } from "@heroicons/react/24/solid";

import { BookmarkAction } from "src/components/site/Action/Bookmark";
import { Checkbox } from "src/theme/components/Checkbox";
import {
  MenuContent,
  MenuItem,
  MenuItemGroup,
  MenuItemGroupLabel,
  MenuSeparator,
} from "src/theme/components/Menu";
import { Menu, MenuPositioner, MenuTrigger } from "src/theme/components/Menu";

import { CollectionCreateTrigger } from "../CollectionCreate/CollectionCreateTrigger";

import { Center, HStack } from "@/styled-system/jsx";
import { Box } from "@/styled-system/jsx";

import { Props, useCollectionMenu } from "./useCollectionMenu";

export function CollectionMenu(props: Props) {
  const { ready, collections, multiSelect, isAlreadySaved, handlers } =
    useCollectionMenu(props);

  if (!ready) return null;

  return (
    <Box
      onKeyDown={handlers.handleKeyDown}
      onKeyUp={handlers.handleKeyUp}
      tabIndex={1}
    >
      <Menu
        size="sm"
        onOpenChange={handlers.handleOpenChange}
        closeOnSelect={!multiSelect}
        userSelect="none"
        onSelect={handlers.handleSelect}
      >
        <MenuTrigger asChild>
          <BookmarkAction bookmarked={isAlreadySaved} />
        </MenuTrigger>

        <Portal>
          <MenuPositioner>
            <MenuContent>
              <MenuItemGroup id="group">
                <MenuItemGroupLabel htmlFor="group">
                  Add to collections
                </MenuItemGroupLabel>

                <MenuSeparator />

                {collections.map((c) => (
                  <MenuItem id={c.id} key={c.id}>
                    <HStack>
                      {multiSelect ? (
                        <Checkbox checked={c.hasPost} />
                      ) : (
                        <Center w="5">
                          {c.hasPost ? <MinusIcon /> : <PlusIcon />}
                        </Center>
                      )}
                      {c.name}
                    </HStack>
                  </MenuItem>
                ))}
              </MenuItemGroup>

              <MenuItemGroup id="create">
                <MenuItem id="create-collection" closeOnSelect={false}>
                  <CollectionCreateTrigger kind="blank" />
                </MenuItem>
              </MenuItemGroup>
            </MenuContent>
          </MenuPositioner>
        </Portal>
      </Menu>
    </Box>
  );
}
