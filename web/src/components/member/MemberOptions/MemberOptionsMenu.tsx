import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { PropsWithChildren } from "react";
import { toast } from "sonner";

import { useSession } from "src/auth";
import { Avatar } from "src/components/site/Avatar/Avatar";
import { WithDisclosure } from "src/utils/useDisclosure";

import * as Menu from "@/components/ui/menu";
import { WEB_ADDRESS } from "@/config";
import { VStack, styled } from "@/styled-system/jsx";

import { MemberSuspensionTrigger } from "../MemberSuspension/MemberSuspensionTrigger";

import { Props } from "./useMemberOptionsScreen";

export function MemberOptionsMenu({
  children,
  ...props
}: PropsWithChildren<WithDisclosure<Props>>) {
  const session = useSession();
  const [_, copy] = useCopyToClipboard();

  const permalink = `${WEB_ADDRESS}/m/${props.handle}`;

  const showAdminOptions = session?.admin && props.handle !== session.handle;

  function handleSelect(value: MenuSelectionDetails) {
    switch (value.value) {
      case "copy-link":
        copy(permalink);
        toast("Link copied to clipboard");
        break;
    }
  }

  return (
    <Menu.Root onOpenChange={props.onOpenChange} onSelect={handleSelect}>
      <Menu.Trigger asChild>{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="group">
              <Menu.ItemGroupLabel display="flex" gap="2" alignItems="center">
                <Avatar handle={props.handle} />
                <VStack alignItems="start" gap="0">
                  <styled.h1 color="fg.default">{props.name}</styled.h1>
                  <styled.h2 color="fg.subtle">@{props.handle}</styled.h2>
                </VStack>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link">Copy link</Menu.Item>

              {showAdminOptions && (
                <MemberSuspensionTrigger {...props}>
                  <Menu.Item
                    value="suspend"
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
