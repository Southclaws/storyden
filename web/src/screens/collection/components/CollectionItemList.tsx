import { Divider, List } from "@chakra-ui/react";
import { Fragment } from "react";

import { CollectionItem as CollectionItemSchema } from "src/api/openapi/schemas";

import { CollectionItemListItem } from "./CollectionItemListItem";

type Props = { items: CollectionItemSchema[] };

export function CollectionItemList(props: Props) {
  return (
    <List width="full" display="flex" flexDirection="column">
      {props.items.map((t) => (
        <Fragment key={t.id}>
          <Divider />
          <CollectionItemListItem key={t.id} item={t} />
        </Fragment>
      ))}
    </List>
  );
}
