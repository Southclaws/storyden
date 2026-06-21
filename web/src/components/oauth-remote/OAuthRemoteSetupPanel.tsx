"use client";

import { type ChangeEvent, useEffect, useState } from "react";

import { handle } from "@/api/client";
import {
  oAuthRemoteConnectionAuthorize,
  oAuthRemoteConnectionCreate,
  oAuthRemoteDiscover,
} from "@/api/openapi-client/admin";
import {
  OAuthRemoteConnection,
  OAuthRemoteDiscoveryResult,
  OAuthRemoteMode,
} from "@/api/openapi-schema";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormLabel } from "@/components/ui/FormLabel";
import * as Alert from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button, ButtonGroup } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  resourceUrl: string;
  initialDiscovery?: OAuthRemoteDiscoveryResult | null;
  showResourceUrl?: boolean;
  prepareDisabled?: boolean;
  onAuthorizationUrlChange?: (url: string) => void;
  onResourceUrlChange?: (value: string) => void;
  onAuthorisationReady?: (result: {
    connection: OAuthRemoteConnection;
    authorizationUrl: string;
  }) => void | Promise<void>;
  onAuthorisationReset?: () => void | Promise<void>;
};

type ManualConfig = {
  clientId: string;
  clientSecret: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  authorizationServer: string;
  redirectUri: string;
};

const emptyManualConfig: ManualConfig = {
  clientId: "",
  clientSecret: "",
  authorizationEndpoint: "",
  tokenEndpoint: "",
  authorizationServer: "",
  redirectUri: "",
};

