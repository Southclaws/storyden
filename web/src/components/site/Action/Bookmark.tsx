import { BookmarkIcon as BookmarkNotSavedIcon } from "@heroicons/react/24/outline";
import { BookmarkIcon as BookmarkSavedIcon } from "@heroicons/react/24/solid";

import { Button, ButtonProps } from "src/theme/components/Button";

export function BookmarkAction(props: ButtonProps & { bookmarked: boolean }) {
  const { bookmarked, ...rest } = props;
  return (
    <Button variant="ghost" size="xs" {...rest}>
      {bookmarked ? (
        <BookmarkSavedIcon width="1.4em" />
      ) : (
        <BookmarkNotSavedIcon width="1.4em" />
      )}
    </Button>
  );
}
