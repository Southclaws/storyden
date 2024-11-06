import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { EditIcon } from "@/components/ui/icons/Edit";

export function EditAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <EditIcon width="4" height="4" />
      {children}
    </Button>
  );
}
