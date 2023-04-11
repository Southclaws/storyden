import { useToast } from "@chakra-ui/react";
import { threadCreate } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import { errorToast } from "src/components/ErrorBanner";

export function useComposeScreen() {
  const account = useSession();
  const toast = useToast();

  const loggedIn = !!account;

  async function onCreate(title: string, category: string, md: string) {
    if (!loggedIn) return;

    await threadCreate({
      title: title,
      body: md,
      category: category,
      tags: [],
    }).catch(errorToast(toast));
  }

  return {
    onCreate,
  };
}
