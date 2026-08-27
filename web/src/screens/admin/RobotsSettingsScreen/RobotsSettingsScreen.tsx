"use client";

import { createListCollection } from "@ark-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { type MouseEvent, useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import { RequestError } from "@/api/common";
import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import {
  getRobotMCPServersListKey,
  getRobotModelsListKey,
  getRobotProvidersListKey,
  getRobotToolsListKey,
  getRobotToolsetsListKey,
  getRobotWorkspacesListKey,
  robotMCPServerDelete,
  robotMCPServerRefresh,
  robotProviderModelsRefresh,
  robotProviderUpdate,
  robotWorkspaceCreate,
  robotWorkspaceDelete,
  useRobotMCPServersList,
  useRobotModelsList,
  useRobotProvidersList,
  useRobotWorkspaceInstancesList,
  useRobotWorkspaceProvidersList,
  useRobotWorkspacesList,
} from "@/api/openapi-client/robots";
import {
  AdminSettingsServiceProps,
  RobotMCPServer,
  RobotModelInfo,
  RobotProviderMutableSettings,
  RobotProviderStatus,
  RobotWorkspace,
  RobotWorkspaceCreateBody,
  RobotWorkspaceInstance,
  RobotWorkspaceProvider,
} from "@/api/openapi-schema";
import { RobotMCPOnboardingModal } from "@/components/robots/RobotMCPOnboardingModal";
import { EmptyState } from "@/components/site/EmptyState";
import { InfoTip } from "@/components/site/InfoTip";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Unready } from "@/components/site/Unready";
import * as Alert from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { CheckboxField } from "@/components/ui/checkbox";
import { ComboboxField } from "@/components/ui/combobox";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Input } from "@/components/ui/input";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { SelectField } from "@/components/ui/select";
import { StatusBadge } from "@/components/ui/status-badge";
import { Text } from "@/components/ui/text";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { HStack, LStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { timestamp } from "@/utils/date";
import { useDisclosure } from "@/utils/useDisclosure";

const ProviderFormSchema = z.object({
  enabled: z.boolean(),
  api_key: z.string(),
});

type ProviderForm = z.infer<typeof ProviderFormSchema>;

const DefaultModelFormSchema = z.object({
  default_model: z.string(),
});

type DefaultModelForm = z.infer<typeof DefaultModelFormSchema>;

const WorkspaceFormSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string(),
  provider: z.nativeEnum(RobotWorkspaceProvider),
  allow_untrusted_commands: z.boolean(),
});

type WorkspaceForm = z.infer<typeof WorkspaceFormSchema>;

type ProviderRefreshError = {
  title: string;
  detail?: string;
};

export function RobotsSettingsScreen() {
  const providersQuery = useRobotProvidersList();
  const settingsQuery = useAdminSettingsGet();
  const modelsQuery = useRobotModelsList();
  const workspacesQuery = useRobotWorkspacesList();
  const providers = (providersQuery.data?.providers ?? []).filter(
    (provider) => provider.provider !== "mock",
  );
  const enabledProviders = providers.filter((p) => p.settings.enabled);
  const defaultModel = settingsQuery.data?.services?.robots?.default_model;
  const available = enabledProviders.length > 0 && !!defaultModel;

  return (
    <LStack gap="4">
      <LStack gap="4">
        <LStack gap="1">
          <WStack justifyContent="space-between">
            <PageHeading>Robot settings</PageHeading>

            {available ? (
              <Badge
                size="sm"
                borderColor="status.success.border"
                backgroundColor="status.success.surface"
                color="status.success.content"
              >
                Available
              </Badge>
            ) : enabledProviders.length > 0 ? (
              <Badge
                size="sm"
                borderColor="status.warning.border"
                backgroundColor="status.warning.surface"
                color="status.warning.content"
              >
                Setup incomplete
              </Badge>
            ) : (
              <Badge
                size="sm"
                borderColor="border.strong"
                backgroundColor="background.controlDisabled"
                color="text.subtle"
              >
                Disabled
              </Badge>
            )}
          </WStack>

          <Text variant="supporting">
            Configure model providers for Robots. Robots become available after
            at least one provider is enabled and a default model is selected.
          </Text>
        </LStack>

        {enabledProviders.length > 0 && !defaultModel && (
          <Text variant="metadata" color="status.warning.content">
            Robots are enabled, but no default model has been selected yet.
          </Text>
        )}

        {!providersQuery.data || !settingsQuery.data ? (
          <Unready error={providersQuery.error ?? settingsQuery.error} />
        ) : (
          <LStack gap="3">
            <GlobalDefaultModelForm
              defaultModel={defaultModel}
              enabledProviders={enabledProviders.length}
              models={modelsQuery.data?.models ?? []}
              modelsReady={!!modelsQuery.data}
              services={settingsQuery.data.services}
            />

            {providers.map((provider) => (
              <RobotProviderItem key={provider.provider} provider={provider} />
            ))}
          </LStack>
        )}
      </LStack>

      <RobotMCPServersSettings />
    </LStack>
  );
}

