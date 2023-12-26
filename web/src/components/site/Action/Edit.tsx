import { PencilIcon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "src/theme/components/Button";

export function EditAction(props: ButtonProps) {
  return (
    <Button kind="ghost" size="xs" {...props}>
      <PencilIcon width="0.5em" height="0.5em" />
    </Button>
  );
}
