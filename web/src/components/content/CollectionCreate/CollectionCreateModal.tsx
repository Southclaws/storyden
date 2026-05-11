import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { useI18n } from "@/i18n/provider";

import { CollectionCreateScreen } from "./CollectionCreateScreen";
import { Props } from "./useCollectionCreate";

export function CollectionCreateModal({ session, ...props }: Props) {
  const { t } = useI18n();

  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title={t("Create collection")}
      >
        <CollectionCreateScreen
          id={props.id}
          session={session}
          onClose={props.onClose}
        />
      </ModalDrawer>
    </>
  );
}
