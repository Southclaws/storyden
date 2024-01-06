import { slice } from "lodash/fp";

import { Item } from "src/api/openapi/schemas";
import { DirectoryPath } from "src/screens/directory/datagraph/useDirectoryPath";

import { Grid } from "@/styled-system/jsx";

import { ItemCard } from "./ItemCard";

type Props = {
  items: Item[];
  directoryPath: DirectoryPath;
};

export function ItemGrid(props: Props) {
  const items = slice(0, 4, props.items);

  return (
    <Grid
      w="full"
      gridTemplateColumns={{
        base: "2",
        sm: "3",
        md: "3",
        lg: "4",
      }}
      gridTemplateRows="1"
    >
      {items.map((item) => (
        <ItemCard
          key={item.id}
          directoryPath={props.directoryPath}
          item={item}
        />
      ))}
    </Grid>
  );
}
