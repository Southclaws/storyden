import { EllipsisHorizontalIcon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "src/theme/components/Button";

export function MoreAction(props: ButtonProps) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <EllipsisHorizontalIcon width="1.4em" />
    </Button>
  );
}
