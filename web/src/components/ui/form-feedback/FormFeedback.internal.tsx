import { PropsWithChildren } from "react";

import { FormErrorText } from "../form-error-text";
import { FormHelperText } from "../form-helper-text";

export function FormFeedback(props: PropsWithChildren<{ error?: string }>) {
  if (props.error) return <FormErrorText>{props.error}</FormErrorText>;

  return <FormHelperText>{props.children}</FormHelperText>;
}
