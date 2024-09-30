import { DefaultSeo } from "next-seo";

import "./globals.css";
import "../fonts.css";

import localFont from "next/font/local";
import Head from "next/head";

import { seo } from "../seo";
import Script from "next/script";

const monasans = localFont({
  src: "./mona-sans.woff2",
  display: "swap",
  variable: "--font-mona-sans",
});

export default function MyApp({ Component, pageProps }) {
  return (
    <>
      <Head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width" />
        <meta name="msapplication-TileColor" content="#303030" />
        <meta name="theme-color" content="#303030" />

        {/* Icons */}
        <link
          rel="apple-touch-icon"
          sizes="180x180"
          href="/apple-touch-icon.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="32x32"
          href="/favicon-32x32.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="16x16"
          href="/favicon-16x16.png"
        />
        <link rel="manifest" href="/site.webmanifest" />
        <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#303030" />

        <Script type="text/javascript">
          {`(function(c,l,a,r,i,t,y){
                c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
                t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
                y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
            })(window, document, "clarity", "script", "obgioniw76");`}
        </Script>
      </Head>

      <style jsx global>
        {`
          :root {
            --font-mona-sans: ${monasans.style.fontFamily};
          }
        `}
      </style>

      <DefaultSeo {...seo} />

      <Component {...pageProps} />
    </>
  );
}
