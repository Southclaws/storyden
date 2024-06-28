import { useRouter } from "next/navigation";

import {
  nodeDelete,
  nodeUpdate,
  nodeUpdateVisibility,
  useNodeGet,
} from "src/api/openapi/nodes";
import {
  Node,
  NodeInitialProps,
  NodeWithChildren,
  Visibility,
} from "src/api/openapi/schemas";

import { replaceDirectoryPath } from "../directory-path";
import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  slug: string;
  node: NodeWithChildren;
};

export function useNodeViewerScreen(props: Props) {
  const router = useRouter();
  const { data, mutate, error } = useNodeGet(props.slug, {
    swr: {
      fallbackData: props.node,
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
      properties: node.properties,
    });
    await mutate();

    // Handle slug changes properly by redirecting to the new path.
    if (newNode.slug !== slug) {
      const newPath = replaceDirectoryPath(directoryPath, slug, newNode.slug);
      router.push(newPath);
    }
  }

  // TODO: Provide a way to set the new parent node for child nodes/items.
  async function handleDelete(node: Node) {
    const { destination } = await nodeDelete(node.slug);

    if (destination) {
      const newPath = replaceDirectoryPath(
        directoryPath,
        slug,
        destination.slug,
      );
      router.push(newPath);
    } else {
      router.push("/directory");
    }
  }

  return {
    ready: true as const,
    data,
    handlers: { handleSave, handleVisibilityChange, handleDelete },
    directoryPath,
    mutate,
  };
}
