import { debounce } from "lodash";
import { useEffect, useMemo } from "react";

import { handle } from "@/api/client";
import {
  nodeUpdate,
  nodeUpdateChildrenPropertySchema,
  nodeVersionUpdate,
} from "@/api/openapi-client/nodes";
import { deriveMutationFromDifference } from "@/lib/library/diff";
import { hydrateNode } from "@/lib/library/metadata";
import { deepEqual } from "@/utils/equality";

import { useLibraryPageContext } from "./Context";
import { LibraryPageEditMode } from "./editMode";
import { useEditState } from "./useEditState";
import { buildNodeVersionMutation } from "./versionedEdit";

export function LibraryPageAutosaveController() {
  const {
    isAutosaveSuppressed,
    nodeID,
    revalidate,
    setSaving,
    store,
    suppressAutosave,
  } = useLibraryPageContext();
  const { editMode, proposalVersion, setProposalVersion } = useEditState();

  const saveDraft = useMemo(
    () =>
      debounce(async () => {
        const state = store.getState();
        const mutation = deriveMutationFromDifference(
          state.original,
          state.draft,
        );

        if (mutation.clean) {
          console.debug("skipping autosave: no changes");
          return;
        }

        if (editMode === LibraryPageEditMode.proposal) {
          if (!proposalVersion) {
            console.debug("skipping version draft save: no version loaded");
            return;
          }

          const patch = buildNodeVersionMutation(mutation);
          if (Object.keys(patch).length === 0) {
            console.debug("skipping version draft save: no supported changes");
            return;
          }

          try {
            setSaving(true);

            const updated = await handle(
              () => nodeVersionUpdate(nodeID, proposalVersion.id, patch),
              {
                errorToast: true,
              },
            );

            if (updated) {
              setProposalVersion(updated);
            }
          } finally {
            setTimeout(() => {
              setSaving(false);
            }, 500);
          }

          return;
        }

        try {
          setSaving(true);

          const updated = await handle(
            async () => {
              if (mutation.childPropertySchemaMutation) {
                await nodeUpdateChildrenPropertySchema(
                  nodeID,
                  mutation.childPropertySchemaMutation,
                );
              }

              const childOperations = Object.entries(
                mutation.childMutation,
              ).map(
                ([childNodeID, changes]) =>
                  [childNodeID, changes.at(-1)!] as const,
              );

              if (childOperations.length > 0) {
                console.debug("Updating child nodes", childOperations);

                await Promise.all(
                  childOperations.map(([childNodeID, child]) =>
                    nodeUpdate(childNodeID, child),
                  ),
                );
              }

              return nodeUpdate(nodeID, mutation.nodeMutation);
            },
            {
              errorToast: true,
            },
          );

          if (!updated) {
            return;
          }

          await revalidate(updated);

          const slugChanged = updated.slug !== state.original.slug;
          if (slugChanged) {
            window.history.replaceState(
              null,
              "",
              `/l/${updated.slug}?edit=true`,
            );
          }

          suppressAutosave(() => {
            const hydrated = hydrateNode(updated);
            store.setState((state) => {
              state.original = hydrated;
              state.draft = hydrated;
            });
          });
        } finally {
          setTimeout(() => {
            setSaving(false);
          }, 500);
        }
      }, 500),
    [
      editMode,
      nodeID,
      proposalVersion,
      revalidate,
      setProposalVersion,
      setSaving,
      store,
      suppressAutosave,
    ],
  );

  useEffect(() => {
    const unsubscribe = store.subscribe((state, prev) => {
      if (isAutosaveSuppressed()) {
        return;
      }

      if (!deepEqual(state.draft, prev.draft)) {
        saveDraft();
      }
    });

    return () => {
      unsubscribe();
      saveDraft.cancel();
    };
  }, [isAutosaveSuppressed, saveDraft, store]);

  return null;
}
