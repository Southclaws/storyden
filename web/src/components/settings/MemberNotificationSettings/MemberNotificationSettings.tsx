import { Button } from "@/components/ui/button";
import { Checkbox, CheckboxField } from "@/components/ui/checkbox";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
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
  const form = useMemberNotificationSettings(session);
  const groups = getNotificationGroups(session, settings);

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
          Choose how notifications behave and which ones reach your devices.
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
          <SectionHeading>Push notifications</SectionHeading>
          <Text variant="supporting">
            These preferences apply to all your devices. Your browser must also
            allow notifications. Your inbox is unaffected.
          </Text>
        </LStack>

        {settings.web_push_public_key ? (
          <>
            <FormControl>
              <CheckboxField control={form.control} name="pushEnabled">
                Receive push notifications
              </CheckboxField>
            </FormControl>

            <LStack gap="3" pt="1">
              <SectionHeading as="h3" fontSize="sm">
                Notify me about
              </SectionHeading>
              {groups.map((group) => (
                <FormControl key={group.id}>
                  <Checkbox
                    checked={getNotificationGroupState(
                      form.events,
                      group.events,
                    )}
                    disabled={!form.pushEnabled}
                    onCheckedChange={({ checked }) => {
                      form.setGroupEnabled(group.events, checked === true);
                    }}
                  >
                    {group.label}
                  </Checkbox>
                  <FormHelperText>{group.description}</FormHelperText>
                </FormControl>
              ))}
            </LStack>
          </>
        ) : (
          <Text variant="supporting">
            Push notifications are not enabled on this site.
          </Text>
        )}
      </LStack>
    </styled.form>
  );
}
