"use client";

import { Portal } from "@ark-ui/react";
import { forwardRef, useEffect, useRef, useState } from "react";
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
    const handleRef = useRef<HTMLDivElement>(null);
    const contentRef = useRef<HTMLDivElement>(null);
    const menuOpen = open ?? internalOpen;

    function setMenuOpen(nextOpen: boolean) {
      if (open === undefined) {
        setInternalOpen(nextOpen);
      }
      onOpenChange?.(nextOpen);
    }

    useEffect(() => {
      if (!menuOpen) return;

      function closeOnOutsidePointerDown(event: PointerEvent) {
        const target = event.target;

        if (
          !(target instanceof Node) ||
          handleRef.current?.contains(target) ||
          contentRef.current?.contains(target)
        ) {
          return;
        }

        setMenuOpen(false);
      }

      document.addEventListener("pointerdown", closeOnOutsidePointerDown);
      return () =>
        document.removeEventListener("pointerdown", closeOnOutsidePointerDown);
    }, [menuOpen]);

    return (
      <Handle ref={handleRef}>
        <Menu.Root
          lazyMount
          open={menuOpen}
          onOpenChange={({ open: nextOpen }) => setMenuOpen(nextOpen)}
          positioning={{ placement: "right-start", gutter: 0 }}
        >
          <Menu.Trigger asChild>
            <MenuTrigger
              {...triggerProps}
              ref={ref}
              aria-label={ariaLabel}
              data-dragging={dragging ? "" : undefined}
            >
              <DragHandleIcon />
            </MenuTrigger>
          </Menu.Trigger>

          <Portal>
            <Menu.Positioner>
              <MenuContent ref={contentRef}>{children}</MenuContent>
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
