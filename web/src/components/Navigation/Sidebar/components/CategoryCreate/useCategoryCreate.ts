import { useDisclosure, useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { categoryCreate } from "src/api/openapi/categories";
import { getCollectionListKey } from "src/api/openapi/collections";
import { APIError } from "src/api/openapi/schemas";
import { errorToast } from "src/components/ErrorBanner";

export const FormSchema = z.object({
  name: z.string(),
  description: z.string(),
  colour: z.string().default("#fff"), // not implemented yet
  admin: z.boolean().default(false), // not implemented yet
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryCreate() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    try {
      const collection = await categoryCreate(data);
      onClose();
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
    isOpen,
    onOpen,
    onClose,
    onSubmit,
    register,
    errors,
  };
}
