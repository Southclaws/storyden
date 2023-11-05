import { Portal } from "@ark-ui/react";

import { Bookmark } from "src/components/site/Action/Action";
import { Menu, MenuPositioner, MenuTrigger } from "src/theme/components/Menu";

import { Box } from "@/styled-system/jsx";

import { CollectionMenuItems } from "./CollectionMenuItems";
import { Props, useCollectionMenu } from "./useCollectionMenu";

export function CollectionMenu(props: Props) {
  const {
    ready,
    collections,
    multiSelect,
    onKeyDown,
    onKeyUp,
    isOpen,
    onOpen,
    onClose,
    isAlreadySaved,
  } = useCollectionMenu(props);

  if (!ready) return null;

  return (
    <Box onKeyDown={onKeyDown} onKeyUp={onKeyUp} tabIndex={1}>
      <Menu
        size="sm"
        isOpen={isOpen}
        onOpen={onOpen}
        onClose={onClose}
        closeOnSelect={!multiSelect}
        userSelect="none"
      >
        <MenuTrigger asChild>
          <Bookmark bookmarked={isAlreadySaved} />
        </MenuTrigger>

        <Portal>
          <MenuPositioner>
            <CollectionMenuItems
              initialCollections={collections}
              thread={props.thread}
              multiSelect={multiSelect}
            />
          </MenuPositioner>
        </Portal>
      </Menu>
    </Box>
  );
}
