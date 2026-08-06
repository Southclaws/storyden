"use client";

import type { Assign } from "@ark-ui/react";
import { Select } from "@ark-ui/react/select";

import { type SelectVariantProps, select } from "@/styled-system/recipes";
import type { ComponentProps, HTMLStyledProps } from "@/styled-system/types";
import { createStyleContext } from "@/utils/create-style-context";

const { withProvider, withContext } = createStyleContext(select);

const selectPositioning = {
  fitViewport: true,
  flip: true,
  gutter: 4,
  overflowPadding: 8,
  slide: true,
} satisfies NonNullable<
  Select.RootBaseProps<Select.CollectionItem>["positioning"]
>;

export type RootProviderProps = ComponentProps<typeof RootProvider>;
export const RootProvider = withProvider<
  HTMLDivElement,
  Assign<
    Assign<
      HTMLStyledProps<"div">,
      Select.RootProviderBaseProps<Select.CollectionItem>
    >,
    SelectVariantProps
  >
>(Select.RootProvider, "root");

const RootBase = withProvider<
  HTMLDivElement,
  Assign<
    Assign<HTMLStyledProps<"div">, Select.RootBaseProps<Select.CollectionItem>>,
    SelectVariantProps
  >
>(Select.Root, "root");

export type RootProps = ComponentProps<typeof RootBase>;

export function Root({ positioning, ...props }: RootProps) {
  return (
    <RootBase
      {...props}
      positioning={{ ...selectPositioning, ...positioning }}
    />
  );
}

export const ClearTrigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, Select.ClearTriggerBaseProps>
>(Select.ClearTrigger, "clearTrigger");

export const Content = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ContentBaseProps>
>(Select.Content, "content");

export const Control = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ControlBaseProps>
>(Select.Control, "control");

export const Indicator = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.IndicatorBaseProps>
>(Select.Indicator, "indicator");

export const ItemGroupLabel = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ItemGroupLabelBaseProps>
>(Select.ItemGroupLabel, "itemGroupLabel");

export const ItemGroup = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ItemGroupBaseProps>
>(Select.ItemGroup, "itemGroup");

export const ItemIndicator = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ItemIndicatorBaseProps>
>(Select.ItemIndicator, "itemIndicator");

export const Item = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ItemBaseProps>
>(Select.Item, "item");

export const ItemText = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"span">, Select.ItemTextBaseProps>
>(Select.ItemText, "itemText");

export const Label = withContext<
  HTMLLabelElement,
  Assign<HTMLStyledProps<"label">, Select.LabelBaseProps>
>(Select.Label, "label");

export const List = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.ListBaseProps>
>(Select.List, "list");

export const Positioner = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, Select.PositionerBaseProps>
>(Select.Positioner, "positioner");

export const Trigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, Select.TriggerBaseProps>
>(Select.Trigger, "trigger");

export const ValueText = withContext<
  HTMLSpanElement,
  Assign<HTMLStyledProps<"span">, Select.ValueTextBaseProps>
>(Select.ValueText, "valueText");

export {
  SelectContext as Context,
  SelectHiddenSelect as HiddenSelect,
} from "@ark-ui/react/select";