function RobotMCPServersSettings() {
  const serversQuery = useRobotMCPServersList();
  const workspacesQuery = useRobotWorkspacesList();
  const disclosure = useDisclosure();
  const servers = serversQuery.data?.servers ?? [];

  return (
    <LStack gap="4">
      <WStack justifyContent="space-between">
        <LStack gap="1">
          <SectionHeading>MCP servers</SectionHeading>
          <Text variant="supporting">
            Connect external MCP servers. Denbot can discover each server as a
            Toolset, which can also be assigned to custom Robots.
          </Text>
        </LStack>

        <Button type="button" variant="subtle" onClick={disclosure.onOpen}>
          Connect
        </Button>
      </WStack>

      {!serversQuery.data ? (
        <Unready error={serversQuery.error} />
      ) : servers.length === 0 ? (
        <VStack w="full">
          <EmptyState hideContributionLabel>
            No MCP servers configured yet.
          </EmptyState>
        </VStack>
      ) : (
        <LStack gap="3">
          {servers.map((server) => (
            <RobotMCPServerItem key={server.id} server={server} />
          ))}
        </LStack>
      )}

      <RobotWorkspaceTemplatesSection
        workspaces={workspacesQuery.data?.workspaces ?? []}
        ready={!!workspacesQuery.data}
        error={workspacesQuery.error}
      />

      <RobotMCPOnboardingModal
        isOpen={disclosure.isOpen}
        onClose={disclosure.onClose}
        onOpen={disclosure.onOpen}
        onOpenChange={disclosure.onOpenChange}
      />
    </LStack>
  );
}

function RobotWorkspaceTemplatesSection({
  workspaces,
  ready,
  error,
}: {
  workspaces: RobotWorkspace[];
  ready: boolean;
  error: unknown;
}) {
  const disclosure = useDisclosure();

  return (
    <LStack gap="3" pt="3">
      <WStack justifyContent="space-between" alignItems="start">
        <LStack gap="1">
          <SectionHeading>Workspace templates</SectionHeading>
          <Text variant="supporting">
            Create reusable workspace templates. Robot sessions create live
            instances from these templates when they mount a workspace.
          </Text>
        </LStack>

        <Button type="button" variant="subtle" onClick={disclosure.onOpen}>
          Create
        </Button>
      </WStack>

      {!ready ? (
        <Unready error={error} />
      ) : workspaces.length === 0 ? (
        <VStack w="full">
          <EmptyState hideContributionLabel>
            No workspace templates yet.
          </EmptyState>
        </VStack>
      ) : (
        <LStack gap="2">
          {workspaces.map((workspace) => (
            <RobotWorkspaceTemplateItem
              key={workspace.id}
              workspace={workspace}
            />
          ))}
        </LStack>
      )}

      <RobotWorkspaceCreateDrawer
        isOpen={disclosure.isOpen}
        onOpen={disclosure.onOpen}
        onClose={disclosure.onClose}
      />
    </LStack>
  );
}

