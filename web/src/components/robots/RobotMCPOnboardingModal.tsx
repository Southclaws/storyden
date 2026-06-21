"use client";

import { useEffect, useMemo, useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { oAuthRemoteDiscover } from "@/api/openapi-client/admin";
import {
  getRobotMCPServersListKey,
  robotMCPServerCreate,
  robotMCPServerDelete,
  robotMCPServerProbe,
  useRobotMCPServersList,
} from "@/api/openapi-client/robots";
import {
  OAuthRemoteDiscoveryResult,
  RobotMCPServer,
  RobotMCPServerProbeResult,
} from "@/api/openapi-schema";
import { OAuthRemoteSetupPanel } from "@/components/oauth-remote/OAuthRemoteSetupPanel";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { FormControl } from "@/components/ui/FormControl";
import { FormLabel } from "@/components/ui/FormLabel";
import * as Alert from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button, ButtonGroup } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { UseDisclosureProps } from "@/utils/useDisclosure";

type AuthMode = "oauth" | "bearer";

type Props = UseDisclosureProps;

export function RobotMCPOnboardingModal(props: Props) {
  return (
    <ModalDrawer
      title="Connect MCP server"
      size="wide"
      isOpen={props.isOpen}
      onClose={props.onClose}
      onOpen={props.onOpen}
      onOpenChange={props.onOpenChange}
    >
      <RobotMCPOnboardingScreen onClose={props.onClose} />
    </ModalDrawer>
  );
}

