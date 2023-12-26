import { ark } from "@ark-ui/react";
import type { ComponentPropsWithoutRef } from "react";
import { styled } from "styled-system/jsx";
import { type TitleInputVariantProps, titleInput } from "styled-system/recipes";

export type TitleInputProps = TitleInputVariantProps &
  ComponentPropsWithoutRef<typeof ark.input>;
export const TitleInput = styled("span", titleInput, {
  defaultProps: {
    contentEditable: true,
  },
});
