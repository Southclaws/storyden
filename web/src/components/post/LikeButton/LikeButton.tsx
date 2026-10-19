import { useId } from "react";

import { LikeAction } from "@/components/site/Action/Like";
import { LikeIcon } from "@/components/ui/icons/Like";
import { Text } from "@/components/ui/text";
import { HStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

import { Props, useLikeButton } from "./useLikeButton";

type LikeButtonProps = Props & {
  showCount?: boolean;
  label?: string;
};

export function LikeButton({
  showCount = false,
  label = "Like",
  ...props
}: LikeButtonProps) {
  const { enabled, pending, handleClick } = useLikeButton(props);
  const countId = useId();
  const likeCount = props.thread.likes.likes;

  if (showCount) {
    return (
      <HStack
        gap="1"
        title={enabled ? undefined : "Sign in to like discussions"}
      >
        {enabled ? (
          <LikeAction
            liked={props.thread.likes.liked}
            aria-label={label}
            aria-describedby={countId}
            aria-disabled={pending || undefined}
            aria-busy={pending || undefined}
            onClick={handleClick}
          />
        ) : (
          <LikeIcon width="4" color="text.muted" aria-hidden="true" />
        )}
        <Text
          id={countId}
          as="span"
          variant="supporting"
          color="text.default"
          fontWeight="medium"
          fontVariantNumeric="tabular-nums"
          fontVariant="tabular-nums"
        >
          <span aria-hidden="true">{likeCount}</span>
          <styled.span srOnly>
            {likeCount} {pluralise(likeCount, "like")}.
            {!enabled && " Sign in to like discussions."}
          </styled.span>
        </Text>
      </HStack>
    );
  }

  return (
    <LikeAction
      variant="subtle"
      liked={props.thread.likes.liked}
      aria-label={label}
      aria-disabled={pending || undefined}
      aria-busy={pending || undefined}
      onClick={handleClick}
      disabled={!enabled}
    />
  );
}
