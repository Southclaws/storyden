import { type RecipeVariantProps, cva } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

const formHelperText = cva({
  base: {
    color: "fg.muted",
    fontSize: "xs",
  },
});

export type FormHelperTextVariants = RecipeVariantProps<typeof formHelperText>;

export const FormHelperText = styled("p", formHelperText);
