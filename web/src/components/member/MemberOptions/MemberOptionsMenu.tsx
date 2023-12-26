import { Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { useSession } from "src/auth";
import { Avatar } from "src/components/site/Avatar/Avatar";
import {
  MenuContent,
  MenuItem,
  MenuItemGroup,
  MenuItemGroupLabel,
  MenuSeparator,
} from "src/theme/components/Menu";
import { Menu, MenuPositioner, MenuTrigger } from "src/theme/components/Menu";
import { WithDisclosure } from "src/utils/useDisclosure";

import { MemberSuspensionTrigger } from "../MemberSuspension/MemberSuspensionTrigger";

import { VStack, styled } from "@/styled-system/jsx";

import { Props } from "./useMemberOptionsScreen";

export function MemberOptionsMenu({
  children,
  ...props
}: PropsWithChildren<WithDisclosure<Props>>) {
  const session = useSession();

  const showAdminOptions = session?.admin && props.handle !== session.handle;

  return (
    <Menu
      size="sm"
      userSelect="none"
      isOpen={props.isOpen}
      onClose={props.onClose}
    >
      <MenuTrigger asChild>{children}</MenuTrigger>

      <Portal>
        <MenuPositioner>
          <MenuContent minW="48">
            <MenuItemGroup id="group">
              <MenuItemGroupLabel
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
              </MenuItemGroupLabel>

              <MenuSeparator />

              {showAdminOptions && (
                <MemberSuspensionTrigger {...props}>
                  <MenuItem
                    id="suspend"
                    color="fg.destructive"
                    _hover={{
                      color: "fg.destructive",
                      background: "bg.destructive",
                    }}
                  >
                    {props.deletedAt ? "Reinstate" : "Suspend"}
                  </MenuItem>
                </MemberSuspensionTrigger>
              )}
            </MenuItemGroup>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
