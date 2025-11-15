import { useSession } from "@/auth";
import { MemberInterfaceSettings } from "@/components/settings/MemberInterfaceSettings/MemberInterfaceSettings";

export function MemberInterfaceSettingsScreen() {
  const session = useSession();
  if (!session) {
    return null;
  }

  return <MemberInterfaceSettings session={session} />;
}
