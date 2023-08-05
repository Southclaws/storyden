import { Divider, OrderedList, VStack } from "@chakra-ui/react";
import { Fragment } from "react";

import { Collection } from "src/api/openapi/schemas";

import { useProfileContext } from "../../context";

import { CollectionCreate } from "./CollectionCreate/CollectionCreate";
import { CollectionListItem } from "./CollectionListItem";

type Props = {
  collections: Collection[];
};
export function CollectionList(props: Props) {
  const { isSelf } = useProfileContext();

  return (
    <VStack alignItems="start">
      {/* TODO: Actually design this lol */}
      {isSelf && <CollectionCreate />}

      <OrderedList gap={4} display="flex" flexDir="column" width="full" m={0}>
        {props.collections.map((c) => (
          <Fragment key={c.id}>
            <Divider />
            <CollectionListItem {...c} />
          </Fragment>
        ))}
      </OrderedList>
    </VStack>
  );
}
