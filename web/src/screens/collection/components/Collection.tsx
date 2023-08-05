import { Heading, Text } from "@chakra-ui/react";

import { Collection, CollectionWithItems } from "src/api/openapi/schemas";
import { Byline } from "src/screens/thread/components/Byline";
import { CollectionItemList } from "./CollectionItemList";

export function Collection(props: CollectionWithItems) {
  return (
    <>
      <Heading>{props.name}</Heading>
      <Byline
        author={props.owner.handle}
        time={new Date(props.createdAt)}
        updated={new Date(props.updatedAt)}
        href={`/p/${props.owner.handle}/collections/${props.id}`}
        />
      <Text>{props.description}</Text>

      <CollectionItemList items={props.items} />
    </>
  );
}
