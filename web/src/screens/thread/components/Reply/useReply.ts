import { useSession } from "src/auth";

export function useReply() {
  const account = useSession();

  const loggedIn = !!account;

  return {
    loggedIn,
  };
}
