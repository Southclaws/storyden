import { Collection } from "src/api/openapi/schemas";
import { Heading1 } from "src/theme/components/Heading/Index";

import { Flex, LinkBox, LinkOverlay, styled } from "@/styled-system/jsx";

export function CollectionListItem(props: Collection) {
  return (
    <LinkBox key={props.id}>
      <styled.li key={props.id} listStyleType="none">
        <Flex id={props.id} flexDir="column" gap="1">
          <LinkOverlay
            href={`/p/${props.owner.handle}/collections/${props.id}`}
          >
            <Heading1 size="md">{props.name}</Heading1>
          </LinkOverlay>
          <p>{props.description}</p>
        </Flex>
      </styled.li>
    </LinkBox>
  );
}
