import { PublicProfile } from "src/api/openapi/schemas";
import { useContent } from "./useContent";
import { Tabs, TabList, Tab, TabPanels, TabPanel, Box } from "@chakra-ui/react";
import { ThreadList } from "src/screens/home/components/ThreadList";
import { Unready } from "src/components/Unready";
import { PostListView } from "src/screens/thread/components/PostList";

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
            <PostListView posts={content.data.posts} />
          </Box>
        </TabPanel>
      </TabPanels>
    </Tabs>
  );
}
