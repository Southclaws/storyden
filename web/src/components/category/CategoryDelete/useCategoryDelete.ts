import { zodResolver } from "@hookform/resolvers/zod";
import { usePathname, useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";
import { categoryGet } from "@/api/openapi-client/categories";
import { useCategoryMutations } from "@/lib/category/mutation";

const FormSchema = z.object({
  move_to: z.string().min(1, "Please select a category to move posts to."),
});

export type FormData = z.infer<typeof FormSchema>;

export interface CategoryDeleteProps extends UseDisclosureProps {
  categorySlug: string;
  categoryName: string;
}

export function useCategoryDelete(props: CategoryDeleteProps) {
  const { deleteCategory } = useCategoryMutations();
  const router = useRouter();
  const pathname = usePathname();

  const form = useForm<FormData>({
    resolver: zodResolver(FormSchema),
  });

  const handleDelete = form.handleSubmit(async (data) => {
    handle(async () => {
      await deleteCategory(props.categorySlug, {
        move_to: data.move_to,
      });

      props.onClose?.();

      const isOnCategoryPage = pathname?.includes(`/d/${props.categorySlug}`);
      if (!isOnCategoryPage) {
        return;
      }

      // TODO: It would be nice to redirect to the move_to category page. But
      // we don't currently have a way to fetch a category by its ID.

      router.push("/d");
    });
  });

  return {
    form,
    handleDelete,
  };
}
