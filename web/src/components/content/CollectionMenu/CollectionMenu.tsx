import { Portal } from "@ark-ui/react";
import { MinusIcon, PlusIcon } from "@heroicons/react/24/solid";

import { BookmarkAction } from "src/components/site/Action/Bookmark";

import { Checkbox } from "@/components/ui/checkbox";
import * as Menu from "@/components/ui/menu";
import { Box, Center, HStack } from "@/styled-system/jsx";

import { CollectionCreateTrigger } from "../CollectionCreate/CollectionCreateTrigger";

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
      <Menu.Root
        size="sm"
        onOpenChange={handlers.handleOpenChange}
        closeOnSelect={!multiSelect}
        onSelect={handlers.handleSelect}
      >
        <Menu.Trigger asChild>
          <BookmarkAction bookmarked={isAlreadySaved} />
        </Menu.Trigger>

        <Portal>
          <Menu.Positioner>
            <Menu.Content userSelect="none">
              <Menu.ItemGroup id="group">
                <Menu.ItemGroupLabel>Add to collections</Menu.ItemGroupLabel>

                <Menu.Separator />

                {collections.map((c) => (
                  <Menu.Item key={c.id} value={c.id}>
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
                  </Menu.Item>
                ))}
              </Menu.ItemGroup>

              <Menu.ItemGroup id="create">
                <Menu.Item value="create-collection" closeOnSelect={false}>
                  <CollectionCreateTrigger variant="ghost" />
                </Menu.Item>
              </Menu.ItemGroup>
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
    </Box>
  );
}
