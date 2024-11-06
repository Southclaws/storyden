import { Collection } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/rich-card";
import { HStack, WStack } from "@/styled-system/jsx";

import { CollectionIcon } from "../ui/icons/Collection";

import { CollectionMenu } from "./CollectionMenu/CollectionMenu";

type Props = {
  collection: Collection;
  hideOwner?: boolean;
};

export function CollectionCard({ collection, hideOwner }: Props) {
  const url = `/c/${collection.slug}`;

  return (
    <Card
      key={collection.id}
      id={collection.id}
      url={url}
      shape="responsive"
      title={collection.name}
      text={collection.description}
      controls={
        <WStack>
          <HStack>
            {hideOwner ? null : (
              <MemberBadge profile={collection.owner} size="sm" name="handle" />
            )}

            <CollectionItems collection={collection} />
          </HStack>

          <HStack>
            <CollectionMenu collection={collection} />
          </HStack>
        </WStack>
      }
    />
  );
}

function CollectionItems(props: Props) {
  const itemsLabel = props.collection.item_count === 1 ? "item" : "items";
  return (
    <Badge size="sm">
      <CollectionIcon />{" "}
      <span>
        {props.collection.item_count} {itemsLabel}
      </span>
    </Badge>
  );
}
