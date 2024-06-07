import { EllipsisHorizontalIcon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "@/components/ui/button";

export function MoreAction(props: ButtonProps) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <EllipsisHorizontalIcon width="1.4em" />
    </Button>
  );
}
