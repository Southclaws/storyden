import { useId } from "react";

import { HStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

import { Button } from "../button";
import { LikeIcon, LikeSavedIcon } from "../icons/Like";

export type LikeButtonProps = {
  size?: "small" | "large";
  count: number;
  liked: boolean;
  enabled?: boolean;
  pending?: boolean;
  label?: string;
  onClick?: () => void;
};

export function LikeButton({
  size = "small",
  count,
  liked,
  enabled = true,
  pending = false,
  label = "Like",
  onClick,
}: LikeButtonProps) {
  const countId = useId();
  const isDisabled = !enabled || pending;
  let description = `${count} ${pluralise(count, "like")}.`;
  let clickHandler = onClick;

  if (!enabled) {
    description += " Sign in to like discussions.";
  }

  if (pending) {
    clickHandler = undefined;
  }

  const buttonProps = {
    type: "button" as const,
    variant: "subtle" as const,
    "aria-label": label,
    "aria-describedby": countId,
    "aria-pressed": liked,
    "aria-disabled": pending || undefined,
    "aria-busy": pending || undefined,
    title: enabled ? label : "Sign in to like discussions",
    onClick: clickHandler,
  };
  let icon = <LikeIcon aria-hidden="true" />;

  if (liked) {
    icon = <LikeSavedIcon aria-hidden="true" />;
  }

  return (
    <Button
      {...buttonProps}
      size={size === "large" ? "md" : "sm"}
      disabled={!enabled}
    >
      <HStack gap="1" aria-hidden="true">
        {icon}
        <styled.span
          color={isDisabled ? "inherit" : "text.default"}
          fontWeight="normal"
          fontVariantNumeric="tabular-nums"
        >
          {count}
        </styled.span>
      </HStack>
      {size === "large" && (
        <styled.span
          aria-hidden="true"
          color={isDisabled ? "inherit" : "text.subtle"}
        >
          Like
        </styled.span>
      )}
      <styled.span id={countId} srOnly>
        {description}
      </styled.span>
    </Button>
  );
}
