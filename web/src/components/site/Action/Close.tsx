import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CloseIcon } from "@/components/ui/icons/Close";

export function CloseAction(props: ButtonProps) {
  return (
    <IconButton variant="ghost" size="sm" {...props}>
      <CloseIcon />
    </IconButton>
  );
}
