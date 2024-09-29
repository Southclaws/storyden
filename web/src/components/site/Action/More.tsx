import { EllipsisHorizontalIcon } from "@heroicons/react/24/outline";

import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";

export function MoreAction(props: ButtonProps) {
  return (
    <IconButton variant="ghost" {...props}>
      <EllipsisHorizontalIcon />
    </IconButton>
  );
}
