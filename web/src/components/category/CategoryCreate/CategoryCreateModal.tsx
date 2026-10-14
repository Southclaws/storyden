import dynamic from "next/dynamic";

import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { UnreadyBanner } from "@/components/site/Unready";

import type { CategoryCreateProps } from "./useCategoryCreate";

const CategoryCreateScreen = dynamic(
  () =>
    import("./CategoryCreateScreen").then(
      (module) => module.CategoryCreateScreen,
    ),
  { ssr: false, loading: () => <UnreadyBanner /> },
);

export function CategoryCreateModal(props: CategoryCreateProps) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        onOpenChange={props.onOpenChange}
        title="Create category"
      >
        <CategoryCreateScreen {...props} />
      </ModalDrawer>
    </>
  );
}
