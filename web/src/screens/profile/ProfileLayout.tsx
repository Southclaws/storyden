import { PropsWithChildren } from "react";

import { Unready } from "src/components/site/Unready";

import { VStack } from "@/styled-system/jsx";

import { Header } from "./components/Header";
import { Metadata } from "./components/Metadata";
import { Props, useProfileScreen } from "./useProfileScreen";

export function ProfileLayout(props: PropsWithChildren<Props>) {
  const profile = useProfileScreen(props);
  if (!profile.ready) return <Unready {...profile.error} />;

  return (
    <VStack py="4" width="full" alignItems="start" gap="2">
      <Header {...profile.data} />
      <Metadata {...profile.data} />
      <p>{profile.data.bio}</p>
      {props.children}
    </VStack>
  );
}
