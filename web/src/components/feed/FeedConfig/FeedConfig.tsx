"use client";

import {
  ListCollection,
  SelectValueChangeDetails,
  createListCollection,
} from "@ark-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useQueryState } from "nuqs";
import { Controller, ControllerProps, useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { EditAction } from "@/components/site/Action/Edit";
import { SaveAction } from "@/components/site/Action/Save";
import {
  Editing,
  EditingSchema,
} from "@/components/site/SiteContextPane/useSiteContextPane";
import { Unready } from "@/components/site/Unready";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { CheckIcon } from "@/components/ui/icons/Check";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import {
  FeedLayoutConfigSchema,
  FeedSourceConfigSchema,
} from "@/lib/settings/feed";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { HStack, styled } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { refreshFeed } from "../../../lib/feed/refresh";

type Props = {
  initialSession?: Account;
  initialSettings: Settings;
};

export const FormSchema = z.object({
  layout: FeedLayoutConfigSchema,
  source: FeedSourceConfigSchema,
});
export type Form = z.infer<typeof FormSchema>;

export function useFeedConfig({ initialSession, initialSettings }: Props) {
  const router = useRouter();
  const session = useSession(initialSession);
  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const { updateSettings, revalidate } = useSettingsMutation(initialSettings);

  const { ready, error, settings } = useSettings(initialSettings);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: settings?.metadata.feed,
  });

  if (!ready) {
    return {
      ready: false as const,
      error,
    };
  }

  const isEditingEnabled = hasPermission(session, Permission.MANAGE_SETTINGS);
  const isEditing = editing === "feed";
  const source = settings.metadata.feed.source.type;

  function handleSetEditing() {
    setEditing("feed");
  }

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateSettings({
          metadata: {
            feed: data,
          },
        });

        setEditing(null);

        await refreshFeed();
        router.refresh();
      },
      {
        promiseToast: {
          loading: "Updating feed configuration...",
          success: "Updated!",
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  });

  return {
    ready: true as const,
    form,
    data: {
      isEditingEnabled,
      isEditing,
      source,
    },
    handlers: {
      handleSetEditing,
      handleSave,
    },
  };
}

const sources = [
  {
    label: "Threads",
    value: "threads" as const,
    icon: <DiscussionIcon width="4" />,
  },
  {
    label: "Library",
    value: "library" as const,
    icon: <LibraryIcon width="4" />,
  },
  {
    label: "Categories",
    value: "categories" as const,
    icon: <CategoryIcon width="4" />,
  },
];

export function FeedConfig(props: Props) {
  const { ready, error, form, data, handlers } = useFeedConfig(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { isEditingEnabled, isEditing, source } = data;

  if (!isEditingEnabled) {
    return null;
  }

  const collection = createListCollection({ items: sources });

  return (
    <HStack w="full" justify="end">
      {isEditing ? (
        <>
          <SelectField
            collection={collection}
            defaultValue={source}
            control={form.control}
            name="source.type"
          />

          <SaveAction onClick={handlers.handleSave}>Save feed</SaveAction>
        </>
      ) : (
        <EditAction onClick={handlers.handleSetEditing}>
          Configure feed
        </EditAction>
      )}
    </HStack>
  );
}

function SelectField<T = any>({
  collection,
  defaultValue,
  ...props
}: Omit<ControllerProps<Form>, "render"> & {
  collection: ListCollection<T>;
  defaultValue: string;
}) {
  return (
    <Controller
      {...props}
      render={({ field, formState, fieldState }) => {
        function handleChange({ value }: SelectValueChangeDetails) {
          const [v] = value;
          if (!v) return;

          field.onChange(v);
        }

        return (
          <Select.Root
            w="fit"
            size="xs"
            defaultValue={[defaultValue]}
            collection={collection}
            positioning={{ sameWidth: false }}
            onValueChange={handleChange}
          >
            <Select.Control>
              <Select.Trigger>
                <Select.ValueText placeholder="Select a Source" />
                <SelectIcon />
              </Select.Trigger>
            </Select.Control>
            <Select.Positioner>
              <Select.Content>
                {sources.map((item) => (
                  <Select.Item key={item.value} item={item}>
                    <Select.ItemText mr="2">
                      <HStack gap="1">
                        <styled.span w="4">{item.icon}</styled.span>
                        <styled.span>{item.label}</styled.span>
                      </HStack>
                    </Select.ItemText>
                    <Select.ItemIndicator>
                      <CheckIcon />
                    </Select.ItemIndicator>
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Positioner>
          </Select.Root>
        );
      }}
    />
  );
}
