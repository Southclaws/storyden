import dynamic from "next/dynamic";
import React, { ReactNode, useState } from "react";

import { ProfileReference } from "@/api/openapi-schema";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { UnreadyBanner } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { useDisclosure } from "@/utils/useDisclosure";

const MemberWarningScreen = dynamic(
  () =>
    import("./MemberWarningScreen").then(
      (module) => module.MemberWarningScreen,
    ),
  { ssr: false, loading: () => <UnreadyBanner /> },
);

type MemberWarningTriggerProps = {
  children?: ReactNode;
  profile: ProfileReference;
};

export function MemberWarningTrigger({
  children,
  profile,
}: MemberWarningTriggerProps) {
  const { onOpen, onClose, isOpen } = useDisclosure();
  const [reason, setReason] = useState("");

  const triggerNode = React.isValidElement<{
    onClick?: React.MouseEventHandler;
  }>(children) ? (
    React.cloneElement(children, { onClick: onOpen })
  ) : (
    <Button intent="warning" onClick={onOpen}>
      Warn
    </Button>
  );

  return (
    <>
      {triggerNode}
      <ModalDrawer
        isOpen={isOpen}
        onClose={onClose}
        title={`Issue warning to ${profile.name}`}
      >
        <MemberWarningScreen
          profile={profile}
          reason={reason}
          onReasonChange={setReason}
          onClose={onClose}
        />
      </ModalDrawer>
    </>
  );
}
