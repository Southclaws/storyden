import data from "@emoji-mart/data";
import Picker from "@emoji-mart/react";
import { PlusIcon } from "@heroicons/react/24/outline";

import { Button } from "@/components/ui/button";
import * as Popover from "@/components/ui/popover";
import { styled } from "@/styled-system/jsx";

import { Props, useReactList } from "./useReactList";

export function ReactList(props: Props) {
  const { authenticated, ref, isOpen, handlers } = useReactList(props);
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
        <Popover.Root open={isOpen} lazyMount closeOnInteractOutside={false}>
          <Popover.Anchor>
            <Button size="xs" aria-label="add" onClick={handlers.handleTrigger}>
              <PlusIcon width="1.25em" />
            </Button>
          </Popover.Anchor>

          <Popover.Positioner>
            <Popover.Content ref={ref}>
              <Picker
                data={data}
                onEmojiSelect={handlers.handleSelect}
                // TODO: When we do dark mode, this needs to be updated!
                theme="light"
              />
            </Popover.Content>
          </Popover.Positioner>
        </Popover.Root>
      )}
    </styled.ul>
  );
}
