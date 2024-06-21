"use client";

import { useSession } from "src/auth";
import { AdminAction } from "src/components/site/Navigation/Anchors/Admin";
import {
  LoginAction,
  RegisterAction,
} from "src/components/site/Navigation/Anchors/Login";
import { SettingsAction } from "src/components/site/Navigation/Anchors/Settings";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";

import { Account } from "@/api/openapi/schemas";
import { HStack } from "@/styled-system/jsx";

import { ComposeAction } from "../Anchors/Compose";
import { DraftsAction } from "../Anchors/Drafts";

type Props = {
  session: Account | undefined;
};

export function Toolbar({ session }: Props) {
  const account = useSession(session);
  return (
    <HStack w="full" gap="2" alignItems="center">
      {account ? (
        <HStack w="full" alignItems="center" justify="end">
          <ComposeAction>New</ComposeAction>
          {account.admin && (
            <>
              <AdminAction />
              {/* TODO: Move public drafts for admin review to /queue */}
              {/* <QueueAction /> */}
            </>
          )}
          <DraftsAction />
          <SettingsAction />

          <ProfilePill
            profileReference={account}
            size="lg"
            showHandle={false}
          />
        </HStack>
      ) : (
        <HStack>
          <RegisterAction w="full" />
          <LoginAction flexShrink={0} />
        </HStack>
      )}
    </HStack>
  );
}
