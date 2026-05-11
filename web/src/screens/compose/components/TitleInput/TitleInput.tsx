import { Controller } from "react-hook-form";

import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { HeadingInput } from "@/components/ui/heading-input";
import { useI18n } from "@/i18n/provider";

import { useTitleInput } from "./useTitleInput";

export function TitleInput() {
  const { control, fieldError } = useTitleInput();
  const { t } = useI18n();

  return (
    <>
      <FormControl>
        <Controller
          render={({ field: { onChange, ...field }, formState }) => {
            return (
              <HeadingInput
                id="title-input"
                placeholder={t("Thread title...")}
                onValueChange={onChange}
                defaultValue={formState.defaultValues?.["title"]}
                {...field}
              />
            );
          }}
          control={control}
          name="title"
        />

        <FormErrorText>{fieldError?.message?.toString()}</FormErrorText>
      </FormControl>
    </>
  );
}
