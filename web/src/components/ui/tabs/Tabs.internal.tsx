"use client";

import type { Assign } from "@ark-ui/react";
import { Tabs } from "@ark-ui/react/tabs";
import { forwardRef, useImperativeHandle } from "react";

import { type TabsVariantProps, tabs } from "@/styled-system/recipes";
import type { ComponentProps, HTMLStyledProps } from "@/styled-system/types";
import { createStyleContext } from "@/utils/create-style-context";
import { useOverflowGradient } from "@/utils/useOverflowGradient";

const { withProvider, withContext } = createStyleContext(tabs);

export type RootProviderProps = ComponentProps<typeof RootProvider>;
export const RootProvider = withProvider<
  HTMLDivElement,
  Assign<
    Assign<HTMLStyledProps<"div">, Tabs.RootProviderBaseProps>,
    TabsVariantProps
  >
>(Tabs.RootProvider, "root");

export type RootProps = ComponentProps<typeof Root>;
export const Root = withProvider<
  HTMLDivElement,
  Assign<Assign<HTMLStyledProps<"div">, Tabs.RootBaseProps>, TabsVariantProps>
>(Tabs.Root, "root");

export const Content = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Tabs.ContentBaseProps>
>(Tabs.Content, "content");

export const Indicator = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Tabs.IndicatorBaseProps>
>(Tabs.Indicator, "indicator");

export const Bar = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "bar",
);

export const Controls = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "controls",
);

const StyledList = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Tabs.ListBaseProps>
>(Tabs.List, "list");

export const List = forwardRef<
  HTMLDivElement,
  ComponentProps<typeof StyledList>
>(function List(props, ref) {
  const viewportRef = useOverflowGradient<HTMLDivElement>();

  useImperativeHandle(ref, () => viewportRef.current!);

  return <StyledList {...props} ref={viewportRef} />;
});

export const Trigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, Tabs.TriggerBaseProps>
>(Tabs.Trigger, "trigger");

export { TabsContext as Context } from "@ark-ui/react/tabs";
