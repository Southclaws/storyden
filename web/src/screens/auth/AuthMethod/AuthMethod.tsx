import { Spinner } from "@chakra-ui/react";
import { SignUp } from "./Password/SignUp";

interface Props {
  method?: string;
}
export function AuthMethod({ method }: Props) {
  if (!method) {
    return <Spinner />;
  }

  switch (method) {
    case "password":
      return <SignUp />;
  }

  // TODO: Status/error component.
  return <Spinner />;
}
