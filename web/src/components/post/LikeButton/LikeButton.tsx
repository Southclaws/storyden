import { Box } from "@/styled-system/jsx";

import { Props, useLikeButton } from "./useLikeButton";
import { LikeAction } from "@/components/site/Action/Like";

export function LikeButton(props: Props) {

  const { handleClick } = useLikeButton(
    {thread: props.thread},
  );

  return (
    <Box>
      <LikeAction
        variant="subtle"
        size="xs"
        liked={props.thread.likes.liked}
        onClick={handleClick}
      />
    </Box>
  );
}