import data from "@emoji-mart/data";
import Picker from "@emoji-mart/react";
import { PlusIcon } from "@heroicons/react/24/outline";

import { Button } from "src/theme/components/Button";
import {
  Popover,
  PopoverContent,
  PopoverPositioner,
  PopoverTrigger,
} from "src/theme/components/Popover";

import { styled } from "@/styled-system/jsx";

import { Props, useReactList } from "./useReactList";

export function ReactList(props: Props) {
  const { authenticated, handlers } = useReactList(props);
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
        <Popover lazyMount>
          <PopoverTrigger asChild>
            <Button size="xs" aria-label="add">
              <PlusIcon width="1.25em" />
            </Button>
          </PopoverTrigger>

          <PopoverPositioner>
            <PopoverContent>
              <Picker
                data={data} //
                onEmojiSelect={handlers.onSelect}
              />
            </PopoverContent>
          </PopoverPositioner>
        </Popover>
      )}
    </styled.ul>
  );
}
