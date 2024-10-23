import { Controller, useFormContext } from "react-hook-form";

import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { ContentComposerProps } from "src/components/content/ContentComposer/useContentComposer";

import { FormControl } from "@/components/ui/form/FormControl";

import { Form } from "./useLibraryPageScreen";

type Props = ContentComposerProps;

export function ContentInput({
  disabled,
  initialValue,
  value,
  onAssetUpload,
}: Props) {
  const { control } = useFormContext<Form>();

  return (
    <FormControl>
      <Controller
        render={({ field }) => (
          <ContentComposer
            disabled={disabled}
            onChange={field.onChange}
            onAssetUpload={onAssetUpload}
            initialValue={initialValue}
            value={value}
          />
        )}
        control={control}
        name="content"
      />
    </FormControl>
  );
}
