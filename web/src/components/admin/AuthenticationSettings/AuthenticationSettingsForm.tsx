import { createListCollection } from "@ark-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { OAuthAutonomousRegistrationMode, Role } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import { RadioGroupCardField } from "@/components/ui/radio-group";
import { SectionHeading } from "@/components/ui/section-heading";
import { SelectField } from "@/components/ui/select";
import { Text } from "@/components/ui/text";
import {
  AuthenticationModeDetail,
  AuthenticationModeList,
  AuthenticationModeSchema,
} from "@/lib/auth/mode";
import {
  RegistrationModeList,
  RegistrationModeSchema,
} from "@/lib/auth/registration-mode";
import { isGuestRole } from "@/lib/role/defaults";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

export type Props = {
  settings: AdminSettings;
  roles: Role[];
  canManageRoles: boolean;
  canApproveRegistrations: boolean;
};

export const FormSchema = z.object({
  authentication_mode: AuthenticationModeSchema,
  registration_mode: RegistrationModeSchema,
  autonomous_registration_mode: z
    .enum(OAuthAutonomousRegistrationMode)
    .optional(),
  autonomous_registration_role_id: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

type AuthenticationModeDetailEnabled = AuthenticationModeDetail & {
  enabled: boolean;
};

export function useAuthenticationSettingsForm(props: Props) {
  const { updateSettings } = useSettingsMutation();
  const capabilities = new Set(props.settings.capabilities);
  const oauthEnabled = capabilities.has("oauth");
  const oauth = props.settings.services?.oauth;
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      authentication_mode: props.settings.authentication_mode,
      registration_mode: props.settings.registration_mode,
      autonomous_registration_mode:
        oauth?.autonomous_registration_mode ?? "disabled",
      autonomous_registration_role_id:
        oauth?.autonomous_registration_role_id ?? "",
    },
  });

  const { dirtyFields } = form.formState;
  const handleSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateSettings({
          authentication_mode: data.authentication_mode,
          registration_mode: data.registration_mode,
          ...(oauthEnabled && {
            services: {
              oauth: {
                autonomous_registration_mode: data.autonomous_registration_mode,
                ...(props.canManageRoles &&
                  dirtyFields.autonomous_registration_role_id && {
                    autonomous_registration_role_id:
                      data.autonomous_registration_role_id || null,
                  }),
              },
            },
          }),
        });
        form.reset(data);
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
            enabled: capabilities.has("email_client"),
          };
        case "phone":
          return {
            ...m,
            enabled: capabilities.has("sms_client"),
          };
      }
    });

  return {
    form,
    availableModes,
    oauthEnabled,
    handleSubmit,
  };
}

export function AuthenticationSettingsForm(props: Props) {
  const { form, availableModes, oauthEnabled, handleSubmit } =
    useAuthenticationSettingsForm(props);
  const registrationEnabled =
    props.settings.services?.oauth?.dynamic_registration_enabled === true;
  const configuredRole =
    props.settings.services?.oauth?.autonomous_registration_role_id;
  const missingRole =
    configuredRole && !props.roles.some((role) => role.id === configuredRole);
  const roles = createListCollection({
    items: [
      { label: "Member only (no additional role)", value: "" },
      ...props.roles
        .filter((role) => !isGuestRole(role))
        .map((role) => ({ label: role.name, value: role.id })),
      ...(missingRole
        ? [{ label: `Deleted role (${configuredRole})`, value: configuredRole }]
        : []),
    ],
  });

  return (
    <styled.form className={lstack({ gap: "4" })} onSubmit={handleSubmit}>
      <LStack gap="1">
        <WStack>
          <PageHeading>Authentication settings</PageHeading>
          <Button type="submit" loading={form.formState.isSubmitting}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Configure how members sign in and who can register.
        </Text>
      </LStack>

      <FormControl>
        <FormLabel>Authentication mode</FormLabel>
        <RadioGroupCardField
          ariaLabel="Authentication mode"
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

      <FormControl>
        <FormLabel>Registration mode</FormLabel>
        <RadioGroupCardField
          ariaLabel="Registration mode"
          control={form.control}
          name="registration_mode"
          items={RegistrationModeList.map((m) => ({
            value: m.value,
            label: m.name,
            description: m.description,
          }))}
        />
        <FormErrorText>
          {form.formState.errors["registration_mode"]?.message}{" "}
        </FormErrorText>
      </FormControl>

      <LStack gap="1">
        <SectionHeading>Agent registration</SectionHeading>
        <Text variant="supporting">
          Choose how autonomous agents create their own bot accounts.
          Applications connecting on behalf of a member still use that member’s
          consent.
        </Text>
        {!oauthEnabled && (
          <Text variant="supporting">
            Enable OAuth on the server to configure agent registration.
          </Text>
        )}
      </LStack>

      {oauthEnabled && !registrationEnabled && (
        <Text variant="supporting">
          Enable dynamic client registration in OAuth settings to allow agent
          registration.
        </Text>
      )}

      <FormControl>
        <FormLabel>Agent registration policy</FormLabel>
        <RadioGroupCardField
          control={form.control}
          name="autonomous_registration_mode"
          ariaLabel="Agent registration policy"
          items={[
            {
              value: "disabled",
              label: "Disabled",
              description: "Agents cannot create new accounts.",
            },
            {
              value: "approval",
              label: "Admin approval",
              description:
                "Review each new account request. Agents with a registration token can register immediately.",
            },
            {
              value: "protected",
              label: "Registration token required",
              description:
                "Only agents with an admin-issued registration token can create accounts.",
            },
            {
              value: "open",
              label: "Open",
              description: "Any agent can create an account without approval.",
            },
          ].map((item) => ({
            ...item,
            disabled: !oauthEnabled || !registrationEnabled,
          }))}
        />
        <FormErrorText>
          {form.formState.errors.autonomous_registration_mode?.message}
        </FormErrorText>
      </FormControl>

      <FormControl>
        <FormLabel>Default agent role</FormLabel>
        <SelectField
          control={form.control}
          name="autonomous_registration_role_id"
          collection={roles}
          placeholder="Member only (no additional role)"
          ariaLabel="Default agent role"
          disabled={!oauthEnabled || !props.canManageRoles}
        />
        <Text variant="supporting">
          New bot accounts receive this role in addition to Member. Their roles
          determine their permissions. Changes apply when an account is created,
          including requests awaiting approval.
        </Text>
        {missingRole && (
          <Text variant="supporting">
            The saved role was deleted. New accounts receive Member permissions
            until you select another role.
          </Text>
        )}
        {!props.canManageRoles && (
          <Text variant="supporting">
            Managing the default role requires permission to manage roles.
          </Text>
        )}
        <FormErrorText>
          {form.formState.errors.autonomous_registration_role_id?.message}
        </FormErrorText>
      </FormControl>

      {oauthEnabled && props.canApproveRegistrations && (
        <LinkButton href="/admin/authentication/agents">
          Manage agent registrations
        </LinkButton>
      )}

      <FormErrorText>{form.formState.errors["root"]?.message} </FormErrorText>

      <WStack justifyContent="end">
        <Button type="submit" loading={form.formState.isSubmitting}>
          Save
        </Button>
      </WStack>
    </styled.form>
  );
}
