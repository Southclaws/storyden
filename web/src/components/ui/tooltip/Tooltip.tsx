import { Portal } from "@ark-ui/react";
import type { ReactElement, ReactNode } from "react";

import * as TooltipPrimitive from "./Tooltip.internal";

export interface TooltipProps extends Omit<
  TooltipPrimitive.RootProps,
  "children"
> {
  children: ReactElement;
  content: ReactNode;
  contentProps?: Omit<TooltipPrimitive.ContentProps, "children">;
  portalled?: boolean;
  positionerProps?: Omit<TooltipPrimitive.PositionerProps, "children">;
  showArrow?: boolean;
}

export function Tooltip({
  children,
  content,
  contentProps,
  portalled = true,
  positionerProps,
  showArrow = true,
  ...rootProps
}: TooltipProps) {
  return (
    <TooltipPrimitive.Root {...rootProps}>
      <TooltipPrimitive.Trigger asChild>{children}</TooltipPrimitive.Trigger>
      <Portal disabled={!portalled}>
        <TooltipPrimitive.Positioner {...positionerProps}>
          {showArrow && (
            <TooltipPrimitive.Arrow>
              <TooltipPrimitive.ArrowTip />
            </TooltipPrimitive.Arrow>
          )}
          <TooltipPrimitive.Content {...contentProps}>
            {content}
          </TooltipPrimitive.Content>
        </TooltipPrimitive.Positioner>
      </Portal>
    </TooltipPrimitive.Root>
  );
}
