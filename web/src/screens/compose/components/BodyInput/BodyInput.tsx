import { PropsWithChildren } from "react";
import { Controller } from "react-hook-form";

import { Asset } from "src/api/openapi/schemas";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { FormControl } from "src/theme/components/FormControl";

import { useBodyInput } from "./useBodyInput";

type Props = {
  onAssetUpload: (asset: Asset) => void;
};

export function BodyInput({
  children,
  onAssetUpload,
}: PropsWithChildren<Props>) {
  const { control } = useBodyInput();

  return (
    <FormControl>
      <Controller
        render={({ field, formState }) => (
          <ContentComposer
            onChange={field.onChange}
            onAssetUpload={onAssetUpload}
            initialValue={formState.defaultValues?.["body"]}
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
