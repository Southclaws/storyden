import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import {
  categoryUpdate,
  useCategoryList,
} from "src/api/openapi-client/categories";
import { useGetInfo } from "src/api/openapi-client/misc";
import { APIError, Category } from "src/api/openapi-schema";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";

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
    await handle(async () => {
      const collection = await categoryUpdate(props.category.id, data);
      const updated = [...(existing?.categories ?? []), collection];

      mutateInfoStatus();
      mutate(
        { categories: updated },
        { populateCache: true, rollbackOnError: true },
      );

      props.onClose?.();
    });
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
