import { Bars3Icon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "src/theme/components/Button";

export function DashboardAction(props: ButtonProps) {
  return (
    <Button title="Main navigation menu" variant="ghost" size="sm" {...props}>
      <Bars3Icon />
    </Button>
  );
}
