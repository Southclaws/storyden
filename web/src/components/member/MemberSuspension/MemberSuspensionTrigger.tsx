import React, { PropsWithChildren } from "react";

import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { useDisclosure } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";

import { MemberSuspensionConfirmation } from "./MemberSuspensionConfirmation";
import { Props } from "./useMemberSuspension";

export function MemberSuspensionTrigger({
  children,
  profile,
}: PropsWithChildren<Props>) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  const title = profile.suspended
    ? `Reinstate account ${profile.name}`
    : `Suspend account ${profile.name}`;

  return (
    <>
      {children ? (
        React.cloneElement(
          // not sure why types broken here, but it works fine.
          children as any,
          {
            onClick: onOpen,
          },
        )
      ) : (
        <Button colorPalette="red" onClick={onOpen}>
          {profile.suspended ? "Reinstate" : "Suspend"}
        </Button>
      )}

      <ModalDrawer isOpen={isOpen} onClose={onClose} title={title}>
        <MemberSuspensionConfirmation onClose={onClose} profile={profile} />
      </ModalDrawer>
    </>
  );
}
