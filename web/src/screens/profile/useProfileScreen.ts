import { useProfileGet } from "src/api/openapi/profiles";
import { APIError, PublicProfile } from "src/api/openapi/schemas";

type ProfileScreen =
  | { ready: false; error: void | APIError }
  | {
      ready: true;
      data: PublicProfile;
    };

export type Props = { handle: string };

export function useProfileScreen({ handle }: Props): ProfileScreen {
  const { data, error } = useProfileGet(handle);

  if (!data) {
    return { error, ready: false };
  }

  return { data, ready: true };
}
