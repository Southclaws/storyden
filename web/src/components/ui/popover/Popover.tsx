import { Portal } from "@ark-ui/react";
import type { ReactElement, ReactNode } from "react";

import * as PopoverPrimitive from "./Popover.internal";

export interface PopoverProps extends Omit<
  PopoverPrimitive.RootProps,
  "children"
> {
  children: ReactNode;
  contentProps?: Omit<PopoverPrimitive.ContentProps, "children">;
  portalled?: boolean;
  positionerProps?: Omit<PopoverPrimitive.PositionerProps, "children">;
  showArrow?: boolean;
  trigger: ReactElement;
}

export function Popover({
  children,
  contentProps,
  portalled = true,
  positionerProps,
  showArrow = true,
  trigger,
  ...rootProps
}: PopoverProps) {
  return (
    <PopoverPrimitive.Root {...rootProps}>
      <PopoverPrimitive.Trigger asChild>{trigger}</PopoverPrimitive.Trigger>
      <Portal disabled={!portalled}>
        <PopoverPrimitive.Positioner {...positionerProps}>
          <PopoverPrimitive.Content {...contentProps}>
            {showArrow && (
              <PopoverPrimitive.Arrow>
                <PopoverPrimitive.ArrowTip />
              </PopoverPrimitive.Arrow>
            )}
            {children}
          </PopoverPrimitive.Content>
        </PopoverPrimitive.Positioner>
      </Portal>
    </PopoverPrimitive.Root>
  );
}
