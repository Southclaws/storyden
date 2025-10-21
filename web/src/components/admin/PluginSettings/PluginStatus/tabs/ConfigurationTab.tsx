import { FormEvent, useEffect, useState } from "react";
import { match } from "ts-pattern";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getPluginGetConfigurationKey,
  usePluginGetConfiguration,
  usePluginGetConfigurationSchema,
  usePluginUpdateConfiguration,
} from "@/api/openapi-client/plugins";
import {
  PluginConfiguration,
  PluginConfigurationFieldUnion,
} from "@/api/openapi-schema";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { deriveError } from "@/utils/error";
import { Box, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  pluginID: string;
};

export function ConfigurationTab({ pluginID }: Props) {
  const { mutate } = useSWRConfig();
  const { data: schema } = usePluginGetConfigurationSchema(pluginID);
  const { data: configuration } = usePluginGetConfiguration(pluginID);
  const { trigger: updateConfiguration } =
    usePluginUpdateConfiguration(pluginID);

  const [values, setValues] = useState<PluginConfiguration>({});
  const [dirty, setDirty] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (configuration) {
      setValues(configuration);
      setDirty(false);
    }
  }, [configuration]);

  function handleFieldChange(id: string, value: unknown) {
    setValues((prev) => ({ ...prev, [id]: value }));
    setDirty(true);
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSaving(true);
    setError(null);

    await handle(
      async () => {
        await mutateTransaction(
          mutate,
          [
            {
              key: getPluginGetConfigurationKey(pluginID),
              optimistic: () => values,
              commit: (_current, result) => result,
            },
          ],
          () => updateConfiguration(values),
          { revalidate: true },
        );

        setDirty(false);
      },
      {
        errorToast: false,
        onError: async (err) => {
          setError(deriveError(err));
        },
        cleanup: async () => {
          setIsSaving(false);
        },
      },
    );
  }

  const fields = schema?.fields ?? [];

  if (fields.length === 0) {
    return (
      <Box minH="12">
        <styled.p fontSize="sm" color="fg.muted">
          This plugin has no configurable fields.
        </styled.p>
      </Box>
    );
  }

  return (
    <styled.form
      onSubmit={handleSubmit}
      display="flex"
      flexDir="column"
      gap="4"
      w="full"
    >
      <styled.p fontSize="sm" color="fg.muted">
        Configure the plugin&apos;s settings below.
      </styled.p>

      <LStack gap="3">
        {fields.map((field) => (
          <ConfigurationField
            key={fieldKey(field)}
            field={field}
            value={values[fieldKey(field)]}
            onChange={(value) => handleFieldChange(fieldKey(field), value)}
          />
        ))}
      </LStack>

      <WStack justifyContent="flex-end">
        <Button
          type="submit"
          size="sm"
          variant="subtle"
          disabled={!dirty || isSaving}
          loading={isSaving}
        >
          Save Configuration
        </Button>
      </WStack>

      <Admonition
        value={!!error}
        kind="failure"
        title="Configuration Error"
        onChange={() => setError(null)}
      >
        {error && <styled.p fontSize="sm">{error}</styled.p>}
      </Admonition>
    </styled.form>
  );
}

type FieldProps = {
  field: PluginConfigurationFieldUnion;
  value: unknown;
  onChange: (value: unknown) => void;
};

function ConfigurationField({ field, value, onChange }: FieldProps) {
  const label = field.label ?? fieldKey(field);
  const description = field.description;

  const input = match(field)
    .with({ type: "string" }, (f) => (
      <styled.input
        type="text"
        value={typeof value === "string" ? value : (f.default ?? "")}
        onChange={(e) => onChange(e.currentTarget.value)}
        borderWidth="thin"
        borderColor="border.default"
        borderRadius="md"
        px="3"
        py="2"
        fontSize="sm"
        w="full"
        bgColor="bg.default"
      />
    ))
    .with({ type: "number" }, (f) => (
      <styled.input
        type="number"
        value={typeof value === "number" ? value : (f.default ?? "")}
        onChange={(e) => onChange(e.currentTarget.valueAsNumber)}
        borderWidth="thin"
        borderColor="border.default"
        borderRadius="md"
        px="3"
        py="2"
        fontSize="sm"
        w="full"
        bgColor="bg.default"
      />
    ))
    .with({ type: "boolean" }, (f) => (
      <styled.input
        type="checkbox"
        checked={typeof value === "boolean" ? value : (f.default ?? false)}
        onChange={(e) => onChange(e.currentTarget.checked)}
        w="4"
        h="4"
        cursor="pointer"
      />
    ))
    .exhaustive();

  return (
    <LStack gap="1">
      <WStack alignItems="center" gap="2">
        <styled.label fontSize="sm" fontWeight="medium">
          {label}
        </styled.label>
        {input}
      </WStack>
      {description && (
        <styled.p fontSize="xs" color="fg.muted">
          {description}
        </styled.p>
      )}
    </LStack>
  );
}

function fieldKey(field: PluginConfigurationFieldUnion): string {
  return field.id ?? field.label ?? "";
}
