import { Text, VStack } from "@chakra-ui/react";
import { PublicProfile } from "src/api/openapi/schemas";
import { Header } from "./Header";
import { Content } from "./Content/Content";
import { Metadata } from "./Metadata";

export function Profile(props: PublicProfile) {
  return (
    <VStack alignItems="start" gap={2}>
      <Header {...props} />
      <Metadata {...props} />
      <Text>{props.bio}</Text>
      <Content {...props} />
    </VStack>
  );
}