export function OAuthRemoteSetupPanel({
  resourceUrl,
  initialDiscovery,
  showResourceUrl = true,
  prepareDisabled = false,
  onAuthorizationUrlChange,
  onResourceUrlChange,
  onAuthorisationReady,
  onAuthorisationReset,
}: Props) {
  const [discovery, setDiscovery] = useState<OAuthRemoteDiscoveryResult | null>(
    null,
  );
  const [mode, setMode] = useState<OAuthRemoteMode | "">("");
  const [manual, setManual] = useState<ManualConfig>(emptyManualConfig);
  const [scope, setScope] = useState("");
  const [authorizationUrl, setAuthorizationUrl] = useState("");
  const [isDiscovering, setIsDiscovering] = useState(false);
  const [isPreparing, setIsPreparing] = useState(false);
  const [isResetting, setIsResetting] = useState(false);

  const selectedMode = mode || discovery?.mode || OAuthRemoteMode.manual;

  useEffect(() => {
    if (!initialDiscovery) {
      return;
    }

    applyDiscovery(initialDiscovery);
  }, [initialDiscovery]);

  function applyDiscovery(result: OAuthRemoteDiscoveryResult) {
    setDiscovery(result);
    setMode(result.mode);
    setManual((current) => ({
      ...current,
      authorizationEndpoint:
        result.authorization_server_metadata.authorization_endpoint ?? "",
      tokenEndpoint: result.authorization_server_metadata.token_endpoint ?? "",
      authorizationServer: result.authorization_server,
      redirectUri: result.redirect_uri,
    }));
  }

  async function resetPreparedAuthorisation() {
    if (!authorizationUrl) {
      return;
    }

    setIsResetting(true);
    setAuthorizationUrl("");
    onAuthorizationUrlChange?.("");

    try {
      await onAuthorisationReset?.();
    } finally {
      setIsResetting(false);
    }
  }

  async function selectMode(nextMode: OAuthRemoteMode) {
    if (nextMode === selectedMode) {
      return;
    }

    setMode(nextMode);
    await resetPreparedAuthorisation();
  }

  function updateScope(event: ChangeEvent<HTMLInputElement>) {
    setScope(event.target.value);
    void resetPreparedAuthorisation();
  }

  function updateManual(value: ManualConfig) {
    setManual(value);
    void resetPreparedAuthorisation();
  }

  async function discover() {
    await resetPreparedAuthorisation();
    setIsDiscovering(true);

    try {
      const result = await handle(
        () => oAuthRemoteDiscover({ resource_url: resourceUrl }),
        {
          promiseToast: {
            loading: "Discovering OAuth configuration...",
            success: "OAuth configuration discovered",
          },
        },
      );

      if (!result) {
        return;
      }

      applyDiscovery(result);
    } finally {
      setIsDiscovering(false);
    }
  }

  async function prepareAuthorization() {
    setIsPreparing(true);

    try {
      const connection = await handle(
        () =>
          oAuthRemoteConnectionCreate({
            resource_url: resourceUrl,
            mode: selectedMode,
            scope: scope || undefined,
            manual:
              selectedMode === OAuthRemoteMode.manual
                ? {
                    client_id: manual.clientId,
                    client_secret: manual.clientSecret || undefined,
                    authorization_endpoint: manual.authorizationEndpoint,
                    token_endpoint: manual.tokenEndpoint,
                    authorization_server:
                      manual.authorizationServer || undefined,
                    redirect_uri: manual.redirectUri || undefined,
                    scope: scope || undefined,
                  }
                : undefined,
          }),
        {
          promiseToast: {
            loading: "Saving OAuth connection...",
            success: "OAuth connection saved",
          },
        },
      );

      if (!connection) {
        return;
      }

      const authorize = await handle(
        () => oAuthRemoteConnectionAuthorize(connection.id),
        {
          promiseToast: {
            loading: "Preparing authorisation link...",
            success: "Authorisation link ready",
          },
        },
      );

      if (!authorize) {
        return;
      }

      await onAuthorisationReady?.({
        connection: authorize.connection,
        authorizationUrl: authorize.authorization_url,
      });
      setAuthorizationUrl(authorize.authorization_url);
      onAuthorizationUrlChange?.(authorize.authorization_url);
    } finally {
      setIsPreparing(false);
    }
  }

  return (
    <LStack gap="2">
      <LStack gap="1">
        <Heading size="sm">OAuth setup</Heading>
        <styled.p color="fg.muted" fontSize="sm">
          Discover the provider configuration, then prepare an authorisation
          link using authorisation code and PKCE.
        </styled.p>
      </LStack>

      {showResourceUrl && (
        <FormControl>
          <FormLabel>Resource URL</FormLabel>
          <WStack gap="2">
            <Input
              value={resourceUrl}
              onChange={(event) => onResourceUrlChange?.(event.target.value)}
              placeholder="https://mcp.example.com/mcp"
              disabled={!onResourceUrlChange}
            />
            <Button
              type="button"
              size="sm"
              variant="subtle"
              loading={isDiscovering}
              onClick={discover}
            >
              Discover
            </Button>
          </WStack>
        </FormControl>
      )}

      {!showResourceUrl && (
        <Button
          type="button"
          size="sm"
          variant="subtle"
          alignSelf="start"
          loading={isDiscovering}
          disabled={!resourceUrl}
          onClick={discover}
        >
          Rediscover OAuth settings
        </Button>
      )}

      {discovery && (
        <LStack
          gap="3"
          borderWidth="thin"
          borderStyle="solid"
          borderColor="border.default"
          borderRadius="md"
          p="3"
        >
          <WStack alignItems="start">
            <LStack gap="1">
              <Heading size="xs">
                {discovery.protected_resource_metadata.resource_name ||
                  discovery.authorization_server}
              </Heading>
              <styled.p color="fg.muted" fontSize="xs" wordBreak="break-word">
                {discovery.authorization_server}
              </styled.p>
            </LStack>
            <ModeBadge mode={selectedMode} />
          </WStack>

          <ButtonGroup attached size="xs" variant="outline">
            <Button
              type="button"
              variant={
                selectedMode === OAuthRemoteMode.cimd ? "subtle" : "outline"
              }
              onClick={() => void selectMode(OAuthRemoteMode.cimd)}
            >
              CIMD
            </Button>
            <Button
              type="button"
              variant={
                selectedMode === OAuthRemoteMode.dcr ? "subtle" : "outline"
              }
              onClick={() => void selectMode(OAuthRemoteMode.dcr)}
            >
              DCR
            </Button>
            <Button
              type="button"
              variant={
                selectedMode === OAuthRemoteMode.manual ? "subtle" : "outline"
              }
              onClick={() => void selectMode(OAuthRemoteMode.manual)}
            >
              Manual
            </Button>
          </ButtonGroup>

          {selectedMode === OAuthRemoteMode.cimd && (
            <DetailRow label="Client ID" value={discovery.client_id} />
          )}

          {selectedMode === OAuthRemoteMode.dcr && (
            <DetailRow
              label="Registration endpoint"
              value={
                discovery.authorization_server_metadata.registration_endpoint ||
                "Not advertised"
              }
            />
          )}
        </LStack>
      )}

      {selectedMode === OAuthRemoteMode.manual && (
        <ManualFields manual={manual} onChange={updateManual} />
      )}

      <FormControl>
        <FormLabel>Scopes</FormLabel>
        <Input
          value={scope}
          onChange={updateScope}
          placeholder="Optional space-separated scopes"
        />
      </FormControl>

      {authorizationUrl ? (
        <Alert.Root>
          <Alert.Content>
            <Alert.Title>Authorisation link ready</Alert.Title>
            <Alert.Description>
              Use the main action to open the provider authorisation flow.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>
      ) : (
        <Button
          type="button"
          alignSelf="end"
          loading={isPreparing || isResetting}
          onClick={prepareAuthorization}
          disabled={!resourceUrl || prepareDisabled || isResetting}
        >
          Prepare authorisation
        </Button>
      )}
    </LStack>
  );
}

