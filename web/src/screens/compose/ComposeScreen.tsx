import { Unready } from "src/components/site/Unready";
import { Skeleton } from "src/theme/components/Skeleton";

import { VStack } from "@/styled-system/jsx";

import { ComposeForm } from "./components/ComposeForm/ComposeForm";
import { Props, useComposeScreen } from "./useComposeScreen";

export function ComposeScreen(props: Props) {
  const { loadingDraft, draft } = useComposeScreen(props);

  if (loadingDraft)
    return (
      <Unready>
        <Skeleton />
      </Unready>
    );

  return (
    <VStack alignItems="start" gap="2" w="full" h="full">
      <ComposeForm {...props} initialDraft={draft} />
    </VStack>
  );
}
