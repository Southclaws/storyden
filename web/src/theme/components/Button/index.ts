import { ark } from "@ark-ui/react";
import { ComponentPropsWithoutRef } from "react";

import { StyledComponent, styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";

export type ButtonProps = ButtonVariantProps &
  ComponentPropsWithoutRef<"button"> &
  StyledComponent<"button"> &
  // 'any' type because idk how the fuck to make these props work properly yet.
  any;

export const Button = styled(ark.button, button);
