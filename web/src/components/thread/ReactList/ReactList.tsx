import data from "@emoji-mart/data";
import Picker from "@emoji-mart/react";
import { PlusIcon } from "@heroicons/react/24/outline";

import { Button } from "src/theme/components/Button";
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
  PopoverPositioner,
  PopoverTrigger,
} from "src/theme/components/Popover";

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
        <Popover open={isOpen} lazyMount closeOnInteractOutside={false}>
          <PopoverAnchor>
            <Button size="xs" aria-label="add" onClick={handlers.handleTrigger}>
              <PlusIcon width="1.25em" />
            </Button>
          </PopoverAnchor>

          <PopoverPositioner ref={ref}>
            <PopoverContent>
              <Picker
                data={data}
                onEmojiSelect={handlers.handleSelect}
                // TODO: When we do dark mode, this needs to be updated!
                theme="light"
              />
            </PopoverContent>
          </PopoverPositioner>
        </Popover>
      )}
    </styled.ul>
  );
}
