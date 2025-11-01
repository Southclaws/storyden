import { debounce, entries, toPairs } from "lodash";
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
  NodeListResult,
  NodeWithChildren,
} from "src/api/openapi-schema";

import {
  nodeUpdate,
  nodeUpdateChildrenPropertySchema,
} from "@/api/openapi-client/nodes";
import { MutationSet } from "@/lib/library/diff";
import { useLibraryMutation } from "@/lib/library/library";
import { WithMetadata, hydrateNode } from "@/lib/library/metadata";
import { deepEqual } from "@/utils/equality";
import { deriveError } from "@/utils/error";

import { NodeStoreAPI, createNodeStore } from "./store";

type LibraryPageContext = {
  nodeID: Identifier;
  initialNode: WithMetadata<NodeWithChildren>;
  initialChildren?: NodeListResult;
  store: NodeStoreAPI;
  saving: boolean;
  revalidate: () => void;
};

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
  childNodes?: NodeListResult;
};

export function LibraryPageProvider({
  node,
  childNodes,
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

  const saveDraft = useRef(
    debounce(() => {
      if (!storeRef.current) {
        return;
      }

      const state = storeRef.current.getState();

      state.commit(async (mutation: MutationSet) => {
        try {
          setSaving(() => true);

          if (mutation.childPropertySchemaMutation) {
            await nodeUpdateChildrenPropertySchema(
              node.id,
              mutation.childPropertySchemaMutation,
            );
          }

          if (mutation.childMutation) {
            const collapsed = Object.fromEntries(
              Object.entries(mutation.childMutation).map(([id, changes]) => [
                id,
                changes.at(-1)!,
              ]),
            );

            const operations = entries(collapsed);

            console.debug("Updating child nodes", operations);

            await Promise.all(
              operations.map(([childNodeID, child]) =>
                nodeUpdate(childNodeID, child),
              ),
            );
          }

          const updated = await nodeUpdate(node.id, mutation.nodeMutation);
          await revalidate(updated);

          const slugChanged = updated.slug !== state.original.slug;
          if (slugChanged) {
            window.history.replaceState(
              null,
              "",
              `/l/${updated.slug}?edit=true`,
            );
          }

          return updated;
        } catch (error) {
          throw new Error(deriveError(error), { cause: error });
        } finally {
          setTimeout(() => {
            setSaving(() => false);
          }, 500);
        }
      });
    }, 500),
  ).current;

  useEffect(() => {
    if (!storeRef.current) {
      return;
    }

    const unsub = storeRef.current.subscribe((state, prev) => {
      if (!deepEqual(state.draft, prev.draft)) {
        saveDraft();
      }
    });

    return unsub;
  }, [saveDraft]);

  // Cancel the saveDraft debounce when the component unmounts.
  useEffect(() => {
    return () => {
      saveDraft.cancel();
    };
  }, []);

  // Handle external changes to the original node state. This happens if another
  // source triggers a mutation+revalidation via SWR and the initial must update
  // the store state. This hook must run after the store subscription is set up.
  useEffect(() => {
    if (!storeRef.current) {
      return;
    }

    const { original, draft } = storeRef.current.getState();

    // We compare the un-hydrated node for original comparison, because the
    // nodeWithMeta object is potentially mutated by the hydration function to
    // set up default values for new nodes. This includes the page's layout.
    const equalToOriginal = deepEqual(original, node);
    const equalToDraft = deepEqual(draft, nodeWithMeta);

    storeRef.current.setState((state) => {
      if (!equalToOriginal) {
        state.original = node;
      }

      if (!equalToDraft) {
        state.draft = nodeWithMeta;
      }
    });
  }, [node, nodeWithMeta]);

  return (
    <Context.Provider
      value={{
        nodeID: node.id,
        initialNode: nodeWithMeta,
        initialChildren: childNodes,
        store: storeRef.current,
        saving,
        revalidate,
      }}
    >
      {children}
    </Context.Provider>
  );
}
