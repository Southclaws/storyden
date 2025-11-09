import { hasCapability } from "@/lib/settings/capabilities";
import { getSettings } from "@/lib/settings/settings-server";

import { AskAnchor } from "./Ask";

export async function AskServer() {
  const { capabilities } = await getSettings();
  const isSemdexEnabled = hasCapability("semdex", capabilities);

  if (!isSemdexEnabled) {
    return null;
  }

  return <AskAnchor />;
}
