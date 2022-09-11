import { Box, Spinner } from "@chakra-ui/react";
import { Password } from "./Password/Password";

interface Props {
  method?: string;
}
export function AuthMethod({ method }: Props) {
  if (!method) {
    return <Spinner />;
  }

  switch (method) {
    case "password":
      return <Password />;
  }

  // TODO: Status/error component.
  return <Spinner />;
}
