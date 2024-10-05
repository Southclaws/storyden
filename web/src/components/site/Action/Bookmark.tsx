import { BookmarkIcon as BookmarkNotSavedIcon } from "@heroicons/react/24/outline";
import { BookmarkIcon as BookmarkSavedIcon } from "@heroicons/react/24/solid";

import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";

type Props = ButtonProps & { bookmarked: boolean };

export function BookmarkAction(props: Props) {
  const { bookmarked, ...rest } = props;
  return (
    <IconButton variant="subtle" size="xs" {...rest}>
      {bookmarked ? <BookmarkSavedIcon /> : <BookmarkNotSavedIcon />}
    </IconButton>
  );
}
