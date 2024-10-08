import { Bars3Icon } from "@heroicons/react/24/outline";

import { Button, ButtonProps } from "@/components/ui/button";

export const MenuIcon = <Bars3Icon />;

export function MenuAction(props: ButtonProps) {
  return (
    <Button title="Main navigation menu" variant="ghost" size="sm" {...props}>
      <Bars3Icon />
    </Button>
  );
}
