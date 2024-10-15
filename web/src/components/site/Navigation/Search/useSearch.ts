import { zodResolver } from "@hookform/resolvers/zod";
import { usePathname, useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useGetInfo } from "src/api/openapi-client/misc";

export const FormSchema = z.object({
  q: z.string().min(1, { message: "Please enter a search term" }),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  query?: string;
};

export function useSearch(props: Props) {
  const { data: infoResult } = useGetInfo();
  const router = useRouter();
  const pathname = usePathname();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });

  const title = infoResult?.title ?? "Storyden";
  const { q } = form.watch();

  const handleSearch = form.handleSubmit((data) => {
    router.push(`/search?q=${data.q}`);
  });

  const handleReset = async () => {
    form.reset();
    if (pathname === "/search") {
      router.push("/search");
    }
  };

  return {
    form,
    data: {
      q,
      title,
    },
    handlers: {
      handleSearch,
      handleReset,
    },
  };
}
