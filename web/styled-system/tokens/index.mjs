const tokens = {
  "animations.backdrop-in": {
    "value": "fade-in 250ms var(--easings-emphasized-in)",
    "variable": "var(--animations-backdrop-in)"
  },
  "animations.backdrop-out": {
    "value": "fade-out 200ms var(--easings-emphasized-out)",
    "variable": "var(--animations-backdrop-out)"
  },
  "animations.dialog-in": {
    "value": "slide-in 400ms var(--easings-emphasized-in)",
    "variable": "var(--animations-dialog-in)"
  },
  "animations.dialog-out": {
    "value": "slide-out 200ms var(--easings-emphasized-out)",
    "variable": "var(--animations-dialog-out)"
  },
  "animations.drawer-in-left": {
    "value": "slide-in-left 400ms var(--easings-emphasized-in)",
    "variable": "var(--animations-drawer-in-left)"
  },
  "animations.drawer-out-left": {
    "value": "slide-out-left 200ms var(--easings-emphasized-out)",
    "variable": "var(--animations-drawer-out-left)"
  },
  "animations.drawer-in-right": {
    "value": "slide-in-right 400ms var(--easings-emphasized-in)",
    "variable": "var(--animations-drawer-in-right)"
  },
  "animations.drawer-out-right": {
    "value": "slide-out-right 200ms var(--easings-emphasized-out)",
    "variable": "var(--animations-drawer-out-right)"
  },
  "animations.skeleton-pulse": {
    "value": "skeleton-pulse 2s var(--easings-pulse) infinite",
    "variable": "var(--animations-skeleton-pulse)"
  },
  "animations.fade-in": {
    "value": "fade-in 400ms var(--easings-emphasized-in)",
    "variable": "var(--animations-fade-in)"
  },
  "animations.collapse-in": {
    "value": "collapse-in 250ms var(--easings-emphasized-in)",
    "variable": "var(--animations-collapse-in)"
  },
  "animations.collapse-out": {
    "value": "collapse-out 200ms var(--easings-emphasized-out)",
    "variable": "var(--animations-collapse-out)"
  },
  "animations.spin": {
    "value": "spin 1s linear infinite",
    "variable": "var(--animations-spin)"
  },
  "animations.target-pulse": {
    "value": "targetPulse 1s var(--easings-pulse) 2",
    "variable": "var(--animations-target-pulse)"
  },
  "blurs.sm": {
    "value": "4px",
    "variable": "var(--blurs-sm)"
  },
  "blurs.base": {
    "value": "8px",
    "variable": "var(--blurs-base)"
  },
  "blurs.md": {
    "value": "12px",
    "variable": "var(--blurs-md)"
  },
  "blurs.lg": {
    "value": "16px",
    "variable": "var(--blurs-lg)"
  },
  "blurs.xl": {
    "value": "24px",
    "variable": "var(--blurs-xl)"
  },
  "blurs.2xl": {
    "value": "40px",
    "variable": "var(--blurs-2xl)"
  },
  "blurs.3xl": {
    "value": "64px",
    "variable": "var(--blurs-3xl)"
  },
  "borders.none": {
    "value": "none",
    "variable": "var(--borders-none)"
  },
  "colors.current": {
    "value": "currentColor",
    "variable": "var(--colors-current)"
  },
  "colors.black": {
    "value": "#000000",
    "variable": "var(--colors-black)"
  },
  "colors.black.a1": {
    "value": "rgba(0, 0, 0, 0.05)",
    "variable": "var(--colors-black-a1)"
  },
  "colors.black.a2": {
    "value": "rgba(0, 0, 0, 0.1)",
    "variable": "var(--colors-black-a2)"
  },
  "colors.black.a3": {
    "value": "rgba(0, 0, 0, 0.15)",
    "variable": "var(--colors-black-a3)"
  },
  "colors.black.a4": {
    "value": "rgba(0, 0, 0, 0.2)",
    "variable": "var(--colors-black-a4)"
  },
  "colors.black.a5": {
    "value": "rgba(0, 0, 0, 0.3)",
    "variable": "var(--colors-black-a5)"
  },
  "colors.black.a6": {
    "value": "rgba(0, 0, 0, 0.4)",
    "variable": "var(--colors-black-a6)"
  },
  "colors.black.a7": {
    "value": "rgba(0, 0, 0, 0.5)",
    "variable": "var(--colors-black-a7)"
  },
  "colors.black.a8": {
    "value": "rgba(0, 0, 0, 0.6)",
    "variable": "var(--colors-black-a8)"
  },
  "colors.black.a9": {
    "value": "rgba(0, 0, 0, 0.7)",
    "variable": "var(--colors-black-a9)"
  },
  "colors.black.a10": {
    "value": "rgba(0, 0, 0, 0.8)",
    "variable": "var(--colors-black-a10)"
  },
  "colors.black.a11": {
    "value": "rgba(0, 0, 0, 0.9)",
    "variable": "var(--colors-black-a11)"
  },
  "colors.black.a12": {
    "value": "rgba(0, 0, 0, 0.95)",
    "variable": "var(--colors-black-a12)"
  },
  "colors.white": {
    "value": "#ffffff",
    "variable": "var(--colors-white)"
  },
  "colors.white.a1": {
    "value": "rgba(255, 255, 255, 0.05)",
    "variable": "var(--colors-white-a1)"
  },
  "colors.white.a2": {
    "value": "rgba(255, 255, 255, 0.1)",
    "variable": "var(--colors-white-a2)"
  },
  "colors.white.a3": {
    "value": "rgba(255, 255, 255, 0.15)",
    "variable": "var(--colors-white-a3)"
  },
  "colors.white.a4": {
    "value": "rgba(255, 255, 255, 0.2)",
    "variable": "var(--colors-white-a4)"
  },
  "colors.white.a5": {
    "value": "rgba(255, 255, 255, 0.3)",
    "variable": "var(--colors-white-a5)"
  },
  "colors.white.a6": {
    "value": "rgba(255, 255, 255, 0.4)",
    "variable": "var(--colors-white-a6)"
  },
  "colors.white.a7": {
    "value": "rgba(255, 255, 255, 0.5)",
    "variable": "var(--colors-white-a7)"
  },
  "colors.white.a8": {
    "value": "rgba(255, 255, 255, 0.6)",
    "variable": "var(--colors-white-a8)"
  },
  "colors.white.a9": {
    "value": "rgba(255, 255, 255, 0.7)",
    "variable": "var(--colors-white-a9)"
  },
  "colors.white.a10": {
    "value": "rgba(255, 255, 255, 0.8)",
    "variable": "var(--colors-white-a10)"
  },
  "colors.white.a11": {
    "value": "rgba(255, 255, 255, 0.9)",
    "variable": "var(--colors-white-a11)"
  },
  "colors.white.a12": {
    "value": "rgba(255, 255, 255, 0.95)",
    "variable": "var(--colors-white-a12)"
  },
  "colors.amber.light.1": {
    "value": "#fefdfb",
    "variable": "var(--colors-amber-light-1)"
  },
  "colors.amber.light.2": {
    "value": "#fefbe9",
    "variable": "var(--colors-amber-light-2)"
  },
  "colors.amber.light.3": {
    "value": "#fff7c2",
    "variable": "var(--colors-amber-light-3)"
  },
  "colors.amber.light.4": {
    "value": "#ffee9c",
    "variable": "var(--colors-amber-light-4)"
  },
  "colors.amber.light.5": {
    "value": "#fbe577",
    "variable": "var(--colors-amber-light-5)"
  },
  "colors.amber.light.6": {
    "value": "#f3d673",
    "variable": "var(--colors-amber-light-6)"
  },
  "colors.amber.light.7": {
    "value": "#e9c162",
    "variable": "var(--colors-amber-light-7)"
  },
  "colors.amber.light.8": {
    "value": "#e2a336",
    "variable": "var(--colors-amber-light-8)"
  },
  "colors.amber.light.9": {
    "value": "#ffc53d",
    "variable": "var(--colors-amber-light-9)"
  },
  "colors.amber.light.10": {
    "value": "#ffba18",
    "variable": "var(--colors-amber-light-10)"
  },
  "colors.amber.light.11": {
    "value": "#ab6400",
    "variable": "var(--colors-amber-light-11)"
  },
  "colors.amber.light.12": {
    "value": "#4f3422",
    "variable": "var(--colors-amber-light-12)"
  },
  "colors.amber.light.a1": {
    "value": "#c0800004",
    "variable": "var(--colors-amber-light-a1)"
  },
  "colors.amber.light.a2": {
    "value": "#f4d10016",
    "variable": "var(--colors-amber-light-a2)"
  },
  "colors.amber.light.a3": {
    "value": "#ffde003d",
    "variable": "var(--colors-amber-light-a3)"
  },
  "colors.amber.light.a4": {
    "value": "#ffd40063",
    "variable": "var(--colors-amber-light-a4)"
  },
  "colors.amber.light.a5": {
    "value": "#f8cf0088",
    "variable": "var(--colors-amber-light-a5)"
  },
  "colors.amber.light.a6": {
    "value": "#eab5008c",
    "variable": "var(--colors-amber-light-a6)"
  },
  "colors.amber.light.a7": {
    "value": "#dc9b009d",
    "variable": "var(--colors-amber-light-a7)"
  },
  "colors.amber.light.a8": {
    "value": "#da8a00c9",
    "variable": "var(--colors-amber-light-a8)"
  },
  "colors.amber.light.a9": {
    "value": "#ffb300c2",
    "variable": "var(--colors-amber-light-a9)"
  },
  "colors.amber.light.a10": {
    "value": "#ffb300e7",
    "variable": "var(--colors-amber-light-a10)"
  },
  "colors.amber.light.a11": {
    "value": "#ab6400",
    "variable": "var(--colors-amber-light-a11)"
  },
  "colors.amber.light.a12": {
    "value": "#341500dd",
    "variable": "var(--colors-amber-light-a12)"
  },
  "colors.amber.dark.1": {
    "value": "#16120c",
    "variable": "var(--colors-amber-dark-1)"
  },
  "colors.amber.dark.2": {
    "value": "#1d180f",
    "variable": "var(--colors-amber-dark-2)"
  },
  "colors.amber.dark.3": {
    "value": "#302008",
    "variable": "var(--colors-amber-dark-3)"
  },
  "colors.amber.dark.4": {
    "value": "#3f2700",
    "variable": "var(--colors-amber-dark-4)"
  },
  "colors.amber.dark.5": {
    "value": "#4d3000",
    "variable": "var(--colors-amber-dark-5)"
  },
  "colors.amber.dark.6": {
    "value": "#5c3d05",
    "variable": "var(--colors-amber-dark-6)"
  },
  "colors.amber.dark.7": {
    "value": "#714f19",
    "variable": "var(--colors-amber-dark-7)"
  },
  "colors.amber.dark.8": {
    "value": "#8f6424",
    "variable": "var(--colors-amber-dark-8)"
  },
  "colors.amber.dark.9": {
    "value": "#ffc53d",
    "variable": "var(--colors-amber-dark-9)"
  },
  "colors.amber.dark.10": {
    "value": "#ffd60a",
    "variable": "var(--colors-amber-dark-10)"
  },
  "colors.amber.dark.11": {
    "value": "#ffca16",
    "variable": "var(--colors-amber-dark-11)"
  },
  "colors.amber.dark.12": {
    "value": "#ffe7b3",
    "variable": "var(--colors-amber-dark-12)"
  },
  "colors.amber.dark.a1": {
    "value": "#e63c0006",
    "variable": "var(--colors-amber-dark-a1)"
  },
  "colors.amber.dark.a2": {
    "value": "#fd9b000d",
    "variable": "var(--colors-amber-dark-a2)"
  },
  "colors.amber.dark.a3": {
    "value": "#fa820022",
    "variable": "var(--colors-amber-dark-a3)"
  },
  "colors.amber.dark.a4": {
    "value": "#fc820032",
    "variable": "var(--colors-amber-dark-a4)"
  },
  "colors.amber.dark.a5": {
    "value": "#fd8b0041",
    "variable": "var(--colors-amber-dark-a5)"
  },
  "colors.amber.dark.a6": {
    "value": "#fd9b0051",
    "variable": "var(--colors-amber-dark-a6)"
  },
  "colors.amber.dark.a7": {
    "value": "#ffab2567",
    "variable": "var(--colors-amber-dark-a7)"
  },
  "colors.amber.dark.a8": {
    "value": "#ffae3587",
    "variable": "var(--colors-amber-dark-a8)"
  },
  "colors.amber.dark.a9": {
    "value": "#ffc53d",
    "variable": "var(--colors-amber-dark-a9)"
  },
  "colors.amber.dark.a10": {
    "value": "#ffd60a",
    "variable": "var(--colors-amber-dark-a10)"
  },
  "colors.amber.dark.a11": {
    "value": "#ffca16",
    "variable": "var(--colors-amber-dark-a11)"
  },
  "colors.amber.dark.a12": {
    "value": "#ffe7b3",
    "variable": "var(--colors-amber-dark-a12)"
  },
  "colors.blue.light.1": {
    "value": "#fbfdff",
    "variable": "var(--colors-blue-light-1)"
  },
  "colors.blue.light.2": {
    "value": "#f4faff",
    "variable": "var(--colors-blue-light-2)"
  },
  "colors.blue.light.3": {
    "value": "#e6f4fe",
    "variable": "var(--colors-blue-light-3)"
  },
  "colors.blue.light.4": {
    "value": "#d5efff",
    "variable": "var(--colors-blue-light-4)"
  },
  "colors.blue.light.5": {
    "value": "#c2e5ff",
    "variable": "var(--colors-blue-light-5)"
  },
  "colors.blue.light.6": {
    "value": "#acd8fc",
    "variable": "var(--colors-blue-light-6)"
  },
  "colors.blue.light.7": {
    "value": "#8ec8f6",
    "variable": "var(--colors-blue-light-7)"
  },
  "colors.blue.light.8": {
    "value": "#5eb1ef",
    "variable": "var(--colors-blue-light-8)"
  },
  "colors.blue.light.9": {
    "value": "#0090ff",
    "variable": "var(--colors-blue-light-9)"
  },
  "colors.blue.light.10": {
    "value": "#0588f0",
    "variable": "var(--colors-blue-light-10)"
  },
  "colors.blue.light.11": {
    "value": "#0d74ce",
    "variable": "var(--colors-blue-light-11)"
  },
  "colors.blue.light.12": {
    "value": "#113264",
    "variable": "var(--colors-blue-light-12)"
  },
  "colors.blue.light.a1": {
    "value": "#0080ff04",
    "variable": "var(--colors-blue-light-a1)"
  },
  "colors.blue.light.a2": {
    "value": "#008cff0b",
    "variable": "var(--colors-blue-light-a2)"
  },
  "colors.blue.light.a3": {
    "value": "#008ff519",
    "variable": "var(--colors-blue-light-a3)"
  },
  "colors.blue.light.a4": {
    "value": "#009eff2a",
    "variable": "var(--colors-blue-light-a4)"
  },
  "colors.blue.light.a5": {
    "value": "#0093ff3d",
    "variable": "var(--colors-blue-light-a5)"
  },
  "colors.blue.light.a6": {
    "value": "#0088f653",
    "variable": "var(--colors-blue-light-a6)"
  },
  "colors.blue.light.a7": {
    "value": "#0083eb71",
    "variable": "var(--colors-blue-light-a7)"
  },
  "colors.blue.light.a8": {
    "value": "#0084e6a1",
    "variable": "var(--colors-blue-light-a8)"
  },
  "colors.blue.light.a9": {
    "value": "#0090ff",
    "variable": "var(--colors-blue-light-a9)"
  },
  "colors.blue.light.a10": {
    "value": "#0086f0fa",
    "variable": "var(--colors-blue-light-a10)"
  },
  "colors.blue.light.a11": {
    "value": "#006dcbf2",
    "variable": "var(--colors-blue-light-a11)"
  },
  "colors.blue.light.a12": {
    "value": "#002359ee",
    "variable": "var(--colors-blue-light-a12)"
  },
  "colors.blue.dark.1": {
    "value": "#0d1520",
    "variable": "var(--colors-blue-dark-1)"
  },
  "colors.blue.dark.2": {
    "value": "#111927",
    "variable": "var(--colors-blue-dark-2)"
  },
  "colors.blue.dark.3": {
    "value": "#0d2847",
    "variable": "var(--colors-blue-dark-3)"
  },
  "colors.blue.dark.4": {
    "value": "#003362",
    "variable": "var(--colors-blue-dark-4)"
  },
  "colors.blue.dark.5": {
    "value": "#004074",
    "variable": "var(--colors-blue-dark-5)"
  },
  "colors.blue.dark.6": {
    "value": "#104d87",
    "variable": "var(--colors-blue-dark-6)"
  },
  "colors.blue.dark.7": {
    "value": "#205d9e",
    "variable": "var(--colors-blue-dark-7)"
  },
  "colors.blue.dark.8": {
    "value": "#2870bd",
    "variable": "var(--colors-blue-dark-8)"
  },
  "colors.blue.dark.9": {
    "value": "#0090ff",
    "variable": "var(--colors-blue-dark-9)"
  },
  "colors.blue.dark.10": {
    "value": "#3b9eff",
    "variable": "var(--colors-blue-dark-10)"
  },
  "colors.blue.dark.11": {
    "value": "#70b8ff",
    "variable": "var(--colors-blue-dark-11)"
  },
  "colors.blue.dark.12": {
    "value": "#c2e6ff",
    "variable": "var(--colors-blue-dark-12)"
  },
  "colors.blue.dark.a1": {
    "value": "#004df211",
    "variable": "var(--colors-blue-dark-a1)"
  },
  "colors.blue.dark.a2": {
    "value": "#1166fb18",
    "variable": "var(--colors-blue-dark-a2)"
  },
  "colors.blue.dark.a3": {
    "value": "#0077ff3a",
    "variable": "var(--colors-blue-dark-a3)"
  },
  "colors.blue.dark.a4": {
    "value": "#0075ff57",
    "variable": "var(--colors-blue-dark-a4)"
  },
  "colors.blue.dark.a5": {
    "value": "#0081fd6b",
    "variable": "var(--colors-blue-dark-a5)"
  },
  "colors.blue.dark.a6": {
    "value": "#0f89fd7f",
    "variable": "var(--colors-blue-dark-a6)"
  },
  "colors.blue.dark.a7": {
    "value": "#2a91fe98",
    "variable": "var(--colors-blue-dark-a7)"
  },
  "colors.blue.dark.a8": {
    "value": "#3094feb9",
    "variable": "var(--colors-blue-dark-a8)"
  },
  "colors.blue.dark.a9": {
    "value": "#0090ff",
    "variable": "var(--colors-blue-dark-a9)"
  },
  "colors.blue.dark.a10": {
    "value": "#3b9eff",
    "variable": "var(--colors-blue-dark-a10)"
  },
  "colors.blue.dark.a11": {
    "value": "#70b8ff",
    "variable": "var(--colors-blue-dark-a11)"
  },
  "colors.blue.dark.a12": {
    "value": "#c2e6ff",
    "variable": "var(--colors-blue-dark-a12)"
  },
  "colors.green.light.1": {
    "value": "#fbfefc",
    "variable": "var(--colors-green-light-1)"
  },
  "colors.green.light.2": {
    "value": "#f4fbf6",
    "variable": "var(--colors-green-light-2)"
  },
  "colors.green.light.3": {
    "value": "#e6f6eb",
    "variable": "var(--colors-green-light-3)"
  },
  "colors.green.light.4": {
    "value": "#d6f1df",
    "variable": "var(--colors-green-light-4)"
  },
  "colors.green.light.5": {
    "value": "#c4e8d1",
    "variable": "var(--colors-green-light-5)"
  },
  "colors.green.light.6": {
    "value": "#adddc0",
    "variable": "var(--colors-green-light-6)"
  },
  "colors.green.light.7": {
    "value": "#8eceaa",
    "variable": "var(--colors-green-light-7)"
  },
  "colors.green.light.8": {
    "value": "#5bb98b",
    "variable": "var(--colors-green-light-8)"
  },
  "colors.green.light.9": {
    "value": "#30a46c",
    "variable": "var(--colors-green-light-9)"
  },
  "colors.green.light.10": {
    "value": "#2b9a66",
    "variable": "var(--colors-green-light-10)"
  },
  "colors.green.light.11": {
    "value": "#218358",
    "variable": "var(--colors-green-light-11)"
  },
  "colors.green.light.12": {
    "value": "#193b2d",
    "variable": "var(--colors-green-light-12)"
  },
  "colors.green.light.a1": {
    "value": "#00c04004",
    "variable": "var(--colors-green-light-a1)"
  },
  "colors.green.light.a2": {
    "value": "#00a32f0b",
    "variable": "var(--colors-green-light-a2)"
  },
  "colors.green.light.a3": {
    "value": "#00a43319",
    "variable": "var(--colors-green-light-a3)"
  },
  "colors.green.light.a4": {
    "value": "#00a83829",
    "variable": "var(--colors-green-light-a4)"
  },
  "colors.green.light.a5": {
    "value": "#019c393b",
    "variable": "var(--colors-green-light-a5)"
  },
  "colors.green.light.a6": {
    "value": "#00963c52",
    "variable": "var(--colors-green-light-a6)"
  },
  "colors.green.light.a7": {
    "value": "#00914071",
    "variable": "var(--colors-green-light-a7)"
  },
  "colors.green.light.a8": {
    "value": "#00924ba4",
    "variable": "var(--colors-green-light-a8)"
  },
  "colors.green.light.a9": {
    "value": "#008f4acf",
    "variable": "var(--colors-green-light-a9)"
  },
  "colors.green.light.a10": {
    "value": "#008647d4",
    "variable": "var(--colors-green-light-a10)"
  },
  "colors.green.light.a11": {
    "value": "#00713fde",
    "variable": "var(--colors-green-light-a11)"
  },
  "colors.green.light.a12": {
    "value": "#002616e6",
    "variable": "var(--colors-green-light-a12)"
  },
  "colors.green.dark.1": {
    "value": "#0e1512",
    "variable": "var(--colors-green-dark-1)"
  },
  "colors.green.dark.2": {
    "value": "#121b17",
    "variable": "var(--colors-green-dark-2)"
  },
  "colors.green.dark.3": {
    "value": "#132d21",
    "variable": "var(--colors-green-dark-3)"
  },
  "colors.green.dark.4": {
    "value": "#113b29",
    "variable": "var(--colors-green-dark-4)"
  },
  "colors.green.dark.5": {
    "value": "#174933",
    "variable": "var(--colors-green-dark-5)"
  },
  "colors.green.dark.6": {
    "value": "#20573e",
    "variable": "var(--colors-green-dark-6)"
  },
  "colors.green.dark.7": {
    "value": "#28684a",
    "variable": "var(--colors-green-dark-7)"
  },
  "colors.green.dark.8": {
    "value": "#2f7c57",
    "variable": "var(--colors-green-dark-8)"
  },
  "colors.green.dark.9": {
    "value": "#30a46c",
    "variable": "var(--colors-green-dark-9)"
  },
  "colors.green.dark.10": {
    "value": "#33b074",
    "variable": "var(--colors-green-dark-10)"
  },
  "colors.green.dark.11": {
    "value": "#3dd68c",
    "variable": "var(--colors-green-dark-11)"
  },
  "colors.green.dark.12": {
    "value": "#b1f1cb",
    "variable": "var(--colors-green-dark-12)"
  },
  "colors.green.dark.a1": {
    "value": "#00de4505",
    "variable": "var(--colors-green-dark-a1)"
  },
  "colors.green.dark.a2": {
    "value": "#29f99d0b",
    "variable": "var(--colors-green-dark-a2)"
  },
  "colors.green.dark.a3": {
    "value": "#22ff991e",
    "variable": "var(--colors-green-dark-a3)"
  },
  "colors.green.dark.a4": {
    "value": "#11ff992d",
    "variable": "var(--colors-green-dark-a4)"
  },
  "colors.green.dark.a5": {
    "value": "#2bffa23c",
    "variable": "var(--colors-green-dark-a5)"
  },
  "colors.green.dark.a6": {
    "value": "#44ffaa4b",
    "variable": "var(--colors-green-dark-a6)"
  },
  "colors.green.dark.a7": {
    "value": "#50fdac5e",
    "variable": "var(--colors-green-dark-a7)"
  },
  "colors.green.dark.a8": {
    "value": "#54ffad73",
    "variable": "var(--colors-green-dark-a8)"
  },
  "colors.green.dark.a9": {
    "value": "#44ffa49e",
    "variable": "var(--colors-green-dark-a9)"
  },
  "colors.green.dark.a10": {
    "value": "#43fea4ab",
    "variable": "var(--colors-green-dark-a10)"
  },
  "colors.green.dark.a11": {
    "value": "#46fea5d4",
    "variable": "var(--colors-green-dark-a11)"
  },
  "colors.green.dark.a12": {
    "value": "#bbffd7f0",
    "variable": "var(--colors-green-dark-a12)"
  },
  "colors.orange.light.1": {
    "value": "#fefcfb",
    "variable": "var(--colors-orange-light-1)"
  },
  "colors.orange.light.2": {
    "value": "#fff7ed",
    "variable": "var(--colors-orange-light-2)"
  },
  "colors.orange.light.3": {
    "value": "#ffefd6",
    "variable": "var(--colors-orange-light-3)"
  },
  "colors.orange.light.4": {
    "value": "#ffdfb5",
    "variable": "var(--colors-orange-light-4)"
  },
  "colors.orange.light.5": {
    "value": "#ffd19a",
    "variable": "var(--colors-orange-light-5)"
  },
  "colors.orange.light.6": {
    "value": "#ffc182",
    "variable": "var(--colors-orange-light-6)"
  },
  "colors.orange.light.7": {
    "value": "#f5ae73",
    "variable": "var(--colors-orange-light-7)"
  },
  "colors.orange.light.8": {
    "value": "#ec9455",
    "variable": "var(--colors-orange-light-8)"
  },
  "colors.orange.light.9": {
    "value": "#f76b15",
    "variable": "var(--colors-orange-light-9)"
  },
  "colors.orange.light.10": {
    "value": "#ef5f00",
    "variable": "var(--colors-orange-light-10)"
  },
  "colors.orange.light.11": {
    "value": "#cc4e00",
    "variable": "var(--colors-orange-light-11)"
  },
  "colors.orange.light.12": {
    "value": "#582d1d",
    "variable": "var(--colors-orange-light-12)"
  },
  "colors.orange.light.a1": {
    "value": "#c0400004",
    "variable": "var(--colors-orange-light-a1)"
  },
  "colors.orange.light.a2": {
    "value": "#ff8e0012",
    "variable": "var(--colors-orange-light-a2)"
  },
  "colors.orange.light.a3": {
    "value": "#ff9c0029",
    "variable": "var(--colors-orange-light-a3)"
  },
  "colors.orange.light.a4": {
    "value": "#ff91014a",
    "variable": "var(--colors-orange-light-a4)"
  },
  "colors.orange.light.a5": {
    "value": "#ff8b0065",
    "variable": "var(--colors-orange-light-a5)"
  },
  "colors.orange.light.a6": {
    "value": "#ff81007d",
    "variable": "var(--colors-orange-light-a6)"
  },
  "colors.orange.light.a7": {
    "value": "#ed6c008c",
    "variable": "var(--colors-orange-light-a7)"
  },
  "colors.orange.light.a8": {
    "value": "#e35f00aa",
    "variable": "var(--colors-orange-light-a8)"
  },
  "colors.orange.light.a9": {
    "value": "#f65e00ea",
    "variable": "var(--colors-orange-light-a9)"
  },
  "colors.orange.light.a10": {
    "value": "#ef5f00",
    "variable": "var(--colors-orange-light-a10)"
  },
  "colors.orange.light.a11": {
    "value": "#cc4e00",
    "variable": "var(--colors-orange-light-a11)"
  },
  "colors.orange.light.a12": {
    "value": "#431200e2",
    "variable": "var(--colors-orange-light-a12)"
  },
  "colors.orange.dark.1": {
    "value": "#17120e",
    "variable": "var(--colors-orange-dark-1)"
  },
  "colors.orange.dark.2": {
    "value": "#1e160f",
    "variable": "var(--colors-orange-dark-2)"
  },
  "colors.orange.dark.3": {
    "value": "#331e0b",
    "variable": "var(--colors-orange-dark-3)"
  },
  "colors.orange.dark.4": {
    "value": "#462100",
    "variable": "var(--colors-orange-dark-4)"
  },
  "colors.orange.dark.5": {
    "value": "#562800",
    "variable": "var(--colors-orange-dark-5)"
  },
  "colors.orange.dark.6": {
    "value": "#66350c",
    "variable": "var(--colors-orange-dark-6)"
  },
  "colors.orange.dark.7": {
    "value": "#7e451d",
    "variable": "var(--colors-orange-dark-7)"
  },
  "colors.orange.dark.8": {
    "value": "#a35829",
    "variable": "var(--colors-orange-dark-8)"
  },
  "colors.orange.dark.9": {
    "value": "#f76b15",
    "variable": "var(--colors-orange-dark-9)"
  },
  "colors.orange.dark.10": {
    "value": "#ff801f",
    "variable": "var(--colors-orange-dark-10)"
  },
  "colors.orange.dark.11": {
    "value": "#ffa057",
    "variable": "var(--colors-orange-dark-11)"
  },
  "colors.orange.dark.12": {
    "value": "#ffe0c2",
    "variable": "var(--colors-orange-dark-12)"
  },
  "colors.orange.dark.a1": {
    "value": "#ec360007",
    "variable": "var(--colors-orange-dark-a1)"
  },
  "colors.orange.dark.a2": {
    "value": "#fe6d000e",
    "variable": "var(--colors-orange-dark-a2)"
  },
  "colors.orange.dark.a3": {
    "value": "#fb6a0025",
    "variable": "var(--colors-orange-dark-a3)"
  },
  "colors.orange.dark.a4": {
    "value": "#ff590039",
    "variable": "var(--colors-orange-dark-a4)"
  },
  "colors.orange.dark.a5": {
    "value": "#ff61004a",
    "variable": "var(--colors-orange-dark-a5)"
  },
  "colors.orange.dark.a6": {
    "value": "#fd75045c",
    "variable": "var(--colors-orange-dark-a6)"
  },
  "colors.orange.dark.a7": {
    "value": "#ff832c75",
    "variable": "var(--colors-orange-dark-a7)"
  },
  "colors.orange.dark.a8": {
    "value": "#fe84389d",
    "variable": "var(--colors-orange-dark-a8)"
  },
  "colors.orange.dark.a9": {
    "value": "#fe6d15f7",
    "variable": "var(--colors-orange-dark-a9)"
  },
  "colors.orange.dark.a10": {
    "value": "#ff801f",
    "variable": "var(--colors-orange-dark-a10)"
  },
  "colors.orange.dark.a11": {
    "value": "#ffa057",
    "variable": "var(--colors-orange-dark-a11)"
  },
  "colors.orange.dark.a12": {
    "value": "#ffe0c2",
    "variable": "var(--colors-orange-dark-a12)"
  },
  "colors.pink.light.1": {
    "value": "#fffcfe",
    "variable": "var(--colors-pink-light-1)"
  },
  "colors.pink.light.2": {
    "value": "#fef7fb",
    "variable": "var(--colors-pink-light-2)"
  },
  "colors.pink.light.3": {
    "value": "#fee9f5",
    "variable": "var(--colors-pink-light-3)"
  },
  "colors.pink.light.4": {
    "value": "#fbdcef",
    "variable": "var(--colors-pink-light-4)"
  },
  "colors.pink.light.5": {
    "value": "#f6cee7",
    "variable": "var(--colors-pink-light-5)"
  },
  "colors.pink.light.6": {
    "value": "#efbfdd",
    "variable": "var(--colors-pink-light-6)"
  },
  "colors.pink.light.7": {
    "value": "#e7acd0",
    "variable": "var(--colors-pink-light-7)"
  },
  "colors.pink.light.8": {
    "value": "#dd93c2",
    "variable": "var(--colors-pink-light-8)"
  },
  "colors.pink.light.9": {
    "value": "#d6409f",
    "variable": "var(--colors-pink-light-9)"
  },
  "colors.pink.light.10": {
    "value": "#cf3897",
    "variable": "var(--colors-pink-light-10)"
  },
  "colors.pink.light.11": {
    "value": "#c2298a",
    "variable": "var(--colors-pink-light-11)"
  },
  "colors.pink.light.12": {
    "value": "#651249",
    "variable": "var(--colors-pink-light-12)"
  },
  "colors.pink.light.a1": {
    "value": "#ff00aa03",
    "variable": "var(--colors-pink-light-a1)"
  },
  "colors.pink.light.a2": {
    "value": "#e0008008",
    "variable": "var(--colors-pink-light-a2)"
  },
  "colors.pink.light.a3": {
    "value": "#f4008c16",
    "variable": "var(--colors-pink-light-a3)"
  },
  "colors.pink.light.a4": {
    "value": "#e2008b23",
    "variable": "var(--colors-pink-light-a4)"
  },
  "colors.pink.light.a5": {
    "value": "#d1008331",
    "variable": "var(--colors-pink-light-a5)"
  },
  "colors.pink.light.a6": {
    "value": "#c0007840",
    "variable": "var(--colors-pink-light-a6)"
  },
  "colors.pink.light.a7": {
    "value": "#b6006f53",
    "variable": "var(--colors-pink-light-a7)"
  },
  "colors.pink.light.a8": {
    "value": "#af006f6c",
    "variable": "var(--colors-pink-light-a8)"
  },
  "colors.pink.light.a9": {
    "value": "#c8007fbf",
    "variable": "var(--colors-pink-light-a9)"
  },
  "colors.pink.light.a10": {
    "value": "#c2007ac7",
    "variable": "var(--colors-pink-light-a10)"
  },
  "colors.pink.light.a11": {
    "value": "#b60074d6",
    "variable": "var(--colors-pink-light-a11)"
  },
  "colors.pink.light.a12": {
    "value": "#59003bed",
    "variable": "var(--colors-pink-light-a12)"
  },
  "colors.pink.dark.1": {
    "value": "#191117",
    "variable": "var(--colors-pink-dark-1)"
  },
  "colors.pink.dark.2": {
    "value": "#21121d",
    "variable": "var(--colors-pink-dark-2)"
  },
  "colors.pink.dark.3": {
    "value": "#37172f",
    "variable": "var(--colors-pink-dark-3)"
  },
  "colors.pink.dark.4": {
    "value": "#4b143d",
    "variable": "var(--colors-pink-dark-4)"
  },
  "colors.pink.dark.5": {
    "value": "#591c47",
    "variable": "var(--colors-pink-dark-5)"
  },
  "colors.pink.dark.6": {
    "value": "#692955",
    "variable": "var(--colors-pink-dark-6)"
  },
  "colors.pink.dark.7": {
    "value": "#833869",
    "variable": "var(--colors-pink-dark-7)"
  },
  "colors.pink.dark.8": {
    "value": "#a84885",
    "variable": "var(--colors-pink-dark-8)"
  },
  "colors.pink.dark.9": {
    "value": "#d6409f",
    "variable": "var(--colors-pink-dark-9)"
  },
  "colors.pink.dark.10": {
    "value": "#de51a8",
    "variable": "var(--colors-pink-dark-10)"
  },
  "colors.pink.dark.11": {
    "value": "#ff8dcc",
    "variable": "var(--colors-pink-dark-11)"
  },
  "colors.pink.dark.12": {
    "value": "#fdd1ea",
    "variable": "var(--colors-pink-dark-12)"
  },
  "colors.pink.dark.a1": {
    "value": "#f412bc09",
    "variable": "var(--colors-pink-dark-a1)"
  },
  "colors.pink.dark.a2": {
    "value": "#f420bb12",
    "variable": "var(--colors-pink-dark-a2)"
  },
  "colors.pink.dark.a3": {
    "value": "#fe37cc29",
    "variable": "var(--colors-pink-dark-a3)"
  },
  "colors.pink.dark.a4": {
    "value": "#fc1ec43f",
    "variable": "var(--colors-pink-dark-a4)"
  },
  "colors.pink.dark.a5": {
    "value": "#fd35c24e",
    "variable": "var(--colors-pink-dark-a5)"
  },
  "colors.pink.dark.a6": {
    "value": "#fd51c75f",
    "variable": "var(--colors-pink-dark-a6)"
  },
  "colors.pink.dark.a7": {
    "value": "#fd62c87b",
    "variable": "var(--colors-pink-dark-a7)"
  },
  "colors.pink.dark.a8": {
    "value": "#ff68c8a2",
    "variable": "var(--colors-pink-dark-a8)"
  },
  "colors.pink.dark.a9": {
    "value": "#fe49bcd4",
    "variable": "var(--colors-pink-dark-a9)"
  },
  "colors.pink.dark.a10": {
    "value": "#ff5cc0dc",
    "variable": "var(--colors-pink-dark-a10)"
  },
  "colors.pink.dark.a11": {
    "value": "#ff8dcc",
    "variable": "var(--colors-pink-dark-a11)"
  },
  "colors.pink.dark.a12": {
    "value": "#ffd3ecfd",
    "variable": "var(--colors-pink-dark-a12)"
  },
  "colors.red.light.1": {
    "value": "#fffcfc",
    "variable": "var(--colors-red-light-1)"
  },
  "colors.red.light.2": {
    "value": "#fff7f7",
    "variable": "var(--colors-red-light-2)"
  },
  "colors.red.light.3": {
    "value": "#feebec",
    "variable": "var(--colors-red-light-3)"
  },
  "colors.red.light.4": {
    "value": "#ffdbdc",
    "variable": "var(--colors-red-light-4)"
  },
  "colors.red.light.5": {
    "value": "#ffcdce",
    "variable": "var(--colors-red-light-5)"
  },
  "colors.red.light.6": {
    "value": "#fdbdbe",
    "variable": "var(--colors-red-light-6)"
  },
  "colors.red.light.7": {
    "value": "#f4a9aa",
    "variable": "var(--colors-red-light-7)"
  },
  "colors.red.light.8": {
    "value": "#eb8e90",
    "variable": "var(--colors-red-light-8)"
  },
  "colors.red.light.9": {
    "value": "#e5484d",
    "variable": "var(--colors-red-light-9)"
  },
  "colors.red.light.10": {
    "value": "#dc3e42",
    "variable": "var(--colors-red-light-10)"
  },
  "colors.red.light.11": {
    "value": "#ce2c31",
    "variable": "var(--colors-red-light-11)"
  },
  "colors.red.light.12": {
    "value": "#641723",
    "variable": "var(--colors-red-light-12)"
  },
  "colors.red.light.a1": {
    "value": "#ff000003",
    "variable": "var(--colors-red-light-a1)"
  },
  "colors.red.light.a2": {
    "value": "#ff000008",
    "variable": "var(--colors-red-light-a2)"
  },
  "colors.red.light.a3": {
    "value": "#f3000d14",
    "variable": "var(--colors-red-light-a3)"
  },
  "colors.red.light.a4": {
    "value": "#ff000824",
    "variable": "var(--colors-red-light-a4)"
  },
  "colors.red.light.a5": {
    "value": "#ff000632",
    "variable": "var(--colors-red-light-a5)"
  },
  "colors.red.light.a6": {
    "value": "#f8000442",
    "variable": "var(--colors-red-light-a6)"
  },
  "colors.red.light.a7": {
    "value": "#df000356",
    "variable": "var(--colors-red-light-a7)"
  },
  "colors.red.light.a8": {
    "value": "#d2000571",
    "variable": "var(--colors-red-light-a8)"
  },
  "colors.red.light.a9": {
    "value": "#db0007b7",
    "variable": "var(--colors-red-light-a9)"
  },
  "colors.red.light.a10": {
    "value": "#d10005c1",
    "variable": "var(--colors-red-light-a10)"
  },
  "colors.red.light.a11": {
    "value": "#c40006d3",
    "variable": "var(--colors-red-light-a11)"
  },
  "colors.red.light.a12": {
    "value": "#55000de8",
    "variable": "var(--colors-red-light-a12)"
  },
  "colors.red.dark.1": {
    "value": "#191111",
    "variable": "var(--colors-red-dark-1)"
  },
  "colors.red.dark.2": {
    "value": "#201314",
    "variable": "var(--colors-red-dark-2)"
  },
  "colors.red.dark.3": {
    "value": "#3b1219",
    "variable": "var(--colors-red-dark-3)"
  },
  "colors.red.dark.4": {
    "value": "#500f1c",
    "variable": "var(--colors-red-dark-4)"
  },
  "colors.red.dark.5": {
    "value": "#611623",
    "variable": "var(--colors-red-dark-5)"
  },
  "colors.red.dark.6": {
    "value": "#72232d",
    "variable": "var(--colors-red-dark-6)"
  },
  "colors.red.dark.7": {
    "value": "#8c333a",
    "variable": "var(--colors-red-dark-7)"
  },
  "colors.red.dark.8": {
    "value": "#b54548",
    "variable": "var(--colors-red-dark-8)"
  },
  "colors.red.dark.9": {
    "value": "#e5484d",
    "variable": "var(--colors-red-dark-9)"
  },
  "colors.red.dark.10": {
    "value": "#ec5d5e",
    "variable": "var(--colors-red-dark-10)"
  },
  "colors.red.dark.11": {
    "value": "#ff9592",
    "variable": "var(--colors-red-dark-11)"
  },
  "colors.red.dark.12": {
    "value": "#ffd1d9",
    "variable": "var(--colors-red-dark-12)"
  },
  "colors.red.dark.a1": {
    "value": "#f4121209",
    "variable": "var(--colors-red-dark-a1)"
  },
  "colors.red.dark.a2": {
    "value": "#f22f3e11",
    "variable": "var(--colors-red-dark-a2)"
  },
  "colors.red.dark.a3": {
    "value": "#ff173f2d",
    "variable": "var(--colors-red-dark-a3)"
  },
  "colors.red.dark.a4": {
    "value": "#fe0a3b44",
    "variable": "var(--colors-red-dark-a4)"
  },
  "colors.red.dark.a5": {
    "value": "#ff204756",
    "variable": "var(--colors-red-dark-a5)"
  },
  "colors.red.dark.a6": {
    "value": "#ff3e5668",
    "variable": "var(--colors-red-dark-a6)"
  },
  "colors.red.dark.a7": {
    "value": "#ff536184",
    "variable": "var(--colors-red-dark-a7)"
  },
  "colors.red.dark.a8": {
    "value": "#ff5d61b0",
    "variable": "var(--colors-red-dark-a8)"
  },
  "colors.red.dark.a9": {
    "value": "#fe4e54e4",
    "variable": "var(--colors-red-dark-a9)"
  },
  "colors.red.dark.a10": {
    "value": "#ff6465eb",
    "variable": "var(--colors-red-dark-a10)"
  },
  "colors.red.dark.a11": {
    "value": "#ff9592",
    "variable": "var(--colors-red-dark-a11)"
  },
  "colors.red.dark.a12": {
    "value": "#ffd1d9",
    "variable": "var(--colors-red-dark-a12)"
  },
  "colors.slate.light.1": {
    "value": "hsl(240, 20%, 99%)",
    "variable": "var(--colors-slate-light-1)"
  },
  "colors.slate.light.2": {
    "value": "hsl(240, 20%, 98%)",
    "variable": "var(--colors-slate-light-2)"
  },
  "colors.slate.light.3": {
    "value": "hsl(240, 11.1%, 94.7%)",
    "variable": "var(--colors-slate-light-3)"
  },
  "colors.slate.light.4": {
    "value": "hsl(240, 9.5%, 91.8%)",
    "variable": "var(--colors-slate-light-4)"
  },
  "colors.slate.light.5": {
    "value": "hsl(230, 10.7%, 89%)",
    "variable": "var(--colors-slate-light-5)"
  },
  "colors.slate.light.6": {
    "value": "hsl(240, 10.1%, 86.5%)",
    "variable": "var(--colors-slate-light-6)"
  },
  "colors.slate.light.7": {
    "value": "hsl(233.3, 9.9%, 82.2%)",
    "variable": "var(--colors-slate-light-7)"
  },
  "colors.slate.light.8": {
    "value": "hsl(230.8, 10.2%, 75.1%)",
    "variable": "var(--colors-slate-light-8)"
  },
  "colors.slate.light.9": {
    "value": "hsl(230.8, 5.9%, 57.1%)",
    "variable": "var(--colors-slate-light-9)"
  },
  "colors.slate.light.10": {
    "value": "hsl(226.2, 5.4%, 52.7%)",
    "variable": "var(--colors-slate-light-10)"
  },
  "colors.slate.light.11": {
    "value": "hsl(220, 5.9%, 40%)",
    "variable": "var(--colors-slate-light-11)"
  },
  "colors.slate.light.12": {
    "value": "hsl(210, 12.5%, 12.5%)",
    "variable": "var(--colors-slate-light-12)"
  },
  "colors.slate.light.a1": {
    "value": "hsla(240, 100%, 16.7%, 0)",
    "variable": "var(--colors-slate-light-a1)"
  },
  "colors.slate.light.a2": {
    "value": "hsla(240, 100%, 16.7%, 0)",
    "variable": "var(--colors-slate-light-a2)"
  },
  "colors.slate.light.a3": {
    "value": "hsla(240, 100%, 10%, 0.1)",
    "variable": "var(--colors-slate-light-a3)"
  },
  "colors.slate.light.a4": {
    "value": "hsla(240, 100%, 8.8%, 0.1)",
    "variable": "var(--colors-slate-light-a4)"
  },
  "colors.slate.light.a5": {
    "value": "hsla(229.2, 100%, 9.8%, 0.1)",
    "variable": "var(--colors-slate-light-a5)"
  },
  "colors.slate.light.a6": {
    "value": "hsla(240, 100%, 9.2%, 0.1)",
    "variable": "var(--colors-slate-light-a6)"
  },
  "colors.slate.light.a7": {
    "value": "hsla(232.2, 100%, 9%, 0.2)",
    "variable": "var(--colors-slate-light-a7)"
  },
  "colors.slate.light.a8": {
    "value": "hsla(230, 100%, 9.4%, 0.3)",
    "variable": "var(--colors-slate-light-a8)"
  },
  "colors.slate.light.a9": {
    "value": "hsla(229.7, 100%, 5.7%, 0.5)",
    "variable": "var(--colors-slate-light-a9)"
  },
  "colors.slate.light.a10": {
    "value": "hsla(224.4, 100%, 5.3%, 0.5)",
    "variable": "var(--colors-slate-light-a10)"
  },
  "colors.slate.light.a11": {
    "value": "hsla(219, 100%, 3.9%, 0.6)",
    "variable": "var(--colors-slate-light-a11)"
  },
  "colors.slate.light.a12": {
    "value": "hsla(206.7, 100%, 1.8%, 0.9)",
    "variable": "var(--colors-slate-light-a12)"
  },
  "colors.slate.dark.1": {
    "value": "hsl(240, 5.6%, 7.1%)",
    "variable": "var(--colors-slate-dark-1)"
  },
  "colors.slate.dark.2": {
    "value": "hsl(220, 5.9%, 10%)",
    "variable": "var(--colors-slate-dark-2)"
  },
  "colors.slate.dark.3": {
    "value": "hsl(225, 5.7%, 13.7%)",
    "variable": "var(--colors-slate-dark-3)"
  },
  "colors.slate.dark.4": {
    "value": "hsl(210, 7.1%, 16.5%)",
    "variable": "var(--colors-slate-dark-4)"
  },
  "colors.slate.dark.5": {
    "value": "hsl(214.3, 7.1%, 19.4%)",
    "variable": "var(--colors-slate-dark-5)"
  },
  "colors.slate.dark.6": {
    "value": "hsl(213.3, 7.7%, 22.9%)",
    "variable": "var(--colors-slate-dark-6)"
  },
  "colors.slate.dark.7": {
    "value": "hsl(212.7, 7.6%, 28.4%)",
    "variable": "var(--colors-slate-dark-7)"
  },
  "colors.slate.dark.8": {
    "value": "hsl(212, 7.7%, 38.2%)",
    "variable": "var(--colors-slate-dark-8)"
  },
  "colors.slate.dark.9": {
    "value": "hsl(218.6, 6.3%, 43.9%)",
    "variable": "var(--colors-slate-dark-9)"
  },
  "colors.slate.dark.10": {
    "value": "hsl(221.5, 5.2%, 49.2%)",
    "variable": "var(--colors-slate-dark-10)"
  },
  "colors.slate.dark.11": {
    "value": "hsl(216, 6.8%, 71%)",
    "variable": "var(--colors-slate-dark-11)"
  },
  "colors.slate.dark.12": {
    "value": "hsl(220, 9.1%, 93.5%)",
    "variable": "var(--colors-slate-dark-12)"
  },
  "colors.slate.dark.a1": {
    "value": "hsla(220, 5.9%, 10%, 0.02)",
    "variable": "var(--colors-slate-dark-a1)"
  },
  "colors.slate.dark.a2": {
    "value": "hsla(220, 5.9%, 10%, 0.04)",
    "variable": "var(--colors-slate-dark-a2)"
  },
  "colors.slate.dark.a3": {
    "value": "hsla(220, 5.9%, 10%, 0.08)",
    "variable": "var(--colors-slate-dark-a3)"
  },
  "colors.slate.dark.a4": {
    "value": "hsla(220, 5.9%, 10%, 0.12)",
    "variable": "var(--colors-slate-dark-a4)"
  },
  "colors.slate.dark.a5": {
    "value": "hsla(220, 5.9%, 10%, 0.16)",
    "variable": "var(--colors-slate-dark-a5)"
  },
  "colors.slate.dark.a6": {
    "value": "hsla(220, 5.9%, 10%, 0.24)",
    "variable": "var(--colors-slate-dark-a6)"
  },
  "colors.slate.dark.a7": {
    "value": "hsla(220, 5.9%, 10%, 0.32)",
    "variable": "var(--colors-slate-dark-a7)"
  },
  "colors.slate.dark.a8": {
    "value": "hsla(220, 5.9%, 10%, 0.42)",
    "variable": "var(--colors-slate-dark-a8)"
  },
  "colors.slate.dark.a9": {
    "value": "hsla(220, 5.9%, 10%, 0.52)",
    "variable": "var(--colors-slate-dark-a9)"
  },
  "colors.slate.dark.a10": {
    "value": "hsla(220, 5.9%, 10%, 0.62)",
    "variable": "var(--colors-slate-dark-a10)"
  },
  "colors.slate.dark.a11": {
    "value": "hsla(220, 5.9%, 10%, 0.7)",
    "variable": "var(--colors-slate-dark-a11)"
  },
  "colors.slate.dark.a12": {
    "value": "hsla(220, 5.9%, 10%, 0.9)",
    "variable": "var(--colors-slate-dark-a12)"
  },
  "colors.tomato.light.1": {
    "value": "#fffcfc",
    "variable": "var(--colors-tomato-light-1)"
  },
  "colors.tomato.light.2": {
    "value": "#fff8f7",
    "variable": "var(--colors-tomato-light-2)"
  },
  "colors.tomato.light.3": {
    "value": "#feebe7",
    "variable": "var(--colors-tomato-light-3)"
  },
  "colors.tomato.light.4": {
    "value": "#ffdcd3",
    "variable": "var(--colors-tomato-light-4)"
  },
  "colors.tomato.light.5": {
    "value": "#ffcdc2",
    "variable": "var(--colors-tomato-light-5)"
  },
  "colors.tomato.light.6": {
    "value": "#fdbdaf",
    "variable": "var(--colors-tomato-light-6)"
  },
  "colors.tomato.light.7": {
    "value": "#f5a898",
    "variable": "var(--colors-tomato-light-7)"
  },
  "colors.tomato.light.8": {
    "value": "#ec8e7b",
    "variable": "var(--colors-tomato-light-8)"
  },
  "colors.tomato.light.9": {
    "value": "#e54d2e",
    "variable": "var(--colors-tomato-light-9)"
  },
  "colors.tomato.light.10": {
    "value": "#dd4425",
    "variable": "var(--colors-tomato-light-10)"
  },
  "colors.tomato.light.11": {
    "value": "#d13415",
    "variable": "var(--colors-tomato-light-11)"
  },
  "colors.tomato.light.12": {
    "value": "#5c271f",
    "variable": "var(--colors-tomato-light-12)"
  },
  "colors.tomato.light.a1": {
    "value": "#ff000003",
    "variable": "var(--colors-tomato-light-a1)"
  },
  "colors.tomato.light.a2": {
    "value": "#ff200008",
    "variable": "var(--colors-tomato-light-a2)"
  },
  "colors.tomato.light.a3": {
    "value": "#f52b0018",
    "variable": "var(--colors-tomato-light-a3)"
  },
  "colors.tomato.light.a4": {
    "value": "#ff35002c",
    "variable": "var(--colors-tomato-light-a4)"
  },
  "colors.tomato.light.a5": {
    "value": "#ff2e003d",
    "variable": "var(--colors-tomato-light-a5)"
  },
  "colors.tomato.light.a6": {
    "value": "#f92d0050",
    "variable": "var(--colors-tomato-light-a6)"
  },
  "colors.tomato.light.a7": {
    "value": "#e7280067",
    "variable": "var(--colors-tomato-light-a7)"
  },
  "colors.tomato.light.a8": {
    "value": "#db250084",
    "variable": "var(--colors-tomato-light-a8)"
  },
  "colors.tomato.light.a9": {
    "value": "#df2600d1",
    "variable": "var(--colors-tomato-light-a9)"
  },
  "colors.tomato.light.a10": {
    "value": "#d72400da",
    "variable": "var(--colors-tomato-light-a10)"
  },
  "colors.tomato.light.a11": {
    "value": "#cd2200ea",
    "variable": "var(--colors-tomato-light-a11)"
  },
  "colors.tomato.light.a12": {
    "value": "#460900e0",
    "variable": "var(--colors-tomato-light-a12)"
  },
  "colors.tomato.dark.1": {
    "value": "#181111",
    "variable": "var(--colors-tomato-dark-1)"
  },
  "colors.tomato.dark.2": {
    "value": "#1f1513",
    "variable": "var(--colors-tomato-dark-2)"
  },
  "colors.tomato.dark.3": {
    "value": "#391714",
    "variable": "var(--colors-tomato-dark-3)"
  },
  "colors.tomato.dark.4": {
    "value": "#4e1511",
    "variable": "var(--colors-tomato-dark-4)"
  },
  "colors.tomato.dark.5": {
    "value": "#5e1c16",
    "variable": "var(--colors-tomato-dark-5)"
  },
  "colors.tomato.dark.6": {
    "value": "#6e2920",
    "variable": "var(--colors-tomato-dark-6)"
  },
  "colors.tomato.dark.7": {
    "value": "#853a2d",
    "variable": "var(--colors-tomato-dark-7)"
  },
  "colors.tomato.dark.8": {
    "value": "#ac4d39",
    "variable": "var(--colors-tomato-dark-8)"
  },
  "colors.tomato.dark.9": {
    "value": "#e54d2e",
    "variable": "var(--colors-tomato-dark-9)"
  },
  "colors.tomato.dark.10": {
    "value": "#ec6142",
    "variable": "var(--colors-tomato-dark-10)"
  },
  "colors.tomato.dark.11": {
    "value": "#ff977d",
    "variable": "var(--colors-tomato-dark-11)"
  },
  "colors.tomato.dark.12": {
    "value": "#fbd3cb",
    "variable": "var(--colors-tomato-dark-12)"
  },
  "colors.tomato.dark.a1": {
    "value": "#f1121208",
    "variable": "var(--colors-tomato-dark-a1)"
  },
  "colors.tomato.dark.a2": {
    "value": "#ff55330f",
    "variable": "var(--colors-tomato-dark-a2)"
  },
  "colors.tomato.dark.a3": {
    "value": "#ff35232b",
    "variable": "var(--colors-tomato-dark-a3)"
  },
  "colors.tomato.dark.a4": {
    "value": "#fd201142",
    "variable": "var(--colors-tomato-dark-a4)"
  },
  "colors.tomato.dark.a5": {
    "value": "#fe332153",
    "variable": "var(--colors-tomato-dark-a5)"
  },
  "colors.tomato.dark.a6": {
    "value": "#ff4f3864",
    "variable": "var(--colors-tomato-dark-a6)"
  },
  "colors.tomato.dark.a7": {
    "value": "#fd644a7d",
    "variable": "var(--colors-tomato-dark-a7)"
  },
  "colors.tomato.dark.a8": {
    "value": "#fe6d4ea7",
    "variable": "var(--colors-tomato-dark-a8)"
  },
  "colors.tomato.dark.a9": {
    "value": "#fe5431e4",
    "variable": "var(--colors-tomato-dark-a9)"
  },
  "colors.tomato.dark.a10": {
    "value": "#ff6847eb",
    "variable": "var(--colors-tomato-dark-a10)"
  },
  "colors.tomato.dark.a11": {
    "value": "#ff977d",
    "variable": "var(--colors-tomato-dark-a11)"
  },
  "colors.tomato.dark.a12": {
    "value": "#ffd6cefb",
    "variable": "var(--colors-tomato-dark-a12)"
  },
  "colors.gray.light.1": {
    "value": "#fcfcfc",
    "variable": "var(--colors-gray-light-1)"
  },
  "colors.gray.light.2": {
    "value": "#f9f9f9",
    "variable": "var(--colors-gray-light-2)"
  },
  "colors.gray.light.3": {
    "value": "#f0f0f0",
    "variable": "var(--colors-gray-light-3)"
  },
  "colors.gray.light.4": {
    "value": "#e8e8e8",
    "variable": "var(--colors-gray-light-4)"
  },
  "colors.gray.light.5": {
    "value": "#e0e0e0",
    "variable": "var(--colors-gray-light-5)"
  },
  "colors.gray.light.6": {
    "value": "#d9d9d9",
    "variable": "var(--colors-gray-light-6)"
  },
  "colors.gray.light.7": {
    "value": "#cecece",
    "variable": "var(--colors-gray-light-7)"
  },
  "colors.gray.light.8": {
    "value": "#bbbbbb",
    "variable": "var(--colors-gray-light-8)"
  },
  "colors.gray.light.9": {
    "value": "#8d8d8d",
    "variable": "var(--colors-gray-light-9)"
  },
  "colors.gray.light.10": {
    "value": "#838383",
    "variable": "var(--colors-gray-light-10)"
  },
  "colors.gray.light.11": {
    "value": "#646464",
    "variable": "var(--colors-gray-light-11)"
  },
  "colors.gray.light.12": {
    "value": "#202020",
    "variable": "var(--colors-gray-light-12)"
  },
  "colors.gray.light.a1": {
    "value": "#00000003",
    "variable": "var(--colors-gray-light-a1)"
  },
  "colors.gray.light.a2": {
    "value": "#00000006",
    "variable": "var(--colors-gray-light-a2)"
  },
  "colors.gray.light.a3": {
    "value": "#0000000f",
    "variable": "var(--colors-gray-light-a3)"
  },
  "colors.gray.light.a4": {
    "value": "#00000017",
    "variable": "var(--colors-gray-light-a4)"
  },
  "colors.gray.light.a5": {
    "value": "#0000001f",
    "variable": "var(--colors-gray-light-a5)"
  },
  "colors.gray.light.a6": {
    "value": "#00000026",
    "variable": "var(--colors-gray-light-a6)"
  },
  "colors.gray.light.a7": {
    "value": "#00000031",
    "variable": "var(--colors-gray-light-a7)"
  },
  "colors.gray.light.a8": {
    "value": "#00000044",
    "variable": "var(--colors-gray-light-a8)"
  },
  "colors.gray.light.a9": {
    "value": "#00000072",
    "variable": "var(--colors-gray-light-a9)"
  },
  "colors.gray.light.a10": {
    "value": "#0000007c",
    "variable": "var(--colors-gray-light-a10)"
  },
  "colors.gray.light.a11": {
    "value": "#0000009b",
    "variable": "var(--colors-gray-light-a11)"
  },
  "colors.gray.light.a12": {
    "value": "#000000df",
    "variable": "var(--colors-gray-light-a12)"
  },
  "colors.gray.dark.1": {
    "value": "#111111",
    "variable": "var(--colors-gray-dark-1)"
  },
  "colors.gray.dark.2": {
    "value": "#191919",
    "variable": "var(--colors-gray-dark-2)"
  },
  "colors.gray.dark.3": {
    "value": "#222222",
    "variable": "var(--colors-gray-dark-3)"
  },
  "colors.gray.dark.4": {
    "value": "#2a2a2a",
    "variable": "var(--colors-gray-dark-4)"
  },
  "colors.gray.dark.5": {
    "value": "#313131",
    "variable": "var(--colors-gray-dark-5)"
  },
  "colors.gray.dark.6": {
    "value": "#3a3a3a",
    "variable": "var(--colors-gray-dark-6)"
  },
  "colors.gray.dark.7": {
    "value": "#484848",
    "variable": "var(--colors-gray-dark-7)"
  },
  "colors.gray.dark.8": {
    "value": "#606060",
    "variable": "var(--colors-gray-dark-8)"
  },
  "colors.gray.dark.9": {
    "value": "#6e6e6e",
    "variable": "var(--colors-gray-dark-9)"
  },
  "colors.gray.dark.10": {
    "value": "#7b7b7b",
    "variable": "var(--colors-gray-dark-10)"
  },
  "colors.gray.dark.11": {
    "value": "#b4b4b4",
    "variable": "var(--colors-gray-dark-11)"
  },
  "colors.gray.dark.12": {
    "value": "#eeeeee",
    "variable": "var(--colors-gray-dark-12)"
  },
  "colors.gray.dark.a1": {
    "value": "#00000000",
    "variable": "var(--colors-gray-dark-a1)"
  },
  "colors.gray.dark.a2": {
    "value": "#ffffff09",
    "variable": "var(--colors-gray-dark-a2)"
  },
  "colors.gray.dark.a3": {
    "value": "#ffffff12",
    "variable": "var(--colors-gray-dark-a3)"
  },
  "colors.gray.dark.a4": {
    "value": "#ffffff1b",
    "variable": "var(--colors-gray-dark-a4)"
  },
  "colors.gray.dark.a5": {
    "value": "#ffffff22",
    "variable": "var(--colors-gray-dark-a5)"
  },
  "colors.gray.dark.a6": {
    "value": "#ffffff2c",
    "variable": "var(--colors-gray-dark-a6)"
  },
  "colors.gray.dark.a7": {
    "value": "#ffffff3b",
    "variable": "var(--colors-gray-dark-a7)"
  },
  "colors.gray.dark.a8": {
    "value": "#ffffff55",
    "variable": "var(--colors-gray-dark-a8)"
  },
  "colors.gray.dark.a9": {
    "value": "#ffffff64",
    "variable": "var(--colors-gray-dark-a9)"
  },
  "colors.gray.dark.a10": {
    "value": "#ffffff72",
    "variable": "var(--colors-gray-dark-a10)"
  },
  "colors.gray.dark.a11": {
    "value": "#ffffffaf",
    "variable": "var(--colors-gray-dark-a11)"
  },
  "colors.gray.dark.a12": {
    "value": "#ffffffed",
    "variable": "var(--colors-gray-dark-a12)"
  },
  "colors.neutral.light.1": {
    "value": "#fcfcfc",
    "variable": "var(--colors-neutral-light-1)"
  },
  "colors.neutral.light.2": {
    "value": "#f9f9f9",
    "variable": "var(--colors-neutral-light-2)"
  },
  "colors.neutral.light.3": {
    "value": "#f0f0f0",
    "variable": "var(--colors-neutral-light-3)"
  },
  "colors.neutral.light.4": {
    "value": "#e8e8e8",
    "variable": "var(--colors-neutral-light-4)"
  },
  "colors.neutral.light.5": {
    "value": "#e0e0e0",
    "variable": "var(--colors-neutral-light-5)"
  },
  "colors.neutral.light.6": {
    "value": "#d9d9d9",
    "variable": "var(--colors-neutral-light-6)"
  },
  "colors.neutral.light.7": {
    "value": "#cecece",
    "variable": "var(--colors-neutral-light-7)"
  },
  "colors.neutral.light.8": {
    "value": "#bbbbbb",
    "variable": "var(--colors-neutral-light-8)"
  },
  "colors.neutral.light.9": {
    "value": "#8d8d8d",
    "variable": "var(--colors-neutral-light-9)"
  },
  "colors.neutral.light.10": {
    "value": "#838383",
    "variable": "var(--colors-neutral-light-10)"
  },
  "colors.neutral.light.11": {
    "value": "#646464",
    "variable": "var(--colors-neutral-light-11)"
  },
  "colors.neutral.light.12": {
    "value": "#202020",
    "variable": "var(--colors-neutral-light-12)"
  },
  "colors.neutral.light.a1": {
    "value": "#00000003",
    "variable": "var(--colors-neutral-light-a1)"
  },
  "colors.neutral.light.a2": {
    "value": "#00000006",
    "variable": "var(--colors-neutral-light-a2)"
  },
  "colors.neutral.light.a3": {
    "value": "#0000000f",
    "variable": "var(--colors-neutral-light-a3)"
  },
  "colors.neutral.light.a4": {
    "value": "#00000017",
    "variable": "var(--colors-neutral-light-a4)"
  },
  "colors.neutral.light.a5": {
    "value": "#0000001f",
    "variable": "var(--colors-neutral-light-a5)"
  },
  "colors.neutral.light.a6": {
    "value": "#00000026",
    "variable": "var(--colors-neutral-light-a6)"
  },
  "colors.neutral.light.a7": {
    "value": "#00000031",
    "variable": "var(--colors-neutral-light-a7)"
  },
  "colors.neutral.light.a8": {
    "value": "#00000044",
    "variable": "var(--colors-neutral-light-a8)"
  },
  "colors.neutral.light.a9": {
    "value": "#00000072",
    "variable": "var(--colors-neutral-light-a9)"
  },
  "colors.neutral.light.a10": {
    "value": "#0000007c",
    "variable": "var(--colors-neutral-light-a10)"
  },
  "colors.neutral.light.a11": {
    "value": "#0000009b",
    "variable": "var(--colors-neutral-light-a11)"
  },
  "colors.neutral.light.a12": {
    "value": "#000000df",
    "variable": "var(--colors-neutral-light-a12)"
  },
  "colors.neutral.dark.1": {
    "value": "#111111",
    "variable": "var(--colors-neutral-dark-1)"
  },
  "colors.neutral.dark.2": {
    "value": "#191919",
    "variable": "var(--colors-neutral-dark-2)"
  },
  "colors.neutral.dark.3": {
    "value": "#222222",
    "variable": "var(--colors-neutral-dark-3)"
  },
  "colors.neutral.dark.4": {
    "value": "#2a2a2a",
    "variable": "var(--colors-neutral-dark-4)"
  },
  "colors.neutral.dark.5": {
    "value": "#313131",
    "variable": "var(--colors-neutral-dark-5)"
  },
  "colors.neutral.dark.6": {
    "value": "#3a3a3a",
    "variable": "var(--colors-neutral-dark-6)"
  },
  "colors.neutral.dark.7": {
    "value": "#484848",
    "variable": "var(--colors-neutral-dark-7)"
  },
  "colors.neutral.dark.8": {
    "value": "#606060",
    "variable": "var(--colors-neutral-dark-8)"
  },
  "colors.neutral.dark.9": {
    "value": "#6e6e6e",
    "variable": "var(--colors-neutral-dark-9)"
  },
  "colors.neutral.dark.10": {
    "value": "#7b7b7b",
    "variable": "var(--colors-neutral-dark-10)"
  },
  "colors.neutral.dark.11": {
    "value": "#b4b4b4",
    "variable": "var(--colors-neutral-dark-11)"
  },
  "colors.neutral.dark.12": {
    "value": "#eeeeee",
    "variable": "var(--colors-neutral-dark-12)"
  },
  "colors.neutral.dark.a1": {
    "value": "#00000000",
    "variable": "var(--colors-neutral-dark-a1)"
  },
  "colors.neutral.dark.a2": {
    "value": "#ffffff09",
    "variable": "var(--colors-neutral-dark-a2)"
  },
  "colors.neutral.dark.a3": {
    "value": "#ffffff12",
    "variable": "var(--colors-neutral-dark-a3)"
  },
  "colors.neutral.dark.a4": {
    "value": "#ffffff1b",
    "variable": "var(--colors-neutral-dark-a4)"
  },
  "colors.neutral.dark.a5": {
    "value": "#ffffff22",
    "variable": "var(--colors-neutral-dark-a5)"
  },
  "colors.neutral.dark.a6": {
    "value": "#ffffff2c",
    "variable": "var(--colors-neutral-dark-a6)"
  },
  "colors.neutral.dark.a7": {
    "value": "#ffffff3b",
    "variable": "var(--colors-neutral-dark-a7)"
  },
  "colors.neutral.dark.a8": {
    "value": "#ffffff55",
    "variable": "var(--colors-neutral-dark-a8)"
  },
  "colors.neutral.dark.a9": {
    "value": "#ffffff64",
    "variable": "var(--colors-neutral-dark-a9)"
  },
  "colors.neutral.dark.a10": {
    "value": "#ffffff72",
    "variable": "var(--colors-neutral-dark-a10)"
  },
  "colors.neutral.dark.a11": {
    "value": "#ffffffaf",
    "variable": "var(--colors-neutral-dark-a11)"
  },
  "colors.neutral.dark.a12": {
    "value": "#ffffffed",
    "variable": "var(--colors-neutral-dark-a12)"
  },
  "colors.transparent": {
    "value": "rgb(0 0 0 / 0)",
    "variable": "var(--colors-transparent)"
  },
  "colors.accent.light.1": {
    "value": "var(--accent-colour-flat-fill-50)",
    "variable": "var(--colors-accent-light-1)"
  },
  "colors.accent.light.2": {
    "value": "var(--accent-colour-flat-fill-100)",
    "variable": "var(--colors-accent-light-2)"
  },
  "colors.accent.light.3": {
    "value": "var(--accent-colour-flat-fill-200)",
    "variable": "var(--colors-accent-light-3)"
  },
  "colors.accent.light.4": {
    "value": "var(--accent-colour-flat-fill-300)",
    "variable": "var(--colors-accent-light-4)"
  },
  "colors.accent.light.5": {
    "value": "var(--accent-colour-flat-fill-400)",
    "variable": "var(--colors-accent-light-5)"
  },
  "colors.accent.light.6": {
    "value": "var(--accent-colour-flat-fill-500)",
    "variable": "var(--colors-accent-light-6)"
  },
  "colors.accent.light.7": {
    "value": "var(--accent-colour-flat-fill-600)",
    "variable": "var(--colors-accent-light-7)"
  },
  "colors.accent.light.8": {
    "value": "var(--accent-colour-flat-fill-700)",
    "variable": "var(--colors-accent-light-8)"
  },
  "colors.accent.light.9": {
    "value": "var(--accent-colour-flat-fill-800)",
    "variable": "var(--colors-accent-light-9)"
  },
  "colors.accent.light.10": {
    "value": "var(--accent-colour-flat-fill-900)",
    "variable": "var(--colors-accent-light-10)"
  },
  "colors.accent.light.text.1": {
    "value": "var(--accent-colour-flat-text-50)",
    "variable": "var(--colors-accent-light-text-1)"
  },
  "colors.accent.light.text.2": {
    "value": "var(--accent-colour-flat-text-100)",
    "variable": "var(--colors-accent-light-text-2)"
  },
  "colors.accent.light.text.3": {
    "value": "var(--accent-colour-flat-text-200)",
    "variable": "var(--colors-accent-light-text-3)"
  },
  "colors.accent.light.text.4": {
    "value": "var(--accent-colour-flat-text-300)",
    "variable": "var(--colors-accent-light-text-4)"
  },
  "colors.accent.light.text.5": {
    "value": "var(--accent-colour-flat-text-400)",
    "variable": "var(--colors-accent-light-text-5)"
  },
  "colors.accent.light.text.6": {
    "value": "var(--accent-colour-flat-text-500)",
    "variable": "var(--colors-accent-light-text-6)"
  },
  "colors.accent.light.text.7": {
    "value": "var(--accent-colour-flat-text-600)",
    "variable": "var(--colors-accent-light-text-7)"
  },
  "colors.accent.light.text.8": {
    "value": "var(--accent-colour-flat-text-700)",
    "variable": "var(--colors-accent-light-text-8)"
  },
  "colors.accent.light.text.9": {
    "value": "var(--accent-colour-flat-text-800)",
    "variable": "var(--colors-accent-light-text-9)"
  },
  "colors.accent.light.text.10": {
    "value": "var(--accent-colour-flat-text-900)",
    "variable": "var(--colors-accent-light-text-10)"
  },
  "colors.accent.dark.1": {
    "value": "var(--accent-colour-dark-fill-50)",
    "variable": "var(--colors-accent-dark-1)"
  },
  "colors.accent.dark.2": {
    "value": "var(--accent-colour-dark-fill-100)",
    "variable": "var(--colors-accent-dark-2)"
  },
  "colors.accent.dark.3": {
    "value": "var(--accent-colour-dark-fill-200)",
    "variable": "var(--colors-accent-dark-3)"
  },
  "colors.accent.dark.4": {
    "value": "var(--accent-colour-dark-fill-300)",
    "variable": "var(--colors-accent-dark-4)"
  },
  "colors.accent.dark.5": {
    "value": "var(--accent-colour-dark-fill-400)",
    "variable": "var(--colors-accent-dark-5)"
  },
  "colors.accent.dark.6": {
    "value": "var(--accent-colour-dark-fill-500)",
    "variable": "var(--colors-accent-dark-6)"
  },
  "colors.accent.dark.7": {
    "value": "var(--accent-colour-dark-fill-600)",
    "variable": "var(--colors-accent-dark-7)"
  },
  "colors.accent.dark.8": {
    "value": "var(--accent-colour-dark-fill-700)",
    "variable": "var(--colors-accent-dark-8)"
  },
  "colors.accent.dark.9": {
    "value": "var(--accent-colour-dark-fill-800)",
    "variable": "var(--colors-accent-dark-9)"
  },
  "colors.accent.dark.10": {
    "value": "var(--accent-colour-dark-fill-900)",
    "variable": "var(--colors-accent-dark-10)"
  },
  "colors.accent.dark.text.1": {
    "value": "var(--accent-colour-dark-text-50)",
    "variable": "var(--colors-accent-dark-text-1)"
  },
  "colors.accent.dark.text.2": {
    "value": "var(--accent-colour-dark-text-100)",
    "variable": "var(--colors-accent-dark-text-2)"
  },
  "colors.accent.dark.text.3": {
    "value": "var(--accent-colour-dark-text-200)",
    "variable": "var(--colors-accent-dark-text-3)"
  },
  "colors.accent.dark.text.4": {
    "value": "var(--accent-colour-dark-text-300)",
    "variable": "var(--colors-accent-dark-text-4)"
  },
  "colors.accent.dark.text.5": {
    "value": "var(--accent-colour-dark-text-400)",
    "variable": "var(--colors-accent-dark-text-5)"
  },
  "colors.accent.dark.text.6": {
    "value": "var(--accent-colour-dark-text-500)",
    "variable": "var(--colors-accent-dark-text-6)"
  },
  "colors.accent.dark.text.7": {
    "value": "var(--accent-colour-dark-text-600)",
    "variable": "var(--colors-accent-dark-text-7)"
  },
  "colors.accent.dark.text.8": {
    "value": "var(--accent-colour-dark-text-700)",
    "variable": "var(--colors-accent-dark-text-8)"
  },
  "colors.accent.dark.text.9": {
    "value": "var(--accent-colour-dark-text-800)",
    "variable": "var(--colors-accent-dark-text-9)"
  },
  "colors.accent.dark.text.10": {
    "value": "var(--accent-colour-dark-text-900)",
    "variable": "var(--colors-accent-dark-text-10)"
  },
  "durations.fastest": {
    "value": "50ms",
    "variable": "var(--durations-fastest)"
  },
  "durations.faster": {
    "value": "100ms",
    "variable": "var(--durations-faster)"
  },
  "durations.fast": {
    "value": "150ms",
    "variable": "var(--durations-fast)"
  },
  "durations.normal": {
    "value": "200ms",
    "variable": "var(--durations-normal)"
  },
  "durations.slow": {
    "value": "300ms",
    "variable": "var(--durations-slow)"
  },
  "durations.slower": {
    "value": "400ms",
    "variable": "var(--durations-slower)"
  },
  "durations.slowest": {
    "value": "500ms",
    "variable": "var(--durations-slowest)"
  },
  "easings.pulse": {
    "value": "cubic-bezier(0.4, 0.0, 0.6, 1.0)",
    "variable": "var(--easings-pulse)"
  },
  "easings.default": {
    "value": "cubic-bezier(0.2, 0.0, 0, 1.0)",
    "variable": "var(--easings-default)"
  },
  "easings.emphasized-in": {
    "value": "cubic-bezier(0.05, 0.7, 0.1, 1.0)",
    "variable": "var(--easings-emphasized-in)"
  },
  "easings.emphasized-out": {
    "value": "cubic-bezier(0.3, 0.0, 0.8, 0.15)",
    "variable": "var(--easings-emphasized-out)"
  },
  "fonts.inter": {
    "value": "var(--font-inter)",
    "variable": "var(--fonts-inter)"
  },
  "fonts.interDisplay": {
    "value": "var(--font-inter-display)",
    "variable": "var(--fonts-inter-display)"
  },
  "fonts.mono": {
    "value": "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, \"Liberation Mono\", \"Courier New\", monospace",
    "variable": "var(--fonts-mono)"
  },
  "fontSizes.2xs": {
    "value": "0.5rem",
    "variable": "var(--font-sizes-2xs)"
  },
  "fontSizes.xs": {
    "value": "0.75rem",
    "variable": "var(--font-sizes-xs)"
  },
  "fontSizes.sm": {
    "value": "0.875rem",
    "variable": "var(--font-sizes-sm)"
  },
  "fontSizes.md": {
    "value": "1rem",
    "variable": "var(--font-sizes-md)"
  },
  "fontSizes.lg": {
    "value": "1.125rem",
    "variable": "var(--font-sizes-lg)"
  },
  "fontSizes.xl": {
    "value": "1.25rem",
    "variable": "var(--font-sizes-xl)"
  },
  "fontSizes.2xl": {
    "value": "1.5rem",
    "variable": "var(--font-sizes-2xl)"
  },
  "fontSizes.3xl": {
    "value": "1.875rem",
    "variable": "var(--font-sizes-3xl)"
  },
  "fontSizes.4xl": {
    "value": "2.25rem",
    "variable": "var(--font-sizes-4xl)"
  },
  "fontSizes.5xl": {
    "value": "3rem",
    "variable": "var(--font-sizes-5xl)"
  },
  "fontSizes.6xl": {
    "value": "3.75rem",
    "variable": "var(--font-sizes-6xl)"
  },
  "fontSizes.7xl": {
    "value": "4.5rem",
    "variable": "var(--font-sizes-7xl)"
  },
  "fontSizes.8xl": {
    "value": "6rem",
    "variable": "var(--font-sizes-8xl)"
  },
  "fontSizes.9xl": {
    "value": "8rem",
    "variable": "var(--font-sizes-9xl)"
  },
  "fontWeights.thin": {
    "value": "100",
    "variable": "var(--font-weights-thin)"
  },
  "fontWeights.extralight": {
    "value": "200",
    "variable": "var(--font-weights-extralight)"
  },
  "fontWeights.light": {
    "value": "300",
    "variable": "var(--font-weights-light)"
  },
  "fontWeights.normal": {
    "value": "400",
    "variable": "var(--font-weights-normal)"
  },
  "fontWeights.medium": {
    "value": "500",
    "variable": "var(--font-weights-medium)"
  },
  "fontWeights.semibold": {
    "value": "600",
    "variable": "var(--font-weights-semibold)"
  },
  "fontWeights.bold": {
    "value": "700",
    "variable": "var(--font-weights-bold)"
  },
  "fontWeights.extrabold": {
    "value": "800",
    "variable": "var(--font-weights-extrabold)"
  },
  "fontWeights.black": {
    "value": "900",
    "variable": "var(--font-weights-black)"
  },
  "letterSpacings.tighter": {
    "value": "-0.05em",
    "variable": "var(--letter-spacings-tighter)"
  },
  "letterSpacings.tight": {
    "value": "-0.025em",
    "variable": "var(--letter-spacings-tight)"
  },
  "letterSpacings.normal": {
    "value": "0em",
    "variable": "var(--letter-spacings-normal)"
  },
  "letterSpacings.wide": {
    "value": "0.025em",
    "variable": "var(--letter-spacings-wide)"
  },
  "letterSpacings.wider": {
    "value": "0.05em",
    "variable": "var(--letter-spacings-wider)"
  },
  "letterSpacings.widest": {
    "value": "0.1em",
    "variable": "var(--letter-spacings-widest)"
  },
  "lineHeights.none": {
    "value": "1",
    "variable": "var(--line-heights-none)"
  },
  "lineHeights.tight": {
    "value": "1.25",
    "variable": "var(--line-heights-tight)"
  },
  "lineHeights.normal": {
    "value": "1.5",
    "variable": "var(--line-heights-normal)"
  },
  "lineHeights.relaxed": {
    "value": "1.75",
    "variable": "var(--line-heights-relaxed)"
  },
  "lineHeights.loose": {
    "value": "2",
    "variable": "var(--line-heights-loose)"
  },
  "radii.none": {
    "value": "0",
    "variable": "var(--radii-none)"
  },
  "radii.2xs": {
    "value": "0.0625rem",
    "variable": "var(--radii-2xs)"
  },
  "radii.xs": {
    "value": "0.125rem",
    "variable": "var(--radii-xs)"
  },
  "radii.sm": {
    "value": "0.25rem",
    "variable": "var(--radii-sm)"
  },
  "radii.md": {
    "value": "0.375rem",
    "variable": "var(--radii-md)"
  },
  "radii.lg": {
    "value": "0.5rem",
    "variable": "var(--radii-lg)"
  },
  "radii.xl": {
    "value": "0.75rem",
    "variable": "var(--radii-xl)"
  },
  "radii.2xl": {
    "value": "1rem",
    "variable": "var(--radii-2xl)"
  },
  "radii.3xl": {
    "value": "1.5rem",
    "variable": "var(--radii-3xl)"
  },
  "radii.full": {
    "value": "9999px",
    "variable": "var(--radii-full)"
  },
  "sizes.0": {
    "value": "0rem",
    "variable": "var(--sizes-0)"
  },
  "sizes.1": {
    "value": "0.25rem",
    "variable": "var(--sizes-1)"
  },
  "sizes.2": {
    "value": "0.5rem",
    "variable": "var(--sizes-2)"
  },
  "sizes.3": {
    "value": "0.75rem",
    "variable": "var(--sizes-3)"
  },
  "sizes.4": {
    "value": "1rem",
    "variable": "var(--sizes-4)"
  },
  "sizes.5": {
    "value": "1.25rem",
    "variable": "var(--sizes-5)"
  },
  "sizes.6": {
    "value": "1.5rem",
    "variable": "var(--sizes-6)"
  },
  "sizes.7": {
    "value": "1.75rem",
    "variable": "var(--sizes-7)"
  },
  "sizes.8": {
    "value": "2rem",
    "variable": "var(--sizes-8)"
  },
  "sizes.9": {
    "value": "2.25rem",
    "variable": "var(--sizes-9)"
  },
  "sizes.10": {
    "value": "2.5rem",
    "variable": "var(--sizes-10)"
  },
  "sizes.11": {
    "value": "2.75rem",
    "variable": "var(--sizes-11)"
  },
  "sizes.12": {
    "value": "3rem",
    "variable": "var(--sizes-12)"
  },
  "sizes.14": {
    "value": "3.5rem",
    "variable": "var(--sizes-14)"
  },
  "sizes.16": {
    "value": "4rem",
    "variable": "var(--sizes-16)"
  },
  "sizes.20": {
    "value": "5rem",
    "variable": "var(--sizes-20)"
  },
  "sizes.24": {
    "value": "6rem",
    "variable": "var(--sizes-24)"
  },
  "sizes.28": {
    "value": "7rem",
    "variable": "var(--sizes-28)"
  },
  "sizes.32": {
    "value": "8rem",
    "variable": "var(--sizes-32)"
  },
  "sizes.36": {
    "value": "9rem",
    "variable": "var(--sizes-36)"
  },
  "sizes.40": {
    "value": "10rem",
    "variable": "var(--sizes-40)"
  },
  "sizes.44": {
    "value": "11rem",
    "variable": "var(--sizes-44)"
  },
  "sizes.48": {
    "value": "12rem",
    "variable": "var(--sizes-48)"
  },
  "sizes.52": {
    "value": "13rem",
    "variable": "var(--sizes-52)"
  },
  "sizes.56": {
    "value": "14rem",
    "variable": "var(--sizes-56)"
  },
  "sizes.60": {
    "value": "15rem",
    "variable": "var(--sizes-60)"
  },
  "sizes.64": {
    "value": "16rem",
    "variable": "var(--sizes-64)"
  },
  "sizes.72": {
    "value": "18rem",
    "variable": "var(--sizes-72)"
  },
  "sizes.80": {
    "value": "20rem",
    "variable": "var(--sizes-80)"
  },
  "sizes.96": {
    "value": "24rem",
    "variable": "var(--sizes-96)"
  },
  "sizes.0.5": {
    "value": "0.125rem",
    "variable": "var(--sizes-0\\.5)"
  },
  "sizes.1.5": {
    "value": "0.375rem",
    "variable": "var(--sizes-1\\.5)"
  },
  "sizes.2.5": {
    "value": "0.625rem",
    "variable": "var(--sizes-2\\.5)"
  },
  "sizes.3.5": {
    "value": "0.875rem",
    "variable": "var(--sizes-3\\.5)"
  },
  "sizes.4.5": {
    "value": "1.125rem",
    "variable": "var(--sizes-4\\.5)"
  },
  "sizes.safeBottom": {
    "value": "env(safe-area-inset-bottom)",
    "variable": "var(--sizes-safe-bottom)"
  },
  "sizes.scrollGutter": {
    "value": "var(--spacing-2)",
    "variable": "var(--sizes-scroll-gutter)"
  },
  "sizes.2xs": {
    "value": "16rem",
    "variable": "var(--sizes-2xs)"
  },
  "sizes.xs": {
    "value": "20rem",
    "variable": "var(--sizes-xs)"
  },
  "sizes.sm": {
    "value": "24rem",
    "variable": "var(--sizes-sm)"
  },
  "sizes.md": {
    "value": "28rem",
    "variable": "var(--sizes-md)"
  },
  "sizes.lg": {
    "value": "32rem",
    "variable": "var(--sizes-lg)"
  },
  "sizes.xl": {
    "value": "36rem",
    "variable": "var(--sizes-xl)"
  },
  "sizes.2xl": {
    "value": "42rem",
    "variable": "var(--sizes-2xl)"
  },
  "sizes.3xl": {
    "value": "48rem",
    "variable": "var(--sizes-3xl)"
  },
  "sizes.4xl": {
    "value": "56rem",
    "variable": "var(--sizes-4xl)"
  },
  "sizes.5xl": {
    "value": "64rem",
    "variable": "var(--sizes-5xl)"
  },
  "sizes.6xl": {
    "value": "72rem",
    "variable": "var(--sizes-6xl)"
  },
  "sizes.7xl": {
    "value": "80rem",
    "variable": "var(--sizes-7xl)"
  },
  "sizes.8xl": {
    "value": "90rem",
    "variable": "var(--sizes-8xl)"
  },
  "sizes.full": {
    "value": "100%",
    "variable": "var(--sizes-full)"
  },
  "sizes.min": {
    "value": "min-content",
    "variable": "var(--sizes-min)"
  },
  "sizes.max": {
    "value": "max-content",
    "variable": "var(--sizes-max)"
  },
  "sizes.fit": {
    "value": "fit-content",
    "variable": "var(--sizes-fit)"
  },
  "sizes.breakpoint-sm": {
    "value": "640px",
    "variable": "var(--sizes-breakpoint-sm)"
  },
  "sizes.breakpoint-md": {
    "value": "768px",
    "variable": "var(--sizes-breakpoint-md)"
  },
  "sizes.breakpoint-lg": {
    "value": "1024px",
    "variable": "var(--sizes-breakpoint-lg)"
  },
  "sizes.breakpoint-xl": {
    "value": "1280px",
    "variable": "var(--sizes-breakpoint-xl)"
  },
  "sizes.breakpoint-2xl": {
    "value": "1536px",
    "variable": "var(--sizes-breakpoint-2xl)"
  },
  "spacing.0": {
    "value": "0rem",
    "variable": "var(--spacing-0)"
  },
  "spacing.1": {
    "value": "0.25rem",
    "variable": "var(--spacing-1)"
  },
  "spacing.2": {
    "value": "0.5rem",
    "variable": "var(--spacing-2)"
  },
  "spacing.3": {
    "value": "0.75rem",
    "variable": "var(--spacing-3)"
  },
  "spacing.4": {
    "value": "1rem",
    "variable": "var(--spacing-4)"
  },
  "spacing.5": {
    "value": "1.25rem",
    "variable": "var(--spacing-5)"
  },
  "spacing.6": {
    "value": "1.5rem",
    "variable": "var(--spacing-6)"
  },
  "spacing.7": {
    "value": "1.75rem",
    "variable": "var(--spacing-7)"
  },
  "spacing.8": {
    "value": "2rem",
    "variable": "var(--spacing-8)"
  },
  "spacing.9": {
    "value": "2.25rem",
    "variable": "var(--spacing-9)"
  },
  "spacing.10": {
    "value": "2.5rem",
    "variable": "var(--spacing-10)"
  },
  "spacing.11": {
    "value": "2.75rem",
    "variable": "var(--spacing-11)"
  },
  "spacing.12": {
    "value": "3rem",
    "variable": "var(--spacing-12)"
  },
  "spacing.14": {
    "value": "3.5rem",
    "variable": "var(--spacing-14)"
  },
  "spacing.16": {
    "value": "4rem",
    "variable": "var(--spacing-16)"
  },
  "spacing.20": {
    "value": "5rem",
    "variable": "var(--spacing-20)"
  },
  "spacing.24": {
    "value": "6rem",
    "variable": "var(--spacing-24)"
  },
  "spacing.28": {
    "value": "7rem",
    "variable": "var(--spacing-28)"
  },
  "spacing.32": {
    "value": "8rem",
    "variable": "var(--spacing-32)"
  },
  "spacing.36": {
    "value": "9rem",
    "variable": "var(--spacing-36)"
  },
  "spacing.40": {
    "value": "10rem",
    "variable": "var(--spacing-40)"
  },
  "spacing.44": {
    "value": "11rem",
    "variable": "var(--spacing-44)"
  },
  "spacing.48": {
    "value": "12rem",
    "variable": "var(--spacing-48)"
  },
  "spacing.52": {
    "value": "13rem",
    "variable": "var(--spacing-52)"
  },
  "spacing.56": {
    "value": "14rem",
    "variable": "var(--spacing-56)"
  },
  "spacing.60": {
    "value": "15rem",
    "variable": "var(--spacing-60)"
  },
  "spacing.64": {
    "value": "16rem",
    "variable": "var(--spacing-64)"
  },
  "spacing.72": {
    "value": "18rem",
    "variable": "var(--spacing-72)"
  },
  "spacing.80": {
    "value": "20rem",
    "variable": "var(--spacing-80)"
  },
  "spacing.96": {
    "value": "24rem",
    "variable": "var(--spacing-96)"
  },
  "spacing.0.5": {
    "value": "0.125rem",
    "variable": "var(--spacing-0\\.5)"
  },
  "spacing.1.5": {
    "value": "0.375rem",
    "variable": "var(--spacing-1\\.5)"
  },
  "spacing.2.5": {
    "value": "0.625rem",
    "variable": "var(--spacing-2\\.5)"
  },
  "spacing.3.5": {
    "value": "0.875rem",
    "variable": "var(--spacing-3\\.5)"
  },
  "spacing.4.5": {
    "value": "1.125rem",
    "variable": "var(--spacing-4\\.5)"
  },
  "spacing.safeBottom": {
    "value": "env(safe-area-inset-bottom)",
    "variable": "var(--spacing-safe-bottom)"
  },
  "spacing.scrollGutter": {
    "value": "var(--spacing-2)",
    "variable": "var(--spacing-scroll-gutter)"
  },
  "zIndex.hide": {
    "value": -1,
    "variable": "var(--z-index-hide)"
  },
  "zIndex.base": {
    "value": 0,
    "variable": "var(--z-index-base)"
  },
  "zIndex.docked": {
    "value": 10,
    "variable": "var(--z-index-docked)"
  },
  "zIndex.dropdown": {
    "value": 1000,
    "variable": "var(--z-index-dropdown)"
  },
  "zIndex.sticky": {
    "value": 1100,
    "variable": "var(--z-index-sticky)"
  },
  "zIndex.banner": {
    "value": 1200,
    "variable": "var(--z-index-banner)"
  },
  "zIndex.overlay": {
    "value": 1300,
    "variable": "var(--z-index-overlay)"
  },
  "zIndex.modal": {
    "value": 1400,
    "variable": "var(--z-index-modal)"
  },
  "zIndex.popover": {
    "value": 1500,
    "variable": "var(--z-index-popover)"
  },
  "zIndex.skipLink": {
    "value": 1600,
    "variable": "var(--z-index-skip-link)"
  },
  "zIndex.toast": {
    "value": 1700,
    "variable": "var(--z-index-toast)"
  },
  "zIndex.tooltip": {
    "value": 1800,
    "variable": "var(--z-index-tooltip)"
  },
  "shadows.xs": {
    "value": "0px 1px 2px 0px var(--colors-gray-light-a5), 0px 0px 1px 0px var(--colors-gray-light-a7)",
    "variable": "var(--shadows-xs)"
  },
  "shadows.sm": {
    "value": "0px 2px 4px 0px var(--colors-gray-light-a3), 0px 0px 1px 0px var(--colors-gray-light-a7)",
    "variable": "var(--shadows-sm)"
  },
  "shadows.md": {
    "value": "0px 4px 8px 0px var(--colors-gray-light-a3), 0px 0px 1px 0px var(--colors-gray-light-a7)",
    "variable": "var(--shadows-md)"
  },
  "shadows.lg": {
    "value": "0px 8px 16px 0px var(--colors-gray-light-a3), 0px 0px 1px 0px var(--colors-gray-light-a7)",
    "variable": "var(--shadows-lg)"
  },
  "shadows.xl": {
    "value": "0px 16px 24px 0px var(--colors-gray-light-a3), 0px 0px 1px 0px var(--colors-gray-light-a7)",
    "variable": "var(--shadows-xl)"
  },
  "shadows.2xl": {
    "value": "0px 24px 40px 0px var(--colors-gray-light-a3), 0px 0px 1px 0px var(--colors-gray-light-a7)",
    "variable": "var(--shadows-2xl)"
  },
  "breakpoints.sm": {
    "value": "640px",
    "variable": "var(--breakpoints-sm)"
  },
  "breakpoints.md": {
    "value": "768px",
    "variable": "var(--breakpoints-md)"
  },
  "breakpoints.lg": {
    "value": "1024px",
    "variable": "var(--breakpoints-lg)"
  },
  "breakpoints.xl": {
    "value": "1280px",
    "variable": "var(--breakpoints-xl)"
  },
  "breakpoints.2xl": {
    "value": "1536px",
    "variable": "var(--breakpoints-2xl)"
  },
  "colors.conicGradient": {
    "value": "\nconic-gradient(\n    oklch(80% 0.15 0),\noklch(80% 0.15 10),\noklch(80% 0.15 20),\noklch(80% 0.15 30),\noklch(80% 0.15 40),\noklch(80% 0.15 50),\noklch(80% 0.15 60),\noklch(80% 0.15 70),\noklch(80% 0.15 80),\noklch(80% 0.15 90),\noklch(80% 0.15 100),\noklch(80% 0.15 110),\noklch(80% 0.15 120),\noklch(80% 0.15 130),\noklch(80% 0.15 140),\noklch(80% 0.15 150),\noklch(80% 0.15 160),\noklch(80% 0.15 170),\noklch(80% 0.15 180),\noklch(80% 0.15 190),\noklch(80% 0.15 200),\noklch(80% 0.15 210),\noklch(80% 0.15 220),\noklch(80% 0.15 230),\noklch(80% 0.15 240),\noklch(80% 0.15 250),\noklch(80% 0.15 260),\noklch(80% 0.15 270),\noklch(80% 0.15 280),\noklch(80% 0.15 290),\noklch(80% 0.15 300),\noklch(80% 0.15 310),\noklch(80% 0.15 320),\noklch(80% 0.15 330),\noklch(80% 0.15 340),\noklch(80% 0.15 350),\noklch(80% 0.15 360)\n)\n",
    "variable": "var(--colors-conic-gradient)"
  },
  "colors.cardBackgroundGradient": {
    "value": "linear-gradient(90deg, var(--colors-bg), transparent)",
    "variable": "var(--colors-card-background-gradient)"
  },
  "colors.backgroundGradientH": {
    "value": "linear-gradient(90deg, var(--colors-bg), transparent)",
    "variable": "var(--colors-background-gradient-h)"
  },
  "colors.backgroundGradientV": {
    "value": "linear-gradient(0deg, var(--colors-bg), transparent)",
    "variable": "var(--colors-background-gradient-v)"
  },
  "radii.l1": {
    "value": "var(--radii-xs)",
    "variable": "var(--radii-l1)"
  },
  "radii.l2": {
    "value": "var(--radii-sm)",
    "variable": "var(--radii-l2)"
  },
  "radii.l3": {
    "value": "var(--radii-md)",
    "variable": "var(--radii-l3)"
  },
  "fonts.body": {
    "value": "var(--fonts-inter)",
    "variable": "var(--fonts-body)"
  },
  "fonts.heading": {
    "value": "var(--fonts-inter-display)",
    "variable": "var(--fonts-heading)"
  },
  "blurs.frosted": {
    "value": "10px",
    "variable": "var(--blurs-frosted)"
  },
  "opacity.0": {
    "value": "0",
    "variable": "var(--opacity-0)"
  },
  "opacity.1": {
    "value": "0.1",
    "variable": "var(--opacity-1)"
  },
  "opacity.2": {
    "value": "0.2",
    "variable": "var(--opacity-2)"
  },
  "opacity.3": {
    "value": "0.3",
    "variable": "var(--opacity-3)"
  },
  "opacity.4": {
    "value": "0.4",
    "variable": "var(--opacity-4)"
  },
  "opacity.5": {
    "value": "0.5",
    "variable": "var(--opacity-5)"
  },
  "opacity.6": {
    "value": "0.6",
    "variable": "var(--opacity-6)"
  },
  "opacity.7": {
    "value": "0.7",
    "variable": "var(--opacity-7)"
  },
  "opacity.8": {
    "value": "0.8",
    "variable": "var(--opacity-8)"
  },
  "opacity.9": {
    "value": "0.9",
    "variable": "var(--opacity-9)"
  },
  "opacity.full": {
    "value": "1",
    "variable": "var(--opacity-full)"
  },
  "borderWidths.none": {
    "value": "0",
    "variable": "var(--border-widths-none)"
  },
  "borderWidths.hairline": {
    "value": "0.5px",
    "variable": "var(--border-widths-hairline)"
  },
  "borderWidths.thin": {
    "value": "1px",
    "variable": "var(--border-widths-thin)"
  },
  "borderWidths.medium": {
    "value": "2px",
    "variable": "var(--border-widths-medium)"
  },
  "borderWidths.thick": {
    "value": "3px",
    "variable": "var(--border-widths-thick)"
  },
  "sizes.prose": {
    "value": "65ch",
    "variable": "var(--sizes-prose)"
  },
  "spacing.safeTop": {
    "value": "calc(env(keyboard-inset-height) + 4px)",
    "variable": "var(--spacing-safe-top)"
  },
  "spacing.-1": {
    "value": "calc(var(--spacing-1) * -1)",
    "variable": "var(--spacing-1)"
  },
  "spacing.-2": {
    "value": "calc(var(--spacing-2) * -1)",
    "variable": "var(--spacing-2)"
  },
  "spacing.-3": {
    "value": "calc(var(--spacing-3) * -1)",
    "variable": "var(--spacing-3)"
  },
  "spacing.-4": {
    "value": "calc(var(--spacing-4) * -1)",
    "variable": "var(--spacing-4)"
  },
  "spacing.-5": {
    "value": "calc(var(--spacing-5) * -1)",
    "variable": "var(--spacing-5)"
  },
  "spacing.-6": {
    "value": "calc(var(--spacing-6) * -1)",
    "variable": "var(--spacing-6)"
  },
  "spacing.-7": {
    "value": "calc(var(--spacing-7) * -1)",
    "variable": "var(--spacing-7)"
  },
  "spacing.-8": {
    "value": "calc(var(--spacing-8) * -1)",
    "variable": "var(--spacing-8)"
  },
  "spacing.-9": {
    "value": "calc(var(--spacing-9) * -1)",
    "variable": "var(--spacing-9)"
  },
  "spacing.-10": {
    "value": "calc(var(--spacing-10) * -1)",
    "variable": "var(--spacing-10)"
  },
  "spacing.-11": {
    "value": "calc(var(--spacing-11) * -1)",
    "variable": "var(--spacing-11)"
  },
  "spacing.-12": {
    "value": "calc(var(--spacing-12) * -1)",
    "variable": "var(--spacing-12)"
  },
  "spacing.-14": {
    "value": "calc(var(--spacing-14) * -1)",
    "variable": "var(--spacing-14)"
  },
  "spacing.-16": {
    "value": "calc(var(--spacing-16) * -1)",
    "variable": "var(--spacing-16)"
  },
  "spacing.-20": {
    "value": "calc(var(--spacing-20) * -1)",
    "variable": "var(--spacing-20)"
  },
  "spacing.-24": {
    "value": "calc(var(--spacing-24) * -1)",
    "variable": "var(--spacing-24)"
  },
  "spacing.-28": {
    "value": "calc(var(--spacing-28) * -1)",
    "variable": "var(--spacing-28)"
  },
  "spacing.-32": {
    "value": "calc(var(--spacing-32) * -1)",
    "variable": "var(--spacing-32)"
  },
  "spacing.-36": {
    "value": "calc(var(--spacing-36) * -1)",
    "variable": "var(--spacing-36)"
  },
  "spacing.-40": {
    "value": "calc(var(--spacing-40) * -1)",
    "variable": "var(--spacing-40)"
  },
  "spacing.-44": {
    "value": "calc(var(--spacing-44) * -1)",
    "variable": "var(--spacing-44)"
  },
  "spacing.-48": {
    "value": "calc(var(--spacing-48) * -1)",
    "variable": "var(--spacing-48)"
  },
  "spacing.-52": {
    "value": "calc(var(--spacing-52) * -1)",
    "variable": "var(--spacing-52)"
  },
  "spacing.-56": {
    "value": "calc(var(--spacing-56) * -1)",
    "variable": "var(--spacing-56)"
  },
  "spacing.-60": {
    "value": "calc(var(--spacing-60) * -1)",
    "variable": "var(--spacing-60)"
  },
  "spacing.-64": {
    "value": "calc(var(--spacing-64) * -1)",
    "variable": "var(--spacing-64)"
  },
  "spacing.-72": {
    "value": "calc(var(--spacing-72) * -1)",
    "variable": "var(--spacing-72)"
  },
  "spacing.-80": {
    "value": "calc(var(--spacing-80) * -1)",
    "variable": "var(--spacing-80)"
  },
  "spacing.-96": {
    "value": "calc(var(--spacing-96) * -1)",
    "variable": "var(--spacing-96)"
  },
  "spacing.-0.5": {
    "value": "calc(var(--spacing-0\\.5) * -1)",
    "variable": "var(--spacing-0\\.5)"
  },
  "spacing.-1.5": {
    "value": "calc(var(--spacing-1\\.5) * -1)",
    "variable": "var(--spacing-1\\.5)"
  },
  "spacing.-2.5": {
    "value": "calc(var(--spacing-2\\.5) * -1)",
    "variable": "var(--spacing-2\\.5)"
  },
  "spacing.-3.5": {
    "value": "calc(var(--spacing-3\\.5) * -1)",
    "variable": "var(--spacing-3\\.5)"
  },
  "spacing.-4.5": {
    "value": "calc(var(--spacing-4\\.5) * -1)",
    "variable": "var(--spacing-4\\.5)"
  },
  "spacing.-safeBottom": {
    "value": "calc(var(--spacing-safe-bottom) * -1)",
    "variable": "var(--spacing-safe-bottom)"
  },
  "spacing.-scrollGutter": {
    "value": "calc(var(--spacing-scroll-gutter) * -1)",
    "variable": "var(--spacing-scroll-gutter)"
  },
  "spacing.-safeTop": {
    "value": "calc(var(--spacing-safe-top) * -1)",
    "variable": "var(--spacing-safe-top)"
  },
  "colors.bg.default": {
    "value": "var(--colors-bg-default)",
    "variable": "var(--colors-bg-default)"
  },
  "colors.bg.subtle": {
    "value": "var(--colors-bg-subtle)",
    "variable": "var(--colors-bg-subtle)"
  },
  "colors.bg.muted": {
    "value": "var(--colors-bg-muted)",
    "variable": "var(--colors-bg-muted)"
  },
  "colors.bg.selected": {
    "value": "var(--colors-bg-selected)",
    "variable": "var(--colors-bg-selected)"
  },
  "colors.bg.emphasized": {
    "value": "var(--colors-bg-emphasized)",
    "variable": "var(--colors-bg-emphasized)"
  },
  "colors.bg.disabled": {
    "value": "var(--colors-bg-disabled)",
    "variable": "var(--colors-bg-disabled)"
  },
  "colors.bg.destructive": {
    "value": "var(--colors-bg-destructive)",
    "variable": "var(--colors-bg-destructive)"
  },
  "colors.bg.success": {
    "value": "var(--colors-bg-success)",
    "variable": "var(--colors-bg-success)"
  },
  "colors.bg.warning": {
    "value": "var(--colors-bg-warning)",
    "variable": "var(--colors-bg-warning)"
  },
  "colors.bg.error": {
    "value": "var(--colors-bg-error)",
    "variable": "var(--colors-bg-error)"
  },
  "colors.bg.info": {
    "value": "var(--colors-bg-info)",
    "variable": "var(--colors-bg-info)"
  },
  "colors.bg.accent": {
    "value": "var(--colors-bg-accent)",
    "variable": "var(--colors-bg-accent)"
  },
  "colors.bg.site": {
    "value": "var(--colors-bg-site)",
    "variable": "var(--colors-bg-site)"
  },
  "colors.bg.opaque": {
    "value": "var(--colors-bg-opaque)",
    "variable": "var(--colors-bg-opaque)"
  },
  "colors.fg.default": {
    "value": "var(--colors-fg-default)",
    "variable": "var(--colors-fg-default)"
  },
  "colors.fg.subtle": {
    "value": "var(--colors-fg-subtle)",
    "variable": "var(--colors-fg-subtle)"
  },
  "colors.fg.muted": {
    "value": "var(--colors-fg-muted)",
    "variable": "var(--colors-fg-muted)"
  },
  "colors.fg.selected": {
    "value": "var(--colors-fg-selected)",
    "variable": "var(--colors-fg-selected)"
  },
  "colors.fg.emphasized": {
    "value": "var(--colors-fg-emphasized)",
    "variable": "var(--colors-fg-emphasized)"
  },
  "colors.fg.disabled": {
    "value": "var(--colors-fg-disabled)",
    "variable": "var(--colors-fg-disabled)"
  },
  "colors.fg.destructive": {
    "value": "var(--colors-fg-destructive)",
    "variable": "var(--colors-fg-destructive)"
  },
  "colors.fg.success": {
    "value": "var(--colors-fg-success)",
    "variable": "var(--colors-fg-success)"
  },
  "colors.fg.warning": {
    "value": "var(--colors-fg-warning)",
    "variable": "var(--colors-fg-warning)"
  },
  "colors.fg.error": {
    "value": "var(--colors-fg-error)",
    "variable": "var(--colors-fg-error)"
  },
  "colors.fg.info": {
    "value": "var(--colors-fg-info)",
    "variable": "var(--colors-fg-info)"
  },
  "colors.fg.accent": {
    "value": "var(--colors-fg-accent)",
    "variable": "var(--colors-fg-accent)"
  },
  "colors.border": {
    "value": "var(--colors-border)",
    "variable": "var(--colors-border)"
  },
  "colors.border.default": {
    "value": "var(--colors-border-default)",
    "variable": "var(--colors-border-default)"
  },
  "colors.border.subtle": {
    "value": "var(--colors-border-subtle)",
    "variable": "var(--colors-border-subtle)"
  },
  "colors.border.muted": {
    "value": "var(--colors-border-muted)",
    "variable": "var(--colors-border-muted)"
  },
  "colors.border.destructive": {
    "value": "var(--colors-border-destructive)",
    "variable": "var(--colors-border-destructive)"
  },
  "colors.border.success": {
    "value": "var(--colors-border-success)",
    "variable": "var(--colors-border-success)"
  },
  "colors.border.warning": {
    "value": "var(--colors-border-warning)",
    "variable": "var(--colors-border-warning)"
  },
  "colors.border.error": {
    "value": "var(--colors-border-error)",
    "variable": "var(--colors-border-error)"
  },
  "colors.border.info": {
    "value": "var(--colors-border-info)",
    "variable": "var(--colors-border-info)"
  },
  "colors.border.accent": {
    "value": "var(--colors-border-accent)",
    "variable": "var(--colors-border-accent)"
  },
  "colors.border.disabled": {
    "value": "var(--colors-border-disabled)",
    "variable": "var(--colors-border-disabled)"
  },
  "colors.border.outline": {
    "value": "var(--colors-border-outline)",
    "variable": "var(--colors-border-outline)"
  },
  "colors.visibility.published.bg": {
    "value": "var(--colors-visibility-published-bg)",
    "variable": "var(--colors-visibility-published-bg)"
  },
  "colors.visibility.published.fg": {
    "value": "var(--colors-visibility-published-fg)",
    "variable": "var(--colors-visibility-published-fg)"
  },
  "colors.visibility.published.border": {
    "value": "var(--colors-visibility-published-border)",
    "variable": "var(--colors-visibility-published-border)"
  },
  "colors.visibility.draft.bg": {
    "value": "var(--colors-visibility-draft-bg)",
    "variable": "var(--colors-visibility-draft-bg)"
  },
  "colors.visibility.draft.fg": {
    "value": "var(--colors-visibility-draft-fg)",
    "variable": "var(--colors-visibility-draft-fg)"
  },
  "colors.visibility.draft.border": {
    "value": "var(--colors-visibility-draft-border)",
    "variable": "var(--colors-visibility-draft-border)"
  },
  "colors.visibility.review.bg": {
    "value": "var(--colors-visibility-review-bg)",
    "variable": "var(--colors-visibility-review-bg)"
  },
  "colors.visibility.review.fg": {
    "value": "var(--colors-visibility-review-fg)",
    "variable": "var(--colors-visibility-review-fg)"
  },
  "colors.visibility.review.border": {
    "value": "var(--colors-visibility-review-border)",
    "variable": "var(--colors-visibility-review-border)"
  },
  "colors.visibility.unlisted.bg": {
    "value": "var(--colors-visibility-unlisted-bg)",
    "variable": "var(--colors-visibility-unlisted-bg)"
  },
  "colors.visibility.unlisted.fg": {
    "value": "var(--colors-visibility-unlisted-fg)",
    "variable": "var(--colors-visibility-unlisted-fg)"
  },
  "colors.visibility.unlisted.border": {
    "value": "var(--colors-visibility-unlisted-border)",
    "variable": "var(--colors-visibility-unlisted-border)"
  },
  "colors.accent.1": {
    "value": "var(--colors-accent-1)",
    "variable": "var(--colors-accent-1)"
  },
  "colors.accent.2": {
    "value": "var(--colors-accent-2)",
    "variable": "var(--colors-accent-2)"
  },
  "colors.accent.3": {
    "value": "var(--colors-accent-3)",
    "variable": "var(--colors-accent-3)"
  },
  "colors.accent.4": {
    "value": "var(--colors-accent-4)",
    "variable": "var(--colors-accent-4)"
  },
  "colors.accent.5": {
    "value": "var(--colors-accent-5)",
    "variable": "var(--colors-accent-5)"
  },
  "colors.accent.6": {
    "value": "var(--colors-accent-6)",
    "variable": "var(--colors-accent-6)"
  },
  "colors.accent.7": {
    "value": "var(--colors-accent-7)",
    "variable": "var(--colors-accent-7)"
  },
  "colors.accent.8": {
    "value": "var(--colors-accent-8)",
    "variable": "var(--colors-accent-8)"
  },
  "colors.accent.9": {
    "value": "var(--colors-accent-9)",
    "variable": "var(--colors-accent-9)"
  },
  "colors.accent.10": {
    "value": "var(--colors-accent-10)",
    "variable": "var(--colors-accent-10)"
  },
  "colors.accent.default": {
    "value": "var(--colors-accent-default)",
    "variable": "var(--colors-accent-default)"
  },
  "colors.accent.subtle": {
    "value": "var(--colors-accent-subtle)",
    "variable": "var(--colors-accent-subtle)"
  },
  "colors.accent.muted": {
    "value": "var(--colors-accent-muted)",
    "variable": "var(--colors-accent-muted)"
  },
  "colors.blue.1": {
    "value": "var(--colors-blue-1)",
    "variable": "var(--colors-blue-1)"
  },
  "colors.blue.2": {
    "value": "var(--colors-blue-2)",
    "variable": "var(--colors-blue-2)"
  },
  "colors.blue.3": {
    "value": "var(--colors-blue-3)",
    "variable": "var(--colors-blue-3)"
  },
  "colors.blue.4": {
    "value": "var(--colors-blue-4)",
    "variable": "var(--colors-blue-4)"
  },
  "colors.blue.5": {
    "value": "var(--colors-blue-5)",
    "variable": "var(--colors-blue-5)"
  },
  "colors.blue.6": {
    "value": "var(--colors-blue-6)",
    "variable": "var(--colors-blue-6)"
  },
  "colors.blue.7": {
    "value": "var(--colors-blue-7)",
    "variable": "var(--colors-blue-7)"
  },
  "colors.blue.8": {
    "value": "var(--colors-blue-8)",
    "variable": "var(--colors-blue-8)"
  },
  "colors.blue.9": {
    "value": "var(--colors-blue-9)",
    "variable": "var(--colors-blue-9)"
  },
  "colors.blue.10": {
    "value": "var(--colors-blue-10)",
    "variable": "var(--colors-blue-10)"
  },
  "colors.blue.11": {
    "value": "var(--colors-blue-11)",
    "variable": "var(--colors-blue-11)"
  },
  "colors.blue.12": {
    "value": "var(--colors-blue-12)",
    "variable": "var(--colors-blue-12)"
  },
  "colors.blue.a1": {
    "value": "var(--colors-blue-a1)",
    "variable": "var(--colors-blue-a1)"
  },
  "colors.blue.a2": {
    "value": "var(--colors-blue-a2)",
    "variable": "var(--colors-blue-a2)"
  },
  "colors.blue.a3": {
    "value": "var(--colors-blue-a3)",
    "variable": "var(--colors-blue-a3)"
  },
  "colors.blue.a4": {
    "value": "var(--colors-blue-a4)",
    "variable": "var(--colors-blue-a4)"
  },
  "colors.blue.a5": {
    "value": "var(--colors-blue-a5)",
    "variable": "var(--colors-blue-a5)"
  },
  "colors.blue.a6": {
    "value": "var(--colors-blue-a6)",
    "variable": "var(--colors-blue-a6)"
  },
  "colors.blue.a7": {
    "value": "var(--colors-blue-a7)",
    "variable": "var(--colors-blue-a7)"
  },
  "colors.blue.a8": {
    "value": "var(--colors-blue-a8)",
    "variable": "var(--colors-blue-a8)"
  },
  "colors.blue.a9": {
    "value": "var(--colors-blue-a9)",
    "variable": "var(--colors-blue-a9)"
  },
  "colors.blue.a10": {
    "value": "var(--colors-blue-a10)",
    "variable": "var(--colors-blue-a10)"
  },
  "colors.blue.a11": {
    "value": "var(--colors-blue-a11)",
    "variable": "var(--colors-blue-a11)"
  },
  "colors.blue.a12": {
    "value": "var(--colors-blue-a12)",
    "variable": "var(--colors-blue-a12)"
  },
  "colors.blue.default": {
    "value": "var(--colors-blue-default)",
    "variable": "var(--colors-blue-default)"
  },
  "colors.blue.emphasized": {
    "value": "var(--colors-blue-emphasized)",
    "variable": "var(--colors-blue-emphasized)"
  },
  "colors.blue.fg": {
    "value": "var(--colors-blue-fg)",
    "variable": "var(--colors-blue-fg)"
  },
  "colors.blue.text": {
    "value": "var(--colors-blue-text)",
    "variable": "var(--colors-blue-text)"
  },
  "colors.overflow-fade": {
    "value": "var(--colors-overflow-fade)",
    "variable": "var(--colors-overflow-fade)"
  },
  "colors.colorPalette": {
    "value": "var(--colors-color-palette)",
    "variable": "var(--colors-color-palette)"
  },
  "colors.colorPalette.a1": {
    "value": "var(--colors-color-palette-a1)",
    "variable": "var(--colors-color-palette-a1)"
  },
  "colors.colorPalette.a2": {
    "value": "var(--colors-color-palette-a2)",
    "variable": "var(--colors-color-palette-a2)"
  },
  "colors.colorPalette.a3": {
    "value": "var(--colors-color-palette-a3)",
    "variable": "var(--colors-color-palette-a3)"
  },
  "colors.colorPalette.a4": {
    "value": "var(--colors-color-palette-a4)",
    "variable": "var(--colors-color-palette-a4)"
  },
  "colors.colorPalette.a5": {
    "value": "var(--colors-color-palette-a5)",
    "variable": "var(--colors-color-palette-a5)"
  },
  "colors.colorPalette.a6": {
    "value": "var(--colors-color-palette-a6)",
    "variable": "var(--colors-color-palette-a6)"
  },
  "colors.colorPalette.a7": {
    "value": "var(--colors-color-palette-a7)",
    "variable": "var(--colors-color-palette-a7)"
  },
  "colors.colorPalette.a8": {
    "value": "var(--colors-color-palette-a8)",
    "variable": "var(--colors-color-palette-a8)"
  },
  "colors.colorPalette.a9": {
    "value": "var(--colors-color-palette-a9)",
    "variable": "var(--colors-color-palette-a9)"
  },
  "colors.colorPalette.a10": {
    "value": "var(--colors-color-palette-a10)",
    "variable": "var(--colors-color-palette-a10)"
  },
  "colors.colorPalette.a11": {
    "value": "var(--colors-color-palette-a11)",
    "variable": "var(--colors-color-palette-a11)"
  },
  "colors.colorPalette.a12": {
    "value": "var(--colors-color-palette-a12)",
    "variable": "var(--colors-color-palette-a12)"
  },
  "colors.colorPalette.light.1": {
    "value": "var(--colors-color-palette-light-1)",
    "variable": "var(--colors-color-palette-light-1)"
  },
  "colors.colorPalette.1": {
    "value": "var(--colors-color-palette-1)",
    "variable": "var(--colors-color-palette-1)"
  },
  "colors.colorPalette.light.2": {
    "value": "var(--colors-color-palette-light-2)",
    "variable": "var(--colors-color-palette-light-2)"
  },
  "colors.colorPalette.2": {
    "value": "var(--colors-color-palette-2)",
    "variable": "var(--colors-color-palette-2)"
  },
  "colors.colorPalette.light.3": {
    "value": "var(--colors-color-palette-light-3)",
    "variable": "var(--colors-color-palette-light-3)"
  },
  "colors.colorPalette.3": {
    "value": "var(--colors-color-palette-3)",
    "variable": "var(--colors-color-palette-3)"
  },
  "colors.colorPalette.light.4": {
    "value": "var(--colors-color-palette-light-4)",
    "variable": "var(--colors-color-palette-light-4)"
  },
  "colors.colorPalette.4": {
    "value": "var(--colors-color-palette-4)",
    "variable": "var(--colors-color-palette-4)"
  },
  "colors.colorPalette.light.5": {
    "value": "var(--colors-color-palette-light-5)",
    "variable": "var(--colors-color-palette-light-5)"
  },
  "colors.colorPalette.5": {
    "value": "var(--colors-color-palette-5)",
    "variable": "var(--colors-color-palette-5)"
  },
  "colors.colorPalette.light.6": {
    "value": "var(--colors-color-palette-light-6)",
    "variable": "var(--colors-color-palette-light-6)"
  },
  "colors.colorPalette.6": {
    "value": "var(--colors-color-palette-6)",
    "variable": "var(--colors-color-palette-6)"
  },
  "colors.colorPalette.light.7": {
    "value": "var(--colors-color-palette-light-7)",
    "variable": "var(--colors-color-palette-light-7)"
  },
  "colors.colorPalette.7": {
    "value": "var(--colors-color-palette-7)",
    "variable": "var(--colors-color-palette-7)"
  },
  "colors.colorPalette.light.8": {
    "value": "var(--colors-color-palette-light-8)",
    "variable": "var(--colors-color-palette-light-8)"
  },
  "colors.colorPalette.8": {
    "value": "var(--colors-color-palette-8)",
    "variable": "var(--colors-color-palette-8)"
  },
  "colors.colorPalette.light.9": {
    "value": "var(--colors-color-palette-light-9)",
    "variable": "var(--colors-color-palette-light-9)"
  },
  "colors.colorPalette.9": {
    "value": "var(--colors-color-palette-9)",
    "variable": "var(--colors-color-palette-9)"
  },
  "colors.colorPalette.light.10": {
    "value": "var(--colors-color-palette-light-10)",
    "variable": "var(--colors-color-palette-light-10)"
  },
  "colors.colorPalette.10": {
    "value": "var(--colors-color-palette-10)",
    "variable": "var(--colors-color-palette-10)"
  },
  "colors.colorPalette.light.11": {
    "value": "var(--colors-color-palette-light-11)",
    "variable": "var(--colors-color-palette-light-11)"
  },
  "colors.colorPalette.11": {
    "value": "var(--colors-color-palette-11)",
    "variable": "var(--colors-color-palette-11)"
  },
  "colors.colorPalette.light.12": {
    "value": "var(--colors-color-palette-light-12)",
    "variable": "var(--colors-color-palette-light-12)"
  },
  "colors.colorPalette.12": {
    "value": "var(--colors-color-palette-12)",
    "variable": "var(--colors-color-palette-12)"
  },
  "colors.colorPalette.light.a1": {
    "value": "var(--colors-color-palette-light-a1)",
    "variable": "var(--colors-color-palette-light-a1)"
  },
  "colors.colorPalette.light.a2": {
    "value": "var(--colors-color-palette-light-a2)",
    "variable": "var(--colors-color-palette-light-a2)"
  },
  "colors.colorPalette.light.a3": {
    "value": "var(--colors-color-palette-light-a3)",
    "variable": "var(--colors-color-palette-light-a3)"
  },
  "colors.colorPalette.light.a4": {
    "value": "var(--colors-color-palette-light-a4)",
    "variable": "var(--colors-color-palette-light-a4)"
  },
  "colors.colorPalette.light.a5": {
    "value": "var(--colors-color-palette-light-a5)",
    "variable": "var(--colors-color-palette-light-a5)"
  },
  "colors.colorPalette.light.a6": {
    "value": "var(--colors-color-palette-light-a6)",
    "variable": "var(--colors-color-palette-light-a6)"
  },
  "colors.colorPalette.light.a7": {
    "value": "var(--colors-color-palette-light-a7)",
    "variable": "var(--colors-color-palette-light-a7)"
  },
  "colors.colorPalette.light.a8": {
    "value": "var(--colors-color-palette-light-a8)",
    "variable": "var(--colors-color-palette-light-a8)"
  },
  "colors.colorPalette.light.a9": {
    "value": "var(--colors-color-palette-light-a9)",
    "variable": "var(--colors-color-palette-light-a9)"
  },
  "colors.colorPalette.light.a10": {
    "value": "var(--colors-color-palette-light-a10)",
    "variable": "var(--colors-color-palette-light-a10)"
  },
  "colors.colorPalette.light.a11": {
    "value": "var(--colors-color-palette-light-a11)",
    "variable": "var(--colors-color-palette-light-a11)"
  },
  "colors.colorPalette.light.a12": {
    "value": "var(--colors-color-palette-light-a12)",
    "variable": "var(--colors-color-palette-light-a12)"
  },
  "colors.colorPalette.dark.1": {
    "value": "var(--colors-color-palette-dark-1)",
    "variable": "var(--colors-color-palette-dark-1)"
  },
  "colors.colorPalette.dark.2": {
    "value": "var(--colors-color-palette-dark-2)",
    "variable": "var(--colors-color-palette-dark-2)"
  },
  "colors.colorPalette.dark.3": {
    "value": "var(--colors-color-palette-dark-3)",
    "variable": "var(--colors-color-palette-dark-3)"
  },
  "colors.colorPalette.dark.4": {
    "value": "var(--colors-color-palette-dark-4)",
    "variable": "var(--colors-color-palette-dark-4)"
  },
  "colors.colorPalette.dark.5": {
    "value": "var(--colors-color-palette-dark-5)",
    "variable": "var(--colors-color-palette-dark-5)"
  },
  "colors.colorPalette.dark.6": {
    "value": "var(--colors-color-palette-dark-6)",
    "variable": "var(--colors-color-palette-dark-6)"
  },
  "colors.colorPalette.dark.7": {
    "value": "var(--colors-color-palette-dark-7)",
    "variable": "var(--colors-color-palette-dark-7)"
  },
  "colors.colorPalette.dark.8": {
    "value": "var(--colors-color-palette-dark-8)",
    "variable": "var(--colors-color-palette-dark-8)"
  },
  "colors.colorPalette.dark.9": {
    "value": "var(--colors-color-palette-dark-9)",
    "variable": "var(--colors-color-palette-dark-9)"
  },
  "colors.colorPalette.dark.10": {
    "value": "var(--colors-color-palette-dark-10)",
    "variable": "var(--colors-color-palette-dark-10)"
  },
  "colors.colorPalette.dark.11": {
    "value": "var(--colors-color-palette-dark-11)",
    "variable": "var(--colors-color-palette-dark-11)"
  },
  "colors.colorPalette.dark.12": {
    "value": "var(--colors-color-palette-dark-12)",
    "variable": "var(--colors-color-palette-dark-12)"
  },
  "colors.colorPalette.dark.a1": {
    "value": "var(--colors-color-palette-dark-a1)",
    "variable": "var(--colors-color-palette-dark-a1)"
  },
  "colors.colorPalette.dark.a2": {
    "value": "var(--colors-color-palette-dark-a2)",
    "variable": "var(--colors-color-palette-dark-a2)"
  },
  "colors.colorPalette.dark.a3": {
    "value": "var(--colors-color-palette-dark-a3)",
    "variable": "var(--colors-color-palette-dark-a3)"
  },
  "colors.colorPalette.dark.a4": {
    "value": "var(--colors-color-palette-dark-a4)",
    "variable": "var(--colors-color-palette-dark-a4)"
  },
  "colors.colorPalette.dark.a5": {
    "value": "var(--colors-color-palette-dark-a5)",
    "variable": "var(--colors-color-palette-dark-a5)"
  },
  "colors.colorPalette.dark.a6": {
    "value": "var(--colors-color-palette-dark-a6)",
    "variable": "var(--colors-color-palette-dark-a6)"
  },
  "colors.colorPalette.dark.a7": {
    "value": "var(--colors-color-palette-dark-a7)",
    "variable": "var(--colors-color-palette-dark-a7)"
  },
  "colors.colorPalette.dark.a8": {
    "value": "var(--colors-color-palette-dark-a8)",
    "variable": "var(--colors-color-palette-dark-a8)"
  },
  "colors.colorPalette.dark.a9": {
    "value": "var(--colors-color-palette-dark-a9)",
    "variable": "var(--colors-color-palette-dark-a9)"
  },
  "colors.colorPalette.dark.a10": {
    "value": "var(--colors-color-palette-dark-a10)",
    "variable": "var(--colors-color-palette-dark-a10)"
  },
  "colors.colorPalette.dark.a11": {
    "value": "var(--colors-color-palette-dark-a11)",
    "variable": "var(--colors-color-palette-dark-a11)"
  },
  "colors.colorPalette.dark.a12": {
    "value": "var(--colors-color-palette-dark-a12)",
    "variable": "var(--colors-color-palette-dark-a12)"
  },
  "colors.colorPalette.light.text.1": {
    "value": "var(--colors-color-palette-light-text-1)",
    "variable": "var(--colors-color-palette-light-text-1)"
  },
  "colors.colorPalette.text.1": {
    "value": "var(--colors-color-palette-text-1)",
    "variable": "var(--colors-color-palette-text-1)"
  },
  "colors.colorPalette.light.text.2": {
    "value": "var(--colors-color-palette-light-text-2)",
    "variable": "var(--colors-color-palette-light-text-2)"
  },
  "colors.colorPalette.text.2": {
    "value": "var(--colors-color-palette-text-2)",
    "variable": "var(--colors-color-palette-text-2)"
  },
  "colors.colorPalette.light.text.3": {
    "value": "var(--colors-color-palette-light-text-3)",
    "variable": "var(--colors-color-palette-light-text-3)"
  },
  "colors.colorPalette.text.3": {
    "value": "var(--colors-color-palette-text-3)",
    "variable": "var(--colors-color-palette-text-3)"
  },
  "colors.colorPalette.light.text.4": {
    "value": "var(--colors-color-palette-light-text-4)",
    "variable": "var(--colors-color-palette-light-text-4)"
  },
  "colors.colorPalette.text.4": {
    "value": "var(--colors-color-palette-text-4)",
    "variable": "var(--colors-color-palette-text-4)"
  },
  "colors.colorPalette.light.text.5": {
    "value": "var(--colors-color-palette-light-text-5)",
    "variable": "var(--colors-color-palette-light-text-5)"
  },
  "colors.colorPalette.text.5": {
    "value": "var(--colors-color-palette-text-5)",
    "variable": "var(--colors-color-palette-text-5)"
  },
  "colors.colorPalette.light.text.6": {
    "value": "var(--colors-color-palette-light-text-6)",
    "variable": "var(--colors-color-palette-light-text-6)"
  },
  "colors.colorPalette.text.6": {
    "value": "var(--colors-color-palette-text-6)",
    "variable": "var(--colors-color-palette-text-6)"
  },
  "colors.colorPalette.light.text.7": {
    "value": "var(--colors-color-palette-light-text-7)",
    "variable": "var(--colors-color-palette-light-text-7)"
  },
  "colors.colorPalette.text.7": {
    "value": "var(--colors-color-palette-text-7)",
    "variable": "var(--colors-color-palette-text-7)"
  },
  "colors.colorPalette.light.text.8": {
    "value": "var(--colors-color-palette-light-text-8)",
    "variable": "var(--colors-color-palette-light-text-8)"
  },
  "colors.colorPalette.text.8": {
    "value": "var(--colors-color-palette-text-8)",
    "variable": "var(--colors-color-palette-text-8)"
  },
  "colors.colorPalette.light.text.9": {
    "value": "var(--colors-color-palette-light-text-9)",
    "variable": "var(--colors-color-palette-light-text-9)"
  },
  "colors.colorPalette.text.9": {
    "value": "var(--colors-color-palette-text-9)",
    "variable": "var(--colors-color-palette-text-9)"
  },
  "colors.colorPalette.light.text.10": {
    "value": "var(--colors-color-palette-light-text-10)",
    "variable": "var(--colors-color-palette-light-text-10)"
  },
  "colors.colorPalette.text.10": {
    "value": "var(--colors-color-palette-text-10)",
    "variable": "var(--colors-color-palette-text-10)"
  },
  "colors.colorPalette.dark.text.1": {
    "value": "var(--colors-color-palette-dark-text-1)",
    "variable": "var(--colors-color-palette-dark-text-1)"
  },
  "colors.colorPalette.dark.text.2": {
    "value": "var(--colors-color-palette-dark-text-2)",
    "variable": "var(--colors-color-palette-dark-text-2)"
  },
  "colors.colorPalette.dark.text.3": {
    "value": "var(--colors-color-palette-dark-text-3)",
    "variable": "var(--colors-color-palette-dark-text-3)"
  },
  "colors.colorPalette.dark.text.4": {
    "value": "var(--colors-color-palette-dark-text-4)",
    "variable": "var(--colors-color-palette-dark-text-4)"
  },
  "colors.colorPalette.dark.text.5": {
    "value": "var(--colors-color-palette-dark-text-5)",
    "variable": "var(--colors-color-palette-dark-text-5)"
  },
  "colors.colorPalette.dark.text.6": {
    "value": "var(--colors-color-palette-dark-text-6)",
    "variable": "var(--colors-color-palette-dark-text-6)"
  },
  "colors.colorPalette.dark.text.7": {
    "value": "var(--colors-color-palette-dark-text-7)",
    "variable": "var(--colors-color-palette-dark-text-7)"
  },
  "colors.colorPalette.dark.text.8": {
    "value": "var(--colors-color-palette-dark-text-8)",
    "variable": "var(--colors-color-palette-dark-text-8)"
  },
  "colors.colorPalette.dark.text.9": {
    "value": "var(--colors-color-palette-dark-text-9)",
    "variable": "var(--colors-color-palette-dark-text-9)"
  },
  "colors.colorPalette.dark.text.10": {
    "value": "var(--colors-color-palette-dark-text-10)",
    "variable": "var(--colors-color-palette-dark-text-10)"
  },
  "colors.colorPalette.default": {
    "value": "var(--colors-color-palette-default)",
    "variable": "var(--colors-color-palette-default)"
  },
  "colors.colorPalette.subtle": {
    "value": "var(--colors-color-palette-subtle)",
    "variable": "var(--colors-color-palette-subtle)"
  },
  "colors.colorPalette.muted": {
    "value": "var(--colors-color-palette-muted)",
    "variable": "var(--colors-color-palette-muted)"
  },
  "colors.colorPalette.selected": {
    "value": "var(--colors-color-palette-selected)",
    "variable": "var(--colors-color-palette-selected)"
  },
  "colors.colorPalette.emphasized": {
    "value": "var(--colors-color-palette-emphasized)",
    "variable": "var(--colors-color-palette-emphasized)"
  },
  "colors.colorPalette.disabled": {
    "value": "var(--colors-color-palette-disabled)",
    "variable": "var(--colors-color-palette-disabled)"
  },
  "colors.colorPalette.destructive": {
    "value": "var(--colors-color-palette-destructive)",
    "variable": "var(--colors-color-palette-destructive)"
  },
  "colors.colorPalette.success": {
    "value": "var(--colors-color-palette-success)",
    "variable": "var(--colors-color-palette-success)"
  },
  "colors.colorPalette.warning": {
    "value": "var(--colors-color-palette-warning)",
    "variable": "var(--colors-color-palette-warning)"
  },
  "colors.colorPalette.error": {
    "value": "var(--colors-color-palette-error)",
    "variable": "var(--colors-color-palette-error)"
  },
  "colors.colorPalette.info": {
    "value": "var(--colors-color-palette-info)",
    "variable": "var(--colors-color-palette-info)"
  },
  "colors.colorPalette.accent": {
    "value": "var(--colors-color-palette-accent)",
    "variable": "var(--colors-color-palette-accent)"
  },
  "colors.colorPalette.site": {
    "value": "var(--colors-color-palette-site)",
    "variable": "var(--colors-color-palette-site)"
  },
  "colors.colorPalette.opaque": {
    "value": "var(--colors-color-palette-opaque)",
    "variable": "var(--colors-color-palette-opaque)"
  },
  "colors.colorPalette.outline": {
    "value": "var(--colors-color-palette-outline)",
    "variable": "var(--colors-color-palette-outline)"
  },
  "colors.colorPalette.published.bg": {
    "value": "var(--colors-color-palette-published-bg)",
    "variable": "var(--colors-color-palette-published-bg)"
  },
  "colors.colorPalette.bg": {
    "value": "var(--colors-color-palette-bg)",
    "variable": "var(--colors-color-palette-bg)"
  },
  "colors.colorPalette.published.fg": {
    "value": "var(--colors-color-palette-published-fg)",
    "variable": "var(--colors-color-palette-published-fg)"
  },
  "colors.colorPalette.fg": {
    "value": "var(--colors-color-palette-fg)",
    "variable": "var(--colors-color-palette-fg)"
  },
  "colors.colorPalette.published.border": {
    "value": "var(--colors-color-palette-published-border)",
    "variable": "var(--colors-color-palette-published-border)"
  },
  "colors.colorPalette.border": {
    "value": "var(--colors-color-palette-border)",
    "variable": "var(--colors-color-palette-border)"
  },
  "colors.colorPalette.draft.bg": {
    "value": "var(--colors-color-palette-draft-bg)",
    "variable": "var(--colors-color-palette-draft-bg)"
  },
  "colors.colorPalette.draft.fg": {
    "value": "var(--colors-color-palette-draft-fg)",
    "variable": "var(--colors-color-palette-draft-fg)"
  },
  "colors.colorPalette.draft.border": {
    "value": "var(--colors-color-palette-draft-border)",
    "variable": "var(--colors-color-palette-draft-border)"
  },
  "colors.colorPalette.review.bg": {
    "value": "var(--colors-color-palette-review-bg)",
    "variable": "var(--colors-color-palette-review-bg)"
  },
  "colors.colorPalette.review.fg": {
    "value": "var(--colors-color-palette-review-fg)",
    "variable": "var(--colors-color-palette-review-fg)"
  },
  "colors.colorPalette.review.border": {
    "value": "var(--colors-color-palette-review-border)",
    "variable": "var(--colors-color-palette-review-border)"
  },
  "colors.colorPalette.unlisted.bg": {
    "value": "var(--colors-color-palette-unlisted-bg)",
    "variable": "var(--colors-color-palette-unlisted-bg)"
  },
  "colors.colorPalette.unlisted.fg": {
    "value": "var(--colors-color-palette-unlisted-fg)",
    "variable": "var(--colors-color-palette-unlisted-fg)"
  },
  "colors.colorPalette.unlisted.border": {
    "value": "var(--colors-color-palette-unlisted-border)",
    "variable": "var(--colors-color-palette-unlisted-border)"
  },
  "colors.colorPalette.text": {
    "value": "var(--colors-color-palette-text)",
    "variable": "var(--colors-color-palette-text)"
  }
}

export function token(path, fallback) {
  return tokens[path]?.value || fallback
}

function tokenVar(path, fallback) {
  return tokens[path]?.variable || fallback
}

token.var = tokenVar