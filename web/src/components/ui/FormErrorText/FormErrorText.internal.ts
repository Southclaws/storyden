import { type RecipeVariantProps, cva } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

const formErrorText = cva({
  base: {
    color: "fg.destructive",
    fontSize: "xs",
  },
});

export type FormErrorTextVariants = RecipeVariantProps<typeof formErrorText>;

export const FormErrorText = styled("p", formErrorText);
