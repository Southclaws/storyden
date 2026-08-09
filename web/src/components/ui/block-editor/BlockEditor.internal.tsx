"use client";

import { Portal } from "@ark-ui/react";
import {
  type ComponentProps,
  forwardRef,
  useEffect,
  useRef,
  useState,
} from "react";
import { blockEditor, button } from "styled-system/recipes";

import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import * as Menu from "@/components/ui/menu";
import { cx } from "@/styled-system/css";
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

type MenuTriggerProps = ComponentProps<typeof Menu.Trigger>;

const MenuTriggerBase = forwardRef<HTMLButtonElement, MenuTriggerProps>(
  function MenuTriggerBase({ className, ...props }, ref) {
    return (
      <Menu.Trigger
        {...props}
        ref={ref}
        className={cx(button({ size: "sm", variant: "subtle" }), className)}
      />
    );
  },
);

const MenuTrigger = withContext<HTMLButtonElement, MenuTriggerProps>(
  MenuTriggerBase,
  "menuTrigger",
);

const MenuContent = withContext<HTMLDivElement, Menu.ContentProps>(
  Menu.Content,
  "menuContent",
);

export interface MenuHandleProps extends Omit<
  MenuTriggerProps,
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
          <MenuTrigger
            {...triggerProps}
            ref={ref}
            aria-label={ariaLabel}
            title={triggerProps.title ?? ariaLabel}
            data-dragging={dragging ? "" : undefined}
          >
            <DragHandleIcon />
          </MenuTrigger>

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
