import slugify from "@sindresorhus/slugify";
import { dequal } from "dequal/lite";
import { omit, uniqueId } from "lodash";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";
import { Xid } from "xid-ts";

import {
  getNodeGetKey,
  getNodeListKey,
  nodeAddAsset,
  nodeCreate,
  nodeDelete,
  nodeRemoveAsset,
  nodeUpdate,
  nodeUpdatePosition,
  nodeUpdateVisibility,
} from "@/api/openapi-client/nodes";
import {
  Asset,
  AssetID,
  Node,
  NodeGetOKResponse,
  NodeListOKResponse,
  NodeMutableProps,
  NodeUpdatePositionBody,
  NodeWithChildren,
  Property,
  PropertyMutation,
  PropertyType,
  TagReference,
  Visibility,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import {
  joinLibraryPath,
  replaceLibraryPath,
} from "@/screens/library/library-path";
import { useLibraryPath } from "@/screens/library/useLibraryPath";

import { CoverImage, NodeMetadata } from "./metadata";

export type CreateNodeArgs = {
  initialName?: string;
  parentSlug?: string;
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
  const { mutate } = useSWRConfig();
  const router = useRouter();
  const libraryPath = useLibraryPath();

  // for revalidating all node list queries (published and private)
  const nodeListKey = getNodeListKey();
  const nodeListAllKeyFn = (key: Arguments) => {
    return (
      Array.isArray(key) &&
      key[0].startsWith(nodeListKey) &&
      // Don't pass for /nodes/<slug> keys
      // NOTE: This may be buggy for cases with query params...
      !key[0].startsWith(nodeListKey + "/")
    );
  };
  // for revalidating only private node list queries
  const nodeListPrivateKey = getNodeListKey({
    // NOTE: The order here matters.
    visibility: [Visibility.draft, Visibility.review, Visibility.unlisted],
  });
  const nodeListPrivateKeyFn = (key: Arguments) => {
    return dequal(key, nodeListPrivateKey);
  };

  // For revalidating one specific node.
  const nodeKey = node && getNodeGetKey(node.slug);
  const nodeKeyFn =
    node &&
    ((key: Arguments) => {
      return Array.isArray(key) && key[0].startsWith(nodeKey);
    });

  const createNode = async ({ initialName, parentSlug }: CreateNodeArgs) => {
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
    const name = initialName ?? `untitled-${new Xid().toString()}`;

    const initial: NodeWithChildren = {
      id: "optimistic_node_" + uniqueId(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      name,
      slug: slugify(name),
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
      tag_suggestions: [],
    };

    const mutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      return {
        ...data,
        nodes: [initial, ...data.nodes],
      };
    };

    await mutate(nodeListPrivateKeyFn, mutator, { revalidate: false });

    const created = await nodeCreate({ name: name, parent: parentSlug });
    const newPath = joinLibraryPath(libraryPath, created.slug);

    router.push(`/l/${newPath}?edit=true`);
  };

  const updateNode = async (
    slug: string,
    newNode: NodeMutableProps,
    cover?: CoverImageArgs,
  ) => {
    // if moving from hidden children to displayed children, the actual data is
    // not present in the swr cache, so we need to trigger a full revalidation.
    const nonOptimisticMutation =
      newNode.hide_child_tree == false && node?.hide_child_tree === true;

    const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
      if (!data) return;

      const nodeProps = omit(newNode, "parent");

      const withProperties = {
        properties:
          newNode.properties?.map((p: PropertyMutation) => {
            return {
              fid: p.fid ?? uniqueId("new_field_"),
              sort: p.sort ?? "",
              type: p.type ?? PropertyType.text,
              ...p,
            } satisfies Property;
          }) ?? data.properties,
      };

      const withNewCover = cover?.asset && mergePrimaryImageAsset(data, cover);

      const withTags = {
        tags:
          newNode.tags?.map(
            (t) =>
              ({
                name: t,
                colour: "white",
                item_count: 1,
              }) satisfies TagReference,
          ) ?? [],
      };

      const withHiddenChildren = {
        children: newNode.hide_child_tree ? [] : data.children,
      };

      const updated = {
        ...data,
        ...nodeProps,
        ...withProperties,
        ...withNewCover,
        ...withTags,
        ...withHiddenChildren,
      } satisfies NodeWithChildren;

      return updated;
    };

    const listMutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data || !data.nodes) return;

      const newNodes = data.nodes.map((n) => {
        if (n.slug === slug) {
          const updated = nodeMutator(n);
          return { ...n, ...updated };
        }

        return n;
      });

      return {
        ...data,
        nodes: newNodes,
      };
    };

    const slugChanged = newNode.slug !== undefined && newNode.slug !== slug;

    await mutate(nodeListAllKeyFn, listMutator, { revalidate: false });
    await mutate(nodeKeyFn, nodeMutator, { revalidate: false });

    const newMeta = {
      ...node?.meta,
      ...newNode.meta,
      ...(cover && !cover.isReplacement
        ? { coverImage: cover.config }
        : undefined),
    };
    await nodeUpdate(slug, {
      ...newNode,
      primary_image_asset_id: cover?.asset.id,
      // NOTE: We don't have access to the original node's meta, so we have to
      // fully replace it. Right now no other features use metadata, but this
      // will need to be fixed eventually. Probably by either calling the API
      // within this hook to fetch the latest version of the node and spreading.
      meta: newMeta,
    });

    // Handle slug changes properly by redirecting to the new path.
    if (slugChanged && newNode.slug /* Needed for TS narrowing */) {
      const newPath = replaceLibraryPath(libraryPath, slug, newNode.slug);
      router.push(newPath);
    }

    if (nonOptimisticMutation) {
      await revalidate();
    }

    return slugChanged;
  };

  const suggestTags = async (slug: string) => {
    const { tag_suggestions } = await nodeUpdate(
      slug,
      {},
      { tag_fill_rule: "query" },
    );

    return tag_suggestions;
  };

  const suggestTitle = async (slug: string) => {
    const { title_suggestion } = await nodeUpdate(
      slug,
      {},
      { title_fill_rule: "query" },
    );

    return title_suggestion;
  };

  const suggestSummary = async (slug: string) => {
    const { title_suggestion } = await nodeUpdate(
      slug,
      {},
      { content_fill_rule: "query" },
    );

    return title_suggestion;
  };

  const importFromLink = async (slug: string, url: string) => {
    const { title_suggestion, tag_suggestions, content_suggestion } =
      await nodeUpdate(
        slug,
        { url },
        {
          // generate all these fields from the provided URL:
          fill_source: "url",
          // title_fill_rule: "query",
          tag_fill_rule: "query",
          content_fill_rule: "query",
        },
      );

    return { title_suggestion, tag_suggestions, content_suggestion };
  };

  const removeNodeCoverImage = async (slug: string) => {
    const nodeKey = getNodeGetKey(slug);
    const nodeKeyFn = (key: Arguments) => {
      return Array.isArray(key) && key[0].startsWith(nodeKey);
    };

    const mutator: MutatorCallback<NodeGetOKResponse> = (data) => {
      if (!data) return;

      const newNode = { ...data, primary_image: undefined };

      return newNode;
    };

    await mutate(nodeKeyFn, mutator, { revalidate: false });

    await nodeUpdate(slug, {
      primary_image_asset_id: null,
      // NOTE: We don't have access to the original node's meta, so we can't
      // remove the old cover config, but it doesn't really matter.
    });
  };

  const updateNodeVisibility = async (slug: string, visibility: Visibility) => {
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

    await mutate(nodeListAllKeyFn, mutator, { revalidate: false });

    await nodeUpdateVisibility(slug, { visibility });
  };

  const addAsset = async (slug: string, asset: Asset) => {
    const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
      if (!data) return;

      const assets = [...data.assets, asset];

      const updated = {
        ...data,
        assets,
      } satisfies NodeWithChildren;

      return updated;
    };

    const nodeKey = getNodeGetKey(slug);
    const nodeKeyFn = (key: Arguments) => {
      return Array.isArray(key) && key[0].startsWith(nodeKey);
    };

    await mutate(nodeKeyFn, nodeMutator, { revalidate: false });

    await nodeAddAsset(slug, asset.id);
  };

  const removeAsset = async (slug: string, assetID: AssetID) => {
    const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
      if (!data) return;

      const assets = data.assets.filter((a) => a.id !== assetID);

      const updated = {
        ...data,
        assets,
      } satisfies NodeWithChildren;

      return updated;
    };

    const nodeKey = getNodeGetKey(slug);
    const nodeKeyFn = (key: Arguments) => {
      return Array.isArray(key) && key[0].startsWith(nodeKey);
    };

    await mutate(nodeKeyFn, nodeMutator, { revalidate: false });

    await nodeRemoveAsset(slug, assetID);
  };

  const deleteNode = async (slug: string, newParent?: string) => {
    const mutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      const newNodes = data.nodes.filter((node) => node.slug !== slug);

      return {
        ...data,
        nodes: newNodes,
      };
    };

    await mutate(nodeListAllKeyFn, mutator, { revalidate: false });

    await nodeDelete(slug, { target_node: newParent });

    // TODO: Ensure redirect only happens if you're viewing this actual page.
    if (newParent) {
      const newPath = replaceLibraryPath(libraryPath, slug, newParent);
      router.push(newPath);
    } else {
      router.push("/l");
    }
  };

  const moveNode = async (
    draggingNodeId: string,
    dropTargetId: string,
    dropPosition: "above" | "below" | "inside",
    newParent: string | null,
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

    await mutate(nodeListAllKeyFn, mutator, { revalidate: false });

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

  const revalidate = async (data?: MutatorCallback<NodeListOKResponse>) => {
    await mutate(nodeListAllKeyFn, data);
    if (node) {
      await mutate(nodeKeyFn);
    }
  };

  return {
    createNode,
    updateNode,
    suggestTitle,
    suggestSummary,
    suggestTags,
    importFromLink,
    removeNodeCoverImage,
    updateNodeVisibility,
    addAsset,
    removeAsset,
    deleteNode,
    moveNode,
    revalidate,
  };
}

// Used purely for optimistic mutation where the asset is swapped out.
function mergePrimaryImageAsset(
  oldNode: Node,
  coverConfig: CoverImageArgs,
): Pick<Node, "primary_image" | "meta"> {
  // If replacing the asset, we don't want to keep the old parent.
  if (coverConfig.isReplacement) {
    return {
      primary_image: coverConfig.asset,
      meta: { ...oldNode.meta, coverImage: null },
    };
  }

  const parentAsset = oldNode.primary_image?.parent ?? oldNode.primary_image;

  return {
    primary_image: { ...coverConfig.asset, parent: parentAsset },
    meta: { ...oldNode.meta, coverImage: coverConfig.config },
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
