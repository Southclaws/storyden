import { PlusIcon } from "@heroicons/react/24/outline";

import {
  Box,
  Button,
  IconButton,
  List,
  ListItem,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from "src/theme/components";

import { Props, emojiPickerContainerID, useReactList } from "./useReactList";

export function ReactList(props: Props) {
  const { onOpen, authenticated } = useReactList(props);
  return (
    <List
      display="flex"
      flexDirection="row"
      gap={1}
      alignItems="center"
      alignContent="center"
      flexWrap="wrap"
      margin={0}
    >
      {props.reacts?.map((r) => (
        <ListItem key={r.id}>
          <Button size="xs">{r.emoji}</Button>
        </ListItem>
      ))}

      {authenticated && (
        <Popover onOpen={onOpen}>
          <PopoverTrigger>
            <IconButton
              size="xs"
              aria-label="add"
              icon={<PlusIcon width="1.25em" />}
            />
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
    </List>
  );
}
