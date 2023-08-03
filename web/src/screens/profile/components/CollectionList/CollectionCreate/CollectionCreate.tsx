import {
  Box,
  Button,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from "@chakra-ui/react";
import { FolderPlusIcon } from "@heroicons/react/24/outline";

import { useCollectionCreate } from "./useCollectionCreate";

export function CollectionCreate() {
  const { onOpen, onClose, isOpen, register, onSubmit } = useCollectionCreate();
  return (
    <>
      <Button leftIcon={<FolderPlusIcon width="1rem" />} onClick={onOpen}>
        New collection
      </Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Create collection</ModalHeader>
          <ModalCloseButton />
          <Box as="form" onSubmit={onSubmit}>
            <ModalBody>
              <FormControl>
                <FormLabel>Name</FormLabel>
                <Input {...register("name")} type="text" />
                <FormHelperText>The name for your collection</FormHelperText>
              </FormControl>
              <FormControl>
                <FormLabel>Description</FormLabel>
                <Input {...register("description")} type="text" />
                <FormHelperText>Describe your collection</FormHelperText>
              </FormControl>
            </ModalBody>

            <ModalFooter>
              <Button type="submit">Create</Button>
            </ModalFooter>
          </Box>
        </ModalContent>
      </Modal>
    </>
  );
}
