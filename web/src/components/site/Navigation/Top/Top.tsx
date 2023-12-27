import { navbarStyles } from "../common";
import { Title } from "../components/Title";
import { Toolbar } from "../components/Toolbar";
import { useNavigation } from "../useNavigation";

import { HStack } from "@/styled-system/jsx";

export function Top() {
  const { title } = useNavigation();

  return (
    <HStack
      className={navbarStyles}
      justify="space-between"
      alignItems="center"
      height="16"
      px="4"
    >
      <Title>{title}</Title>
      <Toolbar />
    </HStack>
  );
}
