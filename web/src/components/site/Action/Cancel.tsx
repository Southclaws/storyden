import { XCircleIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";

export function CancelAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <XCircleIcon width="1.4em" /> {children}
    </Button>
  );
}
