import {
  Box,
  CloseButton,
  HStack,
  Heading,
  UseDisclosureProps,
  VStack,
} from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Drawer } from "vaul";

type Props = {
  title?: string;
} & UseDisclosureProps;

export function ModalDrawer({ children, ...props }: PropsWithChildren<Props>) {
  const onOpenChange = (open: boolean) => {
    if (open) props.onOpen?.();
    else props.onClose?.();
  };

  return (
    <>
      <Drawer.Root
        open={props.isOpen}
        onOpenChange={onOpenChange}
        shouldScaleBackground
      >
        <Drawer.Portal>
          <Drawer.Overlay className="modaldrawer__overlay" />
          <Drawer.Content className="modaldrawer__content">
            <VStack
              height={{ base: "full", md: "unset" }}
              borderTopRadius="1em"
              borderBottomRadius={{ base: "0", md: "1em" }}
              bgColor="gray.100"
              p={4}
            >
              <HStack w="full" justify="space-between">
                <Heading size="md">{props.title}</Heading>
                <CloseButton onClick={props.onClose} />
              </HStack>

              <Box h="full" w="full" pb={3}>
                {children}
              </Box>
            </VStack>
          </Drawer.Content>
        </Drawer.Portal>
      </Drawer.Root>

      <style jsx global>{`
        .modaldrawer__overlay {
          position: fixed;
          inset: 0;
          background-color: var(--chakra-colors-blackAlpha-600);
        }

        /* Modal mode - on desktop screens */
        @media screen and (min-width: 48em) {
          .modaldrawer__content {
            height: 100%;
            width: 100%;
            top: 0;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            position: fixed;
          }
        }

        /* Drawer mode - on mobile screens */
        @media screen and (max-width: 48em) {
          .modaldrawer__content {
            height: 100%;
            width: 100%;
            top: 0;
            display: flex;
            flex-direction: column;
            position: fixed;
            margin-top: 3rem;
            max-height: 96%;
          }
        }
      `}</style>
    </>
  );
}
