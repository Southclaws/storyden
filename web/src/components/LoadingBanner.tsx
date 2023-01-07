import { Flex, Spinner } from "@chakra-ui/react";
import { FC } from "react";

const LoadingBanner: FC = () => {
  return (
    <Flex minHeight="100vh" alignItems="center" justifyContent="center">
      <Spinner
        thickness="4px"
        speed="0.65s"
        emptyColor="gray.200"
        color="purple.500"
        size="xl"
      />
    </Flex>
  );
};

export default LoadingBanner;
