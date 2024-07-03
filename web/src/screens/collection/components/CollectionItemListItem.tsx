import Link from "next/link";

import { CollectionItem } from "src/api/openapi/schemas";
import { Byline } from "src/components/content/Byline";

import { Heading } from "@/components/ui/heading";
import { Flex, styled } from "@/styled-system/jsx";

type Props = { item: CollectionItem };

export function CollectionItemListItem(props: Props) {
  const permalink = `/t/${props.item.slug}`;

  return (
    <styled.section display="flex" flexDir="column" py="2" width="full" gap="2">
      <styled.article>
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <Link href={permalink}>{props.item.name}</Link>
          </Heading>
        </Flex>

        <styled.p lineClamp={3}>{props.item.description}</styled.p>
      </styled.article>

      <Flex justifyContent="space-between">
        <Byline
          href={permalink}
          author={props.item.owner}
          time={new Date(props.item.createdAt)}
          updated={new Date(props.item.updatedAt)}
        />
      </Flex>
    </styled.section>
  );
}
