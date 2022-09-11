import { Box, Button, HStack, Link, Text } from "@chakra-ui/react";

interface Props {
  name: string;
  icon: JSX.Element;
  method: string;
}

export function AuthSelectionOption({ name, icon, method }: Props) {
  return (
    <Button width="full">
      <Link href={`/login/${method}`}>
        <HStack>
          <Box overflow="clip" height="1rem">
            {icon}
          </Box>{" "}
          <Text>{name}</Text>
        </HStack>
      </Link>
    </Button>
  );
}
