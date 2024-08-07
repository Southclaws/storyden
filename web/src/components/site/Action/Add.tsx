import { PlusIcon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "@/components/ui/button";

export function AddAction(props: ButtonProps) {
  return (
    <Button variant="ghost" size="sm" {...props}>
      <PlusIcon width="1.4em" />
    </Button>
  );
}
