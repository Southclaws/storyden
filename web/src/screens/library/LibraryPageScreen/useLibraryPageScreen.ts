import { zodResolver } from "@hookform/resolvers/zod";
import slugify from "@sindresorhus/slugify";
import { dequal } from "dequal";
import { omit, values } from "lodash/fp";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useEffect, useMemo, useRef, useState } from "react";
import { FixedCropperRef } from "react-advanced-cropper";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { match } from "ts-pattern";
import { z } from "zod";

import { nodeAddAsset, nodeRemoveAsset } from "src/api/openapi-client/nodes";
import {
  Asset,
  InstanceCapability,
  LinkReference,
  NodeWithChildren,
  Permission,
  PropertyList,
  PropertyMutation,
  PropertyMutationList,
  Visibility,
} from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { linkCreate } from "@/api/openapi-client/links";
import { useLibraryMutation } from "@/lib/library/library";
import {
  CoverImage,
  CoverImageSchema,
  parseNodeMetadata,
} from "@/lib/library/metadata";
import { useCapability } from "@/lib/settings/capabilities";
import { getAssetURL } from "@/utils/asset";
import { hasPermissionOr } from "@/utils/permissions";

import { useLibraryPath } from "../useLibraryPath";

export const CROP_STENCIL_WIDTH = 1536;
export const CROP_STENCIL_HEIGHT = 384;

const CoverImageFormSchema = z.union([
  CoverImageSchema,
  z.object({
    asset_id: z.string(),
  }),
]);

