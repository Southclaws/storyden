import { Text, VStack } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { Unready } from "src/components/site/Unready";

import { Header } from "./components/Header";
import { Metadata } from "./components/Metadata";
import { Props, useProfileScreen } from "./useProfileScreen";

export function ProfileLayout(props: PropsWithChildren<Props>) {
  const profile = useProfileScreen(props);
  if (!profile.ready) return <Unready {...profile.error} />;

  return (
    <VStack py={4} width="full" alignItems="start" gap={2}>
      <Header {...profile.data} />
      <Metadata {...profile.data} />
      <Text>{profile.data.bio}</Text>
      {props.children}
    </VStack>
  );
}
