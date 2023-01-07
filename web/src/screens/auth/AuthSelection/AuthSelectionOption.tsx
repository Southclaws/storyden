import { Box, Button, HStack, Link, Text } from "@chakra-ui/react";

interface Props {
  name: string;
  icon: JSX.Element;
  method: string;
  link?: string | undefined;
}

export function AuthSelectionOption({ name, icon, method, link }: Props) {
  return (
    <Button width="full">
      <Link href={link ?? `/auth/${method}`}>
        <HStack>
          <Box overflow="clip" height="1rem">
            {icon}
          </Box>
          <Text>{name}</Text>
        </HStack>
      </Link>
    </Button>
  );
}
