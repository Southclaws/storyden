import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import {
  collectionCreate,
  getCollectionListKey,
} from "src/api/openapi/collections";
import { APIError } from "src/api/openapi/schemas";
import { errorToast } from "src/components/site/ErrorBanner";
import { UseDisclosureProps } from "src/utils/useDisclosure";
import { useToast } from "src/utils/useToast";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the collection."),
  description: z.string().min(1, "Please enter a short description."),
});
export type Form = z.infer<typeof FormSchema>;

export function useCollectionCreate(props: UseDisclosureProps) {
  const toast = useToast();
  const { register, handleSubmit } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    try {
      const collection = await collectionCreate(data);
      props.onClose?.();
      mutate(getCollectionListKey());
      toast({
        title: "Collection created",
        description: `${collection.name} is now ready to be filled with stuff!`,
      });
    } catch (e: unknown) {
      errorToast(toast)(e as APIError);
    }
  });

  return {
    onSubmit,
    register,
  };
}
