import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { useI18n } from "@/i18n/provider";

import { CreateAccessKeyScreen } from "./CreateAccessKeyScreen";
import { Props, WithDisclosure } from "./useCreateAccessKeyScreen";

export function CreateAccessKeyModal(props: WithDisclosure<Props>) {
  const { t } = useI18n();

  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title={t("Create Access Key")}
    >
      <CreateAccessKeyScreen onClose={props.onClose || (() => {})} />
    </ModalDrawer>
  );
}
