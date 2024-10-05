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
    <Box onKeyDown={handlers.handleKeyDown} onKeyUp={handlers.handleKeyUp}>
      <Menu.Root
        closeOnSelect={!multiSelect}
        onOpenChange={handlers.handleOpenChange}
        onSelect={handlers.handleSelect}
        positioning={{
          slide: true,
          fitViewport: true,
        }}
      >
        <Menu.Trigger asChild>
          <BookmarkAction
            variant="solid"
            bgColor="bg.muted"
            color="fg.default"
            size="xs"
            bookmarked={isAlreadySaved}
          />
        </Menu.Trigger>

        <Portal>
          <Menu.Positioner>
            <Menu.Content userSelect="none" overflowY="scroll" maxH="60">
              <Menu.ItemGroup>
                <Menu.Item
                  value="create-collection"
                  closeOnSelect={false}
                  asChild
                >
                  <CollectionCreateTrigger variant="ghost" />
                </Menu.Item>

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
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
    </Box>
  );
}
