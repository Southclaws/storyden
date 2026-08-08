import { Portal } from "@ark-ui/react";
import type { ReactElement, ReactNode } from "react";

import * as MenuPrimitive from "./Menu.internal";

export interface MenuProps extends Omit<MenuPrimitive.RootProps, "children"> {
  children: ReactNode;
  contentProps?: Omit<MenuPrimitive.ContentProps, "children">;
  portalled?: boolean;
  positionerProps?: Omit<MenuPrimitive.PositionerProps, "children">;
  trigger: ReactElement;
}

export function Menu({
  children,
  contentProps,
  portalled = true,
  positionerProps,
  trigger,
  ...rootProps
}: MenuProps) {
  return (
    <MenuPrimitive.Root {...rootProps}>
      <MenuPrimitive.Trigger asChild>{trigger}</MenuPrimitive.Trigger>
      <Portal disabled={!portalled}>
        <MenuPrimitive.Positioner {...positionerProps}>
          <MenuPrimitive.Content {...contentProps}>
            {children}
          </MenuPrimitive.Content>
        </MenuPrimitive.Positioner>
      </Portal>
    </MenuPrimitive.Root>
  );
}
