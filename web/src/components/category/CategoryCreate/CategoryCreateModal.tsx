import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { useI18n } from "@/i18n/provider";

import {
  CategoryCreateProps,
  CategoryCreateScreen,
} from "./CategoryCreateScreen";

export function CategoryCreateModal(props: CategoryCreateProps) {
  const { t } = useI18n();

  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        onOpenChange={props.onOpenChange}
        title={t("Create category")}
      >
        <CategoryCreateScreen {...props} />
      </ModalDrawer>
    </>
  );
}
