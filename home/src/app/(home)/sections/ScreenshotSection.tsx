import { Box, VStack } from "@/styled-system/jsx";

export function ScreenshotSection() {
  return (
    <Box
      maxW="100vw"
      w="full"
      maxH={{
        base: "calc(min(100vw - 32px, 390px) * 1.68)",
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
        <Box as="picture" maxW={{ base: "390px", md: "none" }}>
          <source
            media="(max-width: 767px)"
            type="image/webp"
            srcSet="/2026_app_screenshot_mobile_390.webp 390w, /2026_app_screenshot_mobile.webp 780w"
            sizes="(min-width: 422px) 390px, calc(100vw - 32px)"
            width={780}
            height={1544}
          />
          <source
            media="(max-width: 767px)"
            srcSet="/2026_app_screenshot_mobile.png"
            width={780}
            height={1544}
          />
          <source
            type="image/webp"
            srcSet="/2026_app_screenshot_1024.webp 1024w, /2026_app_screenshot.webp 1469w"
            sizes="(min-width: 1597px) 1469px, (min-width: 1280px) calc(100vw - 128px), calc(100vw - 96px)"
          />
          {/* oxlint-disable-next-line next/no-img-element -- Pre-optimized picture sources provide separate mobile and desktop crops. */}
          <img
            src="/2026_app_screenshot.png"
            alt=""
            role="presentation"
            width={1469}
            height={961}
            decoding="async"
          />
        </Box>
      </VStack>
    </Box>
  );
}
