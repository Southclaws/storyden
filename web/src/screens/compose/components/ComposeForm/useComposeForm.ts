import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useCategoryList } from "src/api/openapi/categories";
import { Asset, Thread, ThreadStatus } from "src/api/openapi/schemas";
import { threadCreate, threadUpdate } from "src/api/openapi/threads";
import { errorToast } from "src/components/ErrorBanner";

export type Props = { editing?: string; initialDraft?: Thread };

export const FormShapeSchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  category: z.string(),
  tags: z.string().array().optional(),
  assets: z.array(z.string()),
});
export type FormShape = z.infer<typeof FormShapeSchema>;

export function useComposeForm({ initialDraft, editing }: Props) {
  const router = useRouter();
  const toast = useToast();
  const { data } = useCategoryList();
  const formContext = useForm<FormShape>({
    resolver: zodResolver(FormShapeSchema),
    reValidateMode: "onChange",
    defaultValues: initialDraft
      ? {
          title: initialDraft.title,
          body: initialDraft.posts[0]?.body,
          tags: initialDraft.tags,
          assets: initialDraft.assets.map((v) => v.id),
        }
      : {
          // hack: the underlying category list select component can't do this.
          category: data?.categories[0]?.id,
          assets: [],
        },
  });

  const doSave = async (data: FormShape) => {
    const payload = {
      ...data,
      status: ThreadStatus.draft,
    };

    if (editing) {
      await threadUpdate(editing, payload);
    } else {
      const { id } = await threadCreate(payload);

      router.push(`/new?id=${id}`);
    }
  };

  const doPublish = async ({ title, body, category }: FormShape) => {
    if (editing) {
      const { slug } = await threadUpdate(editing, {
        status: ThreadStatus.published,
      });
      router.push(`/t/${slug}`);
    } else {
      const { slug } = await threadCreate({
        title,
        body,
        category,
        status: ThreadStatus.published,
        tags: [],
      });
      router.push(`/t/${slug}`);
    }
  };

  const onAssetUpload = async (asset: Asset) => {
    const state: FormShape = formContext.getValues();

    const newAssets = [...state.assets, asset.id];

    const newState = {
      ...state,
      assets: newAssets,
    };

    await doSave(newState).catch(errorToast(toast));

    formContext.setValue("assets", newAssets);
  };

  function onBack() {
    router.back();
  }

  const onSave = formContext.handleSubmit((data) =>
    doSave(data).catch(errorToast(toast))
  );

  const onPublish = formContext.handleSubmit((data) =>
    doPublish(data).catch(errorToast(toast))
  );

  return {
    onBack,
    onSave,
    onPublish,
    onAssetUpload,
    formContext,
  };
}
