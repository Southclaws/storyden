"use client";

import { Controller, ControllerProps } from "react-hook-form";

import { FormErrorText } from "@/components/ui/FormErrorText";
import { Input } from "@/components/ui/input";
import { Form } from "@/screens/library/LibraryPageScreen/useLibraryPageScreen";
import { HStack, LStack } from "@/styled-system/jsx";

function useLibraryPageImportFromURL() {
  //
}

export function LibraryPageImportFromURL(
  props: Omit<ControllerProps<Form>, "render">,
) {
  return (
    <Controller<Form>
      control={props.control}
      name={props.name}
      render={(form) => {
        //

        return (
          <LStack gap="0">
            <HStack w="full" justify="space-between">
              <Input
                w="full"
                size="sm"
                variant="ghost"
                color="fg.muted"
                placeholder="External URL..."
                onChange={form.field.onChange}
              />

              {/* <HStack>
                <InfoTip title="Generating a page from a URL">
                  Importing a URL will fetch the content and store it in this
                  page.
                </InfoTip>
                <Button type="button" size="xs" variant="subtle">
                  Import
                </Button>
              </HStack> */}
            </HStack>
            <FormErrorText>{form.fieldState.error?.message}</FormErrorText>
          </LStack>
        );
      }}
    />
  );
}
