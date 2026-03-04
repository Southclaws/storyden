import { PropsWithChildren } from "react";

import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormHelperText } from "@/components/ui/FormHelperText";

export function FormFeedback(props: PropsWithChildren<{ error?: string }>) {
  if (props.error) return <FormErrorText>{props.error}</FormErrorText>;

  return <FormHelperText>{props.children}</FormHelperText>;
}
