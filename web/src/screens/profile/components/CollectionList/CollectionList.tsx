import { Fragment } from "react";

import { Collection } from "src/api/openapi/schemas";
import { CollectionCreateTrigger } from "src/components/content/CollectionCreate/CollectionCreateTrigger";

import { useProfileContext } from "../../context";

import { Divider, VStack, styled } from "@/styled-system/jsx";

import { CollectionListItem } from "./CollectionListItem";

type Props = {
  collections: Collection[];
};
export function CollectionList(props: Props) {
  const { isSelf } = useProfileContext();

  return (
    <VStack alignItems="start">
      {/* TODO: Actually design this lol */}
      {isSelf && <CollectionCreateTrigger />}

      <styled.ol gap="4" display="flex" flexDir="column" width="full" m="0">
        {props.collections.map((c) => (
          <Fragment key={c.id}>
            <Divider />
            <CollectionListItem {...c} />
          </Fragment>
        ))}
      </styled.ol>
    </VStack>
  );
}
