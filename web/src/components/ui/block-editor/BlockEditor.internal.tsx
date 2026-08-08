"use client";

import { Portal } from "@ark-ui/react";
import { forwardRef, useState } from "react";
import { blockEditor } from "styled-system/recipes";

import { IconButton, IconButtonProps } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import * as Menu from "@/components/ui/menu";
import type { HTMLStyledProps } from "@/styled-system/types";
import { createStyleContext } from "@/utils/create-style-context";

const { withProvider, withContext } = createStyleContext(blockEditor);

export const Root = withProvider<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "root",
);
export const Gutter = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "gutter",
);
export const Handle = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "handle",
);

const MenuTrigger = withContext<HTMLButtonElement, IconButtonProps>(
  IconButton,
  "menuTrigger",
);

const MenuContent = withContext<HTMLDivElement, Menu.ContentProps>(
  Menu.Content,
  "menuContent",
);

const MenuAnchor = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "menuAnchor",
);

export interface MenuHandleProps extends Omit<
  IconButtonProps,
  "children" | "onClick"
> {
  children: React.ReactNode;
  dragging?: boolean;
  onOpenChange?: (open: boolean) => void;
  open?: boolean;
}

export const MenuHandle = forwardRef<HTMLButtonElement, MenuHandleProps>(
  function MenuHandle(
    {
      children,
      dragging = false,
      onOpenChange,
      open,
      "aria-label": ariaLabel = "Move or configure block",
      ...triggerProps
    },
    ref,
  ) {
    const [internalOpen, setInternalOpen] = useState(false);
    const menuOpen = open ?? internalOpen;

    function setMenuOpen(nextOpen: boolean) {
      if (open === undefined) {
        setInternalOpen(nextOpen);
      }
      onOpenChange?.(nextOpen);
    }

    return (
      <Handle>
        <Menu.Root
          lazyMount
          open={menuOpen}
          onOpenChange={({ open: nextOpen }) => setMenuOpen(nextOpen)}
          positioning={{ placement: "right-start", gutter: 0 }}
        >
          <MenuTrigger
            {...triggerProps}
            ref={ref}
            aria-expanded={menuOpen}
            aria-haspopup="menu"
            aria-label={ariaLabel}
            data-dragging={dragging ? "" : undefined}
            onClick={() => setMenuOpen(!menuOpen)}
          >
            <DragHandleIcon />
          </MenuTrigger>

          <Menu.Trigger asChild>
            <MenuAnchor aria-hidden="true" tabIndex={-1} />
          </Menu.Trigger>

          <Portal>
            <Menu.Positioner>
              <MenuContent>{children}</MenuContent>
            </Menu.Positioner>
          </Portal>
        </Menu.Root>
      </Handle>
    );
  },
);

export const Content = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "content",
);
