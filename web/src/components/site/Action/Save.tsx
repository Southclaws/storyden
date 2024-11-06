import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { SaveIcon } from "@/components/ui/icons/Save";

export function SaveAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button variant="subtle" size="xs" {...props}>
      <SaveIcon /> {children}
    </Button>
  );
}
