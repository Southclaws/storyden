import { FormEvent, useState } from "react";
import { mutate } from "swr";
import { parse as parseYAML } from "yaml";

import { fetcher, handle } from "@/api/client";
import { getPluginListKey } from "@/api/openapi-client/plugins";
import {
  PluginGetOKResponse,
  PluginInitialExternal,
  PluginInitialExternalMode,
} from "@/api/openapi-schema";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { deriveError } from "@/utils/error";
import { UseDisclosureProps } from "@/utils/useDisclosure";

const defaultPayload = `{
  "id": "my-external-plugin",
  "name": "My External Plugin",
  "description": "External plugin",
  "version": "0.0.1",
  "author": "you",
  "events_consumed": []
}`;

export function PluginAddExternal({ onClose }: UseDisclosureProps) {
  const [payload, setPayload] = useState(defaultPayload);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSubmitting(true);
    setError(null);

    let jsonPayload: PluginInitialExternal;

    try {
      jsonPayload = parseExternalPayload(payload);
    } catch (err) {
      setError(deriveError(err));
      setIsSubmitting(false);
      return;
    }

    await handle(
      async () => {
        await fetcher<PluginGetOKResponse>({
          url: "/plugins",
          method: "POST",
          headers: { "Content-Type": "application/json" },
          data: jsonPayload,
        });
        await mutate(getPluginListKey());
        onClose?.();
      },
      {
        async onError(err) {
          const error = deriveError(err);
          setError(error);
        },
        async cleanup() {
          setIsSubmitting(false);
        },
      },
    );
  };

  const handleClose = () => {
    if (!isSubmitting) {
      setError(null);
      onClose?.();
    }
  };

  return (
    <form className={lstack({ gap: "4" })} onSubmit={handleSubmit}>
      <styled.p color="fg.muted">
        Register an external plugin that connects to Storyden over authenticated
        RPC. Storyden will not manage this plugin process.
      </styled.p>

      <Alert.Root>
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>Security notice</Alert.Title>
          <Alert.Description>
            External plugin tokens grant full access over RPC. Keep them secret
            and only run plugins from trusted sources.
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>

      <styled.label fontSize="sm" fontWeight="medium">
        Manifest (YAML or JSON)
      </styled.label>

      <styled.textarea
        value={payload}
        onChange={(event) => setPayload(event.currentTarget.value)}
        disabled={isSubmitting}
        minH="72"
        w="full"
        resize="vertical"
        borderWidth="thin"
        borderColor="border.default"
        borderRadius="md"
        p="3"
        bgColor="bg.default"
        color="fg.default"
        fontFamily="mono"
        fontSize="xs"
        lineHeight="tight"
        spellCheck={false}
      />

      {error && (
        <Alert.Root colorPalette="red">
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>Configuration Error</Alert.Title>
            <Alert.Description>{error}</Alert.Description>
          </Alert.Content>
        </Alert.Root>
      )}

      <WStack justifyContent="end" gap="2">
        <Button
          type="button"
          variant="outline"
          onClick={handleClose}
          disabled={isSubmitting}
        >
          {isSubmitting ? "Adding..." : "Cancel"}
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Adding..." : "Add External Plugin"}
        </Button>
      </WStack>
    </form>
  );
}

function parseExternalPayload(raw: string): PluginInitialExternal {
  const value = raw.trim();
  let parsed: unknown;

  try {
    parsed = JSON.parse(value) as unknown;
  } catch {
    parsed = parseYAML(value) as unknown;
  }

  if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
    throw new Error("Manifest payload must be an object");
  }

  return {
    mode: PluginInitialExternalMode.external,
    manifest: parsed as Record<string, unknown>,
  };
}
