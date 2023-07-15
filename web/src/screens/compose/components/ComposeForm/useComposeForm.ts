import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useCategoryList } from "src/api/openapi/categories";
import {
  Asset,
  Thread,
  ThreadCreateOKResponse,
  ThreadStatus,
} from "src/api/openapi/schemas";
import { threadCreate, threadUpdate } from "src/api/openapi/threads";
import { errorToast } from "src/components/ErrorBanner";

export type Props = { editing?: string; draft?: Thread };

export const ThreadMutationSchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  category: z.string(),
  tags: z.string().array().optional(),
  assets: z.array(z.string()),
});
export type ThreadMutation = z.infer<typeof ThreadMutationSchema>;

export function useComposeForm({ draft, editing }: Props) {
  const router = useRouter();
  const toast = useToast();
  const { data } = useCategoryList();
  const formContext = useForm<ThreadMutation>({
    resolver: zodResolver(ThreadMutationSchema),
    reValidateMode: "onChange",
    defaultValues: draft
      ? {
          title: draft.title,
          body: draft.posts[0]?.body,
          tags: draft.tags,
        }
      : {
          // hack: the underlying category list select component can't do this.
          category: data?.categories[0]?.id,
        },
  });

  function onBack() {
    router.back();
  }

  const doSave = async (props: ThreadMutation) => {
    const payload = {
      title: props.title,
      category: props.category,
      body: props.body,
      tags: [],
      status: ThreadStatus.draft,
      assets: props.assets,
    };

    if (editing) {
      await threadUpdate(editing, payload)
        .then((thread: ThreadCreateOKResponse) => thread.id)
        .catch(errorToast(toast));
    } else {
      const id = await threadCreate(payload)
        .then((thread: ThreadCreateOKResponse) => thread.id)
        .catch(errorToast(toast));

      router.push(`/new?id=${id}`);
    }
  };

  const onAssetUpload = async (asset: Asset) => {
    const state: ThreadMutation = formContext.getValues();

    return await doSave({
      ...state,
      assets: [...state.assets, asset.id],
    });
  };

  const onSave = formContext.handleSubmit(doSave);

  const onPublish = formContext.handleSubmit(
    async ({ title, body, category }: ThreadMutation) => {
      if (editing) {
        threadUpdate(editing, { status: ThreadStatus.published })
          .then((thread: ThreadCreateOKResponse) =>
            router.push(`/t/${thread.slug}`)
          )
          .catch(errorToast(toast));
      } else {
        await threadCreate({
          title,
          body,
          category,
          status: ThreadStatus.published,
          tags: [],
        })
          .then((thread: ThreadCreateOKResponse) =>
            router.push(`/t/${thread.slug}`)
          )
          .catch(errorToast(toast));
      }
    }
  );

  return {
    onBack,
    onSave,
    onPublish,
    onAssetUpload,
    formContext,
  };
}
