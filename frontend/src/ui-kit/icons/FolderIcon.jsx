import React, { memo } from "react";

const FolderIcon = ({ className = "w-6 h-6", color = "#EF4444" }) => (
  <svg
    className={className}
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M3 5C3 3.89543 3.89543 3 5 3H8.17157C8.70201 3 9.21071 3.21071 9.58579 3.58579L10.4142 4.41421C10.7893 4.78929 11.298 5 11.8284 5H19C20.1046 5 21 5.89543 21 7V19C21 20.1046 20.1046 21 19 21H5C3.89543 21 3 20.1046 3 19V5Z"
      fill={color}
    />
  </svg>
);
export default memo(FolderIcon);
