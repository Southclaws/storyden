import { dequal } from "dequal";
import { debounce } from "lodash";
import {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

import {
  Identifier,
  NodeMutableProps,
  NodeWithChildren,
} from "src/api/openapi-schema";

import { nodeUpdate } from "@/api/openapi-client/nodes";
import { useLibraryMutation } from "@/lib/library/library";
import { WithMetadata, hydrateNode } from "@/lib/library/metadata";

import { createNodeStore } from "./store";

type LibraryPageContext = {
  nodeID: Identifier;
  currentNode: WithMetadata<NodeWithChildren>;
  store: ReturnType<typeof createNodeStore>;
  saving: boolean;
};

type NodeStoreAPI = ReturnType<typeof createNodeStore>;

const Context = createContext<LibraryPageContext | null>(null);

export function useLibraryPageContext() {
  const context = useContext(Context);
  if (!context) {
    throw new Error(
      "useLibraryPageContext must be used within a LibraryPageProvider",
    );
  }

  return context;
}

export type Props = {
  node: NodeWithChildren;
};

export function LibraryPageProvider({
  node,
  children,
}: PropsWithChildren<Props>) {
  const [saving, setSaving] = useState(false);
  const nodeWithMeta = useMemo(() => hydrateNode(node), [node]);
  const { revalidate } = useLibraryMutation(node);

  const storeRef = useRef<NodeStoreAPI | null>(null);
  if (storeRef.current === null) {
    storeRef.current = createNodeStore({
      original: nodeWithMeta,
      draft: nodeWithMeta,
    });
  }

  // Handle external changes to the original node state. This happens if another
  // source triggers a mutation+revalidation via SWR and the initial must update
  useEffect(() => {
    if (!storeRef.current) {
      return;
    }

    const { original, draft } = storeRef.current.getState();

    const equalToOriginal = dequal(original, nodeWithMeta);
    const equalToDraft = dequal(draft, nodeWithMeta);

    storeRef.current.setState((state) => {
      if (!equalToOriginal) {
        state.original = nodeWithMeta;
      }

      if (!equalToDraft) {
        state.draft = nodeWithMeta;
      }
    });
  }, [nodeWithMeta]);

  const saveDraft = useRef(
    debounce(() => {
      if (!storeRef.current) {
        return;
      }

      const state = storeRef.current.getState();

      state.commit(async (patch: NodeMutableProps) => {
        setSaving(() => true);

        const updated = await nodeUpdate(node.id, patch);
        await revalidate(updated);

        const slugChanged = updated.slug !== state.original.slug;
        if (slugChanged) {
          window.history.replaceState(null, "", `/l/${updated.slug}?edit=true`);
        }

        setTimeout(() => {
          setSaving(() => false);
        }, 500);

        return updated;
      });
    }, 500),
  ).current;

  useEffect(() => {
    if (!storeRef.current) {
      return;
    }

    const unsub = storeRef.current.subscribe((state, prev) => {
      if (!dequal(state.draft, prev.draft)) {
        saveDraft();
      }
    });

    return unsub;
  }, [saveDraft]);

  return (
    <Context.Provider
      value={{
        nodeID: node.id,
        currentNode: nodeWithMeta,
        store: storeRef.current,
        saving,
      }}
    >
      {children}
    </Context.Provider>
  );
}
