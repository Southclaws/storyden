import { useFormContext } from "react-hook-form";

import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { HeadingInputField } from "@/components/ui/heading-input";

import { FormShape } from "../ComposeForm/useComposeForm";

export function TitleInput() {
  const form = useFormContext<FormShape>();

  return (
    <FormControl>
      <HeadingInputField
        id="title-input"
        placeholder="Thread title..."
        control={form.control}
        name="title"
      />

      <FormErrorText>
        {form.formState.errors.title?.message?.toString()}
      </FormErrorText>
    </FormControl>
  );
}
