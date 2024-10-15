import { FileUploadFileAcceptDetails } from "@ark-ui/react";

import { Asset } from "@/api/openapi-schema";
import { ImageEditor } from "@/components/site/ImageEditor/ImageEditor";
import { LStack } from "@/styled-system/jsx";

export type Props = {
  initialValue?: Asset;
  onFinish: (a: Asset) => Promise<void>;
};

export function useMediaEditScreen(props: Props) {
  //
}

export function MediaEditScreen(props: Props) {
  useMediaEditScreen(props);

  function handleUpload(asset: Asset) {
    props.onFinish(asset);
  }

  return (
    <LStack w="breakpoint-lg" maxW="breakpoint-lg">
      <ImageEditor onUpload={handleUpload} initialValue={props.initialValue} />
    </LStack>
  );
}
