import { dequal } from "dequal";
import { uniqueId } from "lodash";
import { immer } from "zustand/middleware/immer";
import { useStoreWithEqualityFn } from "zustand/traditional";
import { createStore } from "zustand/vanilla";

import {
  Identifier,
  Link,
  LinkReference,
  NodeMutableProps,
  NodeWithChildren,
  NodeWithChildrenAllOf,
  PropertyMutationList,
  PropertyName,
  PropertySchema,
  PropertySchemaList,
  PropertyType,
} from "@/api/openapi-schema";
import { deriveMutationFromDifference } from "@/lib/library/diff";
import { CoverImageArgs } from "@/lib/library/library";
import {
  DefaultLayout,
  LibraryPageBlock,
  LibraryPageBlockType,
  NodeMetadata,
  WithMetadata,
} from "@/lib/library/metadata";
import { applyNodeChanges } from "@/lib/library/mutators";

import { useLibraryPageContext } from "./Context";

type NodeDraftEvent = {
  type: "patch";
  data: NodeMutableProps;
};

export type State = {
  original: WithMetadata<NodeWithChildren>;
  draft: WithMetadata<NodeWithChildren>;
};

export type Actions = {
  // Simple mutations - direct draft changes and mutation events.
  setName: (v: string) => void;
  setSlug: (v: string) => void;
  setContent: (v: string) => void;
  setTags: (tags: string[]) => void;

  // Complex mutations - require slightly more logic to change the draft.
  setPrimaryImage(args: CoverImageArgs): void;
  removePrimaryImage(): void;
  setLink: (link: LinkReference) => void;

  // Properties
  addProperty: (name: PropertyName, type: PropertyType, value?: string) => void;
  removePropertyByName: (name: PropertyName) => void;
  removePropertyByID: (id: Identifier) => void;
  setPropertyName: (name: PropertyName, newName: PropertyName) => void;
  setPropertyValue: (name: PropertyName, value: string) => void;
  // Properties, from parent perspective
  addChildProperty(p: PropertySchema): void;
  removeChildPropertyByID: (id: Identifier) => void;
  setChildPropertyName: (name: PropertyName, newName: PropertyName) => void;
  setChildPropertyHiddenState: (fid: string, hidden: boolean) => void;

  // Layout blocks
  moveBlock: (type: LibraryPageBlockType, newIndex: number) => void;
  addBlock: (type: LibraryPageBlockType) => void;
  removeBlock: (type: LibraryPageBlockType) => void;

  commit: (
    callback: (draft: NodeMutableProps) => Promise<NodeWithChildrenAllOf>,
  ) => Promise<void>;
};

export type Store = State & Actions;

