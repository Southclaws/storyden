import {
  Box,
  Button,
  HStack,
  Heading,
  Image,
  List,
  ListItem,
  Text,
} from "@chakra-ui/react";

import { Unready } from "src/components/site/Unready";
import { AuthSelectionOption } from "src/screens/auth/components/AuthSelectionOption";

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
                {v.logo_url && (
                  <Box overflow="clip" height="1rem">
                    <Image src={v.logo_url} height="1rem" alt="" />
                  </Box>
                )}
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
            <AuthSelectionOption
              name={v.name}
              method={v.provider}
              icon={v.logo_url}
              link={v.link || undefined}
            />
          </ListItem>
        ))}
      </List>
    </SettingsSection>
  );
}