export const FormNodePropertySchema = z.object({
  fid: z.string().optional(),
  name: z.string(),
  type: z.string(),
  sort: z.string(),
  value: z.string(),
});
export type FormNodeProperty = z.infer<typeof FormNodePropertySchema>;

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  slug: z.string().optional(),
  properties: z.array(FormNodePropertySchema),
  tags: z.string().array().optional(),
  link: z.preprocess((v) => {
    if (typeof v === "string" && v === "") {
      return undefined;
    }

    return v;
  }, z.string().url("Invalid URL").optional()),
  coverImage: CoverImageFormSchema.optional(),
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
  const {
    revalidate,
    updateNode,
    updateNodeVisibility,
    suggestTitle,
    suggestSummary,
    importFromLink,
    deleteNode,
  } = useLibraryMutation(node);
  const isTitleSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const cropperRef = useRef<FixedCropperRef>(null);

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

  const defaults = useMemo<Form>(
    () =>
      ({
        name: node.name,
        slug: node.slug,
        properties: node.properties.map((p, i) => ({
          fid: p.fid,
          name: p.name ?? `Field ${i}`,
          type: p.type ?? "text",
          sort: p.sort,
          value: p.value ?? "",
        })),
        tags: node.tags.map((t) => t.name),
        link: node.link?.url,
        content: node.content,
      }) satisfies Form,
    [node],
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: defaults,
  });

  const { name, content } = form.watch();

  useEffect(() => {
    if (!form.getFieldState("slug").isDirty) {
      const autoSlug = slugify(name);
      form.setValue("slug", autoSlug);
    }
  }, [form, name]);

  // Summary suggestion logic
  const [generatedContent, setGeneratedContent] = useState<string | undefined>(
    undefined,
  );

  function handleResetGeneratedContent() {
    setGeneratedContent(undefined);
  }

  useEffect(() => {
    if (content && generatedContent) {
      // if the content field changes, and it was previously using a controlled
      // value, reset this controlled value to move it back to uncontrolled.
      setGeneratedContent(undefined);
    }
  }, [form, content, generatedContent]);

  // Title suggestion logic
  const [generatedTitle, setGeneratedTitle] = useState<string | undefined>(
    undefined,
  );
  const [isLoadingSuggestTitle, setLoadingSuggestTitle] = useState(false);

  function handleResetGeneratedTitle() {
    setGeneratedTitle(undefined);
  }

  async function handleSuggestTitle() {
    await handle(
      async () => {
        setLoadingSuggestTitle(true);

        const title = await suggestTitle(node.id);
        if (!title) {
          throw new Error("No title could be suggested for this content.");
        }

        form.setValue("name", title);
        setGeneratedTitle(title);
      },
      {
        cleanup: async () => setLoadingSuggestTitle(false),
      },
    );
  }

  // This URL is used for the crop editor, it will always be the original image
  // depending on whether the current primary image has any new versions set.
  // The parent is always set to the originally uploaded image while the actual
  // `primary_image` field has whichever version is currently set as the cover.
  const primaryAssetEditingURL = getAssetURL(
    node.primary_image?.parent?.path ?? node.primary_image?.path,
  );

  const primaryAssetURL = getAssetURL(node.primary_image?.path);

  const initialCoverCoordinates = parseNodeMetadata(node.meta).coverImage;

  function handleEditMode() {
    if (editing) {
      setEditing(false);

      form.reset(defaults);
    } else {
      if (!isAllowedToEdit) return;

      setEditing(true);

      form.reset(defaults);
    }
  }

  const uploadCroppedCover = async () => {
    if (!cropperRef.current) {
      return;
    }

    const canvas = cropperRef.current.getCanvas();
    if (!canvas) {
      throw new Error("An unexpected error occurred with the image editor.");
    }

    const coordinates =
      cropperRef.current.getCoordinates() satisfies CoverImage | null;
    if (!coordinates) {
      throw new Error(
        "An unexpected error occurred with the image editor: unable to get crop coordinates.",
      );
    }

    const blob = await new Promise<Blob>((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (blob == null) {
          reject("An unexpected error occurred with the image editor.");
          return;
        }

        resolve(blob);
      });
    });

    if (node.primary_image) {
      // TODO: Delete the original asset maybe?
    }

    // The cover image is determined to be a copy of an original if it has a
    // parent asset associated with it. Original assets do not have parents.
    const isCopy = node.primary_image?.parent?.id !== undefined;

    const parent_asset_id = isCopy
      ? // If the primary image is already a copy, use the existing parent.
        node.primary_image?.parent?.id
      : // Otherwise, use the primary image asset ID, which will result in
        // this ID becoming the parent.
        node.primary_image?.id;

    // The result of this is that when the asset revalidates, the original
    // asset will be present in the primary image asset object. This means
    // that users will always edit the originally uploaded image and when
    // they stop editing, the saved asset will be the cropped version where
    // the original is still present within the asset's parent property.
    const asset = await assetUpload(blob, {
      // TODO: Split filename from ID in API side (Marks) use original name.
      filename: "cropped-cover",
      parent_asset_id,
    });

    return {
      isReplacement: false,
      config: coordinates,
      asset,
    };
  };

  async function handleImportFromLink(link: LinkReference) {
    await handle(
      async () => {
        const { title_suggestion, tag_suggestions, content_suggestion } =
          await importFromLink(node.slug, link.url);

        setGeneratedTitle(title_suggestion);
        setGeneratedContent(content_suggestion);

        if (title_suggestion) {
          form.setValue("name", title_suggestion);
        }
        form.setValue("tags", tag_suggestions);
        form.setValue("content", content_suggestion);
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(
      async () => {
        const coverConfig = await uploadCroppedCover();

        const isRedirecting = await updateNode(
          node.slug,
          {
            ...payload,
            properties: payload.properties.map((p) => {
              if (p.fid?.startsWith("new_field_")) {
                return omit("fid", p);
              }
              return p;
            }),
            url: payload.link,
          },
          coverConfig,
        );

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
      handleSuggestTitle,
      handleResetGeneratedContent,
      handleResetGeneratedTitle,
      handleAssetUpload,
      handleAssetRemove,
      handleImportFromLink,
    },
    libraryPath,
    editing,
    node,
    generatedContent,
    generatedTitle,
    cropperRef,
    primaryAssetURL,
    primaryAssetEditingURL,
    initialCoverCoordinates,
    isAllowedToEdit,
    isSaving: form.formState.isSubmitting,
    isAllowedToDelete,
    isTitleSuggestEnabled,
    isLoadingSuggestTitle,
  };
}
