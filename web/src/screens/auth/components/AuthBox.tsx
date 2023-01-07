import { Box, BoxProps } from "@chakra-ui/react";

export function AuthBox(props: BoxProps) {
  return (
    <Box id="AuthBox" p={6} borderRadius={8} width="full" maxW="xs" {...props}>
      {props.children}
    </Box>
  );
}
