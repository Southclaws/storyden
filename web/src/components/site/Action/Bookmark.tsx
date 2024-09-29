import { BookmarkIcon as BookmarkNotSavedIcon } from "@heroicons/react/24/outline";
import { BookmarkIcon as BookmarkSavedIcon } from "@heroicons/react/24/solid";

import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";

export function BookmarkAction(props: ButtonProps & { bookmarked: boolean }) {
  const { bookmarked, ...rest } = props;
  return (
    <IconButton variant="ghost" size="xs" {...rest}>
      {bookmarked ? <BookmarkSavedIcon /> : <BookmarkNotSavedIcon />}
    </IconButton>
  );
}
