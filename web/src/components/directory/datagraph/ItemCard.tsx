import { Item } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/directory-path";
import { Card } from "src/theme/components/Card";

import { DirectoryBadge } from "../DirectoryBadge";

import { HStack } from "@/styled-system/jsx";
import { CardVariantProps } from "@/styled-system/recipes";

export type Props = {
  item: Item;
  directoryPath: DirectoryPath;
} & CardVariantProps;

export function ItemCard({ item, directoryPath, ...rest }: Props) {
  const slug = joinDirectoryPath(directoryPath, item.slug);
  const asset = item.assets?.[0];
  const url = `/directory/${slug}`;

  return (
    <Card
      id={item.id}
      title={item.name}
      text={item.description}
      url={url}
      image={asset?.url}
      {...rest}
    >
      <HStack color="fg.muted">
        <DirectoryBadge />

        <Timestamp
          created={item.createdAt}
          updated={item.updatedAt}
          href={url}
        />
      </HStack>
    </Card>
  );
}
