import Link from "next/link";

import { CollectionItem } from "src/api/openapi/schemas";
import { Byline } from "src/components/content/Byline";
import { Heading1 } from "src/theme/components/Heading/Index";

import { Flex, styled } from "@/styled-system/jsx";

type Props = { item: CollectionItem };

export function CollectionItemListItem(props: Props) {
  const permalink = `/t/${props.item.slug}`;

  return (
    <styled.section display="flex" flexDir="column" py="2" width="full" gap="2">
      <styled.article>
        <Flex justifyContent="space-between">
          <Heading1 size="sm">
            <Link href={permalink}>{props.item.title}</Link>
          </Heading1>
        </Flex>

        <styled.p lineClamp={3}>{props.item.short}</styled.p>
      </styled.article>

      <Flex justifyContent="space-between">
        <Byline
          href={permalink}
          author={props.item.author}
          time={new Date(props.item.createdAt)}
          updated={new Date(props.item.updatedAt)}
        />
      </Flex>
    </styled.section>
  );
}
