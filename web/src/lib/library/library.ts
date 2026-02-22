import { uniqueId } from "lodash";
import { usePathname, useRouter } from "next/navigation";
import { toast } from "sonner";
import { MutatorCallback, useSWRConfig } from "swr";

import { linkCreate } from "@/api/openapi-client/links";
import {
  nodeCreate,
  nodeDelete,
  nodeGenerateContent,
  nodeGenerateTags,
  nodeGenerateTitle,
  nodeUpdate,
  nodeUpdatePosition,
  nodeUpdateVisibility,
} from "@/api/openapi-client/nodes";
import {
  Asset,
  Identifier,
  Node,
  NodeGetOKResponse,
  NodeListOKResponse,
  NodeUpdatePositionBody,
  NodeWithChildren,
  Visibility,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import {
  joinLibraryPath,
  replaceLibraryPath,
} from "@/screens/library/library-path";
import { useLibraryPath } from "@/screens/library/useLibraryPath";
import { slugify } from "@/utils/slugify";
import { generateXid } from "@/utils/xid";

import { useCapability } from "../settings/capabilities";

import { CoverImage } from "./metadata";
import { nodeListMutator, nodeMutator } from "./mutator-functions";
import {
  buildNodeChildrenListKey,
  buildNodeKey,
  buildNodeListKey,
  nodeListPrivateKeyFn,
} from "./mutator-keys";

export type CreateNodeArgs = {
  initialName?: string;
  parentID?: string;
  parentSlug?: string;
  disableRedirect?: boolean;
};

export type CoverImageArgs =
  | {
      /**
       * The asset to use as the cover image.
       */
      asset: Asset;

      /**
       * The configuration for the cropper, this is used to store the crop coords
       * for when the user re-enters the edit mode and loads the original image.
       */
      config: CoverImage;

      /**
       * Is this cover image a full replacement or a crop of the original?
       */
      isReplacement: false;
    }
  | {
      asset: Asset;
      isReplacement: true;
    };

// TODO: Remove slug params from API calls and use the node object instead.
export function useLibraryMutation(node?: Node) {
  const session = useSession();
  const genaiAvailable = useCapability("gen_ai");
  const { mutate } = useSWRConfig();
  const router = useRouter();
  const pathname = usePathname();
  const libraryPath = useLibraryPath();

  const isOnNodePage = node ? pathname.includes(`/l/${node?.slug}`) : false;

  const createNode = async ({
    initialName,
    parentID,
    parentSlug,
    disableRedirect,
  }: CreateNodeArgs) => {
    if (!session) return;

    // NOTE: This is a stopgap until the API deals with initial empty states in
    // a nicer way. For now we simply generate a dumb name which in turn results
    // in a unique slug. Eventually, the API should handle empty names and slugs
    // which it will generate a suitable unique mark for, like how Notion works.
    //
    // NOTE 2: We use the Xid library to generate a unique ID for the new page
    // however, the way that marks work is XID-format strings are assumed to be
    // node IDs not slugs. So we need to prefix the random name to prevent this.
    //
    const name = initialName?.trim() || `untitled`;
    const slug = slugify(`${name}-${generateXid()}`);

    const initial: NodeWithChildren = {
      id: "optimistic_node_" + uniqueId(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      name,
      slug: slug,
      description: "",
      owner: session,
      properties: [],
      child_property_schema: [],
      hide_child_tree: false,
      meta: {},
      children: [],
      assets: [],
      tags: [],
      visibility: "draft",
      recomentations: [],
    };

    const mutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      return {
        ...data,
        nodes: [initial, ...data.nodes],
      };
    };

    await mutate(nodeListPrivateKeyFn, mutator, { revalidate: false });

    const parent = parentID ?? parentSlug;
    if (parent) {
      const childListKeyFn = buildNodeChildrenListKey(parent);
      await mutate(childListKeyFn, mutator, { revalidate: false });
    }

    const created = await nodeCreate({
      name,
      slug,
      parent,
    });

    if (!disableRedirect) {
      const newPath = joinLibraryPath(libraryPath, created.slug);
      router.push(`/l/${newPath}?edit=true`);
    }
  };

  const suggestTags = async (slug: string, content: string) => {
    const { tags } = await nodeGenerateTags(slug, { content });

    return tags;
  };

  const suggestTitle = async (slug: string, content: string) => {
    const { title } = await nodeGenerateTitle(slug, { content });

    return title;
  };

  const suggestSummary = async (slug: string, currentContent: string) => {
    const { content } = await nodeGenerateContent(slug, {
      content: currentContent,
    });

    return content;
  };

  const importFromLink = async (id: string, url: string) => {
    const { title, description, primary_image } = await linkCreate({ url });

    if (genaiAvailable && description) {
      const [tag_suggestions, title_suggestion, content_suggestion] =
        await Promise.all([
          suggestTags(id, description).catch(() => undefined),
          suggestTitle(id, description).catch(() => undefined),
          suggestSummary(id, description).catch(() => undefined),
        ]);
      return {
        title_suggestion: title_suggestion || title,
        tag_suggestions: tag_suggestions || [],
        content_suggestion: content_suggestion || description,
        primary_image,
      };
    }

    return {
      title_suggestion: title,
      tag_suggestions: [] as string[],
      content_suggestion: description,
      primary_image,
    };
  };

  const updateNodeVisibility = async (
    slug: string,
    visibility: Visibility,
    parentID?: Identifier,
  ) => {
    const mutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      const newNodes = data.nodes.map((node) => {
        if (node.slug === slug) {
          if (
            node.parent &&
            visibility === Visibility.published &&
            node.parent.visibility !== Visibility.published
          ) {
            toast.warning(
              "Page is staged for publishing but has not been published yet because this page's parent is not published. When the parent is published, this page be visible on the site.",
              {
                duration: 15000,
                dismissible: true,
                closeButton: true,
              },
            );
          }
          const newNode = { ...node, visibility };
          return newNode;
        }
        return node;
      });

      return {
        ...data,
        nodes: newNodes,
      };
    };

    const keyFn = buildNodeListKey();
    await mutate(keyFn, mutator, { revalidate: false });

    // Only optimistically update non-root moves.
    if (parentID) {
      const childListKeyFn = buildNodeChildrenListKey(parentID);
      await mutate(childListKeyFn, mutator, { revalidate: false });
    }

    await nodeUpdateVisibility(slug, { visibility });
  };

  const updateNodeChildVisibility = async (
    slug: string,
    hideChildTree: boolean,
  ) => {
    const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
      if (!data) return;

      const updated = {
        ...data,
        hide_child_tree: hideChildTree,
      } satisfies NodeWithChildren;

      return updated;
    };

    const nodeListMutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      const newNodes = data.nodes.map((node) => {
        if (node.slug === slug) {
          const newNode = { ...node, hide_child_tree: hideChildTree };
          return newNode;
        }
        return node;
      });

      return {
        ...data,
        nodes: newNodes,
      };
    };

    const listKeyFn = buildNodeListKey();
    await mutate(listKeyFn, nodeListMutator, { revalidate: false });

    const nodeKeyFn = buildNodeKey(slug);
    await mutate(nodeKeyFn, nodeMutator, { revalidate: false });

    const updated = await nodeUpdate(slug, {
      hide_child_tree: hideChildTree,
    });

    revalidate();

    return updated;
  };

  const deleteNode = async (
    slug: string,
    oldParent?: string,
    newParent?: string,
  ) => {
    const mutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      const newNodes = data.nodes.filter((node) => node.slug !== slug);

      return {
        ...data,
        nodes: newNodes,
      };
    };

    const listKeyFn = buildNodeListKey();
    await mutate(listKeyFn, mutator, { revalidate: false });

    if (oldParent) {
      const childListKeyFn = buildNodeChildrenListKey(oldParent);
      await mutate(childListKeyFn, mutator, { revalidate: false });
    }

    await nodeDelete(slug, { target_node: newParent });

    if (isOnNodePage) {
      if (newParent) {
        const newPath = replaceLibraryPath(libraryPath, slug, newParent);
        router.push(newPath);
      } else {
        router.push("/l");
      }
    }
  };

  const moveNode = async (
    draggingNodeId: string,
    dropTargetId: string,
    dropPosition: "above" | "below" | "inside",
    newParent:
      | string // Set a new parent
      | null // Set no parent (root)
      | undefined, // Keep current parent
    oldParent:
      | Identifier // node being moved has a parent
      | undefined, // node being moved is root
  ) => {
    const mutator: MutatorCallback<NodeListOKResponse> = (prevData) => {
      if (!prevData) return prevData;

      const newNodes = moveNodeInTree({
        tree: prevData.nodes,
        draggingNodeId,
        dropTargetId,
        dropPosition,
      });

      return { ...prevData, nodes: newNodes };
    };

    const listKeyFn = buildNodeListKey();
    await mutate(listKeyFn, mutator, { revalidate: false });

    // Only optimistically update non-root moves.
    if (oldParent) {
      const childListKeyFn = buildNodeChildrenListKey(oldParent);
      await mutate(childListKeyFn, mutator, { revalidate: false });
    }

    const params: NodeUpdatePositionBody = (() => {
      switch (dropPosition) {
        case "above":
          return {
            before: dropTargetId,
            parent: newParent,
          };

        case "below":
          return {
            after: dropTargetId,
            parent: newParent,
          };

        case "inside":
          return {
            parent: dropTargetId,
          };
      }
    })();

    await nodeUpdatePosition(draggingNodeId, params);
  };

  const revalidate = async (updated?: NodeWithChildren) => {
    const listKeyFn = buildNodeListKey();
    await mutate<NodeListOKResponse>(
      listKeyFn,
      updated ? nodeListMutator(updated) : undefined,
    );

    if (node) {
      const nodeKeyFn = buildNodeKey(updated?.slug ?? node.slug);
      await mutate<NodeGetOKResponse>(
        nodeKeyFn,
        updated ? nodeMutator(updated) : undefined,
      );
    }
  };

  return {
    createNode,
    suggestTitle,
    suggestSummary,
    suggestTags,
    importFromLink,
    updateNodeVisibility,
    updateNodeChildVisibility,
    deleteNode,
    moveNode,
    revalidate,
  };
}

