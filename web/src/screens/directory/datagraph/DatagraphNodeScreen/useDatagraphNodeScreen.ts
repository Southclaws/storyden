import { zodResolver } from "@hookform/resolvers/zod";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { clusterAddAsset, clusterRemoveAsset } from "src/api/openapi/clusters";
import { itemAddAsset } from "src/api/openapi/items";
import { Asset, Visibility } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import {
  DatagraphNode,
  DatagraphNodeInitialProps,
  DatagraphNodeWithRelations,
} from "src/components/directory/datagraph/DatagraphNode";

import { useDirectoryPath } from "../useDirectoryPath";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  slug: z.string().min(1, "Please enter a slug."),
  description: z.string().min(1, "Please enter a short description."),
  content: z.string().optional(),
  asset_ids: z.array(z.string()),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  node: DatagraphNodeWithRelations;
  initialEditingState?: boolean;
  editable?: boolean;
  onVisibilityChange?: (v: Visibility) => Promise<void>;
  onSave: (c: DatagraphNodeInitialProps) => Promise<void>;
  onDelete?: (c: DatagraphNode) => Promise<void>;
};

export function useDatagraphNodeScreen({
  node,
  initialEditingState = false,
  editable = true,
  onVisibilityChange,
  onSave,
  onDelete,
}: Props) {
  const directoryPath = useDirectoryPath();
  const account = useSession();
  const [editing, setEditing] = useState(initialEditingState);
  const [isSaving, setIsSaving] = useState(false);

  const isAllowedToEdit =
    editable && Boolean(account?.id === node.owner.id || account?.admin);

  const isAllowedToDelete =
    editable &&
    Boolean(account?.id === node.owner.id || account?.admin) &&
    onDelete;

  const defaults: Form = useMemo(
    () => ({
      name: node.name,
      slug: node.slug,
      description: node.description,
      content: node.content,
      asset_ids: node.assets.map((a) => a.id),
    }),
    [node],
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: defaults,
  });

  function triggerSavingPopover() {
    setIsSaving(true);
    setTimeout(() => setIsSaving(false), 3000);
  }

  function handleEditMode() {
    if (editing) {
      form.reset();
      setEditing(false);

      return;
    }
    if (!isAllowedToEdit) return;

    setEditing(true);
    form.reset(node);
  }

  function handleSave(payload: Form) {
    if (!editing) return;

    triggerSavingPopover();
    setEditing(false);
    onSave({ type: node.type, ...payload });
  }

  async function handleVisibilityChange(v: Visibility) {
    triggerSavingPopover();
    onVisibilityChange?.(v);
  }

  function handleDelete() {
    if (editing) return;

    onDelete?.(node);
  }

  async function handleAssetUpload(asset: Asset) {
    if (!editing) return;

    // We only want to run these updates for edits of existing clusters.
    if (!node.id) return;

    triggerSavingPopover();
    if (node.type === "cluster") {
      await clusterAddAsset(node.slug, asset.id);
    } else {
      await itemAddAsset(node.slug, asset.id);
    }
  }

  async function handleAssetRemove(asset: Asset) {
    if (!editing) return;
    if (!node.id) return;

    triggerSavingPopover();
    if (node.type === "cluster") {
      await clusterRemoveAsset(node.slug, asset.id);
    } else {
      await itemAddAsset(node.slug, asset.id);
    }
  }

  const handleSubmit = form.handleSubmit(handleSave);

  return {
    form,
    handlers: {
      handleEditMode,
      handleSubmit,
      handleVisibilityChange,
      handleDelete,
      handleAssetUpload,
      handleAssetRemove,
    },
    directoryPath,
    editing,
    node,
    isAllowedToEdit,
    isSaving,
    isAllowedToDelete,
  };
}
