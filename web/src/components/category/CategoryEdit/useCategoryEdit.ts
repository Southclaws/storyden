import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Category } from "src/api/openapi-schema";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";
import { useCategoryMutations } from "@/lib/category/mutation";

export type Props = {
  category: Category;
} & UseDisclosureProps;

export const FormSchema = z.object({
  name: z.string().min(1),
  description: z.string().min(1),
  colour: z.string().default("#fff"),
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryEdit(props: Props) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    values: {
      name: props.category.name,
      description: props.category.description,
      colour: props.category.colour,
    },
  });

  const { revalidateList, updateCategory } = useCategoryMutations();

  const handleSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateCategory(props.category.id, data);

        props.onClose?.();
      },
      {
        promiseToast: {
          loading: "Updating category...",
          success: "Category updated.",
        },
        cleanup: () => revalidateList(),
      },
    );
  });

  function handleCancel() {
    form.reset();
    props.onClose?.();
  }

  return {
    form,
    handlers: {
      handleSubmit,
      handleCancel,
    },
  };
}
