import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { useI18n } from "@/i18n/provider";

import { CategoryDeleteProps, CategoryDeleteScreen } from "./CategoryDeleteScreen";

export function CategoryDeleteModal(props: CategoryDeleteProps) {
  const { t } = useI18n();

  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title={t("Delete category")}
    >
      <CategoryDeleteScreen {...props} />
    </ModalDrawer>
  );
}
