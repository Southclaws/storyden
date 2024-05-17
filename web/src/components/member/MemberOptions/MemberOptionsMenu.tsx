import { Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { useSession } from "src/auth";
import { Avatar } from "src/components/site/Avatar/Avatar";
import { WithDisclosure } from "src/utils/useDisclosure";

import { MemberSuspensionTrigger } from "../MemberSuspension/MemberSuspensionTrigger";

import * as Menu from "@/components/ui/menu";
import { VStack, styled } from "@/styled-system/jsx";

import { Props } from "./useMemberOptionsScreen";

export function MemberOptionsMenu({
  children,
  ...props
}: PropsWithChildren<WithDisclosure<Props>>) {
  const session = useSession();

  const showAdminOptions = session?.admin && props.handle !== session.handle;

  return (
    <Menu.Root size="sm" onOpenChange={props.onOpenChange}>
      <Menu.Trigger asChild>{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="group">
              <Menu.ItemGroupLabel
                htmlFor="group"
                display="flex"
                gap="2"
                alignItems="center"
              >
                <Avatar handle={props.handle} />
                <VStack alignItems="start" gap="0">
                  <styled.h1 color="fg.default">{props.name}</styled.h1>
                  <styled.h2 color="fg.subtle">@{props.handle}</styled.h2>
                </VStack>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              {showAdminOptions && (
                <MemberSuspensionTrigger {...props}>
                  <Menu.Item
                    id="suspend"
                    color="fg.destructive"
                    _hover={{
                      color: "fg.destructive",
                      background: "bg.destructive",
                    }}
                  >
                    {props.deletedAt ? "Reinstate" : "Suspend"}
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
