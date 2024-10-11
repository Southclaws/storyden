import { useRouter } from "next/navigation";

import { handle } from "@/api/client";
import { useLibraryMutation } from "@/lib/library/library";

export function useLibraryPageTree(currentNode: string | undefined) {
  const router = useRouter();
  const { deleteNode, revalidate } = useLibraryMutation();

  async function handleDelete(slug: string) {
    handle(
      async () => {
        await deleteNode(slug);
        if (currentNode === slug) {
          router.push("/l");
        }
      },
      {
        promiseToast: {
          loading: "Deleting page...",
          success: "Page deleted.",
        },
        cleanup: async () => {
          revalidate();
        },
      },
    );
  }

  return {
    handleDelete,
  };
}
