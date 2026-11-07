import { useId } from "react";

import { Button } from "@/components/ui/button";
import { CheckboxField } from "@/components/ui/checkbox";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import {
  ToggleMatrix,
  type ToggleMatrixValue,
  type ToggleMatrixValueChangeDetails,
} from "@/components/ui/toggle-matrix";
import type { Member } from "@/lib/settings/member-settings";
import {
  getNotificationGroupState,
  getNotificationGroups,
} from "@/lib/settings/notifications";
import type { Settings } from "@/lib/settings/settings";
import { LStack, WStack, styled } from "@/styled-system/jsx";

import { useMemberNotificationSettings } from "./useMemberNotificationSettings";

type Props = {
  session: Member;
  settings: Settings;
};

export function MemberNotificationSettings({ session, settings }: Props) {
  const descriptionId = useId();
  const form = useMemberNotificationSettings(session);
  const groups = getNotificationGroups(session, settings);
  const value: ToggleMatrixValue = {};

  for (const group of groups) {
    value[group.id] = {
      inbox: getNotificationGroupState(form.events, group.events, "inbox"),
      push: getNotificationGroupState(form.events, group.events, "push"),
    };
  }

  function handleValueChange({
    row,
    column,
    pressed,
  }: ToggleMatrixValueChangeDetails) {
    const group = groups.find((group) => group.id === row);
    if (!group || (column !== "inbox" && column !== "push")) {
      return;
    }

    form.setGroupEnabled(group.events, column, pressed);
  }

  return (
    <styled.form
      width="full"
      display="flex"
      flexDirection="column"
      gap="4"
      onSubmit={form.onSubmit}
    >
      <LStack gap="1">
        <WStack>
          <PageHeading>Notifications</PageHeading>
          <Button type="submit" loading={form.formState.isSubmitting}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Choose which notifications reach your inbox and devices.
        </Text>
      </LStack>

      <LStack gap="3">
        <SectionHeading>Inbox behaviour</SectionHeading>
        <FormControl>
          <CheckboxField control={form.control} name="markReadOnClick">
            Mark notifications as read when opened
          </CheckboxField>
          <FormHelperText>
            Turn this off to keep notifications unread until you mark them as
            read yourself.
          </FormHelperText>
        </FormControl>
      </LStack>

      <LStack gap="3">
        <LStack gap="1">
          <SectionHeading>Notification delivery</SectionHeading>
          <Text variant="supporting" id={descriptionId}>
            Push also enables Inbox; turning Inbox off turns Push off.
          </Text>
        </LStack>

        {settings.web_push_public_key ? (
          <FormControl>
            <CheckboxField control={form.control} name="pushEnabled">
              Receive push notifications
            </CheckboxField>
            <FormHelperText>
              Applies to all your devices. Your browser must also allow
              notifications.
            </FormHelperText>
          </FormControl>
        ) : (
          <Text variant="supporting">
            Push notifications are not enabled on this site.
          </Text>
        )}

        <ToggleMatrix
          label="Notification delivery"
          aria-describedby={descriptionId}
          rowHeading="Notify me about"
          rows={groups.map((group) => ({ ...group, value: group.id }))}
          columns={[
            { value: "inbox", label: "Inbox" },
            {
              value: "push",
              label: "Push",
              disabled: !settings.web_push_public_key || !form.pushEnabled,
            },
          ]}
          value={value}
          onValueChange={handleValueChange}
          disabled={form.formState.isSubmitting}
        />
      </LStack>
    </styled.form>
  );
}
