import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import { accountUpdate } from "@/api/openapi-client/accounts";
import { getProfileGetKey } from "@/api/openapi-client/profiles";
import {
  AccountMutableProps,
  ProfileGetOKResponse,
} from "@/api/openapi-schema";

export function useProfileMutations(handle: string) {
  const { mutate } = useSWRConfig();

  const profileKey = getProfileGetKey(handle);

  function keyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0].startsWith(profileKey);
  }

  const revalidate = async (data?: MutatorCallback<ProfileGetOKResponse>) => {
    await mutate(keyFilterFn, data);
  };

  const update = async (updated: AccountMutableProps) => {
    const mutator: MutatorCallback<ProfileGetOKResponse> = (data) => {
      if (!data) return;

      const newData = {
        ...data,
        ...updated,
      } as ProfileGetOKResponse;

      return newData;
    };

    await mutate(profileKey, mutator, {
      revalidate: false,
    });

    await accountUpdate(updated);
  };

  return {
    revalidate,
    update,
  };
}
