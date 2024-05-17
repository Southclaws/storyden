import { SendIcon } from "src/components/graphics/SendIcon";

import { Button, ButtonProps } from "@/components/ui/button";

export function SendAction(props: ButtonProps) {
  return (
    <Button variant="ghost" size="sm" {...props}>
      <SendIcon width="1.4em" />
    </Button>
  );
}
