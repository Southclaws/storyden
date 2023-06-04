import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useCategoryList } from "src/api/openapi/categories";
import {
  Thread,
  ThreadCreateOKResponse,
  ThreadStatus,
} from "src/api/openapi/schemas";
import { threadCreate, threadUpdate } from "src/api/openapi/threads";
import { errorToast } from "src/components/ErrorBanner";

export type Props = { editing?: string; draft?: Thread };

export const ThreadCreateSchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  category: z.string(),
  tags: z.string().array().optional(),
});
export type ThreadCreate = z.infer<typeof ThreadCreateSchema>;

export function useComposeForm({ draft, editing }: Props) {
  const router = useRouter();
  const toast = useToast();
  const { data } = useCategoryList();
  const formContext = useForm<ThreadCreate>({
    resolver: zodResolver(ThreadCreateSchema),
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

  const onSave = formContext.handleSubmit(async (props: ThreadCreate) => {
    if (editing) {
      await threadUpdate(editing, {
        ...props,
      })
        .then((thread: ThreadCreateOKResponse) => thread.id)
        .catch(errorToast(toast));
    } else {
      const id = await threadCreate({
        ...props,
        status: ThreadStatus.draft,
        tags: [],
      })
        .then((thread: ThreadCreateOKResponse) => thread.id)
        .catch(errorToast(toast));

      if (!id) return;

      router.push(`/new?id=${id}`);
    }
  });

  const onPublish = formContext.handleSubmit(
    async ({ title, body, category }: ThreadCreate) => {
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
    formContext,
  };
}
