import { Skeleton, SkeletonText, VStack } from "@chakra-ui/react";

import { Unready } from "src/components/Unready";

import { ComposeForm } from "./components/ComposeForm/ComposeForm";
import { Props, useComposeScreen } from "./useComposeScreen";

export function ComposeScreen(props: Props) {
  const { loadingDraft, draft } = useComposeScreen(props);

  if (loadingDraft)
    return (
      <Unready>
        <Skeleton height={8} />
        <SkeletonText noOfLines={3} />
      </Unready>
    );

  return (
    <VStack alignItems="start" gap={2} w="full" h="full" py={5}>
      <ComposeForm {...props} initialDraft={draft} />
    </VStack>
  );
}
