import { EditIcon } from "lucide-react";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";

export function EditAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <EditIcon width="0.5em" height="0.5em" />
      {children}
    </Button>
  );
}
