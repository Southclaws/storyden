import { PublicProfile } from "src/api/openapi-schema";
import { Unready } from "src/components/site/Unready";

import { ThreadItemList } from "@/components/feed/ThreadItemList";
import * as Tabs from "@/components/ui/tabs";
import { Box, VStack } from "@/styled-system/jsx";

import { CollectionList } from "../CollectionList/CollectionList";

import { useContent } from "./useContent";

export function Content(props: PublicProfile) {
  const content = useContent(props);

  if (!content.ready) return <Unready {...content.error} />;

  return (
    <VStack alignItems="start" w="full">
      <Tabs.Root width="full" variant="line" defaultValue="posts">
        <Tabs.List>
          <Tabs.Trigger value="posts">Posts</Tabs.Trigger>
          <Tabs.Trigger value="collections">Collections</Tabs.Trigger>
          <Tabs.Indicator />
        </Tabs.List>

        <Tabs.Content value="posts">
          <ThreadItemList threads={content.data.threads} />
        </Tabs.Content>

        <Tabs.Content value="collections">
          <Box>
            <CollectionList collections={content.data.collections} />
          </Box>
        </Tabs.Content>
      </Tabs.Root>
    </VStack>
  );
}
