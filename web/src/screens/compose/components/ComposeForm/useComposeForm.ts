import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { threadCreate, threadUpdate } from "src/api/openapi-client/threads";
import { Thread, ThreadInitialProps, Visibility } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { NO_CATEGORY_VALUE } from "@/components/category/CategorySelect/useCategorySelect";

export type Props = { editing?: string; initialDraft?: Thread };

export const FormShapeSchema = z.object({
  title: z.string().default(""),
  body: z.string().min(1),
  category: z.string().optional(),
  tags: z.string().array().optional(),
  url: z.string().optional(),
});
export type FormShape = z.infer<typeof FormShapeSchema>;

export function useComposeForm({ initialDraft, editing }: Props) {
  const router = useRouter();

  const [isPublishing, setIsPublishing] = useState(false);
  const [isSavingDraft, setIsSavingDraft] = useState(false);

  const form = useForm<FormShape>({
    resolver: zodResolver(FormShapeSchema),
    reValidateMode: "onChange",
    defaultValues: initialDraft
      ? {
          title: initialDraft.title,
          body: initialDraft.body,
          tags: initialDraft.tags.map((t) => t.name),
          url: initialDraft.link?.url,
        }
      : {},
  });

  const saveDraft = async (data: FormShape) => {
    const payload: ThreadInitialProps = {
      ...data,

      // When saving a new draft, these are optional but must be explicitly set.
      title: data.title ?? "",
      body: data.body ?? "",
      url: data.url ?? "",
      tags: data.tags ?? [],
      category: data.category === NO_CATEGORY_VALUE ? undefined : data.category,

      visibility: Visibility.draft,
    };

    if (editing) {
      await threadUpdate(editing, payload);
    } else {
      const { id } = await threadCreate(payload);
      router.push(`/new?id=${id}`);
    }
  };

  const publish = async ({ title, body, category, tags, url }: FormShape) => {
    if (title.length < 1) {
      form.setError("title", {
        message: "Your post must have a title to be published",
      });
      return;
    }

    if (editing) {
      const { slug } = await threadUpdate(editing, {
        title,
        body,
        category: category === NO_CATEGORY_VALUE ? undefined : category,
        visibility: Visibility.published,
        tags,
        url,
      });
      router.push(`/t/${slug}`);
    } else {
      const { slug } = await threadCreate({
        title,
        body,
        category: category === NO_CATEGORY_VALUE ? undefined : category,
        visibility: Visibility.published,
        tags,
        url,
      });
      router.push(`/t/${slug}`);
    }
  };

  const handleSaveDraft = form.handleSubmit((data) =>
    handle(
      async () => {
        setIsSavingDraft(true);
        await saveDraft(data);
      },
      {
        promiseToast: {
          loading: "Saving draft...",
          success: "Draft saved!",
        },
        cleanup: async () => {
          setIsSavingDraft(false);
        },
      },
    ),
  );

  const handlePublish = form.handleSubmit((data) =>
    handle(
      async () => {
        setIsPublishing(true);
        await publish(data);
      },
      {
        promiseToast: {
          loading: "Publishing post...",
          success: "Post published!",
        },
        cleanup: async () => {
          setIsPublishing(false);
        },
      },
    ),
  );

  const handleAssetUpload = async () => {
    await handle(
      async () => {
        setIsSavingDraft(true);
        const state = form.getValues();
        await saveDraft(state);
      },
      {
        promiseToast: {
          loading: "Saving draft...",
          success: "Draft saved!",
        },
        cleanup: async () => {
          setIsSavingDraft(false);
        },
      },
    );
  };

  const handleAssetDelete = async () => {
    await handle(
      async () => {
        setIsSavingDraft(true);
        const state = form.getValues();
        await saveDraft(state);
      },
      {
        promiseToast: {
          loading: "Saving draft...",
          success: "Draft saved!",
        },
        cleanup: async () => {
          setIsSavingDraft(false);
        },
      },
    );
  };

  function handleBack() {
    router.back();
  }

  return {
    form,
    state: {
      isPublishing,
      isSavingDraft,
    },
    handlers: {
      handleSaveDraft,
      handlePublish,
      handleAssetDelete,
      handleAssetUpload,
      handleBack,
    },
  };
}
