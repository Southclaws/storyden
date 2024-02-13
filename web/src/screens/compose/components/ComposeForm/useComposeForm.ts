import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useCategoryList } from "src/api/openapi/categories";
import {
  Thread,
  ThreadInitialProps,
  Visibility,
} from "src/api/openapi/schemas";
import { threadCreate, threadUpdate } from "src/api/openapi/threads";
import { handleError } from "src/components/site/ErrorBanner";

export type Props = { editing?: string; initialDraft?: Thread };

export const FormShapeSchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  category: z.string(),
  tags: z.string().array().optional(),
  url: z.string().optional(),
  // assets: z.array(z.string()),
});
export type FormShape = z.infer<typeof FormShapeSchema>;

export function useComposeForm({ initialDraft, editing }: Props) {
  const router = useRouter();

  const { data } = useCategoryList();
  const formContext = useForm<FormShape>({
    resolver: zodResolver(FormShapeSchema),
    reValidateMode: "onChange",
    defaultValues: initialDraft
      ? {
          title: initialDraft.title,
          body: initialDraft.posts[0]?.body,
          tags: initialDraft.tags,
          url: initialDraft.link?.url,
        }
      : {
          // hack: the underlying category list select component can't do this.
          category: data?.categories[0]?.id,
        },
  });

  const doSave = async (data: FormShape) => {
    const payload: ThreadInitialProps = {
      ...data,

      // When saving a new draft, these are optional but must be explicitly set.
      title: data.title ?? "",
      body: data.body ?? "",
      url: data.url ?? "",
      // assets: data.assets ?? [],

      visibility: Visibility.draft,
    };

    if (editing) {
      await threadUpdate(editing, payload);
    } else {
      const { id } = await threadCreate(payload);

      router.push(`/new?id=${id}`);
    }
  };

  const doPublish = async ({ title, body, category, url }: FormShape) => {
    if (editing) {
      const { slug } = await threadUpdate(editing, {
        title,
        body,
        category,
        visibility: Visibility.published,
        tags: [],
        url,
      });
      router.push(`/t/${slug}`);
    } else {
      const { slug } = await threadCreate({
        title,
        body,
        category,
        visibility: Visibility.published,
        tags: [],
        url,
      });
      router.push(`/t/${slug}`);
    }
  };

  const onAssetUpload = async () => {
    const state = formContext.getValues();
    await doSave(state).catch(handleError);
  };

  const onAssetDelete = async () => {
    const state = formContext.getValues();
    await doSave(state).catch(handleError);
  };

  function onBack() {
    router.back();
  }

  const onSave = formContext.handleSubmit((data) =>
    doSave(data).catch(handleError),
  );

  const onPublish = formContext.handleSubmit((data) =>
    doPublish(data).catch(handleError),
  );

  return {
    onBack,
    onSave,
    onPublish,
    onAssetUpload,
    onAssetDelete,
    formContext,
  };
}
