import { PlusIcon } from "@heroicons/react/24/outline";

import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";

export function AddAction({ children, ...props }: ButtonProps) {
  return (
    <IconButton variant="ghost" size="sm" {...props}>
      <PlusIcon /> {children}
    </IconButton>
  );
}
