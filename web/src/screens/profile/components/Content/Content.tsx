import { PublicProfile } from "src/api/openapi/schemas";
import { useContent } from "./useContent";
import { Tabs, TabList, Tab, TabPanels, TabPanel } from "@chakra-ui/react";
import { ThreadList } from "src/screens/home/components/ThreadList";
import { Unready } from "src/components/Unready";

export function Content(props: PublicProfile) {
  const threads = useContent(props);

  if (!threads.ready) return <Unready {...threads.error} />;

  return (
    <Tabs variant="soft-rounded">
      <TabList>
        <Tab>Threads</Tab>
        <Tab>Posts</Tab>
      </TabList>
      <TabPanels>
        <TabPanel>
          <ThreadList threads={threads.data} />
        </TabPanel>
        <TabPanel>
          <p>No posts found.</p>
        </TabPanel>
      </TabPanels>
    </Tabs>
  );
}