function RobotMCPOnboardingScreen({ onClose }: { onClose?: () => void }) {
  const { mutate } = useSWRConfig();
  const [url, setURL] = useState("");
  const [bearerToken, setBearerToken] = useState("");
  const [authMode, setAuthMode] = useState<AuthMode>("oauth");
  const [probe, setProbe] = useState<RobotMCPServerProbeResult | null>(null);
  const [oauthDiscovery, setOAuthDiscovery] =
    useState<OAuthRemoteDiscoveryResult | null>(null);
  const [server, setServer] = useState<RobotMCPServer | null>(null);
  const [pendingOAuthServerID, setPendingOAuthServerID] = useState("");
  const [authorizationUrl, setAuthorizationUrl] = useState("");
  const [isProbing, setIsProbing] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const serversQuery = useRobotMCPServersList({
    swr: {
      enabled: authMode === "oauth" && pendingOAuthServerID !== "",
      revalidateOnFocus: true,
    },
  });

  const endpointURL =
    oauthDiscovery?.resource_url || probe?.endpoint_url || url;
  const cardTitle = probe?.server_card?.title || probe?.server_card?.name;
  const serverName = useMemo(() => {
    const resourceName =
      oauthDiscovery?.protected_resource_metadata.resource_name;
    if (resourceName) {
      return resourceName;
    }
    if (cardTitle) {
      return cardTitle;
    }
    try {
      return new URL(endpointURL).hostname;
    } catch {
      return "MCP server";
    }
  }, [cardTitle, endpointURL, oauthDiscovery]);

  const oauthReady = authMode === "bearer" || Boolean(oauthDiscovery);
  const mcpReady =
    Boolean(probe?.active) || (authMode === "oauth" && Boolean(probe));
  const canCreateBearer = authMode === "bearer" && probe?.active;
  const waitingForAuthorisation =
    authMode === "oauth" &&
    authorizationUrl !== "" &&
    pendingOAuthServerID !== "";

  useEffect(() => {
    if (pendingOAuthServerID === "") {
      return;
    }

    const connectedServer = serversQuery.data?.servers.find(
      (candidate) =>
        candidate.id === pendingOAuthServerID &&
        candidate.has_oauth_token &&
        candidate.enabled,
    );
    if (!connectedServer) {
      return;
    }

    setServer(connectedServer);
    setAuthorizationUrl("");
    setPendingOAuthServerID("");
  }, [pendingOAuthServerID, serversQuery.data?.servers]);

  async function checkConnection() {
    setIsProbing(true);
    setProbe(null);
    setOAuthDiscovery(null);
    setServer(null);
    setPendingOAuthServerID("");
    setAuthorizationUrl("");

    try {
      let discovered: OAuthRemoteDiscoveryResult | undefined;
      if (authMode === "oauth") {
        discovered = await handle(() =>
          oAuthRemoteDiscover({ resource_url: url }),
        );
        if (discovered) {
          setOAuthDiscovery(discovered);
        }
      }

      const result = await handle(
        () =>
          robotMCPServerProbe({
            url: discovered?.resource_url || url,
            bearer_token: authMode === "bearer" ? bearerToken : undefined,
          }),
        {
          promiseToast: {
            loading: "Checking MCP endpoint...",
            success: "MCP endpoint checked",
          },
        },
      );

      if (result) {
        setProbe(result);
      }
    } finally {
      setIsProbing(false);
    }
  }

  function updateURL(value: string) {
    setURL(value);
    setProbe(null);
    setOAuthDiscovery(null);
    setServer(null);
    void resetPendingAuthorisation();
  }

  function updateAuthMode(value: AuthMode) {
    if (value === authMode) {
      return;
    }

    setAuthMode(value);
    setProbe(null);
    setOAuthDiscovery(null);
    setServer(null);
    void resetPendingAuthorisation();
  }

  async function createBearerServer() {
    setIsCreating(true);

    try {
      const created = await handle(
        () =>
          robotMCPServerCreate({
            name: serverName,
            description: probe?.server_card?.description ?? "",
            endpoint_url: endpointURL,
            bearer_token: bearerToken || undefined,
            enabled: true,
          }),
        {
          promiseToast: {
            loading: "Connecting MCP server...",
            success: "MCP server connected",
          },
          cleanup: async () => {
            await mutate(getRobotMCPServersListKey());
          },
        },
      );

      if (created) {
        setServer(created);
      }
    } finally {
      setIsCreating(false);
    }
  }

  async function resetPendingAuthorisation() {
    const serverID = pendingOAuthServerID;

    setAuthorizationUrl("");
    setPendingOAuthServerID("");

    if (!serverID) {
      return;
    }

    await handle(
      async () => {
        await robotMCPServerDelete(serverID);
      },
      {
        cleanup: async () => {
          await mutate(getRobotMCPServersListKey());
        },
      },
    );
  }

  if (server) {
    return (
      <LStack gap="6">
        <LStack gap="2">
          <Heading size="sm">{server.name} connected</Heading>
          <styled.p color="fg.muted" fontSize="sm">
            {server.tools.length} tools were discovered and added to the Robot
            tool catalogue.
          </styled.p>
        </LStack>

        <Button alignSelf="end" onClick={onClose}>
          Done
        </Button>
      </LStack>
    );
  }

  return (
    <styled.div
      display="grid"
      gridTemplateColumns={{ base: "1fr", lg: "minmax(0, 1fr) minmax(0, 1fr)" }}
      gap="2"
      alignItems="start"
    >
      <LStack gap="2" h="full" justifyContent="space-between">
        <LStack gap="2">
          <FormControl>
            <FormLabel>MCP server URL</FormLabel>
            <HStack gap="0">
              <Input
                value={url}
                onChange={(event) => updateURL(event.target.value)}
                placeholder="https://mcp.example.com"
                borderRightRadius="none"
              />
              <Button
                type="button"
                variant="subtle"
                alignSelf="start"
                borderLeftRadius="none"
                loading={isProbing}
                disabled={!url || waitingForAuthorisation}
                onClick={checkConnection}
              >
                {waitingForAuthorisation ? "Waiting" : "Check"}
              </Button>
            </HStack>
          </FormControl>

          <FormControl>
            <FormLabel>Authentication</FormLabel>
            <ButtonGroup attached size="sm" variant="outline">
              <Button
                type="button"
                variant={authMode === "oauth" ? "subtle" : "outline"}
                onClick={() => updateAuthMode("oauth")}
              >
                OAuth
              </Button>
              <Button
                type="button"
                variant={authMode === "bearer" ? "subtle" : "outline"}
                onClick={() => updateAuthMode("bearer")}
              >
                Bearer token
              </Button>
            </ButtonGroup>
          </FormControl>

          {authMode === "bearer" && (
            <FormControl>
              <FormLabel>Bearer token</FormLabel>
              <Input
                type="password"
                value={bearerToken}
                onChange={(event) => setBearerToken(event.target.value)}
                placeholder="Token used for the MCP probe and server"
              />
            </FormControl>
          )}
        </LStack>

        <LStack gap="2">
          <ReadinessRow
            label="OAuth setup"
            status={
              authMode === "bearer"
                ? "Not required"
                : oauthReady
                  ? "Discovered"
                  : "Not checked"
            }
            ready={oauthReady}
          />
          <MCPStatusCard
            authMode={authMode}
            probe={probe}
            serverName={serverName}
            endpointURL={endpointURL}
            ready={mcpReady}
            waitingForAuthorisation={waitingForAuthorisation}
          />
        </LStack>
      </LStack>

      <LStack gap="2">
        <Heading size="sm">OAuth advanced settings</Heading>

        <styled.div
          bgColor="bg.subtle"
          borderWidth="thin"
          borderStyle="solid"
          borderColor="border.default"
          borderRadius="md"
          p={{ base: "1", md: "2" }}
        >
          {authMode === "oauth" ? (
            <OAuthRemoteSetupPanel
              resourceUrl={endpointURL}
              initialDiscovery={oauthDiscovery}
              showResourceUrl={false}
              prepareDisabled={!mcpReady}
              onAuthorizationUrlChange={setAuthorizationUrl}
              onResourceUrlChange={updateURL}
              onAuthorisationReset={resetPendingAuthorisation}
              onAuthorisationReady={async ({
                connection,
                authorizationUrl,
              }) => {
                const server = await handle(
                  () =>
                    robotMCPServerCreate({
                      name: serverName,
                      description: probe?.server_card?.description ?? "",
                      endpoint_url: endpointURL,
                      oauth_remote_connection_id: connection.id,
                      enabled: false,
                    }),
                  {
                    promiseToast: {
                      loading: "Linking MCP server to OAuth...",
                      success: "MCP server will connect after authorisation",
                    },
                    cleanup: async () => {
                      await mutate(getRobotMCPServersListKey());
                    },
                  },
                );

                if (!server) {
                  throw new Error("MCP server link was not created.");
                }

                setPendingOAuthServerID(server.id);
                setAuthorizationUrl(authorizationUrl);
              }}
            />
          ) : (
            <LStack gap="3">
              <styled.p color="fg.muted" fontSize="sm">
                Bearer token mode skips OAuth. Switch back to OAuth to use CIMD,
                DCR, or manual client credentials.
              </styled.p>
            </LStack>
          )}
        </styled.div>
      </LStack>

      <LStack style={{ gridColumn: "1/-1" }}>
        <Alert.Root>
          <Alert.Content>
            <Alert.Title>External MCP servers introduce risk</Alert.Title>
            <Alert.Description>
              Robots may send user requests to this server and execute returned
              tools. Connect servers you trust.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>

        {authMode === "bearer" && (
          <Button
            type="button"
            alignSelf="end"
            loading={isCreating}
            disabled={!canCreateBearer}
            onClick={createBearerServer}
          >
            Connect server
          </Button>
        )}

        {authorizationUrl && (
          <Button asChild alignSelf="end">
            <styled.a href={authorizationUrl} target="_blank" rel="noreferrer">
              Open authorisation
            </styled.a>
          </Button>
        )}
      </LStack>
    </styled.div>
  );
}

