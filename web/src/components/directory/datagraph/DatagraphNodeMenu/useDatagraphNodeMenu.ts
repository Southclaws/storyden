import { Visibility } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { useDeleteAction } from "src/components/site/Action/Delete";

import { DatagraphNode } from "../DatagraphNode";

export type Props = {
  node: DatagraphNode;
  onVisibilityChange?: (v: Visibility) => Promise<void>;
  onDelete: () => void;
};

export function useDatagraphNodeMenu(props: Props) {
  const account = useSession();
  const deleteProps = useDeleteAction({
    onClick: props.onDelete,
  });

  const isVisibilityChangeEnabled = props.onVisibilityChange ?? false;
  const isAdmin = account?.admin ?? false;
  const isOwner = account?.id === props.node.owner.id;

  // All possible visibility state transitions
  const draftToReview = props.node.visibility === "draft";
  const reviewToPublish = props.node.visibility !== "published" && isAdmin;
  const publishToReview = props.node.visibility === "published" && isAdmin;
  const reviewToDraft = props.node.visibility === "review";
  const draftToPublish = props.node.visibility !== "published" && isAdmin;

  // Only enable visibility changes if the event handler was passed in.
  const reviewFlow = isVisibilityChangeEnabled
    ? {
        draftToReview,
        reviewToPublish,
        publishToReview,
        reviewToDraft,
        draftToPublish,
      }
    : undefined;

  const deleteEnabled = isAdmin || isOwner;

  function handleSelect({ value }: { value: string }) {
    switch (value) {
      case "delete":
        deleteProps.onOpen();
        return;

      case "draft":
        props.onVisibilityChange?.("draft");
        return;

      case "review":
        props.onVisibilityChange?.("review");
        return;

      case "publish":
        props.onVisibilityChange?.("published");
        return;

      default:
        throw new Error(`unknown handler ${value}`);
    }
  }

  return {
    reviewFlow,
    deleteEnabled,
    deleteProps,
    handleSelect,
  };
}
