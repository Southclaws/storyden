import { Heading } from "@chakra-ui/react";

import { Button } from "src/theme/components/Button";
import { Input } from "src/theme/components/Input";

import { VStack, styled } from "@/styled-system/jsx";

export function Password() {
  return (
    <VStack alignItems="start">
      <Heading size="sm">Password</Heading>

      <p>You can change your password here.</p>

      <styled.form display="flex" flexDir="column" gap="2">
        <Input placeholder="current password" />
        <Input placeholder="new password" />
        <Button>Change password</Button>
      </styled.form>
    </VStack>
  );
}