export const createNodeStore = (initState: State) => {
  return createStore<Store>()(
    immer((set, get) => {
      const simplePatch = (data: NodeMutableProps) =>
        set((state) => {
          const newState = applyNodeChanges(state.draft, data);
          Object.assign(state.draft, newState);
        });

      const commit = async (
        callback: (draft: NodeMutableProps) => Promise<NodeWithChildrenAllOf>,
      ) => {
        const current = get().original;
        const draft = get().draft;
        const patch = deriveMutationFromDifference(current, draft);

        const changes = Object.keys(patch).length;

        if (changes === 0) {
          console.debug("skipping commit: no changes");
          return;
        }

        console.debug(`applying commit: ${changes} changes`, patch);

        const updated = await callback(patch);

        set(() => ({
          original: updated,
          draft: updated,
        }));
      };

      return {
        ...initState,

        // -
        // Simple mutations
        // -

        setName: (name) => simplePatch({ name }),
        setSlug: (slug) => simplePatch({ slug }),
        setContent: (content) => simplePatch({ content }),
        setTags: (tags) => simplePatch({ tags }),

        // -
        // Cover image
        // -

        setPrimaryImage: (coverConfig: CoverImageArgs) => {
          if (coverConfig.isReplacement) {
            set((state) => {
              state.draft.primary_image = coverConfig.asset;
              state.draft.meta = {
                ...state.draft.meta,
                coverImage: null,
              };
            });
          } else {
            set((state) => {
              state.draft.primary_image = coverConfig.asset;
              state.draft.meta = {
                ...state.draft.meta,
                coverImage: coverConfig.config,
              };
            });
          }
        },

        removePrimaryImage: () => {
          set((state) => {
            state.draft.primary_image = undefined;
            state.draft.meta = {
              ...state.draft.meta,
              coverImage: null,
            };
          });
        },

        setLink: (link: LinkReference) => {
          set((state) => {
            state.draft.link = link;
          });
        },

        // -
        // Property management
        // -

        addProperty: (
          name: PropertyName,
          type: PropertyType,
          value?: string,
        ) => {
          set((state) => {
            const existingNames = new Set(
              state.draft.properties.map((f) => f.name),
            );
            let newName = name;
            let counter = 1;
            while (existingNames.has(newName)) {
              newName = `${name} ${counter++}`;
            }

            state.draft.properties.push({
              fid: uniqueId("new_field_"),
              name: newName,
              type,
              sort: "5", // TODO: refine later
              value: value ?? "",
            });
          });
        },

        removePropertyByName: (name: PropertyName) => {
          set((state) => {
            state.draft.properties = state.draft.properties.filter(
              (f) => f.name !== name,
            );
          });
        },

        removePropertyByID: (id: Identifier) => {
          set((state) => {
            state.draft.properties = state.draft.properties.filter(
              (f) => f.fid !== id,
            );
          });
        },

        setPropertyName: (name: PropertyName, newName: PropertyName) => {
          set((state) => {
            const target = state.draft.properties.find((f) => f.name === name);
            if (target) {
              target.name = newName;
            }
          });
        },

        setPropertyValue: (name: PropertyName, value: string) => {
          set((state) => {
            const target = state.draft.properties.find((f) => f.name === name);
            if (target) {
              target.value = value;
            }
          });
        },

        // Child properties - used from parent perspective

        addChildProperty(newProperty: PropertySchema): void {
          set((state) => {
            const newColumn = {
              fid: newProperty.fid,
              hidden: false,
            };

            const layout = (state.draft.meta.layout ??= DefaultLayout);

            for (const block of layout.blocks) {
              if (block.type !== "table") continue;

              // config might not be defined yet, it should be in all cases, but
              // typescript is unsure (so am i) so just to be safe, set default.
              if (!block.config) {
                block.config = { columns: [] };
              } else if (!block.config.columns) {
                block.config.columns = [];
              }

              block.config.columns.push(newColumn);
            }
          });
        },

        removeChildPropertyByID: (fid: Identifier) => {
          set((state) => {
            const layout = state.draft.meta.layout;
            if (!layout) return;

            for (const block of layout.blocks) {
              if (block.type !== "table") continue;
              if (!block.config?.columns) continue;

              const columns = block.config.columns;
              const index = columns.findIndex((col) => col.fid === fid);
              if (index !== -1) {
                columns.splice(index, 1);
              }
            }

            const schema = state.draft.child_property_schema;
            const schemaIndex = schema.findIndex((p) => p.fid === fid);
            if (schemaIndex !== -1) {
              schema.splice(schemaIndex, 1);
            }
          });
        },

        setChildPropertyName: (fid: Identifier, newName: PropertyName) => {
          set((state) => {
            for (const p of state.draft.child_property_schema) {
              if (p.fid === fid) {
                p.name = newName;
              }
            }
          });
        },

        setChildPropertyHiddenState: (fid: string, hidden: boolean) => {
          set((state) => {
            const layout = (state.draft.meta.layout ??= DefaultLayout);
            const blocks = layout.blocks;

            for (const block of blocks) {
              if (block.type !== "table") continue;
              if (!block.config?.columns) continue;

              for (const col of block.config.columns) {
                if (col.fid === fid) {
                  col.hidden = hidden;
                }
              }
            }
          });
        },

        // -
        // Block management
        // -

        moveBlock: (type: LibraryPageBlockType, newIndex: number) => {
          set((state) => {
            const layout = (state.draft.meta.layout ??= DefaultLayout);
            const blocks = layout.blocks;

            const activeIndex = blocks.findIndex((b) => b.type === type);
            if (
              activeIndex === -1 ||
              newIndex === -1 ||
              activeIndex === newIndex
            ) {
              return;
            }

            const [moved] = blocks.splice(activeIndex, 1);
            if (!moved) {
              return;
            }

            blocks.splice(newIndex, 0, moved);
          });
        },

        addBlock: (type: LibraryPageBlockType) => {
          set((state) => {
            const layout = (state.draft.meta.layout ??= DefaultLayout);

            // check if the block already exists, if it does, return
            if (layout.blocks.some((b) => b.type === type)) {
              return;
            }

            layout.blocks.push({ type });
          });
        },

        removeBlock: (type: LibraryPageBlockType) => {
          set((state) => {
            const layout = state.draft.meta.layout;
            if (!layout) {
              return;
            }

            const blocks = layout.blocks;

            // check if the block exists, if it does, remove it
            const index = blocks.findIndex((b) => b.type === type);
            if (index === -1) {
              return;
            }

            blocks.splice(index, 1);
          });
        },

        commit,
      };
    }),
  );
};

export function useWatch<T>(selector: (state: State) => T): T {
  const { store } = useLibraryPageContext();
  return useStoreWithEqualityFn(store, selector, dequal);
}
