import { parseAsBoolean, parseAsString, useQueryState } from "nuqs";
import {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { nodeVersionCreate } from "@/api/openapi-client/nodes";
import { NodeVersion } from "@/api/openapi-schema";

import { useLibraryPageContext } from "./Context";
import {
  directEditorSourceKey,
  liveEditorSourceKey,
  versionEditorSourceKey,
} from "./editorSource";
import { LibraryPageEditMode, normaliseLibraryPageEditMode } from "./editMode";
import { useLibraryPagePermissions } from "./permissions";
import { overlayNodeVersion } from "./versionedEdit";

type EditStateContext = {
  editMode: LibraryPageEditMode;
  editing: boolean;
  isDirectEditing: boolean;
  isProposalEditing: boolean;
  editorSourceKey: string;
  saving: boolean;
  proposalVersion?: NodeVersion;
  setProposalVersion: (version?: NodeVersion) => void;
  setEditing: (value: boolean) => void;
  startEditing: (mode: LibraryPageEditMode) => void;
  startDirectEdit: () => void;
  startProposalEdit: (version?: NodeVersion) => Promise<void>;
  stopEditing: () => void;
  handleToggleEditMode: () => void;
};

const Context = createContext<EditStateContext | null>(null);

export function LibraryPageEditProvider({ children }: PropsWithChildren) {
  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });
  const [rawEditMode, setRawEditMode] = useQueryState("editMode", {
    ...parseAsString,
    defaultValue: LibraryPageEditMode.direct,
    clearOnDefault: true,
  });
  const [, setReviewVersionID] = useQueryState("version", {
    ...parseAsString,
    clearOnDefault: true,
  });
  const [proposalVersion, setProposalVersion] = useState<
    NodeVersion | undefined
  >(undefined);

  const {
    initialNode,
    nodeID,
    revalidate,
    saving,
    setSaving,
    store,
    suppressAutosave,
  } = useLibraryPageContext();
  const [editorSourceKey, setEditorSourceKey] = useState(() =>
    liveEditorSourceKey(initialNode),
  );
  const { isAllowedToDirectEdit, isAllowedToProposeEdit } =
    useLibraryPagePermissions();

  const editMode = normaliseLibraryPageEditMode(rawEditMode);
  const isDirectEditing = editing && editMode === LibraryPageEditMode.direct;
  const isProposalEditing =
    editing && editMode === LibraryPageEditMode.proposal;

  function startDirectEdit() {
    if (!isAllowedToDirectEdit) return;

    setProposalVersion(undefined);
    setReviewVersionID(null);
    setEditorSourceKey(directEditorSourceKey(initialNode));
    setRawEditMode(LibraryPageEditMode.direct);
    setEditing(true);
  }

  async function startProposalEdit(existingVersion?: NodeVersion) {
    if (!isAllowedToProposeEdit) return;

    setSaving(true);
    try {
      const version =
        existingVersion ??
        (await handle(() => nodeVersionCreate(nodeID, {}), {
          errorToast: true,
        }));

      if (!version) {
        return;
      }

      const draft = overlayNodeVersion(initialNode, version);

      setProposalVersion(version);
      suppressAutosave(() => {
        store.setState((state) => {
          state.original = initialNode;
          state.draft = draft;
        });
      });
      setEditorSourceKey(versionEditorSourceKey(version));

      setReviewVersionID(null);
      setRawEditMode(LibraryPageEditMode.proposal);
      setEditing(true);
    } finally {
      setSaving(false);
    }
  }

  function startEditing(mode: LibraryPageEditMode) {
    if (mode === LibraryPageEditMode.direct) {
      startDirectEdit();
      return;
    }

    void startProposalEdit();
  }

  function stopEditing() {
    if (isProposalEditing) {
      suppressAutosave(() => {
        store.setState((state) => {
          state.original = initialNode;
          state.draft = initialNode;
        });
      });
    }

    setEditorSourceKey(liveEditorSourceKey(initialNode));
    setEditing(false);
    setRawEditMode(LibraryPageEditMode.direct);
    setProposalVersion(undefined);
    revalidate();
  }

  function handleToggleEditMode() {
    if (editing) {
      stopEditing();
      return;
    }

    if (isAllowedToDirectEdit) {
      startDirectEdit();
      return;
    }

    void startProposalEdit();
  }

  useEffect(() => {
    if (!isProposalEditing || proposalVersion) {
      return;
    }

    setEditing(false);
    setRawEditMode(LibraryPageEditMode.direct);
    toast.error("Start a draft edit from the page controls.");
  }, [isProposalEditing, proposalVersion, setEditing, setRawEditMode]);

  useEffect(() => {
    if (editing) {
      return;
    }

    setEditorSourceKey(liveEditorSourceKey(initialNode));
  }, [editing, initialNode]);

  const value = {
    editMode,
    editing,
    isDirectEditing,
    isProposalEditing,
    editorSourceKey,
    saving,
    proposalVersion,
    setProposalVersion,
    setEditing,
    startEditing,
    startDirectEdit,
    startProposalEdit,
    stopEditing,
    handleToggleEditMode,
  };

  return <Context.Provider value={value}>{children}</Context.Provider>;
}

export function useEditState() {
  const context = useContext(Context);
  if (!context) {
    throw new Error("useEditState must be used within LibraryPageEditProvider");
  }

  return context;
}
