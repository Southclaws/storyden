"use client";

import { Controller } from "react-hook-form";
import { ClipboardIcon } from "lucide-react";

import { useSession } from "@/auth";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import * as Clipboard from "@/components/ui/clipboard";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { Input } from "@/components/ui/input";
import { Permission } from "@/api/openapi-schema";
import { PermissionList } from "@/lib/permission/permission";
import { hasPermission } from "@/utils/permissions";
import { LStack, WStack, styled } from "@/styled-system/jsx";

import {
  Form,
  Props,
  useCreateOAuthClientScreen,
} from "./useCreateOAuthClientScreen";

export function CreateOAuthClientScreen({ onClose }: Props) {
  const session = useSession();
  const { form, createdClient, handleSubmit } = useCreateOAuthClientScreen({
    onClose,
  });

  if (createdClient) {
    return (
      <LStack h="full" gap="8" justifyContent="space-between">
        <LStack gap="5">
          <LStack>
            <Heading>OAuth Client Created</Heading>
            <p>
              <strong>This is the only time you&apos;ll see this secret.</strong>{" "}
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

  const selectablePermissions = PermissionList.filter((permission) =>
    hasPermission(session, permission.value),
  );

  return (
    <styled.form onSubmit={handleSubmit}>
      <LStack gap="5">
        <LStack gap="1">
          <Heading size="sm">Create OAuth Client</Heading>
          <styled.p color="fg.muted" fontSize="sm">
            Choose the permissions this integration may request. OAuth clients
            use client credentials and never gain permissions your account does
            not already hold.
          </styled.p>
        </LStack>

        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input {...form.register("name")} placeholder="e.g. Analytics Sync" />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

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
                        <styled.span fontWeight="medium">
                          {permission.name}
                        </styled.span>
                        <styled.span color="fg.muted" fontSize="xs">
                          {permission.description}
                        </styled.span>
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
      <Heading size="sm" color="fg.muted">
        {label}
      </Heading>
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
