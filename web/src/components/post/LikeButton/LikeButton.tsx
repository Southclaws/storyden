import { LikeAction } from "@/components/site/Action/Like";
import { Button } from "@/components/ui/button";
import { LikeIcon, LikeSavedIcon } from "@/components/ui/icons/Like";
import { Text } from "@/components/ui/text";

import { Props, useLikeButton } from "./useLikeButton";

type LikeButtonProps = Props & {
  showCount?: boolean;
};

export function LikeButton({ showCount = false, ...props }: LikeButtonProps) {
  const { enabled, handleClick } = useLikeButton({ thread: props.thread });
  const likeCount = props.thread.likes.likes;

  if (showCount) {
    return (
      <Button
        type="button"
        variant="subtle"
        display="flex"
        gap="1"
        color="text.subtle"
        aria-label={props.thread.likes.liked ? "Unlike" : "Like"}
        title={props.thread.likes.liked ? "Unlike" : "Like"}
        onClick={handleClick}
        disabled={!enabled}
      >
        <span>
          {props.thread.likes.liked ? (
            <LikeSavedIcon width="4" />
          ) : (
            <LikeIcon width="4" />
          )}
        </span>
        <Text
          as="span"
          variant="supporting"
          color="text.default"
          fontWeight="medium"
          fontVariantNumeric="tabular-nums"
          fontVariant="tabular-nums"
        >
          {likeCount}
        </Text>
      </Button>
    );
  }

  return (
    <LikeAction
      variant="subtle"
      liked={props.thread.likes.liked}
      onClick={handleClick}
      disabled={!enabled}
    />
  );
}
