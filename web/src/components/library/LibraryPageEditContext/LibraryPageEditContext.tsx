import { PropsWithChildren, createContext, useContext, useState } from "react";

import { NodeMutableProps, NodeWithChildren } from "@/api/openapi-schema";

export interface EditState {
  node: NodeWithChildren;

  updateTitle(s: string): void;
}

const LibraryPageContext = createContext<EditState | undefined>(undefined);

export function useLibraryPageContext() {
  const ctx = useContext(LibraryPageContext);
  if (ctx === undefined) {
    throw Error();
  }

  return ctx;
}

type Props = {
  node: NodeWithChildren;
};

export function Provider({ children, node }: PropsWithChildren<Props>) {
  const [nodeDraft, setNodeDraft] = useState<NodeWithChildren>(node);

  const updateTitle = (title: string) => {
    setNodeDraft((current) => {
      const next = {
        ...current,
        name: title,
      };

      return next;
    });
  };

  const state = {
    node: nodeDraft,
    updateTitle,
  };

  return (
    <LibraryPageContext.Provider value={state}>
      {children}
    </LibraryPageContext.Provider>
  );
}
