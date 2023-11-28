import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { linkList } from "src/api/openapi/links";
import { Link } from "src/api/openapi/schemas";

export type Props = {
  links: Link[];
  query?: string;
};

export const FormSchema = z.object({
  q: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useLinkIndexView(props: Props) {
  const router = useRouter();
  const [results, setResults] = useState<Link[]>(props.links);
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });

  const { q } = form.watch();

  const handleSubmission = form.handleSubmit(async (payload) => {
    const { links } = await linkList(payload);

    router.push(`/l?q=${payload.q}`);
    setResults(links);
  });

  const handleReset = async () => {
    const { links } = await linkList();

    form.reset();
    router.push("/l");
    setResults(links);
  };

  return {
    form,
    data: {
      q,
      links: results,
    },
    handlers: {
      handleSubmission,
      handleReset,
    },
  };
}
