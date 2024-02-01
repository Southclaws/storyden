import { Box } from "@/styled-system/jsx";
import { useTheme } from "nextra-theme-docs";
import { seo } from "./seo";

export default {
  logo: (
    <Box id="logo" p="2">
      <Logo />
    </Box>
  ),
  useNextSeoProps() {
    return seo;
  },
  project: {
    link: "https://github.com/Southclaws/storyden",
  },
  docsRepositoryBase: "https://github.com/Southclaws/storyden/blob/main/home",
};

function Logo() {
  const { resolvedTheme } = useTheme();

  const colour = resolvedTheme === "dark" ? "#f8fcfa" : "#303030";

  return (
    <svg
      className="storyden-logo"
      width="150"
      height="50"
      viewBox="0 0 150 50"
      fill={colour}
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill={colour}
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.7131 1.21172C23.9797 -0.0425755 26.0203 -0.0425791 27.2869 1.21172L46.2869 20.0272C46.9032 20.6376 47.25 21.4691 47.25 22.3365V47.4792C47.25 48.7219 46.2426 49.7292 45 49.7292C43.7574 49.7292 42.75 48.7219 42.75 47.4792V23.2749C42.75 23.008 42.6433 22.7522 42.4536 22.5644L25.7036 5.97702C25.3139 5.59108 24.6861 5.59108 24.2964 5.97702L19.0658 11.1568C18.1828 12.0312 16.7582 12.0242 15.8838 11.1413C15.0094 10.2583 15.0164 8.83373 15.8994 7.95935L22.7131 1.21172ZM17.4038 25.1813C17.01 25.5713 17.0085 26.2071 17.4003 26.599L24.2929 33.4915C24.6834 33.8821 25.3166 33.8821 25.7071 33.4915L32.5997 26.599C32.9915 26.2071 32.99 25.5713 32.5962 25.1813L25.7036 18.3557C25.3139 17.9697 24.6861 17.9697 24.2964 18.3557L17.4038 25.1813ZM28.8891 38.0877C28.4986 37.6972 28.4986 37.064 28.8891 36.6735L37.3804 28.1822C38.654 26.9086 38.649 24.8422 37.3692 23.5748L27.2869 13.5904C26.0203 12.3361 23.9797 12.3361 22.7131 13.5904L12.6308 23.5748C11.351 24.8422 11.346 26.9086 12.6196 28.1822L21.1109 36.6735C21.5014 37.064 21.5014 37.6972 21.1109 38.0877L14.2623 44.9363C14.0748 45.1239 13.8204 45.2292 13.5552 45.2292H8.25C7.69772 45.2292 7.25 44.7815 7.25 44.2292V23.2711C7.25 23.0049 7.35615 22.7497 7.54493 22.562L11.271 18.8573C12.1522 17.9812 12.1563 16.5566 11.2801 15.6754C10.404 14.7942 8.97935 14.7901 8.09814 15.6662L3.70852 20.0307C3.09498 20.6407 2.75 21.4702 2.75 22.3354V46.4792C2.75 48.2741 4.20508 49.7292 6 49.7292H14.4872C15.3491 49.7292 16.1758 49.3868 16.7853 48.7773L24.2929 41.2697C24.6834 40.8792 25.3166 40.8792 25.7071 41.2697L33.5076 49.0702C34.3863 49.9489 35.8109 49.9489 36.6896 49.0702C37.5683 48.1915 37.5683 46.7669 36.6896 45.8882L28.8891 38.0877Z"
      />
      <path
        fill={colour}
        d="M68.7406 21.6317C68.3586 21.3214 67.9767 21.0946 67.5948 20.9514C67.2129 20.7963 66.8429 20.7187 66.4848 20.7187C66.0313 20.7187 65.6613 20.8261 65.3749 21.0409C65.0884 21.2558 64.9452 21.5362 64.9452 21.8823C64.9452 22.121 65.0168 22.318 65.16 22.4731C65.3033 22.6283 65.4883 22.7655 65.715 22.8849C65.9537 22.9923 66.2163 23.0878 66.5027 23.1713C66.8011 23.2549 67.0935 23.3444 67.38 23.4399C68.5257 23.8218 69.3612 24.335 69.8863 24.9795C70.4234 25.612 70.6919 26.4415 70.6919 27.4679C70.6919 28.1602 70.5726 28.7868 70.3339 29.3477C70.1071 29.9087 69.767 30.392 69.3134 30.7978C68.8718 31.1917 68.3228 31.496 67.6664 31.7108C67.0219 31.9376 66.2879 32.051 65.4644 32.051C63.7577 32.051 62.1763 31.5438 60.7202 30.5293L62.224 27.7007C62.7492 28.1661 63.2683 28.5123 63.7816 28.739C64.2948 28.9658 64.802 29.0792 65.3033 29.0792C65.8762 29.0792 66.2998 28.9479 66.5744 28.6853C66.8608 28.4227 67.004 28.1244 67.004 27.7902C67.004 27.5873 66.9682 27.4142 66.8966 27.271C66.825 27.1159 66.7056 26.9786 66.5385 26.8593C66.3715 26.728 66.1507 26.6086 65.8762 26.5012C65.6136 26.3938 65.2913 26.2744 64.9094 26.1432C64.4559 25.9999 64.0083 25.8448 63.5667 25.6777C63.1371 25.4987 62.7492 25.2659 62.4031 24.9795C62.0689 24.693 61.7944 24.335 61.5795 23.9053C61.3766 23.4637 61.2752 22.9088 61.2752 22.2404C61.2752 21.572 61.3826 20.9693 61.5974 20.4322C61.8242 19.8832 62.1345 19.4178 62.5284 19.0358C62.9342 18.642 63.4235 18.3376 63.9964 18.1228C64.5812 17.908 65.2317 17.8006 65.9478 17.8006C66.6161 17.8006 67.3143 17.896 68.0424 18.087C68.7704 18.266 69.4686 18.5346 70.137 18.8926L68.7406 21.6317Z"
      />
      <path
        fill={colour}
        d="M75.1371 25.1943V31.675H71.8968V25.1943H70.8226V22.491H71.8968V19.734H75.1371V22.491H76.9811V25.1943H75.1371Z"
      />
      <path
        fill={colour}
        d="M80.7001 27.0383C80.7001 27.3486 80.7538 27.635 80.8612 27.8976C80.9806 28.1482 81.1298 28.369 81.3088 28.56C81.4998 28.751 81.7206 28.9001 81.9712 29.0076C82.2338 29.115 82.5083 29.1687 82.7947 29.1687C83.0812 29.1687 83.3497 29.115 83.6003 29.0076C83.8629 28.9001 84.0837 28.751 84.2627 28.56C84.4537 28.369 84.6029 28.1482 84.7103 27.8976C84.8296 27.635 84.8893 27.3546 84.8893 27.0562C84.8893 26.7697 84.8296 26.5012 84.7103 26.2506C84.6029 25.988 84.4537 25.7612 84.2627 25.5703C84.0837 25.3793 83.8629 25.2301 83.6003 25.1227C83.3497 25.0153 83.0812 24.9616 82.7947 24.9616C82.5083 24.9616 82.2338 25.0153 81.9712 25.1227C81.7206 25.2301 81.4998 25.3793 81.3088 25.5703C81.1298 25.7612 80.9806 25.982 80.8612 26.2327C80.7538 26.4833 80.7001 26.7518 80.7001 27.0383ZM77.2807 27.0025C77.2807 26.3222 77.418 25.6896 77.6925 25.1048C77.967 24.5081 78.3489 23.9948 78.8383 23.5652C79.3276 23.1236 79.9064 22.7775 80.5748 22.5268C81.2551 22.2762 81.9951 22.1509 82.7947 22.1509C83.5824 22.1509 84.3105 22.2762 84.9788 22.5268C85.6591 22.7655 86.2439 23.1057 86.7333 23.5473C87.2346 23.9769 87.6224 24.4961 87.897 25.1048C88.1715 25.7016 88.3087 26.364 88.3087 27.092C88.3087 27.82 88.1655 28.4884 87.879 29.0971C87.6045 29.6938 87.2226 30.213 86.7333 30.6546C86.2439 31.0843 85.6532 31.4184 84.9609 31.6571C84.2806 31.8958 83.5407 32.0152 82.741 32.0152C81.9533 32.0152 81.2253 31.8958 80.5569 31.6571C79.8885 31.4184 79.3097 31.0783 78.8204 30.6367C78.343 30.1951 77.967 29.67 77.6925 29.0613C77.418 28.4406 77.2807 27.7544 77.2807 27.0025Z"
      />
      <path
        fill={colour}
        d="M88.9218 22.491H92.1621V23.9948C92.5083 23.4458 92.9319 23.0281 93.4332 22.7417C93.9345 22.4433 94.5193 22.2941 95.1877 22.2941C95.2712 22.2941 95.3607 22.2941 95.4562 22.2941C95.5636 22.2941 95.683 22.306 95.8143 22.3299V25.4271C95.3846 25.2122 94.9191 25.1048 94.4179 25.1048C93.666 25.1048 93.099 25.3316 92.7171 25.7851C92.3471 26.2267 92.1621 26.8772 92.1621 27.7365V31.675H88.9218V22.491Z"
      />
      <path
        fill={colour}
        d="M100.627 30.4219L96.1875 22.491H99.947L102.4 27.1994L104.781 22.491H108.504L101.129 36.2581H97.4944L100.627 30.4219Z"
      />
      <path
        fill={colour}
        d="M110.64 27.0562C110.64 27.3546 110.694 27.635 110.801 27.8976C110.909 28.1482 111.052 28.369 111.231 28.56C111.422 28.751 111.643 28.9001 111.893 29.0076C112.156 29.115 112.436 29.1687 112.735 29.1687C113.021 29.1687 113.29 29.115 113.54 29.0076C113.803 28.9001 114.024 28.751 114.203 28.56C114.394 28.369 114.543 28.1482 114.65 27.8976C114.77 27.647 114.829 27.3784 114.829 27.092C114.829 26.8055 114.77 26.537 114.65 26.2864C114.543 26.0238 114.394 25.797 114.203 25.6061C114.024 25.4151 113.803 25.2659 113.54 25.1585C113.29 25.0511 113.021 24.9974 112.735 24.9974C112.448 24.9974 112.174 25.0511 111.911 25.1585C111.661 25.2659 111.44 25.4151 111.249 25.6061C111.07 25.797 110.921 26.0178 110.801 26.2685C110.694 26.5072 110.64 26.7697 110.64 27.0562ZM114.722 16.8159H117.98V31.675H114.722V30.6546C114.03 31.5259 113.093 31.9615 111.911 31.9615C111.243 31.9615 110.628 31.8362 110.067 31.5855C109.506 31.3349 109.017 30.9888 108.599 30.5472C108.182 30.1056 107.853 29.5864 107.615 28.9897C107.388 28.3929 107.275 27.7484 107.275 27.0562C107.275 26.3998 107.382 25.7791 107.597 25.1943C107.824 24.5976 108.14 24.0784 108.546 23.6368C108.951 23.1952 109.435 22.8491 109.996 22.5984C110.569 22.3359 111.195 22.2046 111.876 22.2046C113.021 22.2046 113.97 22.6044 114.722 23.4041V16.8159Z"
      />
      <path
        fill={colour}
        d="M125.611 25.7493C125.504 25.2958 125.283 24.9317 124.949 24.6572C124.614 24.3827 124.209 24.2455 123.731 24.2455C123.23 24.2455 122.818 24.3768 122.496 24.6393C122.186 24.9019 121.989 25.2719 121.905 25.7493H125.611ZM121.816 27.6291C121.816 29.0255 122.472 29.7237 123.785 29.7237C124.489 29.7237 125.02 29.4372 125.378 28.8643H128.511C127.879 30.9649 126.297 32.0152 123.767 32.0152C122.991 32.0152 122.281 31.9018 121.637 31.675C120.992 31.4363 120.437 31.1022 119.972 30.6725C119.518 30.2428 119.166 29.7296 118.915 29.1329C118.665 28.5361 118.539 27.8678 118.539 27.1278C118.539 26.364 118.659 25.6777 118.897 25.069C119.136 24.4484 119.476 23.9232 119.918 23.4936C120.36 23.0639 120.891 22.7357 121.511 22.5089C122.144 22.2702 122.854 22.1509 123.642 22.1509C124.417 22.1509 125.116 22.2702 125.736 22.5089C126.357 22.7357 126.882 23.0699 127.312 23.5115C127.741 23.9531 128.07 24.4961 128.296 25.1406C128.523 25.7732 128.636 26.4893 128.636 27.2889V27.6291H121.816Z"
      />
      <path
        fill={colour}
        d="M129.521 22.491H132.762V23.6547C133.203 23.1176 133.651 22.7596 134.104 22.5805C134.558 22.3896 135.089 22.2941 135.698 22.2941C136.342 22.2941 136.891 22.4015 137.345 22.6164C137.81 22.8192 138.204 23.1117 138.526 23.4936C138.789 23.8039 138.968 24.15 139.063 24.5319C139.159 24.9138 139.207 25.3495 139.207 25.8388V31.675H135.966V27.0383C135.966 26.5848 135.931 26.2207 135.859 25.9462C135.799 25.6598 135.686 25.433 135.519 25.2659C135.376 25.1227 135.214 25.0213 135.035 24.9616C134.856 24.9019 134.665 24.8721 134.463 24.8721C133.914 24.8721 133.49 25.0392 133.191 25.3733C132.905 25.6956 132.762 26.1611 132.762 26.7697V31.675H129.521V22.491Z"
      />
    </svg>
  );
}
