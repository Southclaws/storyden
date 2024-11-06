import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import {
  BookmarkIcon,
  BookmarkSavedIcon,
} from "@/components/ui/icons/Bookmark";

type Props = ButtonProps & { bookmarked: boolean };

export function BookmarkAction(props: Props) {
  const { bookmarked, ...rest } = props;
  return (
    <IconButton variant="subtle" size="xs" {...rest}>
      {bookmarked ? <BookmarkSavedIcon /> : <BookmarkIcon />}
    </IconButton>
  );
}
