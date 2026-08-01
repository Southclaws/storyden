import type { ComponentProps } from "react";

import { cva } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

import { FormErrorText } from "../form-error-text";

const formFeedback = cva({
  base: {
    color: "fg.muted",
    fontSize: "xs",
  },
});

const FormFeedbackText = styled("p", formFeedback);

export type FormFeedbackProps = ComponentProps<typeof FormFeedbackText> & {
  error?: string;
};

export function FormFeedback({ error, children, ...props }: FormFeedbackProps) {
  if (error) {
    return <FormErrorText {...props}>{error}</FormErrorText>;
  }

  return <FormFeedbackText {...props}>{children}</FormFeedbackText>;
}
