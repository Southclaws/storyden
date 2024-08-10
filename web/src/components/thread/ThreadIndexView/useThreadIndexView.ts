import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { KeyedMutator } from "swr";
import { z } from "zod";

import { ThreadListResult } from "src/api/openapi-schema";

export type Props = {
  threads: ThreadListResult;
  mutate?: KeyedMutator<ThreadListResult>;
  query?: string;
  page?: number;
};

export const FormSchema = z.object({
  q: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useThreadIndexView(props: Props) {
  const router = useRouter();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });
  const searchParams = useSearchParams();

  const handlePage = async (page: number) => {
    const newParams = new URLSearchParams(searchParams).set(
      "page",
      page.toString(),
    );
    router.push(`/t?${newParams}`);
  };

  const handleSubmission = form.handleSubmit(async (payload) => {
    router.push(`/t?q=${payload.q}`);
  });

  const handleReset = async () => {
    form.reset();
    router.push("/t");
  };

  const handleMutate = async () => {
    await props.mutate?.();
  };

  return {
    form,
    data: {
      threads: props.threads,
    },
    handlers: {
      handlePage,
      handleSubmission,
      handleReset,
      handleMutate,
    },
  };
}
