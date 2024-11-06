"use client";

import type { Assign } from "@ark-ui/react";
import { Checkbox as ArkCheckbox } from "@ark-ui/react/checkbox";
import { forwardRef } from "react";

import { type CheckboxVariantProps, checkbox } from "@/styled-system/recipes";
import type { ComponentProps, HTMLStyledProps } from "@/styled-system/types";
import { createStyleContext } from "@/utils/create-style-context";

import { CheckIcon } from "./icons/Check";
import { RemoveIcon } from "./icons/Remove";

const { withProvider, withContext } = createStyleContext(checkbox);

export type RootProviderProps = ComponentProps<typeof RootProvider>;
export const RootProvider = withProvider<
  HTMLLabelElement,
  Assign<
    Assign<HTMLStyledProps<"label">, ArkCheckbox.RootProviderBaseProps>,
    CheckboxVariantProps
  >
>(ArkCheckbox.RootProvider, "root");

export type RootProps = ComponentProps<typeof Root>;
export const Root = withProvider<
  HTMLLabelElement,
  Assign<
    Assign<HTMLStyledProps<"label">, ArkCheckbox.RootBaseProps>,
    CheckboxVariantProps
  >
>(ArkCheckbox.Root, "root");

export const Control = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkCheckbox.ControlBaseProps>
>(ArkCheckbox.Control, "control");

export const Group = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkCheckbox.GroupBaseProps>
>(ArkCheckbox.Group, "group");

export const Indicator = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkCheckbox.IndicatorBaseProps>
>(ArkCheckbox.Indicator, "indicator");

export const Label = withContext<
  HTMLSpanElement,
  Assign<HTMLStyledProps<"span">, ArkCheckbox.LabelBaseProps>
>(ArkCheckbox.Label, "label");

export {
  CheckboxContext as Context,
  CheckboxHiddenInput as HiddenInput,
} from "@ark-ui/react/checkbox";

export const Checkbox = forwardRef<HTMLLabelElement, RootProps>(
  (props, ref) => {
    const { children, ...rootProps } = props;

    return (
      <Root ref={ref} {...rootProps}>
        <Control>
          <Indicator>
            <CheckIcon />
          </Indicator>
          <Indicator indeterminate>
            <RemoveIcon />
          </Indicator>
        </Control>
        {children && <Label>{children}</Label>}
        <ArkCheckbox.HiddenInput />
      </Root>
    );
  },
);

Checkbox.displayName = "Checkbox";
