import { PostReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Byline } from "@/components/content/Byline";
import { CollectionMenu } from "@/components/content/CollectionMenu/CollectionMenu";
import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";

type Props = {
  item: PostReference;
};

export function PostRef({ item }: Props) {
  const session = useSession();

  const permalink = `/t/${item.slug}#${item.id}`;

  return (
    <Card
      id={item.id}
      title={item.title}
      text={item.description}
      url={permalink}
      shape="row"
      controls={
        session && (
          <HStack>
            <CollectionMenu account={session} thread={item} />
          </HStack>
        )
      }
    >
      <Byline
        href={permalink}
        author={item.author}
        time={new Date(item.createdAt)}
        updated={new Date(item.updatedAt)}
      />
    </Card>
  );
}
