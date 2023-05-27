import { Link } from "@chakra-ui/next-js";
import { Box, Button, HStack, Image, Text } from "@chakra-ui/react";

interface Props {
  name: string;
  icon: string;
  method: string;
  link?: string | undefined;
}

export function AuthSelectionOption({ name, icon, method, link }: Props) {
  return (
    <Link href={link ?? `/auth/${method}`}>
      <Button width="full">
        <HStack>
          {icon && (
            <Box overflow="clip" height="1rem">
              <Image src={icon} height="1rem" alt="" />
            </Box>
          )}
          <Text>{name}</Text>
        </HStack>
      </Button>
    </Link>
  );
}
