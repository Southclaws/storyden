import { PlusIcon } from "@heroicons/react/24/outline";

import {
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from "src/theme/components";
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
        <Popover onOpen={onOpen}>
          <PopoverTrigger>
            <Button size="xs" aria-label="add">
              <PlusIcon width="1.25em" />
            </Button>
          </PopoverTrigger>
          <PopoverContent>
            <PopoverArrow />
            <PopoverBody m={0} p={0}>
              <Box id={`${emojiPickerContainerID}-${props.id}`}>
                [emoji picker]
              </Box>
            </PopoverBody>
          </PopoverContent>
        </Popover>
      )}
    </styled.ul>
  );
}
