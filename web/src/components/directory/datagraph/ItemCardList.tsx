import { slice } from "lodash/fp";

import { Item } from "src/api/openapi/schemas";
import { DirectoryPath } from "src/screens/directory/datagraph/useDirectoryPath";
import { CardGrid } from "src/theme/components/Card";

import { ItemCard } from "./ItemCard";

type Props = {
  items: Item[];
  directoryPath: DirectoryPath;
};

export function ItemCardGrid(props: Props) {
  const items = slice(0, 4, props.items);

  return (
    <CardGrid>
      {items.map((item) => (
        <ItemCard
          key={item.id}
          directoryPath={props.directoryPath}
          item={item}
        />
      ))}
    </CardGrid>
  );
}
