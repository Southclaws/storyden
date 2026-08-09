import { useFormContext } from "react-hook-form";

import { Asset } from "@/api/openapi-schema";
import { ContentComposerField } from "@/components/content/ContentComposer";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";

import { FormShape } from "../ComposeForm/useComposeForm";

type Props = {
  onAssetUpload: (asset: Asset) => void;
};

export function BodyInput({ onAssetUpload }: Props) {
  const form = useFormContext<FormShape>();

  return (
    <FormControl h="full">
      <ContentComposerField
        onAssetUpload={onAssetUpload}
        control={form.control}
        name="body"
      />
      <FormErrorText>{form.formState.errors.body?.message}</FormErrorText>
    </FormControl>
  );
}
