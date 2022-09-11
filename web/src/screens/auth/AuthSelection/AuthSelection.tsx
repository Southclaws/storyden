import {
  Box,
  Divider,
  Heading,
  List,
  ListItem,
  Text,
  VStack,
} from "@chakra-ui/react";
import { AuthSelectionOption } from "./AuthSelectionOption";
import { KeyIcon } from "@heroicons/react/20/solid";

export function AuthSelection() {
  return (
    <VStack gap={4}>
      <Box>
        <Heading size="md">
          Sign up
          <br />
        </Heading>
        <Text size="sm" fontWeight="medium" color="blackAlpha.600">
          or sign in
        </Text>
      </Box>

      <List>
        <ListItem>
          <AuthSelectionOption
            name="Password"
            method="password"
            icon={<KeyIcon height="100%" />}
          />
        </ListItem>
      </List>
    </VStack>
  );
}
