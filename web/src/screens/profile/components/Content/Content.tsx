import { Box, Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/react";
import { PublicProfile } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";
import { ThreadList } from "src/screens/home/components/ThreadList";
import { PostList } from "./PostList/PostList";
import { useContent } from "./useContent";

export function Content(props: PublicProfile) {
  const content = useContent(props);

  if (!content.ready) return <Unready {...content.error} />;

  return (
    <Tabs width="full" variant="soft-rounded">
      <TabList>
        <Tab>Threads</Tab>
        <Tab>Posts</Tab>
      </TabList>
      <TabPanels>
        <TabPanel>
          <ThreadList threads={content.data.threads} />
        </TabPanel>
        <TabPanel>
          <Box>
            <PostList posts={content.data.posts} />
          </Box>
        </TabPanel>
      </TabPanels>
    </Tabs>
  );
}
