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
import { useI18n } from "@/i18n/provider";
import { Box, HStack, VStack, styled } from "@/styled-system/jsx";
import { getColourVariants } from "@/utils/colour";

const DEMO_ACCENT = "#3b82f6";

export default function Page() {
  const { t } = useI18n();
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
                {t("Navigation")}
              </Text>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                <User size={16} />
                {t("Profile")}
              </Button>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                <Settings size={16} />
                {t("Settings")}
              </Button>
              <Button variant="solid" size="sm" justifyContent="start" w="full">
                <Plus size={16} />
                {t("Create Post")}
              </Button>
            </VStack>

            <VStack gap="2" alignItems="start" w="full">
              <Text size="sm" fontWeight="semibold" color="fg.subtle">
                {t("Categories")}
              </Text>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                {t("General Discussion")}
              </Button>
              <Button variant="ghost" size="sm" justifyContent="start" w="full">
                {t("Announcements")}
              </Button>
              <Button
                variant="outline"
                size="sm"
                justifyContent="start"
                w="full"
              >
                {t("Help & Support")}
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
                  {t("Accent Color Palette")}
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
                  {t("Button Variants")}
                </Text>

                <HStack gap="4" flexWrap="wrap">
                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Solid")}
                    </Text>
                    <Button variant="solid" size="sm">
                      {t("Small")}
                    </Button>
                    <Button variant="solid">{t("Medium")}</Button>
                    <Button variant="solid" size="lg">
                      {t("Large")}
                    </Button>
                    <Button variant="solid" disabled>
                      {t("Disabled")}
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Outline")}
                    </Text>
                    <Button variant="outline" size="sm">
                      {t("Small")}
                    </Button>
                    <Button variant="outline">{t("Medium")}</Button>
                    <Button variant="outline" size="lg">
                      {t("Large")}
                    </Button>
                    <Button variant="outline" disabled>
                      {t("Disabled")}
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Ghost")}
                    </Text>
                    <Button variant="ghost" size="sm">
                      {t("Small")}
                    </Button>
                    <Button variant="ghost">{t("Medium")}</Button>
                    <Button variant="ghost" size="lg">
                      {t("Large")}
                    </Button>
                    <Button variant="ghost" disabled>
                      {t("Disabled")}
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Subtle")}
                    </Text>
                    <Button variant="subtle" size="sm">
                      {t("Small")}
                    </Button>
                    <Button variant="subtle">{t("Medium")}</Button>
                    <Button variant="subtle" size="lg">
                      {t("Large")}
                    </Button>
                    <Button variant="subtle" disabled>
                      {t("Disabled")}
                    </Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Link")}
                    </Text>
                    <Button variant="link" size="sm">
                      {t("Small Link")}
                    </Button>
                    <Button variant="link">{t("Medium Link")}</Button>
                    <Button variant="link" size="lg">
                      {t("Large Link")}
                    </Button>
                    <Button variant="link" disabled>
                      {t("Disabled Link")}
                    </Button>
                  </VStack>
                </HStack>

                <HStack gap="4" flexWrap="wrap">
                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Color Palettes")}
                    </Text>
                    <Button>{t("Default")}</Button>
                    <Button colorPalette="red">{t("Destructive")}</Button>
                    <Button colorPalette="green">{t("Success")}</Button>
                    <Button colorPalette="amber">{t("Warning")}</Button>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("With Icons")}
                    </Text>
                    <Button>
                      <Heart size={16} />
                      {t("Like")}
                    </Button>
                    <Button variant="outline">
                      <Share size={16} />
                      {t("Share")}
                    </Button>
                    <Button variant="ghost">
                      <MessageCircle size={16} />
                      {t("Comment")}
                    </Button>
                  </VStack>
                </HStack>
              </VStack>
            </Box>

            {/* Typography Showcase */}
            <Box w="full" bg="bg.subtle" borderRadius="lg" p="6" boxShadow="sm">
              <VStack gap="4" alignItems="start">
                <Text size="lg" fontWeight="semibold">
                  {t("Typography Scale")}
                </Text>

                <VStack gap="3" alignItems="start" w="full">
                  <Text size="7xl" color="fg.default">
                    {t("7XL Heading")}
                  </Text>
                  <Text size="6xl" color="fg.default">
                    {t("6XL Heading")}
                  </Text>
                  <Text size="5xl" color="fg.default">
                    {t("5XL Heading")}
                  </Text>
                  <Text size="4xl" color="fg.default">
                    {t("4XL Heading")}
                  </Text>
                  <Text size="3xl" color="fg.default">
                    {t("3XL Heading")}
                  </Text>
                  <Text size="2xl" color="fg.default">
                    {t("2XL Heading")}
                  </Text>
                  <Text size="xl" color="fg.default">
                    {t("XL Heading")}
                  </Text>
                  <Text size="lg" color="fg.default">
                    {t("Large Text")}
                  </Text>
                  <Text size="md" color="fg.default">
                    {t("Medium Text (Body)")}
                  </Text>
                  <Text size="sm" color="fg.subtle">
                    {t("Small Text (Secondary)")}
                  </Text>
                  <Text size="xs" color="fg.muted">
                    {t("Extra Small Text (Captions)")}
                  </Text>
                </VStack>

                <VStack gap="2" alignItems="start" w="full" mt="4">
                  <Text size="sm" color="fg.subtle">
                    {t("Color Variants")}
                  </Text>
                  <Text color="fg.default">{t("Default foreground text")}</Text>
                  <Text color="fg.subtle">{t("Subtle foreground text")}</Text>
                  <Text color="fg.muted">{t("Muted foreground text")}</Text>
                  <Text color="fg.destructive">{t("Destructive text")}</Text>
                  <Text color="fg.success">{t("Success text")}</Text>
                  <Text color="fg.warning">{t("Warning text")}</Text>
                  <Text color="fg.accent">{t("Accent text")}</Text>
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
                          {t("2 hours ago")}
                        </Text>
                      </VStack>
                    </HStack>
                    <Button variant="ghost" size="sm">
                      <MoreVertical size={16} />
                    </Button>
                  </HStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="lg" fontWeight="semibold">
                      {t("Sample Forum Post Title")}
                    </Text>
                    <Text color="fg.subtle">
                      {t(
                        "This is a sample forum post content that demonstrates how text appears in cards with proper contrast ratios and semantic color usage.",
                      )}
                    </Text>
                  </VStack>

                  <HStack gap="2" flexWrap="wrap">
                    <Badge variant="subtle">{t("Discussion")}</Badge>
                    <Badge variant="outline">{t("Frontend")}</Badge>
                    <Badge variant="solid" colorPalette="green">
                      {t("Solved")}
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
                        {t("Share")}
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
                    {t("Quick Actions")}
                  </Text>
                  <VStack gap="2" w="full">
                    <Button
                      variant="outline"
                      size="sm"
                      w="full"
                      justifyContent="start"
                    >
                      <Plus size={16} />
                      {t("New Post")}
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      w="full"
                      justifyContent="start"
                    >
                      <Filter size={16} />
                      {t("Filter Posts")}
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
                  {t("Form Elements")}
                </Text>

                <HStack gap="8" w="full" alignItems="start">
                  <VStack gap="4" alignItems="start" flex="1">
                    <VStack gap="2" alignItems="start" w="full">
                      <Text size="sm" color="fg.subtle">
                        {t("Input Fields")}
                      </Text>
                      <Input placeholder={t("Default input")} />
                      <Input
                        placeholder={t("Focused input")}
                        value={t("Sample text")}
                      />
                      <Input placeholder={t("Disabled input")} disabled />
                    </VStack>

                    <VStack gap="2" alignItems="start">
                      <Text size="sm" color="fg.subtle">
                        {t("Checkboxes")}
                      </Text>
                      <Checkbox>{t("Unchecked")}</Checkbox>
                      <Checkbox checked>{t("Checked")}</Checkbox>
                      <Checkbox checked="indeterminate">
                        {t("Indeterminate")}
                      </Checkbox>
                      <Checkbox disabled>{t("Disabled")}</Checkbox>
                    </VStack>
                  </VStack>

                  <VStack gap="4" alignItems="start" flex="1">
                    <VStack gap="2" alignItems="start" w="full">
                      <Text size="sm" color="fg.subtle">
                        {t("Badge Variants")}
                      </Text>
                      <HStack gap="2" flexWrap="wrap">
                        <Badge variant="solid">{t("Solid")}</Badge>
                        <Badge variant="subtle">{t("Subtle")}</Badge>
                        <Badge variant="outline">{t("Outline")}</Badge>
                      </HStack>
                      <HStack gap="2" flexWrap="wrap">
                        <Badge size="sm">{t("Small")}</Badge>
                        <Badge size="md">{t("Medium")}</Badge>
                        <Badge size="lg">{t("Large")}</Badge>
                      </HStack>
                      <HStack gap="2" flexWrap="wrap">
                        <Badge colorPalette="red">{t("Error")}</Badge>
                        <Badge colorPalette="green">{t("Success")}</Badge>
                        <Badge colorPalette="amber">{t("Warning")}</Badge>
                        <Badge colorPalette="blue">{t("Info")}</Badge>
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
                  {t("Surface Hierarchy")}
                </Text>

                <HStack gap="4" w="full">
                  <VStack gap="2" alignItems="start" flex="1">
                    <Text size="sm" color="fg.subtle">
                      {t("Default Surface")}
                    </Text>
                    <Box bg="bg.default" p="4" borderRadius="md" w="full">
                      <Text>{t("Content on default background")}</Text>
                      <Button variant="solid" size="sm" mt="2">
                        {t("Action")}
                      </Button>
                    </Box>
                  </VStack>

                  <VStack gap="2" alignItems="start" flex="1">
                    <Text size="sm" color="fg.subtle">
                      {t("Subtle Surface")}
                    </Text>
                    <Box bg="bg.subtle" p="4" borderRadius="md" w="full">
                      <Text>{t("Content on subtle background")}</Text>
                      <Button variant="outline" size="sm" mt="2">
                        {t("Action")}
                      </Button>
                    </Box>
                  </VStack>

                  <VStack gap="2" alignItems="start" flex="1">
                    <Text size="sm" color="fg.subtle">
                      {t("Muted Surface")}
                    </Text>
                    <Box
                      bg="bg.muted"
                      p="4"
                      borderRadius="md"
                      w="full"
                      borderWidth="thin"
                      borderColor="border.default"
                    >
                      <Text>{t("Content on muted background")}</Text>
                      <Button variant="ghost" size="sm" mt="2">
                        {t("Action")}
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
                  {t("State Indicators")}
                </Text>

                <HStack gap="6" w="full">
                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Status Messages")}
                    </Text>
                    <Box
                      bg="bg.success"
                      p="3"
                      borderRadius="md"
                      borderLeftWidth="thick"
                      borderLeftColor="border.success"
                    >
                      <Text size="sm" color="fg.success">
                        {t("Success: Operation completed")}
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
                        {t("Warning: Please review")}
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
                        {t("Error: Something went wrong")}
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
                        {t("Info: Additional context")}
                      </Text>
                    </Box>
                  </VStack>

                  <VStack gap="2" alignItems="start">
                    <Text size="sm" color="fg.subtle">
                      {t("Content States")}
                    </Text>
                    <Box
                      bg="visibility.published.bg"
                      p="3"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="visibility.published.border"
                    >
                      <Text size="sm" color="visibility.published.fg">
                        {t("Published Content")}
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
                        {t("Draft Content")}
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
                        {t("Under Review")}
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
                        {t("Unlisted Content")}
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
                {t("Recent Activity")}
              </Text>
              <VStack gap="2" w="full">
                {[...Array(5)].map((_, i) => (
                  <Box key={i} p="3" bg="bg.subtle" borderRadius="md" w="full">
                    <Text size="xs" color="fg.muted">
                      {t("{{count}}h ago", { count: i + 1 })}
                    </Text>
                    <Text size="sm">
                      {t("User activity {{count}}", { count: i + 1 })}
                    </Text>
                  </Box>
                ))}
              </VStack>
            </VStack>

            <VStack gap="2" alignItems="start" w="full">
              <Text size="sm" fontWeight="semibold" color="fg.subtle">
                {t("Popular Tags")}
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
