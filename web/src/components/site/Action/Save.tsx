import { CloudArrowUpIcon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "src/theme/components/Button";

export function SaveAction(props: ButtonProps) {
  return (
    <Button kind="ghost" size="sm" {...props}>
      <CloudArrowUpIcon width="1.4em" />
    </Button>
  );
}
