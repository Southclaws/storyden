import { Popover } from "@ark-ui/react";
import { PlusIcon } from "@heroicons/react/24/outline";

import { Button } from "src/theme/components/Button";

import { Box, styled } from "@/styled-system/jsx";

import { Props, emojiPickerContainerID, useReactList } from "./useReactList";

export function ReactList(props: Props) {
  const { onOpen, authenticated } = useReactList(props);
  return (
    <styled.ul
      display="flex"
      flexDirection="row"
      gap="1"
      alignItems="center"
      alignContent="center"
      flexWrap="wrap"
      margin="0"
    >
      {props.reacts?.map((r) => (
        <styled.li key={r.id}>
          <Button size="xs">{r.emoji}</Button>
        </styled.li>
      ))}

      {authenticated && (
        <Popover.Root onOpen={onOpen} portalled>
          <Popover.Trigger>
            <Button size="xs" aria-label="add">
              <PlusIcon width="1.25em" />
            </Button>
          </Popover.Trigger>
          <Popover.Positioner>
            <Popover.Content lazyMount>
              <Box
                id={`${emojiPickerContainerID}-${props.id}`}
                zIndex="dropdown"
              >
                [emoji picker]
              </Box>
            </Popover.Content>
          </Popover.Positioner>
        </Popover.Root>
      )}
    </styled.ul>
  );
}
