import { PublicProfile } from "src/api/openapi/schemas";
import { TextPostList } from "src/components/feed/text/TextPostList";
import { Unready } from "src/components/site/Unready";
import {
  Tabs,
  TabsContent,
  TabsIndicator,
  TabsList,
  TabsTrigger,
} from "src/theme/components/Tabs";

import { CollectionList } from "../CollectionList/CollectionList";
import { PostList } from "../PostList/PostList";

import { Box, VStack } from "@/styled-system/jsx";

import { useContent } from "./useContent";

export function Content(props: PublicProfile) {
  const content = useContent(props);

  if (!content.ready) return <Unready {...content.error} />;

  return (
    <VStack alignItems="start" w="full">
      <Tabs width="full" variant="line" defaultValue="posts">
        <TabsList>
          <TabsTrigger value="posts">Posts</TabsTrigger>
          <TabsTrigger value="replies">Replies</TabsTrigger>
          <TabsTrigger value="collections">Collections</TabsTrigger>
          <TabsIndicator />
        </TabsList>

        <TabsContent value="posts">
          <TextPostList
            posts={content.data.threads.threads}
            onDelete={content.handlers.handleDeleteThread}
          />
        </TabsContent>

        <TabsContent value="replies">
          <Box>
            <PostList posts={content.data.posts} />
          </Box>
        </TabsContent>

        <TabsContent value="collections">
          <Box>
            <CollectionList collections={content.data.collections} />
          </Box>
        </TabsContent>
      </Tabs>
    </VStack>
  );
}
