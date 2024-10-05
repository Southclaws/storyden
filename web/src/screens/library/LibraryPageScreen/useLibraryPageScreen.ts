import { zodResolver } from "@hookform/resolvers/zod";
import slugify from "@sindresorhus/slugify";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { P, match } from "ts-pattern";
import { z } from "zod";

import { nodeAddAsset, nodeRemoveAsset } from "src/api/openapi-client/nodes";
import {
  Asset,
  Node,
  NodeInitialProps,
  NodeWithChildren,
  Visibility,
} from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { handle } from "@/api/client";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPath } from "../useLibraryPath";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  slug: z.string().optional(),
  content: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  node: NodeWithChildren;
  initialEditingState?: boolean;
  editable?: boolean;
  onVisibilityChange?: (v: Visibility) => Promise<void>;
  onSave: (c: NodeInitialProps) => Promise<void>;
  onDelete?: (c: Node) => Promise<void>;
};

export function useLibraryPageScreen({
  node,
  initialEditingState = false,
  editable = true,
  onVisibilityChange,
  onSave,
  onDelete,
}: Props) {
  const libraryPath = useLibraryPath();
  const account = useSession();
  const router = useRouter();
  const [editing, setEditing] = useState(initialEditingState);
  const { revalidate } = useLibraryMutation();
  const isNew = !node.id;

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
    }),
    [node],
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: defaults,
  });

  const { name } = form.watch();

  useEffect(() => {
    if (isNew && !form.getFieldState("slug").isDirty) {
      const autoSlug = slugify(name);
      form.setValue("slug", autoSlug);
    }
  }, [isNew, form, name]);

  function handleEditMode() {
    if (editing) {
      setEditing(false);
      form.reset(node);

      if (isNew) {
        router.back();
      }
    } else {
      if (!isAllowedToEdit) return;

      setEditing(true);
      form.reset(node);
    }
  }

  function handleSave(payload: Form) {
    if (!editing) return;

    handle(
      async () => {
        await onSave(payload);
        setEditing(false);
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Page saved!",
        },
        cleanup: () => revalidate(),
      },
    );
  }

  async function handleVisibilityChange(v: Visibility) {
    if (!onVisibilityChange) return;

    handle(
      async () => {
        onVisibilityChange(v);
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: match(v)
            .with(Visibility.published, () => "Published")
            .with(Visibility.draft, () => "Set to draft")
            .with(Visibility.review, () => "Submitted for review")
            .with(Visibility.unlisted, () => "Set to unlisted")
            .exhaustive(),
        },
        cleanup: () => revalidate(),
      },
    );
  }

  function handleDelete() {
    if (editing) return;

    onDelete?.(node);
  }

  async function handleAssetUpload(asset: Asset) {
    if (!editing) return;

    // We only want to run these updates for edits of existing nodes.
    if (!node.id) return;

    await nodeAddAsset(node.slug, asset.id);
  }

  async function handleAssetRemove(asset: Asset) {
    if (!editing) return;
    if (!node.id) return;

    await nodeRemoveAsset(node.slug, asset.id);
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
    libraryPath,
    editing,
    node,
    isAllowedToEdit,
    isSaving: form.formState.isSubmitting,
    isAllowedToDelete,
  };
}
