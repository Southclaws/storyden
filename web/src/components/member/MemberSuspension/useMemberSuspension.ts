import { Arguments, useSWRConfig } from "swr";

import { getProfileListKey } from "src/api/openapi-client/profiles";
import { ProfileReference } from "src/api/openapi-schema";
import { WithDisclosure } from "src/utils/useDisclosure";

import { handle } from "@/api/client";
import {
  adminAccountBanCreate,
  adminAccountBanRemove,
} from "@/api/openapi-client/admin";

export type Props = {
  profile: ProfileReference;
};

export function useMemberSuspension({
  profile,
  ...props
}: WithDisclosure<Props>) {
  const { mutate } = useSWRConfig();

  const profileKey = getProfileListKey()[0];
  const keyFn = (key: Arguments) => {
    return Array.isArray(key) && key[0].startsWith(profileKey);
  };

  async function handleSuspension() {
    await handle(async () => {
      await adminAccountBanCreate(profile.handle);

      mutate(keyFn);
      props.onClose?.();
    });
  }

  async function handleReinstate() {
    await handle(async () => {
      await adminAccountBanRemove(profile.handle);

      mutate(keyFn);
      props.onClose?.();
    });
  }

  return {
    handlers: {
      handleSuspension,
      handleReinstate,
    },
  };
}
