import { dequal } from "dequal";
import { Arguments } from "swr";

import { getNodeGetKey, getNodeListKey } from "@/api/openapi-client/nodes";
import {
  Identifier,
  NodeGetParams,
  NodeListParams,
  Visibility,
} from "@/api/openapi-schema";

type KeyType = ReturnType<typeof getNodeGetKey>;

// for revalidating all node list queries (published and private)
const nodeListKey = getNodeListKey();
const nodeListKeyPath = nodeListKey[0];

// for revalidating only private node list queries
const nodeListPrivateKey = getNodeListKey({
  // NOTE: The order here matters.
  visibility: [Visibility.draft, Visibility.review, Visibility.unlisted],
});

export const nodeListPrivateKeyFn = (key: Arguments) => {
  return dequal(key, nodeListPrivateKey);
};

export function buildNodeListKey(params?: NodeListParams) {
  const nodeListKeyFn = (key: Arguments): key is KeyType => {
    if (!key) return false;

    const path = key[0] as string;

    const isNodeListKey = path.startsWith(nodeListKeyPath);

    // Don't pass for /nodes/<slug> keys
    const notNodeKey = !path.startsWith(nodeListKeyPath + "/");

    const paramsEqual = params === undefined ? true : dequal(key[1], params);

    const matches = isNodeListKey && notNodeKey && paramsEqual;

    return matches;
  };

  return nodeListKeyFn;
}

export function buildNodeKey(slug: Identifier, params?: NodeGetParams) {
  const nodeKey = getNodeGetKey(slug, params);
  const nodeKeyPath = nodeKey[0];

  const nodeKeyFn = (key: Arguments): key is KeyType => {
    if (!key) return false;

    const path = key[0] as string;

    const pathMatches = path === nodeKeyPath;

    return pathMatches;
  };

  return nodeKeyFn;
}
