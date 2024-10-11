import { zodResolver } from "@hookform/resolvers/zod";
import slugify from "@sindresorhus/slugify";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import { match } from "ts-pattern";
import { z } from "zod";

import { nodeAddAsset, nodeRemoveAsset } from "src/api/openapi-client/nodes";
import {
  Asset,
  NodeWithChildren,
  Permission,
  Visibility,
} from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { handle } from "@/api/client";
import { useLibraryMutation } from "@/lib/library/library";
import { hasPermissionOr } from "@/utils/permissions";

import { useLibraryPath } from "../useLibraryPath";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  slug: z.string().optional(),
  link: z.preprocess((v) => {
    if (typeof v === "string" && v === "") {
      return undefined;
    }

    return v;
  }, z.string().url("Invalid URL").optional()),

  content: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  node: NodeWithChildren;
};

export function useLibraryPageScreen({ node }: Props) {
  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });
  const libraryPath = useLibraryPath();
  const account = useSession();
  const { revalidate, updateNode, updateNodeVisibility, deleteNode } =
    useLibraryMutation();

  const isAllowedToEdit = hasPermissionOr(
    account,
    () => account?.id === node.owner.id,
    Permission.MANAGE_LIBRARY,
  );

  const isAllowedToDelete = hasPermissionOr(
    account,
    () => account?.id === node.owner.id,
    Permission.MANAGE_LIBRARY,
  );

  const defaults: Form = useMemo(
    () => ({
      name: node.name,
      slug: node.slug,
      link: node.link?.url,
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
    if (!form.getFieldState("slug").isDirty) {
      const autoSlug = slugify(name);
      form.setValue("slug", autoSlug);
    }
  }, [form, name]);

  function handleEditMode() {
    if (editing) {
      setEditing(false);

      form.reset({
        name: node.name,
        slug: node.slug,
        link: node.link?.url,
        content: node.content,
      });
    } else {
      if (!isAllowedToEdit) return;

      setEditing(true);

      form.reset({
        name: node.name,
        slug: node.slug,
        link: node.link?.url,
        content: node.content,
      });
    }
  }

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(
      async () => {
        const isRedirecting = await updateNode(node.slug, {
          ...payload,
          url: payload.link,
        });

        if (!isRedirecting) {
          // NOTE: This modifies the previous URL state, if updateNode received
          // a new slug, it will redirect to the new path automatically. This
          // causes the page to reload before the new slug is pushed to the URL.
          // So to fix this, we only call setEditing if the slug hasn't changed.
          setEditing(false);
        }
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Page saved!",
        },
        cleanup: () => revalidate(),
      },
    );
  });

  async function handleVisibilityChange(v: Visibility) {
    await handle(
      async () => {
        await updateNodeVisibility(node.id, v);
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

  async function handleDelete() {
    await handle(
      async () => {
        await deleteNode(node.slug);
      },
      {
        promiseToast: {
          loading: "Deleting...",
          success: "Page deleted!",
        },
        cleanup: () => revalidate(),
      },
    );
  }

  async function handleAssetUpload(asset: Asset) {
    await handle(async () => {
      await nodeAddAsset(node.slug, asset.id);
    });
  }

  async function handleAssetRemove(asset: Asset) {
    await handle(async () => {
      await nodeRemoveAsset(node.slug, asset.id);
    });
  }

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
