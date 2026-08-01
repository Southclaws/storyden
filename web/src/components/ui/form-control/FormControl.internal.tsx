import { type RecipeVariantProps, cva } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

const formControl = cva({
  base: {
    display: "flex",
    flexDirection: "column",
    gap: "1",
    width: "full",
  },
});

export type FormControlVariants = RecipeVariantProps<typeof formControl>;

export const FormControl = styled("div", formControl, {
  defaultProps: {
    role: "group",
  },
});
