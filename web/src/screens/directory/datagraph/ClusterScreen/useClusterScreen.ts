import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { ClusterInitialProps, ClusterWithItems } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

import { useDirectoryPath } from "../useDirectoryPath";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  slug: z.string().min(1, "Please enter a slug."),
  description: z.string().min(1, "Please enter a short description."),
  content: z.string().optional(),
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

  const isAllowedToEdit =
    editable && Boolean(account?.id === cluster.owner.id || account?.admin);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: cluster,
  });

  function handleEditMode() {
    if (editing) return;
    if (!isAllowedToEdit) return;

    setEditing(true);
  }

  function handleSave(payload: ClusterInitialProps) {
    if (!editing) return;

    setEditing(false);
    onSave(payload);
  }

  function handleAssetUpload() {
    if (!editing) return;

    onSave(form.getValues());
  }

  return {
    form,
    handlers: {
      handleEditMode,
      handleSave,
      handleAssetUpload,
    },
    directoryPath,
    editing,
    cluster,
    isAllowedToEdit,
  };
}
