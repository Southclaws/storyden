import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import {
  collectionCreate,
  getCollectionListKey,
} from "src/api/openapi-client/collections";
import { APIError } from "src/api/openapi-schema";
import { handleError } from "src/components/site/ErrorBanner";
import { UseDisclosureProps } from "src/utils/useDisclosure";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the collection."),
  description: z.string().min(1, "Please enter a short description."),
});
export type Form = z.infer<typeof FormSchema>;

export function useCollectionCreate(props: UseDisclosureProps) {
  const { register, handleSubmit } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    try {
      await collectionCreate(data);
      props.onClose?.();
      mutate(getCollectionListKey());
    } catch (e: unknown) {
      handleError(e as APIError);
    }
  });

  return {
    onSubmit,
    register,
  };
}
