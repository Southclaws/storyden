import { FormErrorText } from "../form-error-text";
import { Text, type TextProps } from "../text";

export type FormFeedbackProps = Omit<TextProps<"p">, "variant"> & {
  error?: string;
};

export function FormFeedback({ error, children, ...props }: FormFeedbackProps) {
  if (error) {
    return <FormErrorText {...props}>{error}</FormErrorText>;
  }

  return (
    <Text variant="metadata" color="text.subtle" {...props}>
      {children}
    </Text>
  );
}
