import { Box, VStack } from "@/styled-system/jsx";

export function ScreenshotSection() {
  return (
    <Box
      maxW="100vw"
      w="full"
      maxH={{
        base: "30vh",
        sm: "40vh",
        md: "50vh",
        lg: "70vh",
      }}
      overflowY="hidden"
      bgColor="black"
    >
      <VStack
        position="relative"
        zIndex="20"
        w="full"
        paddingX={{
          base: "4",
          sm: "8",
          md: "12",
          xl: "16",
        }}
      >
        <picture>
          <source
            media="(max-width: 768px)"
            srcSet="2025_app_screenshot_viewport.png"
          />
          <source media="(min-width: 768px)" srcSet="2025_app_screenshot.png" />
          <source media="(min-width: 768px)" srcSet="2025_app_screenshot.png" />
          <img
            src="2025_app_screenshot.png"
            alt=""
            role="presentation"
            width={1469}
            height={961}
          />
        </picture>
      </VStack>
    </Box>
  );
}
