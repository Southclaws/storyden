import { XMarkIcon } from "@heroicons/react/24/solid";

import { Button, ButtonProps } from "src/theme/components/Button";

export function CloseAction(props: ButtonProps) {
  return (
    <Button kind="ghost" size="sm" {...props}>
      <XMarkIcon width="1.4em" />
    </Button>
  );
}
