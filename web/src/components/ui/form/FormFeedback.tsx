import { PropsWithChildren } from "react";

import { FormErrorText } from "./FormErrorText";
import { FormHelperText } from "./FormHelperText";

export function FormFeedback(props: PropsWithChildren<{ error?: string }>) {
  if (props.error) return <FormErrorText>{props.error}</FormErrorText>;

  return <FormHelperText>{props.children}</FormHelperText>;
}
