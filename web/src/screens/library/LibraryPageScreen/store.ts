import { dequal } from "dequal";
import { immer } from "zustand/middleware/immer";
import { useStoreWithEqualityFn } from "zustand/traditional";
import { createStore } from "zustand/vanilla";

import {
  NodeMutableProps,
  NodeWithChildren,
  NodeWithChildrenAllOf,
  PropertyMutationList,
} from "@/api/openapi-schema";
import { NodeMetadata, WithMetadata } from "@/lib/library/metadata";

import { useLibraryPageContext } from "./Context";

type NodeDraftEvent = {
  type: "patch";
  data: NodeMutableProps;
};

export type State = {
  draft: WithMetadata<NodeWithChildren>;
  draftEvents: NodeDraftEvent[];
};

export type Actions = {
  setName: (v: string) => void;
  setSlug: (v: string) => void;
  setContent: (v: string) => void;
  setPrimaryImage(assetID: string): void;
  setTags: (tags: string[]) => void;
  setLink: (url: string) => void;
  setProperties: (p: PropertyMutationList) => void;
  // setChildPropertySchema: (p: PropertySchemaList) => void;
  setMeta: (v: NodeMetadata) => void;

  commit: (
    callback: (draft: NodeMutableProps) => Promise<NodeWithChildrenAllOf>,
  ) => Promise<void>;
};

export type Store = State & Actions;

export const createNodeStore = (initState: State) => {
  return createStore<Store>()(
    immer((set, get) => {
      const patchDraft = (data: NodeMutableProps) =>
        set((state) => {
          Object.assign(state.draft, data);
          state.draftEvents.push({ type: "patch", data });
        });

      const commit = async (
        callback: (draft: NodeMutableProps) => Promise<NodeWithChildrenAllOf>,
      ) => {
        const { draftEvents } = get();

        const patch = draftEvents.reduce<NodeMutableProps>((acc, e) => {
          return { ...acc, ...e.data };
        }, {});

        const changes = Object.keys(patch).length;

        if (changes === 0) {
          console.debug("skipping commit: no changes");
          return;
        }

        console.debug(`applying commit: ${changes} changes`);

        const updated = await callback(patch);

        set(() => ({
          draft: updated,
          draftEvents: [],
        }));
      };

      return {
        ...initState,

        setName: (name) => patchDraft({ name }),
        setSlug: (slug) => patchDraft({ slug }),
        setContent: (content) => patchDraft({ content }),
        setPrimaryImage: (assetID) =>
          patchDraft({ primary_image_asset_id: assetID }),
        setTags: (tags) => patchDraft({ tags }),
        setLink: (url) => patchDraft({ url }),
        setProperties: (properties) => patchDraft({ properties }),
        // setChildPropertySchema: (child_property_schema) =>
        //   patchDraft({ child_property_schema }),
        setMeta: (meta) => patchDraft({ meta }),

        commit,
      };
    }),
  );
};

export function useWatch<T>(selector: (state: State) => T): T {
  const { store } = useLibraryPageContext();
  return useStoreWithEqualityFn(store, selector, dequal);
}
