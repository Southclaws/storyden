import { dequal } from "dequal";
import { debounce } from "lodash";
import { parseAsBoolean, useQueryState } from "nuqs";
import {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useRef,
} from "react";

import { NodeMutableProps, NodeWithChildren } from "src/api/openapi-schema";

import { nodeUpdate } from "@/api/openapi-client/nodes";
import { deriveMutationFromDifference } from "@/lib/library/diff";
import { useLibraryMutation } from "@/lib/library/library";
import { WithMetadata, hydrateNode } from "@/lib/library/metadata";

import { createNodeStore } from "./store";

type LibraryPageContext = {
  currentNode: WithMetadata<NodeWithChildren>;
  store: ReturnType<typeof createNodeStore>;
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
  const nodeWithMeta = hydrateNode(node);
  const { revalidate } = useLibraryMutation(node);

  // NOTE: Copied from useEditState - cannot call here though, not in context.
  const [editing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  const storeRef = useRef<NodeStoreAPI | null>(null);
  if (storeRef.current === null) {
    storeRef.current = createNodeStore({
      draft: nodeWithMeta,
      draftEvents: [],
    });
  }

  // Handle external changes to the original node state. This happens if another
  // source triggers a mutation+revalidation via SWR and the initial must update
  useEffect(() => {
    if (!storeRef.current) {
      return;
    }

    storeRef.current.setState((state) => {
      state.draft = nodeWithMeta;
    });
  }, [nodeWithMeta]);

  const saveDraft = useRef(
    debounce(() => {
      if (!storeRef.current) {
        return;
      }

      const current = storeRef.current.getInitialState().draft;
      const updated = storeRef.current.getState().draft;
      const patch = deriveMutationFromDifference(current, updated);
      console.log("experimental patch:", patch);

      storeRef.current.getState().commit(async (patch: NodeMutableProps) => {
        console.log("Saving patch:", patch);

        const updated = await nodeUpdate(node.id, patch);
        await revalidate(updated);

        const slugChanged = updated.slug !== current.slug;
        if (slugChanged) {
          const query = editing ? "?edit=true" : "";
          window.history.replaceState(null, "", `/l/${updated.slug}${query}`);
        }

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
        currentNode: nodeWithMeta,
        store: storeRef.current,
      }}
    >
      {children}
    </Context.Provider>
  );
}
