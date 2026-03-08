import { createContext } from "react";

// TODO: Wrap plugin modal with context so we can easily mark the modal as non-dismissable when uploading without having to thread the isUploading state through all the components
const context = createContext({
  // Mark the dismissable in modal as false if this is true
  isUploading: false,
});
