import { useProfileGet } from "src/api/openapi-client/profiles";
import { APIError, PublicProfile } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { ProfileContextShape } from "./context";

type ProfileScreen =
  | { ready: false; error: void | APIError }
  | {
      ready: true;
      data: PublicProfile;
      state: ProfileContextShape;
    };

export type Props = { handle: string };

export function useProfileScreen({ handle }: Props): ProfileScreen {
  const session = useSession();
  const { data, error } = useProfileGet(handle);

  if (!data) {
    return { error, ready: false };
  }

  const isSelf = session?.id === data.id;

  return {
    data,
    ready: true,
    state: {
      isEditing: false, // TODO: profile editing
      isSelf: isSelf,
    },
  };
}
