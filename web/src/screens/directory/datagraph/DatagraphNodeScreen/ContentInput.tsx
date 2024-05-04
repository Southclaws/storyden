import { Controller, useFormContext } from "react-hook-form";

import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { ContentComposerProps } from "src/components/content/ContentComposer/useContentComposer";
import { FormControl } from "src/theme/components/FormControl";

import { Form } from "./useDatagraphNodeScreen";

type Props = ContentComposerProps;

export function ContentInput({ disabled, initialValue, onAssetUpload }: Props) {
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
          />
        )}
        control={control}
        name="content"
      />
    </FormControl>
  );
}
