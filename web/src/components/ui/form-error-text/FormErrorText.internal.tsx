import { Text, type TextProps } from "../text";

export type FormErrorTextProps = Omit<TextProps<"p">, "variant">;

export function FormErrorText(props: FormErrorTextProps) {
  return <Text variant="metadata" color="status.danger.content" {...props} />;
}
