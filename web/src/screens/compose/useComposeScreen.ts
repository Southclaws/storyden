import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/router";
import { useForm } from "react-hook-form";
import { useCategoryList } from "src/api/openapi/categories";
import { ThreadCreateOKResponse } from "src/api/openapi/schemas";
import { threadCreate } from "src/api/openapi/threads";
import { errorToast } from "src/components/ErrorBanner";
import { z } from "zod";

export const ThreadCreateSchema = z.object({
  title: z.string().min(1),
  body: z.string().min(1),
  category: z.string(),
  tags: z.string().array().optional(),
});
export type ThreadCreate = z.infer<typeof ThreadCreateSchema>;

export function useComposeScreen() {
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

  const onSubmit = async ({ title, body, category }: ThreadCreate) => {
    await threadCreate({
      title,
      body,
      category,
      tags: [],
    })
      .then((thread: ThreadCreateOKResponse) =>
        router.push(`/t/${thread.slug}`)
      )
      .catch(errorToast(toast));
  };

  return {
    isValid,
    onSubmit,
    handleSubmit,
    control,
    register,
    errors,
    isSubmitting,
  };
}
