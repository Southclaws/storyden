import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/router";
import { useForm } from "react-hook-form";
import { useCategoryList } from "src/api/openapi/categories";
import { ThreadCreateOKResponse, ThreadStatus } from "src/api/openapi/schemas";
import { threadCreate, threadUpdate } from "src/api/openapi/threads";
import { errorToast } from "src/components/ErrorBanner";
import { z } from "zod";

export type Props = { editing?: string };

export const ThreadCreateSchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  category: z.string(),
  tags: z.string().array().optional(),
});
export type ThreadCreate = z.infer<typeof ThreadCreateSchema>;

export function useComposeScreen({ editing }: Props) {
  const router = useRouter();
  const toast = useToast();
  const { data } = useCategoryList();
  const {
    handleSubmit,
    control,
    register,
    formState: { isValid, errors, isSubmitting },
  } = useForm<ThreadCreate>({
    resolver: zodResolver(ThreadCreateSchema),
    reValidateMode: "onChange",
    defaultValues: {
      // hack: the underlying category list select component can't do this.
      category: data?.categories[0]?.id,
    },
  });

  const onSave = handleSubmit(async (props: ThreadCreate) => {
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

      router.push(
        {
          pathname: "/new",
          query: { id, edit: "true" },
        },
        `/new?id=${id}`,
        { shallow: true }
      );
    }
  });

  const onPublish = handleSubmit(
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
    isValid,
    onSave,
    onPublish,
    control,
    register,
    errors,
    isSubmitting,
  };
}
