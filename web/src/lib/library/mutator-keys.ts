import { Arguments } from "swr";

import {
  getNodeGetKey,
  getNodeListChildrenKey,
  getNodeListKey,
} from "@/api/openapi-client/nodes";
import {
  Identifier,
  NodeGetParams,
  NodeListParams,
  Visibility,
} from "@/api/openapi-schema";
import { deepEqual } from "@/utils/equality";

type NodeKey = ReturnType<typeof getNodeGetKey>;
type NodeListKey = ReturnType<typeof getNodeListKey>;

// for revalidating all node list queries (published and private)
const nodeListKey = getNodeListKey();
const nodeListKeyPath = nodeListKey[0];

// for revalidating only private node list queries
const nodeListPrivateKey = getNodeListKey({
  // NOTE: The order here matters.
  visibility: [Visibility.draft, Visibility.review, Visibility.unlisted],
});

export const nodeListPrivateKeyFn = (key: Arguments) => {
  return deepEqual(key, nodeListPrivateKey);
};

export function buildNodeListKey(params?: NodeListParams) {
  const nodeListKeyFn = (key: Arguments): key is NodeListKey => {
    if (!key) return false;

    const path = key[0] as string;

    const isNodeListKey = path.startsWith(nodeListKeyPath);

    // Don't pass for /nodes/<slug> keys
    const notNodeKey = !path.startsWith(nodeListKeyPath + "/");

    const paramsEqual = params === undefined ? true : deepEqual(key[1], params);

    const matches = isNodeListKey && notNodeKey && paramsEqual;

    return matches;
  };

  return nodeListKeyFn;
}

export function buildNodeChildrenListKey(
  nodeID: string,
  params?: NodeListParams,
) {
  const nodeListKeyFn = (key: Arguments): key is NodeListKey => {
    if (!key) return false;

    const path = key[0] as string;

    const nodeListChildrenKey = getNodeListChildrenKey(nodeID);
    const nodeListChildrenKeyPath = nodeListChildrenKey[0];

    const isNodeListKey = path.startsWith(nodeListChildrenKeyPath);

    // Don't pass for /nodes/<slug> keys
    const notNodeKey = !path.startsWith(nodeListChildrenKeyPath + "/");

    const paramsEqual = params === undefined ? true : deepEqual(key[1], params);

    const matches = isNodeListKey && notNodeKey && paramsEqual;

    return matches;
  };

  return nodeListKeyFn;
}

export function buildNodeKey(slug: Identifier, params?: NodeGetParams) {
  const nodeKey = getNodeGetKey(slug, params);
  const nodeKeyPath = nodeKey[0];

  const nodeKeyFn = (key: Arguments): key is NodeKey => {
    if (!key) return false;

    const path = key[0] as string;

    const pathMatches = path === nodeKeyPath;

    return pathMatches;
  };

  return nodeKeyFn;
}
