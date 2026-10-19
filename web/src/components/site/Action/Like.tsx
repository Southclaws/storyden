import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { LikeIcon, LikeSavedIcon } from "@/components/ui/icons/Like";

type Props = ButtonProps & { liked: boolean };

export function LikeAction(props: Props) {
  const { liked, ...rest } = props;
  return (
    <IconButton
      variant="subtle"
      aria-pressed={liked}
      aria-label="Like"
      {...rest}
    >
      {liked ? (
        <LikeSavedIcon aria-hidden="true" />
      ) : (
        <LikeIcon aria-hidden="true" />
      )}
    </IconButton>
  );
}
