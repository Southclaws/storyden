import {
  Button,
  HStack,
  Heading,
  List,
  ListItem,
  Text,
} from "@chakra-ui/react";

import { Unready } from "src/components/site/Unready";
import { Link } from "src/theme/components/Link";

import { SettingsSection } from "../SettingsSection/SettingsSection";

import { useAuthMethodSettings } from "./useAuthMethodSettings";

export function AuthMethodSettings() {
  const state = useAuthMethodSettings();

  if (!state.ready) return <Unready {...state.error} />;

  return (
    <SettingsSection>
      <Heading size="sm">Authentication methods</Heading>

      <Text>
        You can add as many authentication methods to your account as you want
        to.
      </Text>

      <Heading size="xs">Active</Heading>

      <List display="flex" flexDir="column" gap={2} w="full">
        {state.active.map((v) => (
          <ListItem key={v.provider}>
            <Button width="full" variant="outline" disabled>
              <HStack>
                <Text>{v.name}</Text>
              </HStack>
            </Button>
          </ListItem>
        ))}
      </List>

      <Heading size="xs">Available</Heading>

      <List display="flex" flexDir="column" gap={2} w="full">
        {state.rest.map((v) => (
          <ListItem key={v.provider}>
            <Link href={v.link}>{v.name}</Link>
          </ListItem>
        ))}
      </List>
    </SettingsSection>
  );
}
