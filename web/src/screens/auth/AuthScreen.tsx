import { Box, Button, Input } from "@chakra-ui/react";

import { useAuthScreen } from "./useAuthScreen";

type Props = {};
export function AuthScreen({}: Props) {
  const { register, onSubmit } = useAuthScreen();

  return (
    <Box>
      <form onSubmit={onSubmit}>
        <Input {...register("handle")} />
        <Button type="submit">Authenticate</Button>
      </form>
    </Box>
  );
}
