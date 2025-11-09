import { getServerSidebarState } from "./server";
import { SidebarToggle } from "./SidebarToggle";

export async function SidebarToggleServer() {
  const initialSidebarState = await getServerSidebarState();

  return <SidebarToggle initialValue={initialSidebarState} />;
}
