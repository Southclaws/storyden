import { zodResolver } from "@hookform/resolvers/zod";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { clusterAddAsset, clusterRemoveAsset } from "src/api/openapi/clusters";
import {
  Asset,
  ClusterInitialProps,
  ClusterWithItems,
} from "src/api/openapi/schemas";
import { useSession } from "src/auth";

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
  cluster: ClusterWithItems;
  initialEditingState?: boolean;
  editable?: boolean;
  onSave: (c: ClusterInitialProps) => Promise<void>;
};

export function useClusterScreen({
  cluster,
  initialEditingState = false,
  editable = true,
  onSave,
}: Props) {
  const directoryPath = useDirectoryPath();
  const account = useSession();
  const [editing, setEditing] = useState(initialEditingState);
  const [isSaving, setIsSaving] = useState(false);

  const isAllowedToEdit =
    editable && Boolean(account?.id === cluster.owner.id || account?.admin);

  const defaults: Form = useMemo(
    () => ({
      name: cluster.name,
      slug: cluster.slug,
      description: cluster.description,
      content: cluster.content,
      asset_ids: cluster.assets.map((a) => a.id),
    }),
    [cluster],
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
    form.reset(cluster);
  }

  function handleSave(payload: ClusterInitialProps) {
    if (!editing) return;

    triggerSavingPopover();
    setEditing(false);
    onSave(payload);
  }

  async function handleAssetUpload(asset: Asset) {
    if (!editing) return;

    // We only want to run these updates for edits of existing clusters.
    if (!cluster.id) return;

    triggerSavingPopover();
    await clusterAddAsset(cluster.slug, asset.id);
  }

  async function handleAssetRemove(asset: Asset) {
    if (!editing) return;
    if (!cluster.id) return;

    triggerSavingPopover();
    await clusterRemoveAsset(cluster.slug, asset.id);
  }

  const handleSubmit = form.handleSubmit(handleSave);

  return {
    form,
    handlers: {
      handleEditMode,
      handleSubmit,
      handleAssetUpload,
      handleAssetRemove,
    },
    directoryPath,
    editing,
    cluster,
    isAllowedToEdit,
    isSaving,
  };
}
