import { LikeAction } from "@/components/site/Action/Like";
import { LikeIcon, LikeSavedIcon } from "@/components/ui/icons/Like";
import { HStack, styled } from "@/styled-system/jsx";

import { Props, useLikeButton } from "./useLikeButton";

type LikeButtonProps = Props & {
  showCount?: boolean;
};

export function LikeButton({ showCount = false, ...props }: LikeButtonProps) {
  const { handleClick } = useLikeButton({ thread: props.thread });
  const likeCount = props.thread.likes.likes;

  if (showCount) {
    return (
      <styled.button
        type="button"
        display="flex"
        gap="1"
        alignItems="center"
        color="fg.muted"
        cursor="pointer"
        background="transparent"
        border="none"
        padding="1"
        borderRadius="sm"
        transition="colors"
        position="relative"
        zIndex="base"
        onClick={handleClick}
        aria-pressed={props.thread.likes.liked}
        aria-label={props.thread.likes.liked ? "Unlike" : "Like"}
        title={props.thread.likes.liked ? "Unlike" : "Like"}
        _hover={{
          color: "fg.default",
          background: "bg.muted",
        }}
      >
        <span>
          {props.thread.likes.liked ? (
            <LikeSavedIcon width="4" />
          ) : (
            <LikeIcon width="4" />
          )}
        </span>
        <styled.span fontSize="sm">{likeCount}</styled.span>
      </styled.button>
    );
  }

  return (
    <LikeAction
      variant="subtle"
      size="xs"
      liked={props.thread.likes.liked}
      onClick={handleClick}
    />
  );
}
