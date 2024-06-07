import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { categoryUpdate, useCategoryList } from "src/api/openapi/categories";
import { useGetInfo } from "src/api/openapi/misc";
import { APIError, Category } from "src/api/openapi/schemas";
import { handleError } from "src/components/site/ErrorBanner";
import { UseDisclosureProps } from "src/utils/useDisclosure";

export type Props = {
  category: Category;
} & UseDisclosureProps;

export const FormSchema = z.object({
  name: z.string().min(1),
  description: z.string().min(1),
  colour: z.string().default("#fff"), // not implemented yet
  admin: z.boolean().default(false), // not implemented yet
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryEdit(props: Props) {
  const {
    reset,
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
    values: {
      name: props.category.name,
      description: props.category.description,
      colour: props.category.colour,
      admin: props.category.admin,
    },
  });
  const { mutate, data: existing } = useCategoryList();
  const { mutate: mutateInfoStatus } = useGetInfo();

  const onSubmit = handleSubmit(async (data) => {
    try {
      const collection = await categoryUpdate(props.category.id, data);
      const updated = [...(existing?.categories ?? []), collection];

      mutateInfoStatus();
      mutate(
        { categories: updated },
        { populateCache: true, rollbackOnError: true },
      );

      props.onClose?.();
    } catch (e: unknown) {
      handleError(e as APIError);
    }
  });

  function onCancel() {
    reset();
    props.onClose?.();
  }

  return {
    onSubmit,
    onCancel,
    register,
    errors,
  };
}
