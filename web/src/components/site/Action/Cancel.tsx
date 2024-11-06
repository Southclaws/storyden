import React, { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";

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
      <CancelIcon /> {children}
    </Button>
  );
}
