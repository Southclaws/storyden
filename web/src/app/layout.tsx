import { PropsWithChildren } from "react";

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
