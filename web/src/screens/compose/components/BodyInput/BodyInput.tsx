import { FormControl } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Controller } from "react-hook-form";

import { ContentComposer } from "src/components/ContentComposer/ContentComposer";

import { useBodyInput } from "./useBodyInput";

export function BodyInput({ children }: PropsWithChildren) {
  const { control } = useBodyInput();

  return (
    <FormControl>
      <Controller
        render={({ field, formState }) => (
          <ContentComposer
            onChange={field.onChange}
            initialValue={formState.defaultValues?.body}
            minHeight="24em"
            height="full"
          >
            {children}
          </ContentComposer>
        )}
        control={control}
        name="body"
      />
    </FormControl>
  );
}
