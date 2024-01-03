import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

export const FormSchema = z.object({
  q: z.string().min(1, { message: "Please enter a search term" }),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  query?: string;
};

export function useSearch(props: Props) {
  const router = useRouter();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });

  const { q } = form.watch();

  const handleSearch = form.handleSubmit((data) => {
    router.push(`/search?q=${data.q}`);
  });

  const handleReset = async () => {
    form.reset();
    router.push("/search");
  };

  return {
    form,
    data: {
      q,
    },
    handlers: {
      handleSearch,
      handleReset,
    },
  };
}
