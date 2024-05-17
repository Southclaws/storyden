import { CloudArrowUpIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";

export function SaveAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <CloudArrowUpIcon width="1.4em" /> {children}
    </Button>
  );
}
