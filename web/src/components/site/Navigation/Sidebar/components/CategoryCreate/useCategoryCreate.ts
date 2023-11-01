import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { categoryCreate, useCategoryList } from "src/api/openapi/categories";
import { useGetInfo } from "src/api/openapi/misc";
import { APIError } from "src/api/openapi/schemas";
import { errorToast } from "src/components/site/ErrorBanner";
import { UseDisclosureProps, useToast } from "src/theme/components";

export const FormSchema = z.object({
  name: z.string(),
  description: z.string(),
  colour: z.string().default("#fff"), // not implemented yet
  admin: z.boolean().default(false), // not implemented yet
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryCreate(props: UseDisclosureProps) {
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { mutate, data: existing } = useCategoryList();
  const { mutate: mutateInfoStatus } = useGetInfo();

  const onSubmit = handleSubmit(async (data) => {
    try {
      const category = await categoryCreate(data);
      const updated = [...(existing?.categories ?? []), category];

      mutateInfoStatus();
      mutate(
        { categories: updated },
        { populateCache: true, rollbackOnError: true },
      );

      props.onClose?.();
      toast({
        title: "Category created",
        description: `${category.name} is now ready to be filled with stuff!`,
      });
    } catch (e: unknown) {
      errorToast(toast)(e as APIError);
    }
  });

  return {
    onSubmit,
    register,
    errors,
  };
}
