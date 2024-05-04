import { useRouter } from "next/navigation";

import { itemDelete, itemUpdate, useItemGet } from "src/api/openapi/items";
import {
  Item,
  ItemInitialProps,
  ItemWithParents,
} from "src/api/openapi/schemas";

import { replaceDirectoryPath } from "../directory-path";
import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  slug: string;
  item: ItemWithParents;
};

export function useItemViewerScreen(props: Props) {
  const router = useRouter();
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

  const { slug } = data;

  async function handleSave(item: ItemInitialProps) {
    await itemUpdate(slug, {
      name: item.name,
      slug: item.slug,
      asset_ids: item.asset_ids,
      url: item.url,
      description: item.description,
      content: item.content,
      properties: item.properties,
    });
    await mutate();

    // Handle slug changes properly by redirecting to the new path.
    if (item.slug !== slug) {
      const newPath = replaceDirectoryPath(directoryPath, slug, item.slug);
      router.push(newPath);
    }
  }

  async function handleDelete(item: Item) {
    await itemDelete(item.slug);
    // TODO: redirect
  }

  return {
    ready: true as const,
    data,
    handlers: { handleSave, handleDelete },
    mutate,
  };
}
