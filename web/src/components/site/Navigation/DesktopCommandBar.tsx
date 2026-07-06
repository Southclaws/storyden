import { getServerSession } from "@/auth/server-session";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";
import { HStack, styled } from "@/styled-system/jsx";

import { SearchAnchor } from "./Anchors/Search";
import { MemberActions } from "./MemberActions";
import { Title } from "./Title";

export async function DesktopCommandBar() {
  const { title, registration_mode } = await getSettings();
  const session = await getServerSession();

  const canRegister = allowsPublicRegistration(registration_mode);

  return (
    <styled.header className="navigation__surface navigation__topbar">
      <HStack className="navigation__topbar-left">
        <SearchAnchor />
      </HStack>

      <HStack className="navigation__topbar-middle" justify="space-around">
        <Title>{title}</Title>
      </HStack>

      <HStack className="navigation__topbar-right">
        <MemberActions session={session} canRegister={canRegister} />
      </HStack>
    </styled.header>
  );
}
