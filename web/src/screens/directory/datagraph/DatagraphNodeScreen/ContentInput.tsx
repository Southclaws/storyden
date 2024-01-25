import { PropsWithChildren } from "react";
import { Controller, useFormContext } from "react-hook-form";

import { Asset } from "src/api/openapi/schemas";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { FormControl } from "src/theme/components/FormControl";

import { Form } from "./useDatagraphNodeScreen";

type Props = {
  onAssetUpload: (asset: Asset) => void;
};

export function ContentInput({
  children,
  onAssetUpload,
}: PropsWithChildren<Props>) {
  const { control } = useFormContext<Form>();

  return (
    <FormControl>
      <Controller
        render={({ field, formState }) => (
          <ContentComposer
            onChange={field.onChange}
            onAssetUpload={onAssetUpload}
            initialValue={formState.defaultValues?.["content"]}
          >
            {children}
          </ContentComposer>
        )}
        control={control}
        name="content"
      />
    </FormControl>
  );
}
