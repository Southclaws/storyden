import { PropsWithChildren } from "react";

import "./global.css";

import { Providers } from "./providers";

export default function RootLayout({ children }: PropsWithChildren) {
  return (
    <html lang="en">
      <body suppressHydrationWarning={true}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}

export const metadata = {
  manifest: "/manifest.json",
};
