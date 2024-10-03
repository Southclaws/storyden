import { useRouter } from "next/navigation";
import { mutate } from "swr";

import { getNodeListKey, nodeDelete } from "@/api/openapi-client/nodes";

export function useLibraryPageTree(currentNode: string | undefined) {
  const router = useRouter();
  async function handleDelete(slug: string) {
    await nodeDelete(slug);
    await mutate(getNodeListKey());

    if (currentNode === slug) {
      router.push("/l");
    }
  }

  return {
    handleDelete,
  };
}