interface MoveNodeParams {
  tree: NodeWithChildren[];
  draggingNodeId: string;
  dropTargetId: string;
  dropPosition: "above" | "below" | "inside";
}

function moveNodeInTree({
  tree,
  draggingNodeId,
  dropTargetId,
  dropPosition,
}: MoveNodeParams): NodeWithChildren[] {
  let draggedNode: NodeWithChildren | null = null;

  function removeNode(nodes: NodeWithChildren[]): NodeWithChildren[] {
    return nodes.reduce<NodeWithChildren[]>((acc, node) => {
      if (node.id === draggingNodeId) {
        draggedNode = { ...node };
        return acc;
      }

      const newChildren = removeNode(node.children || []);
      acc.push({ ...node, children: newChildren });
      return acc;
    }, []);
  }

  function insertNode(nodes: NodeWithChildren[]): NodeWithChildren[] {
    return nodes.reduce<NodeWithChildren[]>((acc, node) => {
      if (node.id === dropTargetId) {
        if (dropPosition === "inside") {
          // Insert as first child
          const newChildren = [draggedNode!, ...(node.children || [])];
          acc.push({ ...node, children: newChildren });
        } else {
          // Insert above or below at sibling level
          if (dropPosition === "above") {
            acc.push(draggedNode!);
            acc.push(node);
          } else if (dropPosition === "below") {
            acc.push(node);
            acc.push(draggedNode!);
          }
        }
      } else {
        // Normal node
        const newChildren = insertNode(node.children || []);
        acc.push({ ...node, children: newChildren });
      }

      return acc;
    }, []);
  }

  const treeWithoutDragged = removeNode(tree);

  if (!draggedNode) {
    console.warn("Dragged node not found");
    return tree;
  }

  const newTree = insertNode(treeWithoutDragged);

  return newTree;
}
