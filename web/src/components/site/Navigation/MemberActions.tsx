"use client";

import { useSession } from "src/auth";
import {
  LoginAnchor,
  RegisterAnchor,
} from "src/components/site/Navigation/Anchors/Login";

import { Account } from "@/api/openapi-schema";
import { NotificationsMenu } from "@/components/notifications/NotificationsMenu";
import { HStack } from "@/styled-system/jsx";

import { AccountMenu } from "./AccountMenu/AccountMenu";
import { ComposeAnchor } from "./Anchors/Compose";
import { LanguageSwitcher } from "./LanguageSwitcher";

type Props = {
  session: Account | undefined;
  canRegister?: boolean;
};

export function MemberActions({ session, canRegister }: Props) {
  const account = useSession(session);

  return (
    <HStack w="full" gap="2" alignItems="center" justify="end" pr="1">
      <LanguageSwitcher />
      {account ? (
        <>
          <ComposeAnchor />
          <NotificationsMenu status="unread" />
          <AccountMenu account={account} />
        </>
      ) : (
        <>
          {canRegister && <RegisterAnchor w="full" />}
          <LoginAnchor flexShrink={0} />
        </>
      )}
    </HStack>
  );
}
