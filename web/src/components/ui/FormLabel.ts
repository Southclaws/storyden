import { type RecipeVariantProps, cva } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

const formLabel = cva({
  base: {
    color: "fg.default",
    fontSize: "sm",
    marginBottom: "2",
  },
});

export type FormLabelVariants = RecipeVariantProps<typeof formLabel>;

export const FormLabel = styled("p", formLabel);
