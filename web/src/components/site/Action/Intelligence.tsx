import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { IntelligenceIcon } from "@/components/ui/icons/Intelligence";

export function IntelligenceAction(props: ButtonProps) {
  return (
    <IconButton type="button" variant="ghost" size="sm" {...props}>
      <IntelligenceIcon />
    </IconButton>
  );
}
