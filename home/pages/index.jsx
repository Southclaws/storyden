import {
  Box,
  Flex,
  Heading,
  HStack,
  Link,
  Text,
  VStack,
} from "@chakra-ui/react";
import Image from "next/image";

export default function Home() {
  return (
    <Flex direction="column">
      <VStack
        bgColor="white"
        color="hsla(0, 0%, 19%, 1)"
        p={8}
        gap={2}
        w="full"
      >
        <Image src="/logo_200x200.png" width={50} height={50} />
        <Heading as="h1" fontWeight="bold" fontSize="lg">
          A forum for the modern age
        </Heading>
        <Heading
          as="h2"
          fontFamily="Halyard Display"
          fontWeight="medium"
          fontSize="md"
          textAlign="center"
          px={8}
        >
          Empower your community with a modern discussion platform
        </Heading>
        <Box bgColor="white">
          <Image src="/graphic.png" width={323} height={297} />
        </Box>
      </VStack>

      <VStack
        bgColor="hsla(0, 0%, 19%, 1)"
        color="hsla(160, 9%, 92%, 1)"
        p={8}
        gap={2}
        w="full"
      >
        <Heading as="h2" fontWeight="bold" fontSize="md" textAlign="center">
          What makes Storyden different?
        </Heading>

        <Flex
          maxW="xl"
          direction={{
            base: "column",
            md: "row",
          }}
        >
          <VStack p={4}>
            <Heading as="h1" fontWeight="bold" fontSize="lg">
              Modern
            </Heading>
            <Text textAlign="center" maxW="sm">
              Simple and accessible with timeless design and modern paradigms.
            </Text>
          </VStack>
          <VStack p={4}>
            <Heading as="h1" fontWeight="bold" fontSize="lg">
              Fast
            </Heading>
            <Text textAlign="center" maxW="sm">
              Built with Go, optimised for performance and user experience.
            </Text>
          </VStack>
          <VStack p={4}>
            <Heading as="h1" fontWeight="bold" fontSize="lg">
              Extensible
            </Heading>
            <Text textAlign="center" maxW="sm">
              The ultimate level of customisation: bring your own frontend or
              app.
            </Text>
          </VStack>
        </Flex>
      </VStack>

      <VStack
        bgColor="hsla(160, 9%, 92%, 1)"
        color="hsla(0, 0%, 19%, 1)"
        p={8}
        gap={2}
        w="full"
        textAlign="center"
      >
        <Heading as="h1" fontWeight="bold" fontSize="lg">
          Interested?
        </Heading>
        <Text>
          Storyden is early in development and is looking for feedback and
          contributors!
        </Text>
        <Text>
          If you have opinions about forum software,{" "}
          <Link
            isExternal
            href="https://airtable.com/shrLY0jDp9CuXPB2X"
            color="hsla(265, 56%, 42%, 1)"
          >
            please click this link!
          </Link>
        </Text>
        <Text>
          If you know Golang or React.js and are interested in contributing to a
          high quality open source project,{" "}
          <Link
            isExternal
            href="https://github.com/Southclaws/storyden"
            color="hsla(265, 56%, 42%, 1)"
          >
            please click this link!
          </Link>
        </Text>
      </VStack>

      <VStack
        bgColor="hsla(0, 0%, 77%, 1)"
        color="hsla(0, 0%, 19%, 1)"
        width="full"
        p={4}
      >
        <HStack>
          <Link href="https://github.com/Southclaws/storyden" isExternal>
            <svg
              width="33"
              height="32"
              viewBox="0 0 33 32"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M16.5 0.340332C7.66 0.340332 0.5 7.50033 0.5 16.3403C0.5 23.4203 5.08 29.4003 11.44 31.5203C12.24 31.6603 12.54 31.1803 12.54 30.7603C12.54 30.3803 12.52 29.1203 12.52 27.7803C8.5 28.5203 7.46 26.8003 7.14 25.9003C6.96 25.4403 6.18 24.0203 5.5 23.6403C4.94 23.3403 4.14 22.6003 5.48 22.5803C6.74 22.5603 7.64 23.7403 7.94 24.2203C9.38 26.6403 11.68 25.9603 12.6 25.5403C12.74 24.5003 13.16 23.8003 13.62 23.4003C10.06 23.0003 6.34 21.6203 6.34 15.5003C6.34 13.7603 6.96 12.3203 7.98 11.2003C7.82 10.8003 7.26 9.16033 8.14 6.96033C8.14 6.96033 9.48 6.54033 12.54 8.60033C13.82 8.24033 15.18 8.06033 16.54 8.06033C17.9 8.06033 19.26 8.24033 20.54 8.60033C23.6 6.52033 24.94 6.96033 24.94 6.96033C25.82 9.16033 25.26 10.8003 25.1 11.2003C26.12 12.3203 26.74 13.7403 26.74 15.5003C26.74 21.6403 23 23.0003 19.44 23.4003C20.02 23.9003 20.52 24.8603 20.52 26.3603C20.52 28.5003 20.5 30.2203 20.5 30.7603C20.5 31.1803 20.8 31.6803 21.6 31.5203C24.7765 30.4483 27.5367 28.407 29.4921 25.6838C31.4474 22.9607 32.4994 19.6928 32.5 16.3403C32.5 7.50033 25.34 0.340332 16.5 0.340332Z"
                fill="#303030"
              />
            </svg>
          </Link>
          <Link href="https://twitter.com/Southclaws" isExternal>
            <svg
              width="33"
              height="27"
              viewBox="0 0 33 27"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M10.552 26.9462C22.628 26.9462 29.234 16.9402 29.234 8.27822C29.234 7.99822 29.234 7.71422 29.222 7.43421C30.5081 6.50319 31.6181 5.35052 32.5 4.03021C31.2986 4.55991 30.0255 4.90902 28.722 5.06621C30.0951 4.24549 31.1234 2.95415 31.616 1.43221C30.326 2.19646 28.914 2.73296 27.442 3.01821C26.4524 1.96426 25.1428 1.26605 23.7162 1.03173C22.2895 0.797405 20.8254 1.04005 19.5507 1.72209C18.2759 2.40412 17.2616 3.48747 16.6649 4.80433C16.0682 6.1212 15.9224 7.59808 16.25 9.00622C13.6395 8.87532 11.0856 8.19715 8.75409 7.0157C6.42254 5.83425 4.36537 4.1759 2.716 2.14821C1.87868 3.59438 1.62315 5.30501 2.00129 6.93275C2.37944 8.56048 3.36291 9.98328 4.752 10.9122C3.71108 10.8769 2.69304 10.5973 1.78 10.0962V10.1862C1.78179 11.7012 2.30661 13.1691 3.26575 14.3417C4.22488 15.5144 5.5595 16.32 7.044 16.6222C6.48053 16.7775 5.89846 16.8548 5.314 16.8522C4.90197 16.8535 4.49076 16.8153 4.086 16.7382C4.50557 18.0423 5.32255 19.1826 6.42249 19.9992C7.52244 20.8158 8.85027 21.2679 10.22 21.2922C7.89311 23.1198 5.0188 24.111 2.06 24.1062C1.53864 24.1084 1.01765 24.0784 0.5 24.0162C3.50303 25.9308 6.9906 26.9473 10.552 26.9462V26.9462Z"
                fill="#303030"
              />
            </svg>
          </Link>
        </HStack>
        <Text fontSize="sm">
          Storyden forum software &copy; Barnaby Keene 2022
        </Text>
      </VStack>
    </Flex>
  );
}
