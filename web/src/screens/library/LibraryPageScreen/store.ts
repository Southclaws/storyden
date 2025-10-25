import { uniqueId } from "lodash";
import { immer } from "zustand/middleware/immer";
import { useStoreWithEqualityFn } from "zustand/traditional";
import { createStore } from "zustand/vanilla";

import { handle } from "@/api/client";
import {
  Asset,
  Identifier,
  LinkReference,
  NodeMutableProps,
  NodeWithChildren,
  NodeWithChildrenAllOf,
  PropertyName,
  PropertySchema,
  PropertyType,
} from "@/api/openapi-schema";
import { MutationSet, deriveMutationFromDifference } from "@/lib/library/diff";
import { CoverImageArgs } from "@/lib/library/library";
import {
  DefaultLayout,
  LibraryPageBlock,
  LibraryPageBlockType,
  WithMetadata,
} from "@/lib/library/metadata";
import { applyNodeChanges } from "@/lib/library/mutators";
import { deepEqual } from "@/utils/equality";

import { useLibraryPageContext } from "./Context";

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
  removeLink: () => void;
  addAsset: (asset: Asset) => void;
  removeAsset: (asset: Asset) => void;

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
  setChildPropertyValue: (
    nodeID: Identifier,
    fid: string,
    value: string,
  ) => void;
  setChildPropertyHiddenState: (fid: string, hidden: boolean) => void;

  // Layout blocks
  moveBlock: (type: LibraryPageBlockType, newIndex: number) => void;
  addBlock: (type: LibraryPageBlockType, index?: number) => void;
  removeBlock: (type: LibraryPageBlockType) => void;
  overwriteBlock: (type: LibraryPageBlock) => void;

  commit: (
    callback: (draft: MutationSet) => Promise<NodeWithChildrenAllOf>,
  ) => Promise<void>;
};

export type Store = State & Actions;

export type NodeStoreAPI = ReturnType<typeof createNodeStore>;

export const createNodeStore = (initState: State) => {
  return createStore<Store>()(
    immer((set, get) => {
      const simplePatch = (data: NodeMutableProps) =>
        set((state) => {
          const newState = applyNodeChanges(state.draft, data);
          Object.assign(state.draft, newState);
        });

      const commit = async (
        callback: (draft: MutationSet) => Promise<NodeWithChildrenAllOf>,
      ) => {
        const current = get().original;
        const draft = get().draft;
        const mutation = deriveMutationFromDifference(current, draft);

        if (mutation.clean) {
          console.debug("skipping commit: no changes");
          return;
        }

        console.debug(`applying commit: `, mutation);

        const updated = await handle(
          async () => {
            return await callback(mutation);
          },
          {
            errorToast: true,
          },
        );

        if (updated) {
          set(() => ({
            original: updated,
            draft: updated,
          }));
        }
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
        removeLink: () => {
          set((state) => {
            state.draft.link = undefined;
          });
        },

        addAsset: (asset: Asset) => {
          set((state) => {
            state.draft.assets.push(asset);
          });
        },

        removeAsset: (asset: Asset) => {
          set((state) => {
            state.draft.assets = state.draft.assets.filter(
              (a) => a.id !== asset.id,
            );
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

            const layout = (state.draft.meta.layout ??=
              structuredClone(DefaultLayout));

            for (const block of layout.blocks) {
              if (block.type !== "directory") continue;

              // config might not be defined yet, it should be in all cases, but
              // typescript is unsure (so am i) so just to be safe, set default.
              if (!block.config) {
                block.config = { layout: "table", columns: [] };
              } else if (!block.config.columns) {
                block.config.columns = [];
              }

              block.config.columns.push(newColumn);
            }

            state.draft.child_property_schema.push(newProperty);
          });
        },

        removeChildPropertyByID: (fid: Identifier) => {
          set((state) => {
            const layout = state.draft.meta.layout;
            if (!layout) return;

            for (const block of layout.blocks) {
              if (block.type !== "directory") continue;
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
            const target = state.draft.child_property_schema.find(
              (f) => f.fid === fid,
            );
            if (target) {
              target.name = newName;
            }
          });
        },

        setChildPropertyValue: (
          nodeID: Identifier,
          fid: string,
          value: string,
        ) => {
          set((state) => {
            const isFixed = fid.startsWith("fixed:");
            const target = state.draft.children.find((f) => f.id === nodeID);
            if (!target) {
              console.warn(
                "Attempting to set property on non-existing child node",
                {
                  nodeID,
                  fid,
                  value,
                },
              );
              return;
            }

            console.debug("Setting child property value", {
              nodeID,
              fid,
              value,
              isFixed,
            });

            if (isFixed) {
              switch (fid) {
                case "fixed:name": {
                  target.name = value;
                  break;
                }
                case "fixed:description": {
                  target.description = value;
                  break;
                }
                case "fixed:link": {
                  // NOTE: The actual value here is only worried about the URL.
                  // see diff.ts projectNodeToMutableProps for why. Since this
                  // is a mutation only, the other fields don't matter.
                  target.link = {
                    url: value,
                  } as LinkReference;
                  break;
                }
              }
            } else {
              const prop = target.properties.find((p) => p.fid === fid);
              if (!prop) {
                console.warn(
                  "Attempting to set value on non-existing property",
                  {
                    nodeID,
                    fid,
                    value,
                  },
                );
                return;
              }
              prop.value = value;
            }
          });
        },

        setChildPropertyHiddenState: (fid: string, hidden: boolean) => {
          set((state) => {
            const layout = (state.draft.meta.layout ??= DefaultLayout);
            const blocks = layout.blocks;

            for (const block of blocks) {
              if (block.type !== "directory") continue;
              if (!block.config?.columns) continue;

              for (const col of block.config.columns) {
                if (col.fid === fid) {
                  col.hidden = hidden;
                  return;
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

        addBlock: (type: LibraryPageBlockType, index?: number) => {
          set((state) => {
            const layout = (state.draft.meta.layout ??= DefaultLayout);

            // check if the block already exists, if it does, return
            if (layout.blocks.some((b) => b.type === type)) {
              return;
            }

            if (index === undefined || index > layout.blocks.length) {
              layout.blocks.push({ type });
            } else {
              layout.blocks.splice(index + 1, 0, { type });
            }
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

        overwriteBlock: (block: LibraryPageBlock) => {
          set((state) => {
            if (state.draft.meta.layout === undefined) {
              state.draft.meta.layout = DefaultLayout;
            }

            // check if the block exists, if not, do nothing.
            const index = state.draft.meta.layout.blocks.findIndex(
              (b) => b.type === block.type,
            );
            if (index === -1) {
              return;
            }

            state.draft.meta.layout.blocks[index] = block;
          });
        },

        commit,
      };
    }),
  );
};

export function useWatch<T>(selector: (state: State) => T): T {
  const { store } = useLibraryPageContext();
  return useStoreWithEqualityFn(store, selector, deepEqual);
}
