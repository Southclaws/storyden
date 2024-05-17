import { XMarkIcon } from "@heroicons/react/24/solid";

import { Button, ButtonProps } from "@/components/ui/button";

export function CloseAction(props: ButtonProps) {
  return (
    <Button variant="ghost" size="sm" {...props}>
      <XMarkIcon width="1.4em" />
    </Button>
  );
}