function ReadinessRow({
  label,
  status,
  ready,
}: {
  label: string;
  status: string;
  ready: boolean;
}) {
  return (
    <WStack
      gap="3"
      borderWidth="thin"
      borderStyle="solid"
      borderColor={ready ? "border.success" : "border.default"}
      borderRadius="md"
      px="3"
      py="2"
    >
      <styled.span fontSize="sm" fontWeight="medium">
        {label}
      </styled.span>
      <Badge size="sm">{status}</Badge>
    </WStack>
  );
}

function MCPStatusCard({
  authMode,
  probe,
  serverName,
  endpointURL,
  ready,
  waitingForAuthorisation,
}: {
  authMode: AuthMode;
  probe: RobotMCPServerProbeResult | null;
  serverName: string;
  endpointURL: string;
  ready: boolean;
  waitingForAuthorisation: boolean;
}) {
  const status = waitingForAuthorisation
    ? "Waiting for auth"
    : probe?.active
      ? "Active"
      : probe
        ? authMode === "oauth"
          ? "Needs auth"
          : "Check failed"
        : "Not checked";

  return (
    <LStack
      gap="1"
      borderWidth="thin"
      borderStyle="solid"
      borderColor={ready ? "border.success" : "border.default"}
      borderRadius="md"
      p="3"
    >
      <WStack gap="1">
        <styled.span fontSize="sm" fontWeight="medium">
          MCP Endpoint
        </styled.span>
        <Badge size="sm">{status}</Badge>
      </WStack>

      {probe && (
        <LStack gap="1">
          <Heading size="xs">{serverName}</Heading>
          <styled.p color="fg.muted" fontSize="xs" wordBreak="break-word">
            {endpointURL}
          </styled.p>

          {probe.server_card?.description && (
            <styled.p color="fg.muted" fontSize="sm">
              {probe.server_card.description}
            </styled.p>
          )}

          {probe.probe_error && !waitingForAuthorisation && (
            <styled.p color="fg.warning" fontSize="xs">
              {probe.probe_error}
            </styled.p>
          )}

          {waitingForAuthorisation && (
            <styled.p color="fg.muted" fontSize="xs">
              Complete authorisation in the new tab, then return here.
            </styled.p>
          )}
        </LStack>
      )}
    </LStack>
  );
}
