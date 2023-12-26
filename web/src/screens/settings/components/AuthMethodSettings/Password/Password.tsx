import { Heading3 } from "src/theme/components/Heading/Index";

import { VStack } from "@/styled-system/jsx";

import { PasswordCreateForm } from "./PasswordCreateForm";
import { PasswordUpdateForm } from "./PasswordUpdateForm";

export type Props = {
  active: boolean;
};

export function Password(props: Props) {
  return (
    <VStack w="full" alignItems="start">
      <Heading3 size="sm">Password</Heading3>
      {props.active ? <PasswordUpdateForm /> : <PasswordCreateForm />}
    </VStack>
  );
}
