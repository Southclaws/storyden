import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { categoryDelete, useCategoryList } from "src/api/openapi-client/categories";
import { UseDisclosureProps } from "src/utils/useDisclosure";

const FormSchema = z.object({
  move_to: z.string().min(1, "Please select a category to move posts to."),
});

export type FormData = z.infer<typeof FormSchema>;

export interface CategoryDeleteProps extends UseDisclosureProps {
  categorySlug: string;
  categoryName: string;
}

export function useCategoryDelete(props: CategoryDeleteProps) {
  const { mutate } = useCategoryList();

  const form = useForm<FormData>({
    resolver: zodResolver(FormSchema),
  });

  const handleDelete = form.handleSubmit(async (data) => {
    handle(async () => {
      await categoryDelete(props.categorySlug, {
        move_to: data.move_to,
      });

      mutate();
      props.onClose?.();
    });
  });

  return {
    form,
    handleDelete,
  };
}