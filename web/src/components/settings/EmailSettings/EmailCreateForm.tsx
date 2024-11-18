import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  accountEmailAdd,
  getAccountGetKey,
} from "@/api/openapi-client/accounts";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { SaveIcon } from "@/components/ui/icons/Save";
import { Input } from "@/components/ui/input";
import { useProfileMutations } from "@/lib/profile/mutation";
import { CardBox, HStack } from "@/styled-system/jsx";
import { wstack } from "@/styled-system/patterns";

export const FormSchema = z.object({
  email: z.string().email(),
});
export type Form = z.infer<typeof FormSchema>;

type Props = {
  onSuccess: () => void;
  onCancel: () => void;
};

export function EmailCreateForm(props: Props) {
  const { mutate } = useSWRConfig();
  const form = useForm<Form>({ resolver: zodResolver(FormSchema) });

  const handleSubmit = form.handleSubmit(async ({ email }) => {
    await handle(
      async () => {
        await accountEmailAdd({ email_address: email });
        props.onSuccess();
      },
      {
        async cleanup() {
          await mutate(getAccountGetKey());
        },
        promiseToast: {
          loading: "Adding email address...",
          success: "Email added",
        },
      },
    );
  });

  return (
    <CardBox>
      <form
        className={wstack({
          alignItems: "center",
          gap: "2",
        })}
        onSubmit={handleSubmit}
      >
        <Input
          size="sm"
          placeholder="Email address..."
          {...form.register("email")}
        />

        <HStack gap="0">
          <Button size="sm" variant="solid" borderRightRadius="none">
            <SaveIcon /> save
          </Button>
          <Button
            type="button"
            size="sm"
            variant="subtle"
            borderLeftRadius="none"
            onClick={props.onCancel}
          >
            <CancelIcon /> cancel
          </Button>
        </HStack>
      </form>
      <FormErrorText>{form.formState.errors["email"]?.message}</FormErrorText>
    </CardBox>
  );
}
