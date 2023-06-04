import { FormControl } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Controller } from "react-hook-form";

import { Compose } from "src/components/Compose/Compose";

import { useBodyInput } from "./useBodyInput";

export function BodyInput({ children }: PropsWithChildren) {
  const { control } = useBodyInput();

  return (
    <FormControl>
      <Controller
        render={({ field, formState }) => (
          <Compose
            onChange={field.onChange}
            initialValue={formState.defaultValues?.body}
            minHeight="24em"
            height="full"
          >
            {children}
          </Compose>
        )}
        control={control}
        name="body"
      />
    </FormControl>
  );
}
