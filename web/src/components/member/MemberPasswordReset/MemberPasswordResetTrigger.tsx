import React, { PropsWithChildren } from "react";

import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { useDisclosure } from "src/utils/useDisclosure";

import { ProfileReference } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";

import { MemberPasswordResetDialog } from "./MemberPasswordResetDialog";

export type Props = {
  profile: ProfileReference;
};

export function MemberPasswordResetTrigger({
  children,
  profile,
}: PropsWithChildren<Props>) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  const title = `Reset password for ${profile.name}`;
  const trigger =
    React.isValidElement<TriggerChildProps>(children) &&
    children.type !== React.Fragment
      ? children
      : null;

  return (
    <>
      {trigger ? (
        React.cloneElement(trigger, {
          onClick: (event) => {
            trigger.props.onClick?.(event);
            onOpen();
          },
        })
      ) : (
        <Button onClick={onOpen}>Reset Password</Button>
      )}

      <ModalDrawer isOpen={isOpen} onClose={onClose} title={title}>
        <MemberPasswordResetDialog onClose={onClose} profile={profile} />
      </ModalDrawer>
    </>
  );
}

type TriggerChildProps = {
  onClick?: React.MouseEventHandler;
};
