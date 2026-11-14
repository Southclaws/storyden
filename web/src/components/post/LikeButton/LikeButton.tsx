import {
  LikeButton as LikeButtonControl,
  LikeButtonProps as LikeButtonControlProps,
} from "@/components/ui/like-button";

import { Props, useLikeButton } from "./useLikeButton";

type LikeButtonProps = Props & Pick<LikeButtonControlProps, "size" | "label">;

export function LikeButton({ size, label, ...props }: LikeButtonProps) {
  const { enabled, pending, handleClick } = useLikeButton(props);

  return (
    <LikeButtonControl
      size={size}
      label={label}
      count={props.thread.likes.likes}
      liked={props.thread.likes.liked}
      enabled={enabled}
      pending={pending}
      onClick={handleClick}
    />
  );
}
