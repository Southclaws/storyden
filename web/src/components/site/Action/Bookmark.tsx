import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import {
  BookmarkIcon,
  BookmarkSavedIcon,
} from "@/components/ui/icons/Bookmark";

type Props = ButtonProps & { bookmarked: boolean };

export function BookmarkAction(props: Props) {
  const { bookmarked, ...rest } = props;
  const label = bookmarked ? "Remove bookmark" : "Bookmark";

  return (
    <IconButton
      variant="subtle"
      aria-label={rest["aria-label"] ?? label}
      {...rest}
    >
      {bookmarked ? <BookmarkSavedIcon /> : <BookmarkIcon />}
    </IconButton>
  );
}
