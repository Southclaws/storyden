import { pull } from "lodash";
import { Controller, useFormContext } from "react-hook-form";

import { Asset } from "src/api/openapi/schemas";
import { EditableAssetWall } from "src/components/directory/datagraph/EditableAssetWall/EditableAssetWall";
import { FormControl } from "src/theme/components/FormControl";
import { FormErrorText } from "src/theme/components/FormErrorText";

import { Form } from "./useClusterScreen";

type Props = {
  editing: boolean;
  initialAssets: Asset[];
  handleAssetUpload: (asset: Asset) => void;
  handleAssetRemove: (asset: Asset) => void;
};

export function AssetsInput({
  editing,
  initialAssets,
  handleAssetUpload,
  handleAssetRemove,
}: Props) {
  const { control, formState } = useFormContext<Form>();

  const fieldError = formState.errors?.["asset_ids"];

  return (
    <FormControl>
      <Controller
        render={({ field }) => {
          function handleUpload(a: Asset) {
            handleAssetUpload(a);
            field.onChange([...(field.value ?? []), a.id]);
          }

          function handleRemove(a: Asset) {
            handleAssetRemove(a);
            field.onChange(pull(field.value ?? [], a.id));
          }

          return (
            <EditableAssetWall
              editing={editing}
              onUpload={handleUpload}
              onRemove={handleRemove}
              initialAssets={initialAssets}
            />
          );
        }}
        control={control}
        name="asset_ids"
      />

      <FormErrorText>{fieldError?.message?.toString()}</FormErrorText>
    </FormControl>
  );
}
