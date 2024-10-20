import { XIcon } from "lucide-react";
import React, { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";

export function CancelAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  const hasLabel = React.Children.count(children) > 0;

  return (
    <Button
      variant="ghost"
      size="xs"
      px={hasLabel ? undefined : "0"}
      {...props}
    >
      <XIcon width="1.4em" /> {children}
    </Button>
  );
}
