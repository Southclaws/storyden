import { useCollapsible } from "@ark-ui/react";
import { useState } from "react";

import { useNodeCommentList } from "@/api/openapi-client/nodes";
import { NodeWithChildren } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import * as Collapsible from "@/components/ui/collapsible";
import { Heading } from "@/components/ui/heading";
import { ChevronRightIcon } from "@/components/ui/icons/Chevron";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { Input } from "@/components/ui/input";
import { Box, HStack, LStack, WStack } from "@/styled-system/jsx";
import { wstack } from "@/styled-system/patterns";

type Props = {
  node: NodeWithChildren;
};
export function LibraryPageCommentsList(props: Props) {
  const collapsible = useCollapsible({
    lazyMount: true,
    unmountOnExit: true,
    onOpenChange(details) {
      console.log("Collapsible open state changed:", details);
    },
  });

  const { data, error } = useNodeCommentList(props.node.slug);
  if (!data) {
    return <Unready error={error} />;
  }

  function handleToggle() {
    collapsible.setOpen(!collapsible.open);
  }

  const pageSize = 50; // Echoes what's in node_comments.go
  const commentCount =
    data.total_pages > 1 ? `${data.total_pages * pageSize}+` : data.results;
  const commentCountLabel = commentCount === 1 ? "comment" : "comments";

  return (
    <LStack>
      <WStack>
        <Heading
          fontSize="sm"
          color="fg.muted"
          display="flex"
          gap="1"
          alignItems="center"
        >
          <DiscussionIcon w="4" />
          {commentCount} {commentCountLabel} on {props.node.name}
        </Heading>
        <Button type="button" size="xs" variant="subtle" onClick={handleToggle}>
          {collapsible.open ? "Hide" : "Show"}
        </Button>
      </WStack>

      <Collapsible.RootProvider value={collapsible}>
        <Collapsible.Content>
          <LStack>
            <NodeCommentCreateForm node={props.node} />

            <ThreadReferenceList threads={data.threads} />
          </LStack>
        </Collapsible.Content>
      </Collapsible.RootProvider>
    </LStack>
  );
}

function NodeCommentCreateForm(props: { node: NodeWithChildren }) {
  return (
    <HStack w="full">
      <Input size="sm" placeholder="Leave a comment..." />
      <Button size="sm" variant="subtle">
        Post
      </Button>
    </HStack>
  );
}
