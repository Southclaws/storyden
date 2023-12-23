import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { WithDisclosure } from "src/theme/components";

import { MemberMenuOptionsScreen } from "./MemberOptionsScreen";
import { Props } from "./useMemberOptionsScreen";

export function MemberOptionsModal(props: WithDisclosure<Props>) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title={props.name}
      >
        <MemberMenuOptionsScreen {...props} />
      </ModalDrawer>
    </>
  );
}
