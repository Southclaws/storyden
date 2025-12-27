"use client";

import { NumberInput as ArkNumberInput } from "@ark-ui/react/number-input";
import {
  ArrowLeftRightIcon,
  ChevronDownIcon,
  ChevronUpIcon,
} from "lucide-react";
import type { ComponentProps } from "react";
import { type RefObject, forwardRef } from "react";
import { createStyleContext } from "styled-system/jsx";

import { numberInput } from "@/styled-system/recipes";

import { InputGroup } from "./input-group";

export { NumberInputContext as Context } from "@ark-ui/react/number-input";

const { withProvider, withContext } = createStyleContext(numberInput);

export type RootProps = ComponentProps<typeof Root>;
export const Root = withProvider(ArkNumberInput.Root, "root");
export const RootProvider = withProvider(ArkNumberInput.RootProvider, "root");
export const DecrementTrigger = withContext(
  ArkNumberInput.DecrementTrigger,
  "decrementTrigger",
  {
    defaultProps: { children: <ChevronDownIcon /> },
  },
);
export const IncrementTrigger = withContext(
  ArkNumberInput.IncrementTrigger,
  "incrementTrigger",
  {
    defaultProps: { children: <ChevronUpIcon /> },
  },
);
export const Input = withContext(ArkNumberInput.Input, "input");
export const Label = withContext(ArkNumberInput.Label, "label");
export const Scrubber = withContext(ArkNumberInput.Scrubber, "scrubber");
export const ValueText = withContext(ArkNumberInput.ValueText, "valueText");
export const Control = withContext(ArkNumberInput.Control, "control", {
  defaultProps: {
    children: (
      <>
        <IncrementTrigger />
        <DecrementTrigger />
      </>
    ),
  },
});

export interface NumberInputProps extends RootProps {
  rootRef?: RefObject<HTMLDivElement | null>;
  inputProps?: ComponentProps<typeof Input>;
  scrubber?: boolean;
}

export const NumberInput = forwardRef<HTMLInputElement, NumberInputProps>(
  function NumberInput(props, ref) {
    const { inputProps, rootRef, scrubber, ...rest } = props;
    return (
      <Root ref={rootRef} {...rest}>
        <Control />
        <InputGroup
          startElement={
            scrubber && (
              <Scrubber pointerEvents="auto">
                <ArrowLeftRightIcon />
              </Scrubber>
            )
          }
        >
          <Input ref={ref} {...inputProps} />
        </InputGroup>
      </Root>
    );
  },
);
