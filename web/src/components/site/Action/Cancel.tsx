import React, { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";

export function CancelAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  const hasLabel = React.Children.count(children) > 0;

  if (!hasLabel) {
    return (
      <IconButton aria-label="Cancel" variant="ghost" {...props}>
        <CancelIcon />
      </IconButton>
    );
  }

  return (
    <Button variant="ghost" {...props}>
      <CancelIcon /> {children}
    </Button>
  );
}
