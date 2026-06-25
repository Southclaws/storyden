import {
  Dispatch,
  PropsWithChildren,
  SetStateAction,
  createContext,
  useCallback,
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
} from "@/api/openapi-schema";

import { useLibraryMutation } from "@/lib/library/library";
import { WithMetadata, hydrateNode } from "@/lib/library/metadata";
import { deepEqual } from "@/utils/equality";

import { NodeStoreAPI, createNodeStore } from "./store";

type LibraryPageContext = {
  nodeID: Identifier;
  initialNode: WithMetadata<NodeWithChildren>;
  initialChildren?: NodeListResult;
  store: NodeStoreAPI;
  saving: boolean;
  setSaving: Dispatch<SetStateAction<boolean>>;
  revalidate: (updated?: NodeWithChildren) => Promise<void>;
  suppressAutosave: (callback: () => void) => void;
  isAutosaveSuppressed: () => boolean;
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
  const suppressAutosaveRef = useRef(false);

  const storeRef = useRef<NodeStoreAPI | null>(null);
  if (storeRef.current === null) {
    storeRef.current = createNodeStore({
      original: nodeWithMeta,
      draft: nodeWithMeta,
    });
  }

  const suppressAutosave = useCallback((callback: () => void) => {
    suppressAutosaveRef.current = true;
    try {
      callback();
    } finally {
      suppressAutosaveRef.current = false;
    }
  }, []);

  const isAutosaveSuppressed = useCallback(
    () => suppressAutosaveRef.current,
    [],
  );

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

    suppressAutosave(() => {
      storeRef.current?.setState((state) => {
        if (!equalToOriginal) {
          state.original = nodeWithMeta;
        }

        if (!equalToDraft) {
          state.draft = nodeWithMeta;
        }
      });
    });
  }, [node, nodeWithMeta, suppressAutosave]);

  return (
    <Context.Provider
      value={{
        nodeID: node.id,
        initialNode: nodeWithMeta,
        initialChildren: childNodes,
        store: storeRef.current,
        saving,
        setSaving,
        revalidate,
        suppressAutosave,
        isAutosaveSuppressed,
      }}
    >
      {children}
    </Context.Provider>
  );
}
