import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import { adminSettingsUpdate } from "@/api/openapi-client/admin";
import { getGetInfoKey } from "@/api/openapi-client/misc";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Button } from "@/components/ui/button";
import { CardGroupRadio } from "@/components/ui/form/CardGroupRadio";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Heading } from "@/components/ui/heading";
import {
  AuthenticationModeDetail,
  AuthenticationModeList,
  AuthenticationModeSchema,
} from "@/lib/auth/mode";
import { Settings } from "@/lib/settings/settings";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

export type Props = {
  settings: Settings;
};

export const FormSchema = z.object({
  authentication_mode: AuthenticationModeSchema,
});
export type Form = z.infer<typeof FormSchema>;

type AuthenticationModeDetailEnabled = AuthenticationModeDetail & {
  enabled: boolean;
};

export function useAuthenticationSettingsForm(props: Props) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      authentication_mode: props.settings.authentication_mode,
    },
  });

  const handleSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await adminSettingsUpdate({
          authentication_mode: data.authentication_mode,
        });
        mutate(getGetInfoKey());
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Authentication settings updated",
        },
      },
    );
  });

  const availableModes: AuthenticationModeDetailEnabled[] =
    AuthenticationModeList.map((m) => {
      switch (m.value) {
        case "handle":
          return { ...m, enabled: true };
        case "email":
          return {
            ...m,
            enabled: props.settings.capabilities.includes("email_client"),
          };
        case "phone":
          return {
            ...m,
            enabled: props.settings.capabilities.includes("sms_client"),
          };
      }
    });

  return {
    form,
    availableModes,
    handleSubmit,
  };
}

export function AuthenticationSettingsForm(props: Props) {
  const { form, availableModes, handleSubmit } =
    useAuthenticationSettingsForm(props);

  return (
    <CardBox>
      <styled.form className={lstack()} onSubmit={handleSubmit}>
        <WStack>
          <Heading size="md">Authentication settings</Heading>
          <Button type="submit" loading={form.formState.isSubmitting}>
            Save
          </Button>
        </WStack>

        <FormControl>
          <FormLabel>Authentication mode</FormLabel>
          <CardGroupRadio
            control={form.control}
            name="authentication_mode"
            items={availableModes.map((m) => ({
              value: m.value,
              label: m.name,
              description: m.description,
              disabled: !m.enabled,
            }))}
          />
          <FormErrorText>
            {form.formState.errors["authentication_mode"]?.message}{" "}
          </FormErrorText>
        </FormControl>

        <FormErrorText>{form.formState.errors["root"]?.message} </FormErrorText>

        <WStack justifyContent="end">
          <Button type="submit" loading={form.formState.isSubmitting}>
            Save
          </Button>
        </WStack>
      </styled.form>
    </CardBox>
  );
}
