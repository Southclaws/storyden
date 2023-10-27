import { Link } from "@chakra-ui/next-js";
import { Button, HStack, Text } from "@chakra-ui/react";

interface Props {
  name: string;
  method: string;
  link?: string | undefined;
}

export function AuthSelectionOption({ name, method, link }: Props) {
  return (
    <Link href={link ?? `/auth/${method}`}>
      <Button width="full">
        <HStack>
          <Text>{name}</Text>
        </HStack>
      </Button>
    </Link>
  );
}
