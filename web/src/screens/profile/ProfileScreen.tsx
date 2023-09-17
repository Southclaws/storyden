import { Unready } from "src/components/site/Unready";

import { Content } from "./components/Content/Content";
import { ProfileContext } from "./context";
import { Props, useProfileScreen } from "./useProfileScreen";

export function ProfileScreen(props: Props) {
  const profile = useProfileScreen(props);

  if (!profile.ready) return <Unready {...profile.error} />;

  return (
    <ProfileContext.Provider value={profile.state}>
      <Content {...profile.data} />
    </ProfileContext.Provider>
  );
}
