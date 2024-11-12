import { zodResolver } from "@hookform/resolvers/zod";
import { usePathname, useRouter } from "next/navigation";
import { parseAsString, useQueryState } from "nuqs";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

export const FormSchema = z.object({
  q: z.string().min(1, { message: "Please enter a search term" }),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  query?: string;
  isLoading?: boolean;
};

export function useSearch(props: Props) {
  const [query, setQuery] = useSearchQueryState();

  // NOTE: This is done via a useEffect because we don't want this to be present
  // on a server-render, only for client side search interactions.
  const [isLoading, setLoading] = useState(false);
  useEffect(() => {
    setLoading(props.isLoading ?? false);
  }, [props.isLoading]);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });

  const { q } = form.watch();

  const handleSearch = form.handleSubmit((data) => {
    setQuery(data.q);
  });

  const handleReset = async () => {
    form.reset();
    setQuery(null);
  };

  return {
    form,
    data: {
      q,
      isLoading,
    },
    handlers: {
      handleSearch,
      handleReset,
    },
  };
}

export function useSearchQueryState() {
  return useQueryState("q", {
    ...parseAsString,
    defaultValue: "",
    clearOnDefault: true,
  });
}