function RobotWorkspaceTemplateItem({
  workspace,
}: {
  workspace: RobotWorkspace;
}) {
  const { mutate } = useSWRConfig();
  const disclosure = useDisclosure();

  async function handleDelete(event: MouseEvent) {
    event.stopPropagation();

    await handle(
      async () => {
        await robotWorkspaceDelete(workspace.id);
      },
      {
        promiseToast: {
          loading: "Deleting workspace...",
          success: "Workspace deleted",
        },
        cleanup: async () => {
          await mutate(getRobotWorkspacesListKey());
        },
      },
    );
  }

  return (
    <>
      <styled.div
        w="full"
        borderWidth="thin"
        borderStyle="solid"
        borderColor="border.default"
        borderRadius="md"
        p="2"
        cursor="pointer"
        role="button"
        tabIndex={0}
        onClick={disclosure.onOpen}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            disclosure.onOpen();
          }
        }}
      >
        <WStack justifyContent="space-between" alignItems="start">
          <LStack gap="1">
            <HStack gap="2" flexWrap="wrap">
              <Text
                variant="supporting"
                color="text.default"
                fontWeight="semibold"
              >
                {workspace.name}
              </Text>
              <Badge
                size="sm"
                borderColor="status.info.border"
                backgroundColor="status.info.surface"
                color="status.info.content"
              >
                {workspace.provider}
              </Badge>
              {workspace.allow_untrusted_commands ? (
                <Badge
                  size="sm"
                  borderColor="status.warning.border"
                  backgroundColor="status.warning.surface"
                  color="status.warning.content"
                >
                  Shell
                </Badge>
              ) : null}
            </HStack>

            {workspace.description ? (
              <Text variant="metadata">{workspace.description}</Text>
            ) : (
              <Text variant="metadata">No description.</Text>
            )}
          </LStack>

          <Button
            type="button"
            variant="ghost"
            intent="destructive"
            onClick={handleDelete}
          >
            Delete
          </Button>
        </WStack>
      </styled.div>

      <RobotWorkspaceInstancesDrawer
        workspace={workspace}
        isOpen={disclosure.isOpen}
        onClose={disclosure.onClose}
      />
    </>
  );
}

function RobotWorkspaceInstancesDrawer({
  workspace,
  isOpen,
  onClose,
}: {
  workspace: RobotWorkspace;
  isOpen: boolean;
  onClose: () => void;
}) {
  const { data, error } = useRobotWorkspaceInstancesList(
    {},
    {
      swr: {
        enabled: isOpen,
      },
    },
  );
  const instances =
    data?.workspace_instances.filter(
      (instance) => instance.workspace_id === workspace.id,
    ) ?? [];

  return (
    <ModalDrawer title={workspace.name} isOpen={isOpen} onClose={onClose}>
      <LStack gap="4">
        <LStack gap="1">
          <Text variant="supporting">
            Live workspace instances created from this template.
          </Text>
          <Text variant="metadata">Template ID: {workspace.id}</Text>
        </LStack>

        {!data ? (
          <Unready error={error} />
        ) : instances.length === 0 ? (
          <CardBox>
            <Text variant="supporting">No live instances yet.</Text>
          </CardBox>
        ) : (
          <LStack gap="2">
            {instances.map((instance) => (
              <RobotWorkspaceInstanceRow
                key={instance.id}
                instance={instance}
              />
            ))}
          </LStack>
        )}
      </LStack>
    </ModalDrawer>
  );
}

function RobotWorkspaceInstanceRow({
  instance,
}: {
  instance: RobotWorkspaceInstance;
}) {
  return (
    <CardBox>
      <LStack gap="2">
        <WStack justifyContent="space-between">
          <Text variant="supporting" color="text.default" fontWeight="semibold">
            {instance.id}
          </Text>
          <Badge
            size="sm"
            borderColor="status.info.border"
            backgroundColor="status.info.surface"
            color="status.info.content"
          >
            {instance.provider}
          </Badge>
        </WStack>

        <LStack gap="1">
          <Text variant="metadata">
            Created {timestamp(instance.createdAt, false)}
          </Text>
          <Text variant="metadata">
            Updated {timestamp(instance.updatedAt, false)}
          </Text>
        </LStack>
      </LStack>
    </CardBox>
  );
}

