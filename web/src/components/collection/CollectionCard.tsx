import { BoxesIcon } from "lucide-react";

import { Collection } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";

import { CollectionMenu } from "./CollectionMenu/CollectionMenu";

type Props = {
  collection: Collection;
  hideOwner?: boolean;
};

export function CollectionCard({ collection, hideOwner }: Props) {
  const url = `/c/${collection.id}`;

  return (
    <Card
      key={collection.id}
      id={collection.id}
      url={url}
      shape="responsive"
      title={collection.name}
      text={collection.description}
      controls={
        <HStack w="full" justify="space-between">
          <HStack>
            {hideOwner ? null : (
              <MemberBadge profile={collection.owner} size="sm" name="handle" />
            )}

            <CollectionItems collection={collection} />
          </HStack>

          <HStack>
            <CollectionMenu collection={collection} />
          </HStack>
        </HStack>
      }
    />
  );
}

function CollectionItems(props: Props) {
  const itemsLabel = props.collection.item_count === 1 ? "item" : "items";
  return (
    <Badge size="sm">
      <BoxesIcon width="1.4rem" />{" "}
      <span>
        {props.collection.item_count} {itemsLabel}
      </span>
    </Badge>
  );
}
