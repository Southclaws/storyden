import { zodResolver } from "@hookform/resolvers/zod";
import { marked } from "marked";
import { use } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { threadCreate } from "@/api/openapi-client/threads";

type DatagraphRef = {
  id: string;
  kp: string;
  href: string;
};

export type Props = {
  contentMarkdown: string;
  sources: DatagraphRef[];
  onFinish?: () => void;
};

export const FormSchema = z.object({
  category: z.string(),
  title: z.string(),
  content: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useCreateThreadFromResult(props: Props) {
  const contentHTML = marked.parse(props.contentMarkdown, {
    async: false,
  });

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      category: undefined,
      title: "",
      content: contentHTML,
    },
  });

  console.log(form.watch());

  const handleSubmit = form.handleSubmit(async (data: Form) => {
    await handle(async () => {
      await threadCreate({
        title: data.title,
        body: data.content,
        category: data.category,
        visibility: "published",
      });

      props.onFinish?.();
    });
  });

  return {
    form,
    contentHTML,
    handlers: {
      handleSubmit,
    },
  };
}