function RobotWorkspaceCreateDrawer({
  isOpen,
  onOpen,
  onClose,
}: {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
}) {
  return (
    <ModalDrawer
      title="New workspace"
      isOpen={isOpen}
      onOpen={onOpen}
      onClose={onClose}
    >
      <RobotWorkspaceCreateForm onClose={onClose} />
    </ModalDrawer>
  );
}

function RobotWorkspaceCreateForm({ onClose }: { onClose: () => void }) {
  const { mutate } = useSWRConfig();
  const providersQuery = useRobotWorkspaceProvidersList();
  const providers = providersQuery.data?.providers ?? [];
  const providerCollection = useMemo(
    () =>
      createListCollection({
        items: providers.map((provider) => ({
          label: provider.name,
          value: provider.provider,
        })),
      }),
    [providers],
  );
  const form = useForm<WorkspaceForm>({
    defaultValues: {
      name: "",
      description: "",
      allow_untrusted_commands: false,
    },
    resolver: zodResolver(WorkspaceFormSchema),
  });

  useEffect(() => {
    const [firstProvider] = providers;

    if (!firstProvider || form.getValues("provider")) {
      return;
    }

    form.setValue("provider", firstProvider.provider, {
      shouldValidate: true,
    });
  }, [form, providers]);

  const handleSave = form.handleSubmit(async (data) => {
    const payload: RobotWorkspaceCreateBody = {
      name: data.name,
      description: data.description,
      provider: data.provider,
      allow_untrusted_commands: data.allow_untrusted_commands,
    };

    await handle(
      async () => {
        await robotWorkspaceCreate(payload);
        form.reset();
        onClose();
      },
      {
        promiseToast: {
          loading: "Creating workspace...",
          success: "Workspace created",
        },
        cleanup: async () => {
          await mutate(getRobotWorkspacesListKey());
        },
      },
    );
  });

  return (
    <styled.form className={lstack()} gap="4" onSubmit={handleSave}>
      <FormControl>
        <FormLabel>Name</FormLabel>
        <Input {...form.register("name")} placeholder="Workspace name" />
        <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
      </FormControl>

      <FormControl>
        <FormLabel>Description</FormLabel>
        <Input
          {...form.register("description")}
          placeholder="What this workspace is for"
        />
        <FormErrorText>
          {form.formState.errors.description?.message}
        </FormErrorText>
      </FormControl>

      <FormControl>
        <FormLabel>Provider</FormLabel>
        <SelectField
          control={form.control}
          name="provider"
          collection={providerCollection}
          placeholder="Select a provider"
          positioning={{ sameWidth: true }}
          disabled={providers.length === 0}
        />
        {providersQuery.error ? <Unready error={providersQuery.error} /> : null}
        <FormErrorText>{form.formState.errors.provider?.message}</FormErrorText>
      </FormControl>

      <CheckboxField
        control={form.control}
        name="allow_untrusted_commands"
        size="sm"
      >
        Allow untrusted commands
      </CheckboxField>

      <WStack justifyContent="end">
        <Button type="button" variant="ghost" onClick={onClose}>
          Cancel
        </Button>
        <Button
          type="submit"
          loading={form.formState.isSubmitting}
          disabled={providers.length === 0}
        >
          Create
        </Button>
      </WStack>
    </styled.form>
  );
}
function RobotMCPServerItem({ server }: { server: RobotMCPServer }) {
  const { mutate } = useSWRConfig();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  async function refresh() {
    setIsRefreshing(true);
    try {
      await handle(() => robotMCPServerRefresh(server.id), {
        promiseToast: {
          loading: "Refreshing MCP tools...",
          success: "MCP tools refreshed",
        },
        cleanup: async () => {
          await Promise.all([
            mutate(getRobotMCPServersListKey()),
            mutate(getRobotToolsListKey()),
            mutate(getRobotToolsetsListKey()),
          ]);
        },
      });
    } finally {
      setIsRefreshing(false);
    }
  }

  async function remove() {
    setIsDeleting(true);
    try {
      await handle(() => robotMCPServerDelete(server.id), {
        promiseToast: {
          loading: "Deleting MCP server...",
          success: "MCP server deleted",
        },
        cleanup: async () => {
          await Promise.all([
            mutate(getRobotMCPServersListKey()),
            mutate(getRobotToolsListKey()),
            mutate(getRobotToolsetsListKey()),
          ]);
        },
      });
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <styled.div
      w="full"
      borderWidth="thin"
      borderStyle="solid"
      borderColor="border.default"
      borderRadius="md"
      p="3"
    >
      <LStack gap="3">
        <WStack alignItems="start">
          <LStack gap="1">
            <HStack gap="2" flexWrap="wrap">
              <Text
                variant="supporting"
                color="text.default"
                fontWeight="semibold"
              >
                {server.name}
              </Text>
              <Badge size="sm">{server.enabled ? "Enabled" : "Disabled"}</Badge>
              {server.has_bearer_token && <Badge size="sm">Bearer</Badge>}
              {server.oauth_remote_connection_id && (
                <Badge size="sm">
                  {server.has_oauth_token ? "OAuth" : "Pending OAuth"}
                </Badge>
              )}
            </HStack>
            <Text variant="metadata" wordBreak="break-word">
              {server.endpoint_url}
            </Text>
          </LStack>

          <HStack gap="2">
            <Button
              type="button"
              variant="outline"
              loading={isRefreshing}
              onClick={refresh}
            >
              Refresh
            </Button>
            <Button
              type="button"
              variant="ghost"
              intent="destructive"
              loading={isDeleting}
              onClick={remove}
            >
              Delete
            </Button>
          </HStack>
        </WStack>

        <Text variant="metadata">
          {server.tools.length} tools available through its MCP Toolset.
        </Text>

        {server.last_error && (
          <Text variant="metadata" color="status.danger.content">
            {server.last_error}
          </Text>
        )}
      </LStack>
    </styled.div>
  );
}

function RobotProviderItem({ provider }: { provider: RobotProviderStatus }) {
  const { mutate } = useSWRConfig();
  const disclosure = useDisclosure();

  async function refreshModels() {
    await handle(
      async () => {
        await robotProviderModelsRefresh(provider.provider);
      },
      {
        promiseToast: {
          loading: "Refreshing models...",
          success: "Models refreshed",
        },
        cleanup: async () => {
          await mutate(getRobotProvidersListKey());
          await mutate(getRobotModelsListKey());
        },
      },
    );
  }

  const lastRefreshed = provider.cache.last_refreshed_at
    ? `${timestamp(provider.cache.last_refreshed_at, false)}`
    : "Never refreshed";

  return (
    <styled.div
      w="full"
      borderWidth="thin"
      borderStyle="solid"
      borderColor="border.default"
      borderRadius="md"
      p="2"
    >
      <LStack gap="3">
        <WStack justifyContent="space-between" alignItems="start">
          <LStack gap="1">
            <HStack gap="2" flexWrap="wrap">
              <Text
                variant="supporting"
                color="text.default"
                fontWeight="semibold"
              >
                {provider.provider}
              </Text>
              <ProviderStatusBadge provider={provider} />
            </HStack>

            <Text variant="metadata">
              {provider.models.length} models available. Last refresh:{" "}
              {lastRefreshed}.
            </Text>
          </LStack>

          <HStack gap="2">
            {provider.settings.has_api_key && (
              <Button type="button" variant="outline" onClick={refreshModels}>
                Refresh
              </Button>
            )}
            <Button type="button" variant="subtle" onClick={disclosure.onOpen}>
              Configure
            </Button>
          </HStack>
        </WStack>

        {provider.cache.last_error && (
          <Text variant="metadata" color="status.danger.content">
            {provider.cache.last_error}
          </Text>
        )}
      </LStack>

      <RobotProviderConfigurationDrawer
        provider={provider}
        isOpen={disclosure.isOpen}
        onOpen={disclosure.onOpen}
        onClose={disclosure.onClose}
      />
    </styled.div>
  );
}

function ProviderStatusBadge({ provider }: { provider: RobotProviderStatus }) {
  if (provider.settings.enabled) {
    return (
      <StatusBadge size="sm" tone="success">
        Enabled
      </StatusBadge>
    );
  }

  if (provider.settings.has_api_key) {
    return (
      <StatusBadge size="sm" tone="info">
        Configured
      </StatusBadge>
    );
  }

  return (
    <StatusBadge size="sm" tone="neutral">
      Disabled
    </StatusBadge>
  );
}

function RobotProviderConfigurationDrawer({
  provider,
  isOpen,
  onOpen,
  onClose,
}: {
  provider: RobotProviderStatus;
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
}) {
  return (
    <ModalDrawer
      title={`Configure ${provider.provider}`}
      isOpen={isOpen}
      onOpen={onOpen}
      onClose={onClose}
    >
      <RobotProviderConfigurationForm
        provider={provider}
        isOpen={isOpen}
        onClose={onClose}
      />
    </ModalDrawer>
  );
}

function RobotProviderConfigurationForm({
  provider,
  isOpen,
  onClose,
}: {
  provider: RobotProviderStatus;
  isOpen: boolean;
  onClose: () => void;
}) {
  const { mutate } = useSWRConfig();
  const [refreshing, setRefreshing] = useState(false);
  const [refreshError, setRefreshError] = useState<
    ProviderRefreshError | undefined
  >();
  const form = useForm<ProviderForm>({
    defaultValues: {
      enabled: provider.settings.enabled,
      api_key: "",
    },
  });
  const enabled = form.watch("enabled");
  const providerName = provider.provider;
  const hasAPIKey = provider.settings.has_api_key;
  const storedRefreshError = provider.cache.last_error
    ? {
        title: "Model refresh failed",
        detail: provider.cache.last_error,
      }
    : undefined;
  const displayedRefreshError = refreshError ?? storedRefreshError;

  useEffect(() => {
    if (!isOpen || !enabled || !hasAPIKey) {
      return;
    }

    let cancelled = false;

    async function refreshOnOpen() {
      setRefreshing(true);
      setRefreshError(undefined);

      try {
        await robotProviderModelsRefresh(providerName);
        if (cancelled) {
          return;
        }

        await mutate(getRobotProvidersListKey());
        await mutate(getRobotModelsListKey());
      } catch (error) {
        if (!cancelled) {
          setRefreshError(parseProviderRefreshError(error));
        }
      } finally {
        if (!cancelled) {
          setRefreshing(false);
        }
      }
    }

    void refreshOnOpen();

    return () => {
      cancelled = true;
    };
  }, [enabled, hasAPIKey, isOpen, mutate, providerName]);

  const handleSave = form.handleSubmit(async (data) => {
    const patch: RobotProviderMutableSettings = {};
    const apiKey = data.api_key.trim();

    if (data.enabled !== provider.settings.enabled) {
      patch.enabled = data.enabled;
    }

    if (!data.enabled) {
      if (provider.settings.has_api_key || apiKey) {
        patch.clear_api_key = true;
      }
    } else if (apiKey) {
      patch.api_key = apiKey;
    }

    if (Object.keys(patch).length === 0) {
      onClose();
      return;
    }

    await handle(
      async () => {
        await robotProviderUpdate(provider.provider, patch);
        onClose();
      },
      {
        promiseToast: {
          loading: "Saving provider...",
          success: "Provider saved",
        },
        cleanup: async () => {
          await mutate(getRobotProvidersListKey());
          await mutate(getRobotModelsListKey());
        },
      },
    );
  });

  return (
    <styled.form className={lstack()} gap="4" onSubmit={handleSave}>
      {displayedRefreshError && (
        <Alert.Root colorPalette="orange">
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>{displayedRefreshError.title}</Alert.Title>
            {displayedRefreshError.detail && (
              <Alert.Description>
                {displayedRefreshError.detail}
              </Alert.Description>
            )}
          </Alert.Content>
        </Alert.Root>
      )}

      <CheckboxField control={form.control} name="enabled" size="sm">
        Enable provider
      </CheckboxField>

      <FormControl>
        <FormLabel>API key</FormLabel>
        <Input
          {...form.register("api_key")}
          type="password"
          placeholder={
            provider.settings.has_api_key
              ? "Leave blank to keep existing key"
              : "Provider API key"
          }
          autoComplete="off"
          disabled={!enabled}
        />
        <FormErrorText>{form.formState.errors.api_key?.message}</FormErrorText>
      </FormControl>

      <WStack justifyContent="end">
        <Button type="button" variant="ghost" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" loading={form.formState.isSubmitting}>
          Save
        </Button>
      </WStack>
    </styled.form>
  );
}

function GlobalDefaultModelForm({
  defaultModel,
  enabledProviders,
  models,
  modelsReady,
  services,
}: {
  defaultModel: string | undefined;
  enabledProviders: number;
  models: RobotModelInfo[];
  modelsReady: boolean;
  services: AdminSettingsServiceProps | undefined;
}) {
  const { updateSettings, revalidate } = useSettingsMutation();
  const form = useForm<DefaultModelForm>({
    defaultValues: {
      default_model: defaultModel ?? "",
    },
  });

  useEffect(() => {
    form.reset({ default_model: defaultModel ?? "" });
  }, [defaultModel, form]);

  const disabled =
    enabledProviders === 0 || !modelsReady || models.length === 0;
  const helperText =
    enabledProviders === 0
      ? "You must configure at least one model provider to set the default model."
      : !modelsReady
        ? "Loading available models..."
        : models.length === 0
          ? "Refresh models for a configured provider before selecting a default model."
          : undefined;

  const onSubmit = form.handleSubmit(async (data) => {
    if (!data.default_model || data.default_model === defaultModel) {
      return;
    }

    await handle(
      async () => {
        await updateSettings({
          services: {
            ...services,
            robots: {
              ...services?.robots,
              default_model: data.default_model,
            },
          },
        });

        form.reset(data);
        await revalidate();
      },
      {
        promiseToast: {
          loading: "Saving default model...",
          success: "Default model saved",
        },
      },
    );
  });

  return (
    <styled.form w="full" onSubmit={onSubmit}>
      <FormControl>
        <HStack gap="1">
          <FormLabel mb="0">Default model</FormLabel>
          <InfoTip title="Robot default model">
            The model used by Denbot and by custom Robots without their own
            model.
          </InfoTip>
        </HStack>

        <WStack alignItems="center">
          <ComboboxField
            control={form.control}
            name="default_model"
            items={models.map((model) => ({
              label: model.ref,
              value: model.ref,
            }))}
            placeholder={
              disabled ? "Configure a provider first" : "Select default model"
            }
            ariaLabel="Select default model"
            disabled={disabled}
          />

          <Button
            type="submit"
            size="sm"
            flexShrink="0"
            loading={form.formState.isSubmitting}
            disabled={disabled || !form.formState.isDirty}
          >
            Save
          </Button>
        </WStack>

        {helperText && (
          <Text variant="metadata" color="status.warning.content">
            {helperText}
          </Text>
        )}
        <FormErrorText>
          {form.formState.errors.default_model?.message}
        </FormErrorText>
      </FormControl>
    </styled.form>
  );
}

function parseProviderRefreshError(error: unknown): ProviderRefreshError {
  if (error instanceof RequestError) {
    return {
      title: error.problem?.title ?? error.message,
      detail: error.problem?.detail ?? error.message,
    };
  }

  if (error instanceof Error) {
    return {
      title: error.message,
    };
  }

  return {
    title: "Model refresh failed",
    detail: String(error),
  };
}
