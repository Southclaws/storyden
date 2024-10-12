import slugify from "@sindresorhus/slugify";
import { dequal } from "dequal/lite";
import { last, uniqueId } from "lodash";
import { useRouter } from "next/navigation";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";
import { Xid } from "xid-ts";

import {
  getNodeListKey,
  nodeCreate,
  nodeDelete,
  nodeUpdate,
  nodeUpdateVisibility,
} from "@/api/openapi-client/nodes";
import {
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

export type CreateNodeArgs = {
  initialName?: string;
  parentSlug?: string;
};

export function useLibraryMutation(params?: NodeListParams) {
  const session = useSession();
  const { mutate } = useSWRConfig();
  const router = useRouter();
  const libraryPath = useLibraryPath();

  // for revalidating all node list queries (published and private)
  const nodeListKey = getNodeListKey(params);
  const nodeListAllKeyFn = (key: Arguments) => {
    return Array.isArray(key) && key[0].startsWith(nodeListKey);
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

  const updateNode = async (slug: string, newNode: NodeMutableProps) => {
    const mutator: MutatorCallback<NodeListOKResponse> = (data) => {
      if (!data) return;

      const newNodes = data.nodes.map((n) => {
        if (n.slug === slug) {
          return { ...n, ...newNode } as NodeWithChildren;
        }

        return n;
      });

      return {
        ...data,
        nodes: newNodes,
      };
    };

    const slugChanged = newNode.slug !== undefined && newNode.slug !== slug;

    await mutate(nodeListAllKeyFn, mutator, { revalidate: false });

    await nodeUpdate(slug, {
      name: newNode.name,
      slug: newNode.slug,
      asset_ids: newNode.asset_ids,
      url: newNode.url,
      content: newNode.content,
      meta: newNode.meta,
    });

    // Handle slug changes properly by redirecting to the new path.
    if (slugChanged && newNode.slug /* Needed for TS narrowing */) {
      const newPath = replaceLibraryPath(libraryPath, slug, newNode.slug);
      router.push(newPath);
    }

    return slugChanged;
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
    updateNodeVisibility,
    deleteNode,
    revalidate,
  };
}
