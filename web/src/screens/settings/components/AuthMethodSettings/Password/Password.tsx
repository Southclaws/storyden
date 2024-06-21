import { Heading } from "@/components/ui/heading";
import { VStack } from "@/styled-system/jsx";

import { PasswordCreateForm } from "./PasswordCreateForm";
import { PasswordUpdateForm } from "./PasswordUpdateForm";

export type Props = {
  active: boolean;
};

export function Password(props: Props) {
  return (
    <VStack w="full" alignItems="start">
      <Heading size="sm">Password</Heading>
      {props.active ? <PasswordUpdateForm /> : <PasswordCreateForm />}
    </VStack>
  );
}
