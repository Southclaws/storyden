import dynamic from "next/dynamic";

import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { UnreadyBanner } from "@/components/site/Unready";

import { type Props } from "./useCategoryEdit";

const CategoryEditForm = dynamic(
  () => import("./CategoryEditForm").then((module) => module.CategoryEditForm),
  { ssr: false, loading: () => <UnreadyBanner /> },
);

export function CategoryEditModal(props: Props) {
  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      onOpenChange={props.onOpenChange}
      title="Edit category"
    >
      <CategoryEditForm {...props} />
    </ModalDrawer>
  );
}
