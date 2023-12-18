import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { profileList } from "src/api/openapi/profiles";
import { PublicProfileListResult } from "src/api/openapi/schemas";

export type Props = {
  profiles: PublicProfileListResult;
  query?: string;
};

export const FormSchema = z.object({
  q: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useMemberIndexView(props: Props) {
  const router = useRouter();
  const [results, setResults] = useState(props.profiles);
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });

  const { q } = form.watch();

  const handleSubmission = form.handleSubmit(async (payload) => {
    const response = await profileList(payload);

    router.push(`/p?q=${payload.q}`);
    setResults(response);
  });

  const handleReset = async () => {
    const response = await profileList();

    form.reset();
    router.push("/p");
    setResults(response);
  };

  return {
    form,
    data: {
      q,
      results,
    },
    handlers: {
      handleSubmission,
      handleReset,
    },
  };
}
