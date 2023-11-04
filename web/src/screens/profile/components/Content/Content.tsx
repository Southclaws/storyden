import {
  Box,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  VStack,
} from "@chakra-ui/react";

import { PublicProfile } from "src/api/openapi/schemas";
import { MixedPostList } from "src/components/feed/mixed/MixedPostList";
import { Unready } from "src/components/site/Unready";

import { CollectionList } from "../CollectionList/CollectionList";
import { PostList } from "../PostList/PostList";

import { useContent } from "./useContent";

export function Content(props: PublicProfile) {
  const content = useContent(props);

  if (!content.ready) return <Unready {...content.error} />;

  return (
    <VStack alignItems="start" w="full">
      <Tabs width="full" variant="soft-rounded">
        <TabList>
          <Tab>Posts</Tab>
          <Tab>Replies</Tab>
          <Tab>Collections</Tab>
        </TabList>
        <TabPanels>
          <TabPanel>
            <MixedPostList
              posts={content.data.threads}
              onDelete={content.handlers.handleDelete}
            />
          </TabPanel>
          <TabPanel>
            <Box>
              <PostList posts={content.data.posts} />
            </Box>
          </TabPanel>
          <TabPanel>
            <Box>
              <CollectionList collections={content.data.collections} />
            </Box>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </VStack>
  );
}
