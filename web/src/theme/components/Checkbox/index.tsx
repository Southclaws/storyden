import * as Ark from "@ark-ui/react/checkbox";
import { UseCheckboxProps } from "@ark-ui/react/checkbox/use-checkbox";
import { PropsWithChildren } from "react";
import { styled } from "styled-system/jsx";

import { createStyleContext } from "src/theme/create-style-context";

import { type CheckboxVariantProps, checkbox } from "@/styled-system/recipes";

const { withProvider, withContext } = createStyleContext(checkbox);

export * from "@ark-ui/react/checkbox";
export type CheckboxProps = Ark.CheckboxProps & CheckboxVariantProps;

const CheckboxRoot = withProvider(styled(Ark.Checkbox.Root), "root");
export const CheckboxControl = withContext(
  styled(Ark.Checkbox.Control),
  "control",
);
export const CheckboxLabel = withContext(styled(Ark.Checkbox.Label), "label");

export const _Checkbox = Object.assign(CheckboxRoot, {
  Root: CheckboxRoot,
  Control: CheckboxControl,
  Label: CheckboxLabel,
});

type Props = CheckboxVariantProps & UseCheckboxProps;

export function Checkbox(props: PropsWithChildren<Props>) {
  return (
    <CheckboxRoot {...props}>
      {(state) => {
        return (
          <>
            <CheckboxControl>
              {state.isChecked && <CheckIcon />}
              {state.isIndeterminate && <MinusIcon />}
            </CheckboxControl>

            {props.children && <CheckboxLabel>{props.children}</CheckboxLabel>}
          </>
        );
      }}
    </CheckboxRoot>
  );
}

const CheckIcon = () => (
  <svg viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path
      d="M11.6666 3.5L5.24992 9.91667L2.33325 7"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const MinusIcon = () => (
  <svg viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path
      d="M2.91675 7H11.0834"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);
