import Link from "next/link";

import { Post } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { HStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { Byline } from "../content/Byline";
import { CollectionMenu } from "../content/CollectionMenu/CollectionMenu";
import { FeedItemMenu } from "../feed/FeedItemMenu/FeedItemMenu";
import { Card } from "../ui/rich-card";

type Props = {
  post: Post;
};

export function PostCard({ post }: Props) {
  const session = useSession();
  const permalink = `/t/${post.slug}`;

  const title = post.title || "Untitled post";

  return (
    <Card
      shape="responsive"
      id={post.id}
      title={title}
      text={post.description}
      url={permalink}
      image={getAssetURL(post.assets?.[0]?.path)}
      controls={
        session && (
          <HStack>
            <CollectionMenu account={session} thread={post} />
            <FeedItemMenu thread={post} />
          </HStack>
        )
      }
    >
      <Byline
        href={permalink}
        author={post.author}
        time={new Date(post.createdAt)}
        updated={new Date(post.updatedAt)}
      />
    </Card>
  );
}
