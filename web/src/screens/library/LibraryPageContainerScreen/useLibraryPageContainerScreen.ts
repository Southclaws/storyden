import { useRouter } from "next/navigation";

import {
  nodeDelete,
  nodeUpdate,
  nodeUpdateVisibility,
  useNodeGet,
} from "src/api/openapi-client/nodes";
import {
  Node,
  NodeInitialProps,
  NodeWithChildren,
  Visibility,
} from "src/api/openapi-schema";

import { replaceLibraryPath } from "../library-path";
import { useLibraryPath } from "../useLibraryPath";

export type Props = {
  slug: string;
  node: NodeWithChildren;
};

export function useLibraryPageContainerScreen(props: Props) {
  const router = useRouter();
  const { data, mutate, error } = useNodeGet(props.slug, {
    swr: {
      fallbackData: props.node,
    },
  });

  const libraryPath = useLibraryPath();

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const { slug } = data;

  async function handleVisibilityChange(visibility: Visibility) {
    await nodeUpdateVisibility(slug, { visibility });
    await mutate();
  }

  async function handleSave(node: NodeInitialProps) {
    const newNode = await nodeUpdate(slug, {
      name: node.name,
      slug: node.slug,
      asset_ids: node.asset_ids,
      url: node.url,
      content: node.content,
      meta: node.meta,
    });
    await mutate();

    // Handle slug changes properly by redirecting to the new path.
    if (newNode.slug !== slug) {
      const newPath = replaceLibraryPath(libraryPath, slug, newNode.slug);
      router.push(newPath);
    }
  }

  // TODO: Provide a way to set the new parent node for child nodes/items.
  async function handleDelete(node: Node) {
    const { destination } = await nodeDelete(node.slug);

    if (destination) {
      const newPath = replaceLibraryPath(libraryPath, slug, destination.slug);
      router.push(newPath);
    } else {
      router.push("/l");
    }
  }

  return {
    ready: true as const,
    data,
    handlers: { handleSave, handleVisibilityChange, handleDelete },
    libraryPath,
    mutate,
  };
}
