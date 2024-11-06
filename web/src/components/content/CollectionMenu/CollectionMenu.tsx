import { Portal } from "@ark-ui/react";
import { KeyboardEvent, useState } from "react";

import { BookmarkAction } from "src/components/site/Action/Bookmark";

import { Unready } from "@/components/site/Unready";
import { Checkbox } from "@/components/ui/checkbox";
import { AddIcon } from "@/components/ui/icons/Add";
import { RemoveIcon } from "@/components/ui/icons/Remove";
import * as Menu from "@/components/ui/menu";
import { Box, Center, HStack } from "@/styled-system/jsx";
import { useDisclosure } from "@/utils/useDisclosure";

import { CollectionCreateTrigger } from "../CollectionCreate/CollectionCreateTrigger";

import { Props, useCollectionMenu } from "./useCollectionMenu";

export function CollectionMenu(props: Props) {
  const [multiSelect, setMultiSelect] = useState(false);
  const [selected, setSelected] = useState(0);

  const handleKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.shiftKey) setMultiSelect(true);
  };

  const handleKeyUp = (e: KeyboardEvent<HTMLDivElement>) => {
    if (!e.shiftKey && multiSelect) {
      setMultiSelect(false);
      if (selected > 0) {
        onToggle();
      }
    }
  };

  const handleReset = () => {
    setMultiSelect(false);
    setSelected(0);
  };

  const { onOpenChange: handleOpenChange, onToggle } = useDisclosure({
    onClose: handleReset,
  });

  return (
    <Box onKeyDown={handleKeyDown} onKeyUp={handleKeyUp}>
      <Menu.Root
        lazyMount
        closeOnSelect={!multiSelect}
        onOpenChange={handleOpenChange}
        positioning={{
          slide: true,
          fitViewport: true,
        }}
      >
        <Menu.Trigger asChild>
          <BookmarkAction
            variant="subtle"
            size="xs"
            bookmarked={props.thread.collections.has_collected}
          />
        </Menu.Trigger>

        <Portal>
          <Menu.Positioner>
            <LazyLoadedMenuContent {...props} multiSelect={multiSelect} />
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
    </Box>
  );
}

type LazyLoadedMenuContentProps = Props & {
  multiSelect: boolean;
};

function LazyLoadedMenuContent(props: LazyLoadedMenuContentProps) {
  const { ready, error, collections, handleSelect } = useCollectionMenu(props);

  if (!ready) {
    return <Unready error={error} />;
  }

  return (
    <Menu.Content userSelect="none" overflowY="scroll" maxH="60">
      <Menu.ItemGroup>
        <Menu.Item value="create-collection" closeOnSelect={false} asChild>
          <CollectionCreateTrigger session={props.account} variant="ghost" />
        </Menu.Item>

        {collections.map((c) => (
          <Menu.Item key={c.id} value={c.id} onClick={handleSelect(c)}>
            <HStack>
              {props.multiSelect ? (
                <Checkbox checked={c.has_queried_item} />
              ) : (
                <Center w="5">
                  {c.has_queried_item ? <RemoveIcon /> : <AddIcon />}
                </Center>
              )}
              {c.name}
            </HStack>
          </Menu.Item>
        ))}
      </Menu.ItemGroup>
    </Menu.Content>
  );
}
