import {
  Box,
  BoxProps,
  Button,
  Checkbox,
  Circle,
  CloseButton,
  CreateToastFnReturn,
  Divider,
  Flex,
  FlexProps,
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  Heading,
  IconButton,
  Image,
  Input,
  Link,
  LinkBox,
  LinkOverlay,
  LinkProps,
  List,
  ListIcon,
  ListItem,
  Menu,
  MenuButton,
  MenuDivider,
  MenuGroup,
  MenuItem,
  MenuList,
  OrderedList,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
  Skeleton,
  SkeletonText,
  SlideFade,
  Spinner,
  StackProps,
  Text,
  UseDisclosureProps,
  VStack,
  useClipboard,
  useDisclosure,
  useOutsideClick,
  useToast,
} from "@chakra-ui/react";

// NOTE: These are being replaced gradually, so these are being re-exported.

// Layout
export { Box, Circle, Flex, HStack, VStack };

// Form

export {
  Button,
  Checkbox,
  FormControl,
  FormErrorMessage,
  FormLabel,
  IconButton,
  Input,
};

// Typography

export { Divider, Heading, Link, List, ListIcon, ListItem, OrderedList, Text };

// Menu

export { Menu, MenuButton, MenuDivider, MenuGroup, MenuItem, MenuList };

// Popover

export { Popover, PopoverArrow, PopoverBody, PopoverContent, PopoverTrigger };

// Other stuff

export {
  CloseButton,
  Image,
  LinkBox,
  LinkOverlay,
  Spinner,
  Skeleton,
  SkeletonText,
  SlideFade,
};

// Disclosure
// TODO: Copy into our codebase:
// https://github.com/chakra-ui/chakra-ui/blob/main/packages/hooks/use-disclosure/src/index.ts

export { useDisclosure };
export type { UseDisclosureProps };
export type WithDisclosure<T> = UseDisclosureProps & T;

// Hooks

export { useClipboard, useOutsideClick, useToast };

// Types

export type { BoxProps, CreateToastFnReturn, FlexProps, LinkProps, StackProps };
