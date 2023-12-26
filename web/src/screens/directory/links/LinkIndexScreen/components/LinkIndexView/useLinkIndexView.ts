import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { KeyedMutator } from "swr";
import { z } from "zod";

import { LinkListResult } from "src/api/openapi/schemas";

export type Props = {
  links: LinkListResult;
  mutate?: KeyedMutator<LinkListResult>;
  query?: string;
  page?: number;
};

export const FormSchema = z.object({
  q: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useLinkIndexView(props: Props) {
  const router = useRouter();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });

  const { q } = form.watch();

  const handlePage = async (page: number) => {
    router.push(`/p?q=${q}&page=${page}`);
  };

  const handleSubmission = form.handleSubmit(async (payload) => {
    router.push(`/l?q=${payload.q}`);
  });

  const handleReset = async () => {
    form.reset();
    router.push("/l");
  };

  const handleMutate = async () => {
    await props.mutate?.();
  };

  return {
    form,
    data: {
      q,
      links: props.links,
    },
    handlers: {
      handlePage,
      handleSubmission,
      handleReset,
      handleMutate,
    },
  };
}
