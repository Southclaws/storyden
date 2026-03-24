import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { handle } from "@/api/client";
import { useCategoryMutations } from "@/lib/category/mutation";
import { HStack, VStack, styled } from "@/styled-system/jsx";

import { ModalDrawer } from "../../site/Modaldrawer/Modaldrawer";
import { Button } from "../../ui/button";
import { Item } from "../../ui/menu";

type Props = {
  category: Category;
};

export function CategoryVisibilityToggleMenuItem({ category }: Props) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { updateCategory, revalidateList } = useCategoryMutations();

  const isHidden = category.visibility === "unlisted";

  async function handleConfirm() {
    await handle(
      async () => {
        await updateCategory(category.slug, {
          visibility: isHidden ? "published" : "unlisted",
        });
        onClose();
      },
      {
        promiseToast: {
          loading: isHidden ? "Unhiding category..." : "Hiding category...",
          success: isHidden ? "Category unhidden." : "Category hidden.",
        },
        cleanup: async () => {
          await revalidateList();
        },
      },
    );
  }

  return (
    <>
      <Item value="toggle-visibility" onClick={onOpen}>
        {isHidden ? "Unhide" : "Hide"}
      </Item>

      <ModalDrawer
        isOpen={isOpen}
        onClose={onClose}
        title={isHidden ? "Unhide category" : "Hide category"}
      >
        <VStack alignItems="start" gap="4" maxW="prose">
          <styled.p>
            {isHidden
              ? `Make ${category.name} and all of its subcategories visible again?`
              : `Hide ${category.name} and all of its subcategories from public category lists and search?`}
          </styled.p>

          <HStack w="full" alignItems="center" justify="end" gap="4">
            <Button size="sm" variant="ghost" onClick={onClose}>
              Cancel
            </Button>
            <Button size="sm" onClick={handleConfirm}>
              {isHidden ? "Unhide" : "Hide"}
            </Button>
          </HStack>
        </VStack>
      </ModalDrawer>
    </>
  );
}
