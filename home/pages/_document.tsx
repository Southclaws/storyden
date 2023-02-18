// pages/_document.js
import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html lang="en">
      <Head>
        <link
          rel="preload"
          as="font"
          href="https://use.typekit.net/gnq7poa.css"
        ></link>
      </Head>
      <body>
        <Main />
        <NextScript />f
      </body>
    </Html>
  );
}
