import React, { PropsWithChildren } from "react";

import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { Button } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

import { MemberSuspensionConfirmation } from "./MemberSuspensionConfirmation";
import { Props } from "./useMemberSuspension";

export function MemberSuspensionTrigger({
  children,
  ...props
}: PropsWithChildren<Props>) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  const title = props.deletedAt
    ? `Reinstate account ${props.name}`
    : `Suspend account ${props.name}`;

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
          {props.deletedAt ? "Reinstate" : "Suspend"}
        </Button>
      )}

      <ModalDrawer isOpen={isOpen} onClose={onClose} title={title}>
        <MemberSuspensionConfirmation onClose={onClose} {...props} />
      </ModalDrawer>
    </>
  );
}