function ManualFields({
  manual,
  onChange,
}: {
  manual: ManualConfig;
  onChange: (value: ManualConfig) => void;
}) {
  function update(key: keyof ManualConfig) {
    return (event: ChangeEvent<HTMLInputElement>) => {
      onChange({ ...manual, [key]: event.target.value });
    };
  }

  return (
    <LStack gap="3">
      <HStack gap="3" alignItems="start">
        <FormControl>
          <FormLabel>Client ID</FormLabel>
          <Input value={manual.clientId} onChange={update("clientId")} />
          <FormErrorText />
        </FormControl>
        <FormControl>
          <FormLabel>Client secret</FormLabel>
          <Input
            type="password"
            value={manual.clientSecret}
            onChange={update("clientSecret")}
          />
        </FormControl>
      </HStack>

      <FormControl>
        <FormLabel>Authorisation endpoint</FormLabel>
        <Input
          value={manual.authorizationEndpoint}
          onChange={update("authorizationEndpoint")}
        />
      </FormControl>

      <FormControl>
        <FormLabel>Token endpoint</FormLabel>
        <Input
          value={manual.tokenEndpoint}
          onChange={update("tokenEndpoint")}
        />
      </FormControl>

      <FormControl>
        <FormLabel>Redirect URI</FormLabel>
        <Input value={manual.redirectUri} onChange={update("redirectUri")} />
      </FormControl>
    </LStack>
  );
}

function ModeBadge({ mode }: { mode: OAuthRemoteMode }) {
  const label =
    mode === OAuthRemoteMode.cimd
      ? "CIMD"
      : mode === OAuthRemoteMode.dcr
        ? "DCR"
        : "Manual";

  return <Badge size="sm">{label}</Badge>;
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <LStack gap="0">
      <styled.span color="fg.muted" fontSize="xs">
        {label}
      </styled.span>
      <styled.span fontSize="xs" wordBreak="break-word">
        {value}
      </styled.span>
    </LStack>
  );
}
