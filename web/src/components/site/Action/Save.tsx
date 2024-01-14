import { CloudArrowUpIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "src/theme/components/Button";

export function SaveAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button kind="ghost" size="xs" {...props}>
      <CloudArrowUpIcon width="1.4em" /> {children}
    </Button>
  );
}
