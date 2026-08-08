import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import { getCollectionListKey } from "@/api/openapi-client/collections";
import { Account } from "@/api/openapi-schema";
import { useCollectionMutations } from "@/lib/collection/mutation";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the collection."),
  description: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = UseDisclosureProps & {
  session: Account;
};

export function useCollectionCreate({ session, ...props }: Props) {
  const { register, handleSubmit } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const { create, revalidate } = useCollectionMutations(session);

  const onSubmit = handleSubmit(async (data) => {
    await handle(
      async () => {
        await create(data);
        props.onClose?.();
        mutate(getCollectionListKey());
      },
      {
        promiseToast: {
          loading: "Creating collection...",
          success: "Collection created!",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  });

  return {
    onSubmit,
    register,
  };
}
