import { numberInputAnatomy } from "@ark-ui/react/anatomy";
import { defineSlotRecipe } from "@pandacss/dev";

import {
  inputBase,
  inputSizeVariants,
  inputVariantVariants,
} from "../input/Input.recipe";

const trigger = {
  alignItems: "center",
  color: "text.subtle",
  cursor: "pointer",
  display: "flex",
  flex: "1",
  justifyContent: "center",
  lineHeight: "1",
  minHeight: "0",
  overflow: "hidden",
  transition: "common",
  userSelect: "none",
  "& svg": {
    flexShrink: "0",
  },
  _hover: {
    bg: "background.controlHover",
  },
  _active: {
    bg: "background.controlSubtle",
  },
};

export const numberInput = defineSlotRecipe({
  className: "number-input",
  slots: numberInputAnatomy.keys(),
  base: {
    root: {
      isolation: "isolate",
      position: "relative",
      width: "full",
      _disabled: {
        layerStyle: "disabled",
      },
    },
    control: {
      borderStartWidth: "1px",
      display: "flex",
      divideY: "1px",
      flexDirection: "column",
      height: "calc(100% - 2px)",
      insetEnd: "0px",
      margin: "1px",
      overflow: "hidden",
      position: "absolute",
      top: "0",
      width: "var(--stepper-width)",
      zIndex: "1",
    },
    input: {
      ...inputBase,
      verticalAlign: "top",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    incrementTrigger: {
      ...trigger,
      borderTopRightRadius: "sm",
    },
    decrementTrigger: {
      ...trigger,
      borderBottomRightRadius: "sm",
    },
  },
  defaultVariants: {
    size: "sm",
    variant: "outline",
  },
  variants: {
    variant: {
      outline: {
        input: inputVariantVariants.outline,
      },
      ghost: {
        input: inputVariantVariants.ghost,
      },
    },
    size: {
      sm: {
        root: {
          "--stepper-width": "sizes.6",
        },
        decrementTrigger: {
          "& svg": {
            height: "2",
            width: "2",
          },
        },
        incrementTrigger: {
          "& svg": {
            height: "2",
            width: "2",
          },
        },
        input: {
          ...inputSizeVariants.sm,
          px: "2",
          pe: "calc(var(--stepper-width) + 0.5rem)!",
        },
      },
      md: {
        root: {
          "--stepper-width": "sizes.8",
        },
        decrementTrigger: {
          "& svg": {
            height: "2.5",
            width: "2.5",
          },
        },
        incrementTrigger: {
          "& svg": {
            height: "2.5",
            width: "2.5",
          },
        },
        input: {
          ...inputSizeVariants.md,
          px: "2.5",
          pe: "calc(var(--stepper-width) + 0.5rem)!",
        },
      },
      lg: {
        root: {
          "--stepper-width": "sizes.10",
        },
        decrementTrigger: {
          "& svg": {
            height: "3",
            width: "3",
          },
        },
        incrementTrigger: {
          "& svg": {
            height: "3",
            width: "3",
          },
        },
        input: {
          ...inputSizeVariants.lg,
          px: "3",
          pe: "calc(var(--stepper-width) + 0.5rem)!",
        },
      },
    },
  },
});
