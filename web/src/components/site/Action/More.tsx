import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { MoreIcon } from "@/components/ui/icons/More";

export function MoreAction(props: ButtonProps) {
  return (
    <IconButton variant="ghost" aria-label="More options" {...props}>
      <MoreIcon />
    </IconButton>
  );
}
