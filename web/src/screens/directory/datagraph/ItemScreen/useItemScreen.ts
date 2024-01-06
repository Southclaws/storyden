import { useItemGet } from "src/api/openapi/items";
import { ItemWithParents } from "src/api/openapi/schemas";

import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  slug: string;
  item: ItemWithParents;
};

export function useItemScreen(props: Props) {
  const { data, mutate, error } = useItemGet(props.slug, {
    swr: {
      fallbackData: props.item,
    },
  });

  const directoryPath = useDirectoryPath();

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    data,
    directoryPath,
    mutate,
  };
}
