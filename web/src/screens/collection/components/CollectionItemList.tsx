import { Fragment } from "react";

import { CollectionItem as CollectionItemSchema } from "src/api/openapi-schema";

import { Divider, styled } from "@/styled-system/jsx";

import { CollectionItemListItem } from "./CollectionItemListItem";

type Props = { items: CollectionItemSchema[] };

export function CollectionItemList(props: Props) {
  return (
    <styled.ul width="full" display="flex" flexDirection="column">
      {props.items.map((t) => (
        <Fragment key={t.id}>
          <Divider />
          <CollectionItemListItem key={t.id} item={t} />
        </Fragment>
      ))}
    </styled.ul>
  );
}
