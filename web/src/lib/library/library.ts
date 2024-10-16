import slugify from "@sindresorhus/slugify";
import { dequal } from "dequal/lite";
import { omit, uniqueId } from "lodash";
import { useRouter } from "next/navigation";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";
import { Xid } from "xid-ts";

import {
  getNodeGetKey,
  getNodeListKey,
  nodeCreate,
  nodeDelete,
  nodeUpdate,
  nodeUpdateVisibility,
} from "@/api/openapi-client/nodes";
import {
  Asset,
  Node,
  NodeGetOKResponse,
  NodeListOKResponse,
  NodeListParams,
  NodeMutableProps,
  NodeWithChildren,
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

export function useLibraryMutation(params?: NodeListParams) {
  const session = useSession();
  const { mutate } = useSWRConfig();
  const router = useRouter();
  const libraryPath = useLibraryPath();

  // for revalidating all node list queries (published and private)
  const nodeListKey = getNodeListKey(params);
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
      meta: {},
      children: [],
      assets: [],
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

    const created = await nodeCreate({ name: name, parent: parentSlug });
    const newPath = joinLibraryPath(libraryPath, created.slug);

    router.push(`/l/${newPath}?edit=true`);
  };

  const updateNode = async (
    slug: string,
    newNode: NodeMutableProps,
    cover?: CoverImageArgs,
  ) => {
    const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
      if (!data) return;

      const nodeProps = omit(newNode, "parent");

      const withNewCover = cover?.asset && mergePrimaryImageAsset(data, cover);

      const updated = {
        ...data,
        ...nodeProps,
        ...withNewCover,
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

    const nodeKey = getNodeGetKey(slug);
    const nodeKeyFn = (key: Arguments) => {
      return Array.isArray(key) && key[0].startsWith(nodeKey);
    };

    const slugChanged = newNode.slug !== undefined && newNode.slug !== slug;

    await mutate(nodeListAllKeyFn, listMutator, { revalidate: false });
    await mutate(nodeKeyFn, nodeMutator, { revalidate: false });

    const newMeta =
      cover && !cover.isReplacement
        ? ({
            // TODO: Spread original node metadata here
            coverImage: cover.config,
          } satisfies NodeMetadata)
        : undefined;

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

    return slugChanged;
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

  const revalidate = async (data?: MutatorCallback<NodeListOKResponse>) => {
    await mutate(nodeListAllKeyFn, data);
  };

  return {
    createNode,
    updateNode,
    removeNodeCoverImage,
    updateNodeVisibility,
    deleteNode,
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
