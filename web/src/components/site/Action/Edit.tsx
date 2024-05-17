import { PencilIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "src/theme/components/Button";

export function EditAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <Button variant="ghost" size="xs" {...props}>
      <PencilIcon width="0.5em" height="0.5em" />
      {children}
    </Button>
  );
}
