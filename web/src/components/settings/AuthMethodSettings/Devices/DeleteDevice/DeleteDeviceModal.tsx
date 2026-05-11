import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { useI18n } from "@/i18n/provider";

import { DeleteDeviceScreen } from "./DeleteDeviceScreen";
import { Props, WithDisclosure } from "./useDeleteDeviceScreen";

export function DeleteDeviceModal(props: WithDisclosure<Props>) {
  const { t } = useI18n();

  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title={t("Delete device")}
      >
        <DeleteDeviceScreen onClose={props.onClose} id={props.id} />
      </ModalDrawer>
    </>
  );
}
