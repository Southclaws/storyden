import { createTreeCollection } from "@ark-ui/react/tree-view";
import {
  type KeyboardEvent,
  type RefObject,
  useEffect,
  useRef,
  useState,
} from "react";

import { handle } from "@/api/client";

export type DragTreePlacement = "above" | "inside" | "below";

export type DragTreeMove<TNode> = {
  node: TNode;
  oldParent: TNode | null;
  target: TNode;
  parent: TNode | null;
  placement: DragTreePlacement;
};

export type DragTreeMoveConfig<TNode> = {
  /** Reject on failure so the tree can keep the move available for retry. */
  onMove: (move: DragTreeMove<TNode>) => Promise<void>;
  /** Refresh/rollback caches after either a successful or failed request. */
  onSettled?: () => Promise<void>;
};

const placements: DragTreePlacement[] = ["above", "inside", "below"];

export function useDragTreeMove<TNode>({
  collection,
  enabled,
  config,
  expand,
  rootRef,
}: {
  collection: ReturnType<typeof createTreeCollection<TNode>>;
  enabled: boolean;
  config?: DragTreeMoveConfig<TNode>;
  expand: (id: string) => void;
  rootRef: RefObject<HTMLDivElement | null>;
}) {
  const [sourceId, setSourceId] = useState<string | null>(null);
  const [focusedId, setFocusedId] = useState<string | null>(null);
  const [placement, setPlacement] = useState<DragTreePlacement>("inside");
  const [keyboard, setKeyboard] = useState(false);
  const [hintId, setHintId] = useState<string | null>(null);
  const [pending, setPending] = useState(false);
  const committing = useRef(false);
  const restoringId = useRef<string | null>(null);
  const [focusRequest, setFocusRequest] = useState<{
    id: string;
    fallback: string;
  } | null>(null);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const source = sourceId ? collection.findNode(sourceId) : undefined;
  const target = focusedId ? collection.findNode(focusedId) : undefined;
  const invalid =
    !!sourceId &&
    (!target ||
      focusedId === sourceId ||
      collection.getDescendantValues(sourceId).includes(focusedId!));
  const parentOf = (id: string): TNode | null => {
    const parent = collection.getParentNode(id);
    return parent && !collection.isRootNode(parent) ? parent : null;
  };
  const parent =
    target && focusedId
      ? placement === "inside"
        ? target
        : parentOf(focusedId)
      : null;
  const location = parent
    ? `inside ${collection.stringifyNode(parent)}`
    : "at the top level";
  const destination = target
    ? placement === "inside"
      ? `inside ${collection.stringifyNode(target)}`
      : `${placement === "above" ? "before" : "after"} ${collection.stringifyNode(target)}, ${location}`
    : "";
  const preview = invalid
    ? "Choose a node outside the carried branch."
    : destination;
  const status = pending
    ? "Placing the node. Please wait."
    : source
      ? `Carrying ${collection.stringifyNode(source)}. ${preview} ${invalid ? "" : "Enter to place. "}Space changes placement; Escape cancels.`
      : notice;

  function focusNode(id: string) {
    for (const ancestor of collection.getParentNodes(id)) {
      if (!collection.isRootNode(ancestor))
        expand(collection.getNodeValue(ancestor));
    }
    requestAnimationFrame(() => {
      const anchor = Array.from(
        rootRef.current?.querySelectorAll<HTMLAnchorElement>(
          "a[data-tree-node-id]",
        ) ?? [],
      ).find((element) => element.dataset["treeNodeId"] === id);
      anchor?.focus();
    });
  }

  function start(id: string) {
    if (!enabled || !config || committing.current) return;
    setError("");
    setNotice("");
    setPlacement("inside");
    setSourceId(id);
    setFocusedId(id);
    setKeyboard(true);
    if (document.activeElement?.getAttribute("data-tree-node-id") !== id)
      focusNode(id);
  }

  function cancel(restoreFocus = true) {
    if (!sourceId || committing.current) return;
    setNotice("Move cancelled. The original position is unchanged.");
    setSourceId(null);
    setError("");
    if (restoreFocus) focusNode(sourceId);
  }

  async function commit() {
    if (
      !enabled ||
      !config ||
      !source ||
      !target ||
      !sourceId ||
      !focusedId ||
      invalid ||
      committing.current
    )
      return;
    committing.current = true;
    setPending(true);
    setError("");
    if (parent) expand(collection.getNodeValue(parent));
    let saved = false;
    await handle(
      async () => {
        await config.onMove({
          node: source,
          oldParent: parentOf(sourceId),
          target,
          parent,
          placement,
        });
        saved = true;
      },
      {
        errorToast: false,
        onError: async () =>
          setError(
            "Could not move the node. Choose a destination and press Enter to retry, or Escape to cancel.",
          ),
        cleanup: async () => {
          await handle(async () => config.onSettled?.(), {
            errorToast: false,
            onError: async () =>
              setError(
                "Could not refresh the tree. Reload to check the saved position.",
              ),
          });
          committing.current = false;
          setPending(false);
        },
      },
    );
    if (saved) {
      setNotice(`Moved ${collection.stringifyNode(source)} ${destination}.`);
      setSourceId(null);
    }
    // Reparenting can replace the focused anchor and temporarily focus body.
    // Recover that focus, but respect a deliberate move to another control.
    if (
      document.activeElement === document.body ||
      rootRef.current?.contains(document.activeElement)
    ) {
      const id = saved ? sourceId : focusedId;
      restoringId.current = id;
      setFocusRequest({ id, fallback: focusedId });
    }
  }

  function onKeyDown(event: KeyboardEvent) {
    const anchor = (event.target as HTMLElement).closest<HTMLAnchorElement>(
      "a[data-tree-node-id]",
    );
    if (
      !anchor ||
      !rootRef.current?.contains(anchor) ||
      event.nativeEvent.isComposing
    )
      return;
    restoringId.current = null;
    setKeyboard(true);
    if (
      enabled &&
      config &&
      event.altKey &&
      event.shiftKey &&
      !event.ctrlKey &&
      !event.metaKey &&
      (event.code === "KeyM" || event.key.toLowerCase() === "m")
    ) {
      event.preventDefault();
      event.stopPropagation();
      if (!event.repeat && !sourceId) start(anchor.dataset["treeNodeId"]!);
      return;
    }
    if (!sourceId) {
      if (event.key === "Escape") setKeyboard(false);
      return;
    }
    if (event.key === "Tab") {
      cancel(false);
      return;
    }
    if (pending || ["Enter", " ", "Escape"].includes(event.key)) {
      event.preventDefault();
      event.stopPropagation();
      if (pending || event.repeat) return;
      if (event.key === "Escape") cancel();
      if (event.key === " ")
        setPlacement(
          (current) =>
            placements[(placements.indexOf(current) + 1) % placements.length]!,
        );
      if (event.key === "Enter") void commit();
    }
  }

  // Cache refreshes may reparent/remount a node more than once. Restore focus
  // against the committed DOM, and stop as soon as the user navigates again.
  useEffect(() => {
    if (!focusRequest || restoringId.current !== focusRequest.id) return;
    for (const ancestor of collection.getParentNodes(focusRequest.id)) {
      if (!collection.isRootNode(ancestor))
        expand(collection.getNodeValue(ancestor));
    }
    const frame = requestAnimationFrame(() => {
      if (restoringId.current !== focusRequest.id) return;
      if (
        document.activeElement !== document.body &&
        !rootRef.current?.contains(document.activeElement)
      )
        return;
      const anchors = Array.from(
        rootRef.current?.querySelectorAll<HTMLAnchorElement>(
          "a[data-tree-node-id]",
        ) ?? [],
      );
      // A destination may deliberately hide its children in the sidebar.
      // Keep focus on that visible parent when the moved node is unavailable.
      const anchor =
        anchors.find(
          (element) => element.dataset["treeNodeId"] === focusRequest.id,
        ) ??
        anchors.find(
          (element) => element.dataset["treeNodeId"] === focusRequest.fallback,
        );
      anchor?.focus();
    });
    return () => cancelAnimationFrame(frame);
  }, [collection, expand, focusRequest, rootRef]);

  useEffect(() => {
    setHintId(null);
    if (!keyboard || sourceId || !enabled) return;
    const timeout = setTimeout(() => setHintId(focusedId), 600);
    return () => clearTimeout(timeout);
  }, [keyboard, sourceId, focusedId, enabled]);

  useEffect(() => {
    if (sourceId && !pending && (!enabled || !source)) {
      setSourceId(null);
      setNotice(
        "Move cancelled because the node or editing permission is no longer available.",
      );
    }
  }, [enabled, pending, source, sourceId]);

  return {
    sourceId,
    focusedId,
    placement,
    pending,
    error,
    status,
    source,
    invalid,
    preview,
    keyboard,
    hintId,
    setKeyboard,
    setFocusedId,
    start,
    cancel,
    commit,
    onKeyDown,
  };
}
