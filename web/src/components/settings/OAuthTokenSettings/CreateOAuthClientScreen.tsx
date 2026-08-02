"use client";

import { ClipboardIcon, PlusIcon, Trash2Icon } from "lucide-react";
import { Controller, useFieldArray } from "react-hook-form";

import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import * as Clipboard from "@/components/ui/clipboard";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { Input } from "@/components/ui/input";
import { PageHeading } from "@/components/ui/page-heading";
import * as RadioGroup from "@/components/ui/radio-group";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { PermissionList } from "@/lib/permission/permission";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import {
  Form,
  OAuthClientPreset,
  Props,
  useCreateOAuthClientScreen,
} from "./useCreateOAuthClientScreen";

export function CreateOAuthClientScreen({ onClose }: Props) {
  const session = useSession();
  const { form, createdClient, handleSubmit } = useCreateOAuthClientScreen({
    onClose,
  });

  const selectablePermissions = PermissionList.filter((permission) =>
    hasPermission(session, permission.value),
  );

  const selectedPreset = form.watch("preset");
  const showRedirectUris =
    selectedPreset === "app_integration" || selectedPreset === "public_app";

  const { fields, append, remove } = useFieldArray<Form>({
    control: form.control,
    name: "redirectUris",
  });

  if (createdClient) {
    return (
      <LStack h="full" gap="8" justifyContent="space-between">
        <LStack gap="5">
          <LStack>
            <PageHeading>OAuth Client Created</PageHeading>
            <p>
              <strong>
                This is the only time you&apos;ll see this secret.
              </strong>{" "}
              Store it somewhere secure before closing this screen.
            </p>
          </LStack>

          <SecretField
            label="Client ID"
            value={createdClient.client.client_id}
          />

          {createdClient.client_secret && (
            <SecretField
              label="Client secret"
              value={createdClient.client_secret}
            />
          )}
        </LStack>

        <WStack>
          <Button onClick={onClose} style={{ width: "100%" }}>
            Done
          </Button>
        </WStack>
      </LStack>
    );
  }

  return (
    <styled.form onSubmit={handleSubmit}>
      <LStack gap="5">
        <LStack gap="1">
          <SectionHeading>Create OAuth Client</SectionHeading>
          <Text variant="supporting">
            Configure your OAuth client for different integration types.
          </Text>
        </LStack>

        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input {...form.register("name")} placeholder="e.g. My Application" />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Client Type</FormLabel>
          <Controller
            control={form.control}
            name="preset"
            render={({ field }) => (
              <RadioGroup.Root
                value={field.value}
                onValueChange={(details) =>
                  field.onChange(details.value as OAuthClientPreset)
                }
              >
                <LStack gap="2">
                  <RadioGroup.Item value="app_integration">
                    <RadioGroup.ItemControl />
                    <RadioGroup.ItemText>
                      <LStack gap="0">
                        <Text
                          as="span"
                          variant="supporting"
                          color="text.default"
                          fontWeight="medium"
                        >
                          App Integration
                        </Text>
                        <Text as="span" variant="metadata">
                          For third-party apps like MCP clients
                          (authorization_code + refresh_token, confidential,
                          PKCE required)
                        </Text>
                      </LStack>
                    </RadioGroup.ItemText>
                  </RadioGroup.Item>

                  <RadioGroup.Item value="public_app">
                    <RadioGroup.ItemControl />
                    <RadioGroup.ItemText>
                      <LStack gap="0">
                        <Text
                          as="span"
                          variant="supporting"
                          color="text.default"
                          fontWeight="medium"
                        >
                          Public App
                        </Text>
                        <Text as="span" variant="metadata">
                          For browser/mobile apps (authorization_code +
                          refresh_token, public, PKCE required)
                        </Text>
                      </LStack>
                    </RadioGroup.ItemText>
                  </RadioGroup.Item>

                  <RadioGroup.Item value="machine">
                    <RadioGroup.ItemControl />
                    <RadioGroup.ItemText>
                      <LStack gap="0">
                        <Text
                          as="span"
                          variant="supporting"
                          color="text.default"
                          fontWeight="medium"
                        >
                          Machine Client
                        </Text>
                        <Text as="span" variant="metadata">
                          For server-to-server (client_credentials,
                          confidential, no redirect URIs)
                        </Text>
                      </LStack>
                    </RadioGroup.ItemText>
                  </RadioGroup.Item>
                </LStack>
              </RadioGroup.Root>
            )}
          />
          <FormErrorText>{form.formState.errors.preset?.message}</FormErrorText>
        </FormControl>

        {showRedirectUris && (
          <FormControl>
            <FormLabel>Redirect URIs</FormLabel>
            <LStack gap="2">
              {fields.map((field, index) => (
                <LStack key={field.id} gap="1">
                  <WStack gap="2">
                    <Input
                      {...form.register(`redirectUris.${index}.value`)}
                      placeholder="https://example.com/callback"
                    />
                    <IconButton
                      type="button"
                      variant="outline"
                      onClick={() => remove(index)}
                    >
                      <Trash2Icon />
                    </IconButton>
                  </WStack>
                  <FormErrorText>
                    {
                      form.formState.errors.redirectUris?.[index]?.value
                        ?.message
                    }
                  </FormErrorText>
                </LStack>
              ))}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => append({ value: "" })}
              >
                <PlusIcon /> Add Redirect URI
              </Button>
            </LStack>
            <FormErrorText>
              {form.formState.errors.redirectUris?.message}
            </FormErrorText>
          </FormControl>
        )}

        <FormControl>
          <FormLabel>Permissions</FormLabel>
          <Controller
            control={form.control}
            name="permissions"
            render={({ field }) => (
              <LStack gap="2">
                {selectablePermissions.map((permission) => {
                  const checked = field.value.includes(permission.value);

                  return (
                    <Checkbox
                      key={permission.value}
                      size="sm"
                      checked={checked}
                      onCheckedChange={({ checked }) => {
                        const next = checked
                          ? Array.from(
                              new Set([...field.value, permission.value]),
                            )
                          : field.value.filter(
                              (value: Permission) => value !== permission.value,
                            );

                        field.onChange(next);
                      }}
                    >
                      <LStack gap="0">
                        <Text
                          as="span"
                          variant="supporting"
                          color="text.default"
                          fontWeight="medium"
                        >
                          {permission.name}
                        </Text>
                        <Text as="span" variant="metadata">
                          {permission.description}
                        </Text>
                      </LStack>
                    </Checkbox>
                  );
                })}
              </LStack>
            )}
          />
          <FormErrorText>
            {form.formState.errors.permissions?.message}
          </FormErrorText>
        </FormControl>

        <WStack>
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            style={{ flex: 1 }}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={form.formState.isSubmitting}
            style={{ flex: 1 }}
          >
            Create Client
          </Button>
        </WStack>
      </LStack>
    </styled.form>
  );
}

function SecretField({ label, value }: { label: string; value: string }) {
  return (
    <LStack gap="2">
      <Text variant="supporting" fontWeight="semibold">
        {label}
      </Text>
      <Clipboard.Root w="full" value={value}>
        <Clipboard.Control>
          <Clipboard.Input asChild>
            <Input />
          </Clipboard.Input>
          <Clipboard.Trigger asChild>
            <IconButton variant="outline">
              <Clipboard.Indicator copied={<CheckIcon />}>
                <ClipboardIcon />
              </Clipboard.Indicator>
            </IconButton>
          </Clipboard.Trigger>
        </Clipboard.Control>
      </Clipboard.Root>
    </LStack>
  );
}
