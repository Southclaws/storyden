"use client";

import { groupBy, keys, partition } from "lodash";
import {
  ArrowUp,
  Filter,
  Heart,
  MessageCircle,
  MoreVertical,
  Plus,
  Settings,
  Share,
  User,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Text } from "@/components/ui/text";
import { Box, HStack, VStack, styled } from "@/styled-system/jsx";
import { getColourVariants } from "@/utils/colour";

const DEMO_ACCENT = "#3b82f6";

export default function Page() {
  const colours = getColourVariants(DEMO_ACCENT);
  const variables = keys(colours);
  const { flat, dark, other } = groupBy(variables, (v) => {
    if (v.includes("flat")) return "flat";
    if (v.includes("dark")) return "dark";
    return "other";
  });
  const [flatText, flatFill] = partition(flat, (v) => v.includes("text"));
  const [darkText, darkFill] = partition(dark, (v) => v.includes("text"));

  return (
    <Box bg="bg.default">
      <HStack gap="0" pt="20">
        {/* Left Sidebar - Fixed */}
        <Box
          pos="fixed"
          left="0"
          top="20"
          bottom="0"
          w="64"
          bg="bg.opaque"
          backdropBlur="frosted"
          backdropFilter="auto"
          borderRightWidth="thin"
          borderRightColor="border.subtle"
          p="4"
          display={{ base: "none", lg: "block" }}
          overflowY="auto"
        >
          <VStack gap="6" alignItems="start">
            <VStack gap="2" alignItems="start" w="full">
              <Text size="sm" fontWeight="semibold" color="fg.subtle">
                Navigation
              </Text>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                <User size={16} />
                Profile
              </Button>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                <Settings size={16} />
                Settings
              </Button>
              <Button variant="solid" size="sm" justifyContent="start" w="full">
                <Plus size={16} />
                Create Post
              </Button>
            </VStack>

            <VStack gap="2" alignItems="start" w="full">
              <Text size="sm" fontWeight="semibold" color="fg.subtle">
                Categories
              </Text>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                General Discussion
              </Button>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                Announcements
              </Button>
              <Button
                variant="outline"
                size="sm"
                justifyContent="start"
                w="full"
              >
                Help & Support
              </Button>
            </VStack>
          </VStack>
        </Box>

        {/* Main Content Area */}
        <Box flex="1" pl={{ base: "0", lg: "64" }} pr={{ base: "0", xl: "80" }}>
          <VStack gap="6" p="6" maxW="4xl" mx="auto">
            {/* Accent Color Palette Tester */}
            <Box w="full" bg="bg.subtle" borderRadius="lg" p="6" boxShadow="sm">
              <VStack gap="4" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  Accent Color Palette
                </Text>

                {other?.map((v) => (
                  <styled.div
                    key={v}
                    p="2"
                    borderRadius="md"
                    style={{ backgroundColor: `var(${v})` }}
                  >
                    <Text size="sm" fontWeight="medium">
                      {v}: {colours[v]}
                    </Text>
                  </styled.div>
                ))}

                <HStack w="full" gap="0">
                  <VStack w="full" gap="0">
                    {flatFill?.map((v, i) => (
                      <styled.div
                        key={v}
                        w="full"
                        p="3"
                        style={{
                          backgroundColor: `var(${v})`,
                          color: `var(${flatText[i]})`,
                        }}
                      >
                        <Text size="sm" fontWeight="medium">
                          {v}
                        </Text>
                      </styled.div>
                    ))}
                  </VStack>

                  <VStack w="full" gap="0">
                    {darkFill?.map((v, i) => (
                      <styled.div
                        key={v}
                        w="full"
                        p="3"
                        style={{
                          backgroundColor: `var(${v})`,
                          color: `var(${darkText[i]})`,
                        }}
                      >
                        <Text size="sm" fontWeight="medium">
                          {v}
                        </Text>
                      </styled.div>
                    ))}
                  </VStack>
                </HStack>
              </VStack>
            </Box>

            {/* Button Variants & States */}
            <Box
              w="full"
              bg="bg.default"
              borderRadius="lg"
              p="6"
              boxShadow="sm"
            >
              <VStack gap="6" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  Button Variants
                </Text>

                <HStack gap="4" flexWrap="wrap">
                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Solid
                    </Text>
                    <Button variant="solid" size="sm">
                      Small
                    </Button>
                    <Button variant="solid">Medium</Button>
                    <Button variant="solid" size="lg">
                      Large
                    </Button>
                    <Button variant="solid" disabled>
                      Disabled
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Outline
                    </Text>
                    <Button variant="outline" size="sm">
                      Small
                    </Button>
                    <Button variant="outline">Medium</Button>
                    <Button variant="outline" size="lg">
                      Large
                    </Button>
                    <Button variant="outline" disabled>
                      Disabled
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Ghost
                    </Text>
                    <Button variant="ghost" size="sm">
                      Small
                    </Button>
                    <Button variant="ghost">Medium</Button>
                    <Button variant="ghost" size="lg">
                      Large
                    </Button>
                    <Button variant="ghost" disabled>
                      Disabled
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Subtle
                    </Text>
                    <Button variant="subtle" size="sm">
                      Small
                    </Button>
                    <Button variant="subtle">Medium</Button>
                    <Button variant="subtle" size="lg">
                      Large
                    </Button>
                    <Button variant="subtle" disabled>
                      Disabled
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Link
                    </Text>
                    <Button variant="link" size="sm">
                      Small Link
                    </Button>
                    <Button variant="link">Medium Link</Button>
                    <Button variant="link" size="lg">
                      Large Link
                    </Button>
                    <Button variant="link" disabled>
                      Disabled Link
                    </Button>
                  </VStack>
                </HStack>

                <HStack gap="4" flexWrap="wrap">
                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Color Palettes
                    </Text>
                    <Button>Default</Button>
                    <Button colorPalette="red">Destructive</Button>
                    <Button colorPalette="green">Success</Button>
                    <Button colorPalette="amber">Warning</Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      With Icons
                    </Text>
                    <Button>
                      <Heart size={16} />
                      Like
                    </Button>
                    <Button variant="outline">
                      <Share size={16} />
                      Share
                    </Button>
                    <Button variant="ghost">
                      <MessageCircle size={16} />
                      Comment
                    </Button>
                  </VStack>
                </HStack>
              </VStack>
            </Box>

            {/* Typography Showcase */}
            <Box w="full" bg="bg.subtle" borderRadius="lg" p="6" boxShadow="sm">
              <VStack gap="4" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  Typography Scale
                </Text>

                <VStack gap="3" alignItems="start" w="full">
                  <Text size="7xl" color="fg.default">
                    7XL Heading
                  </Text>
                  <Text size="6xl" color="fg.default">
                    6XL Heading
                  </Text>
                  <Text size="5xl" color="fg.default">
                    5XL Heading
                  </Text>
                  <Text size="4xl" color="fg.default">
                    4XL Heading
                  </Text>
                  <Text size="3xl" color="fg.default">
                    3XL Heading
                  </Text>
                  <Text size="2xl" color="fg.default">
                    2XL Heading
                  </Text>
                  <Text size="xl" color="fg.default">
                    XL Heading
                  </Text>
                  <Text size="lg" color="fg.default">
                    Large Text
                  </Text>
                  <Text size="md" color="fg.default">
                    Medium Text (Body)
                  </Text>
                  <Text size="sm" color="fg.subtle">
                    Small Text (Secondary)
                  </Text>
                  <Text size="xs" color="fg.muted">
                    Extra Small Text (Captions)
                  </Text>
                </VStack>

                <VStack gap="2" alignItems="start" w="full" mt="4">
                  <Text size="sm" color="fg.subtle">
                    Color Variants
                  </Text>
                  <Text color="fg.default">Default foreground text</Text>
                  <Text color="fg.subtle">Subtle foreground text</Text>
                  <Text color="fg.muted">Muted foreground text</Text>
                  <Text color="fg.destructive">Destructive text</Text>
                  <Text color="fg.success">Success text</Text>
                  <Text color="fg.warning">Warning text</Text>
                  <Text color="fg.accent">Accent text</Text>
                </VStack>
              </VStack>
            </Box>

            {/* Cards & Content Containers */}
            <HStack gap="4" w="full" alignItems="start">
              {/* Post Card */}
              <Box
                flex="1"
                bg="bg.default"
                borderRadius="lg"
                p="6"
                boxShadow="sm"
              >
                <VStack gap="4" alignItems="start">
                  <HStack justify="space-between" w="full">
                    <HStack gap="3">
                      <Box w="10" h="10" bg="bg.muted" borderRadius="full" />
                      <VStack gap="1" alignItems="start">
                        <Text size="sm" fontWeight="semibold">
                          John Doe
                        </Text>
                        <Text size="xs" color="fg.muted">
                          2 hours ago
                        </Text>
                      </VStack>
                    </HStack>
                    <Button variant="ghost" size="sm">
                      <MoreVertical size={16} />
                    </Button>
                  </HStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="lg" fontWeight="semibold">
                      Sample Forum Post Title
                    </Text>
                    <Text color="fg.subtle">
                      This is a sample forum post content that demonstrates how
                      text appears in cards with proper contrast ratios and
                      semantic color usage.
                    </Text>
                  </VStack>

                  <HStack gap="2" flexWrap="wrap">
                    <Badge variant="subtle">Discussion</Badge>
                    <Badge variant="outline">Frontend</Badge>
                    <Badge variant="solid" colorPalette="green">
                      Solved
                    </Badge>
                  </HStack>

                  <HStack
                    justify="space-between"
                    w="full"
                    pt="2"
                    borderTopWidth="thin"
                    borderTopColor="border.subtle"
                  >
                    <HStack gap="4">
                      <Button variant="ghost" size="sm">
                        <Heart size={16} />
                        42
                      </Button>
                      <Button variant="ghost" size="sm">
                        <MessageCircle size={16} />
                        12
                      </Button>
                      <Button variant="ghost" size="sm">
                        <Share size={16} />
                        Share
                      </Button>
                    </HStack>
                    <Button variant="ghost" size="sm">
                      <ArrowUp size={16} />
                    </Button>
                  </HStack>
                </VStack>
              </Box>

              {/* Sidebar Widget */}
              <Box
                flex="1"
                bg="bg.subtle"
                borderRadius="lg"
                p="4"
                boxShadow="sm"
              >
                <VStack gap="3" alignItems="start">
                  <Text size="sm" fontWeight="semibold">
                    Quick Actions
                  </Text>
                  <VStack gap="2" w="full">
                    <Button
                      variant="outline"
                      size="sm"
                      w="full"
                      justifyContent="start"
                    >
                      <Plus size={16} />
                      New Post
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      w="full"
                      justifyContent="start"
                    >
                      <Filter size={16} />
                      Filter Posts
                    </Button>
                  </VStack>
                </VStack>
              </Box>
            </HStack>

            {/* Form Elements */}
            <Box
              w="full"
              bg="bg.default"
              borderRadius="lg"
              p="6"
              boxShadow="sm"
            >
              <VStack gap="6" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  Form Elements
                </Text>

                <HStack gap="8" w="full" alignItems="start">
                  <VStack gap="4" alignItems="start" flex="1">
                    <VStack gap="2" alignItems="start" w="full">
                      <Text size="sm" color="fg.subtle">
                        Input Fields
                      </Text>
                      <Input placeholder="Default input" />
                      <Input placeholder="Focused input" value="Sample text" />
                      <Input placeholder="Disabled input" disabled />
                    </VStack>

                    <VStack gap="2" alignItems="start">
                      <Text size="sm" color="fg.subtle">
                        Checkboxes
                      </Text>
                      <Checkbox>Unchecked</Checkbox>
                      <Checkbox checked>Checked</Checkbox>
                      <Checkbox checked="indeterminate">Indeterminate</Checkbox>
                      <Checkbox disabled>Disabled</Checkbox>
                    </VStack>
                  </VStack>

                  <VStack gap="4" alignItems="start" flex="1">
                    <VStack gap="2" alignItems="start" w="full">
                      <Text size="sm" color="fg.subtle">
                        Badge Variants
                      </Text>
                      <HStack gap="2" flexWrap="wrap">
                        <Badge variant="solid">Solid</Badge>
                        <Badge variant="subtle">Subtle</Badge>
                        <Badge variant="outline">Outline</Badge>
                      </HStack>
                      <HStack gap="2" flexWrap="wrap">
                        <Badge size="sm">Small</Badge>
                        <Badge size="md">Medium</Badge>
                        <Badge size="lg">Large</Badge>
                      </HStack>
                      <HStack gap="2" flexWrap="wrap">
                        <Badge colorPalette="red">Error</Badge>
                        <Badge colorPalette="green">Success</Badge>
                        <Badge colorPalette="amber">Warning</Badge>
                        <Badge colorPalette="blue">Info</Badge>
                      </HStack>
                    </VStack>
                  </VStack>
                </HStack>
              </VStack>
            </Box>

            {/* Background & Surface Testing */}
            <Box w="full" bg="bg.muted" borderRadius="lg" p="6" boxShadow="sm">
              <VStack gap="4" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  Surface Hierarchy
                </Text>

                <HStack gap="4" w="full">
                  <VStack gap="2" alignItems="start" flex="1">
                    <Text size="sm" color="fg.subtle">
                      Default Surface
                    </Text>
                    <Box bg="bg.default" p="4" borderRadius="md" w="full">
                      <Text>Content on default background</Text>
                      <Button variant="solid" size="sm" mt="2">
                        Action
                      </Button>
                    </Box>
                  </VStack>

                  <VStack gap="2" alignItems="start" flex="1">
                    <Text size="sm" color="fg.subtle">
                      Subtle Surface
                    </Text>
                    <Box bg="bg.subtle" p="4" borderRadius="md" w="full">
                      <Text>Content on subtle background</Text>
                      <Button variant="outline" size="sm" mt="2">
                        Action
                      </Button>
                    </Box>
                  </VStack>

                  <VStack gap="2" alignItems="start" flex="1">
                    <Text size="sm" color="fg.subtle">
                      Muted Surface
                    </Text>
                    <Box
                      bg="bg.muted"
                      p="4"
                      borderRadius="md"
                      w="full"
                      borderWidth="thin"
                      borderColor="border.default"
                    >
                      <Text>Content on muted background</Text>
                      <Button variant="ghost" size="sm" mt="2">
                        Action
                      </Button>
                    </Box>
                  </VStack>
                </HStack>
              </VStack>
            </Box>

            {/* State Indicators */}
            <Box
              w="full"
              bg="bg.default"
              borderRadius="lg"
              p="6"
              boxShadow="sm"
            >
              <VStack gap="4" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  State Indicators
                </Text>

                <HStack gap="6" w="full">
                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Status Messages
                    </Text>
                    <Box
                      bg="bg.success"
                      p="3"
                      borderRadius="md"
                      borderLeftWidth="thick"
                      borderLeftColor="border.success"
                    >
                      <Text size="sm" color="fg.success">
                        Success: Operation completed
                      </Text>
                    </Box>
                    <Box
                      bg="bg.warning"
                      p="3"
                      borderRadius="md"
                      borderLeftWidth="thick"
                      borderLeftColor="border.warning"
                    >
                      <Text size="sm" color="fg.warning">
                        Warning: Please review
                      </Text>
                    </Box>
                    <Box
                      bg="bg.destructive"
                      p="3"
                      borderRadius="md"
                      borderLeftWidth="thick"
                      borderLeftColor="border.destructive"
                    >
                      <Text size="sm" color="fg.destructive">
                        Error: Something went wrong
                      </Text>
                    </Box>
                    <Box
                      bg="bg.info"
                      p="3"
                      borderRadius="md"
                      borderLeftWidth="thick"
                      borderLeftColor="border.info"
                    >
                      <Text size="sm" color="fg.info">
                        Info: Additional context
                      </Text>
                    </Box>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      Content States
                    </Text>
                    <Box
                      bg="visibility.published.bg"
                      p="3"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="visibility.published.border"
                    >
                      <Text size="sm" color="visibility.published.fg">
                        Published Content
                      </Text>
                    </Box>
                    <Box
                      bg="visibility.draft.bg"
                      p="3"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="visibility.draft.border"
                    >
                      <Text size="sm" color="visibility.draft.fg">
                        Draft Content
                      </Text>
                    </Box>
                    <Box
                      bg="visibility.review.bg"
                      p="3"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="visibility.review.border"
                    >
                      <Text size="sm" color="visibility.review.fg">
                        Under Review
                      </Text>
                    </Box>
                    <Box
                      bg="visibility.unlisted.bg"
                      p="3"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="visibility.unlisted.border"
                    >
                      <Text size="sm" color="visibility.unlisted.fg">
                        Unlisted Content
                      </Text>
                    </Box>
                  </VStack>
                </HStack>
              </VStack>
            </Box>
          </VStack>
        </Box>

        {/* Right Sidebar - Fixed */}
        <Box
          pos="fixed"
          right="0"
          top="20"
          bottom="0"
          w="80"
          bg="bg.opaque"
          backdropBlur="frosted"
          backdropFilter="auto"
          borderLeftWidth="thick"
          borderLeftColor="border.subtle"
          p="4"
          display={{ base: "none", xl: "block" }}
          overflowY="auto"
        >
          <VStack gap="6" alignItems="start">
            <VStack gap="2" alignItems="start" w="full">
              <Text size="sm" fontWeight="semibold" color="fg.subtle">
                Recent Activity
              </Text>
              <VStack gap="2" w="full">
                {[...Array(5)].map((_, i) => (
                  <Box key={i} p="3" bg="bg.subtle" borderRadius="md" w="full">
                    <Text size="xs" color="fg.muted">
                      {i + 1}h ago
                    </Text>
                    <Text size="sm">User activity {i + 1}</Text>
                  </Box>
                ))}
              </VStack>
            </VStack>

            <VStack gap="2" alignItems="start" w="full">
              <Text size="sm" fontWeight="semibold" color="fg.subtle">
                Popular Tags
              </Text>
              <HStack gap="1" flexWrap="wrap">
                <Badge size="sm" variant="subtle">
                  React
                </Badge>
                <Badge size="sm" variant="subtle">
                  TypeScript
                </Badge>
                <Badge size="sm" variant="subtle">
                  Design
                </Badge>
                <Badge size="sm" variant="subtle">
                  Backend
                </Badge>
              </HStack>
            </VStack>
          </VStack>
        </Box>
      </HStack>
    </Box>
  );
}
