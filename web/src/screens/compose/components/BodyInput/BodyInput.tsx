import { PropsWithChildren } from "react";
import { Controller } from "react-hook-form";

import { Asset } from "src/api/openapi-schema";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { FormControl } from "@/components/ui/form/FormControl";

import { useBodyInput } from "./useBodyInput";

type Props = {
  onAssetUpload: (asset: Asset) => void;
};

export function BodyInput({ onAssetUpload }: PropsWithChildren<Props>) {
  const { control } = useBodyInput();

  return (
    <FormControl h="full">
      <Controller
        render={({ field, formState }) => (
          <ContentComposer
            onChange={field.onChange}
            onAssetUpload={onAssetUpload}
            initialValue={formState.defaultValues?.["body"]}
          />
        )}
        control={control}
        name="body"
      />
    </FormControl>
  );
}
