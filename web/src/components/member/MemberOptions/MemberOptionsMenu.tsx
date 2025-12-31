import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import Link from "next/link";
import { PropsWithChildren } from "react";
import { toast } from "sonner";

import { useSession } from "src/auth";

import { ProfileReference } from "@/api/openapi-schema";
import { ReportMemberMenuItem } from "@/components/report/ReportMemberMenuItem";
import * as Menu from "@/components/ui/menu";
import { WEB_ADDRESS } from "@/config";
import { hasPermission } from "@/utils/permissions";
import { useCopyToClipboard } from "@/utils/useCopyToClipboard";

import { MemberIdent } from "../MemberBadge/MemberIdent";
import { MemberRoleMenu } from "../MemberRoleMenu/MemberRoleMenu";
import { MemberSuspensionTrigger } from "../MemberSuspension/MemberSuspensionTrigger";

export type Props = {
  profile: ProfileReference;
  asChild?: boolean;
};

export function MemberOptionsMenu({
  children,
  profile,
  ...props
}: PropsWithChildren<Props>) {
  const session = useSession();
  const [_, copy] = useCopyToClipboard();

  const permalink = `${WEB_ADDRESS}/m/${profile.handle}`;

  const isSelf = session?.id === profile.id;

  const isSuspendEnabled =
    !isSelf && hasPermission(session, "MANAGE_SUSPENSIONS");

  const isRoleChangeEnabled = hasPermission(session, "MANAGE_ROLES");

  function handleSelect(value: MenuSelectionDetails) {
    switch (value.value) {
      case "copy-link":
        copy(permalink);
        toast("Link copied to clipboard");
        break;
    }
  }

  return (
    <Menu.Root onSelect={handleSelect}>
      <Menu.Trigger
        className="member-options-menu__trigger"
        maxW="full"
        cursor="pointer"
        asChild={props.asChild}
        textAlign="start"
      >
        {children}
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="group">
              <Menu.ItemGroupLabel display="flex" gap="2" alignItems="center">
                <MemberIdent profile={profile} size="md" name="full-vertical" />
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.ItemGroup>
                <Link href={permalink}>
                  <Menu.Item value="view">View profile</Menu.Item>
                </Link>
              </Menu.ItemGroup>

              <Menu.Item value="copy-link">Copy link</Menu.Item>

              <ReportMemberMenuItem profile={profile} />

              {isRoleChangeEnabled && <MemberRoleMenu profile={profile} />}

              {isSuspendEnabled && (
                <MemberSuspensionTrigger profile={profile}>
                  <Menu.Item
                    value="suspend"
                    color="fg.destructive"
                    _hover={{
                      color: "fg.destructive",
                      background: "bg.destructive",
                    }}
                  >
                    {profile.suspended ? "Reinstate" : "Suspend"}
                  </Menu.Item>
                </MemberSuspensionTrigger>
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
