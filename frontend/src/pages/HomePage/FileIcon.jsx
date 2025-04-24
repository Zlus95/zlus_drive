import React, { memo } from "react";

const FileIcon = ({ type, size = 24, className = "", name }) => {
  const iconStyle = {
    width: size,
    height: size,
  };

  const getIcon = () => {
    const fileType = type?.split("/")[0] || "unknown";
    const fileExtension = type?.split("/")[1] || "";

    if (fileType === "image") {
      return (
        <svg viewBox="0 0 24 24" style={iconStyle} className={className}>
          <path
            fill="currentColor"
            d="M8.5 13.5l2.5 3 3.5-4.5 4.5 6H5m16 1V5a2 2 0 00-2-2H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2z"
          />
        </svg>
      );
    }

    if (fileExtension === "pdf" || type === "application/pdf") {
      return (
        <svg
          viewBox="0 0 24 24"
          style={iconStyle}
          className={`text-red-500 ${className}`}
        >
          <path
            fill="currentColor"
            d="M5 4v16h14V8h-4V4H5zm0-2h9l5 5v11a2 2 0 01-2 2H5a2 2 0 01-2-2V4a2 2 0 012-2zm2 8h2v4H7v-4zm4-2h2v6h-2V8zm4 1h2v5h-2V9z"
          />
        </svg>
      );
    }

    if (
      fileExtension.includes("word") ||
      type.includes("msword") ||
      type.includes("wordprocessingml") ||
      name.includes("docx")
    ) {
      return (
        <svg
          viewBox="0 0 24 24"
          style={iconStyle}
          className={`text-blue-500 ${className}`}
        >
          <path
            fill="currentColor"
            d="M6 2h8l6 6v12a2 2 0 01-2 2H6a2 2 0 01-2-2V4a2 2 0 012-2zm8 0v6h6l-6-6zm-2 9h-1v5h1v-5zm-2 0H9v5h1v-5zm-2 0H7v5h1v-5zm6 0h-1v5h1v-5z"
          />
        </svg>
      );
    }

    if (fileExtension.includes("excel") || type.includes("spreadsheetml")) {
      return (
        <svg
          viewBox="0 0 24 24"
          style={iconStyle}
          className={`text-green-600 ${className}`}
        >
          <path
            fill="currentColor"
            d="M6 2h8l6 6v12a2 2 0 01-2 2H6a2 2 0 01-2-2V4a2 2 0 012-2zm8 0v6h6l-6-6zm-6 8l2 3-2 3h3l-2 3 2 3h-3l-2-3-2 3H5l2-3-2-3h3l2-3-2-3h3z"
          />
        </svg>
      );
    }

    if (
      fileExtension.includes("zip") ||
      fileExtension.includes("rar") ||
      fileExtension.includes("7z")
    ) {
      return (
        <svg
          viewBox="0 0 24 24"
          style={iconStyle}
          className={`text-yellow-500 ${className}`}
        >
          <path
            fill="currentColor"
            d="M14 2h-4v2h4V2zm-2 4h-2v2h2V6zm2 0h-2v2h2V6zm-2 4h-2v2h2v-2zm2 0h-2v2h2v-2zm-2 4h-2v2h2v-2zm2 0h-2v2h2v-2zm-2 4h-2v2h2v-2zm2 0h-2v2h2v-2zM6 4v16h12V4H6z"
          />
        </svg>
      );
    }

    if (fileType === "audio") {
      return (
        <svg
          viewBox="0 0 24 24"
          style={iconStyle}
          className={`text-purple-500 ${className}`}
        >
          <path
            fill="currentColor"
            d="M12 3v9.28a4.39 4.39 0 00-1.5-.28C8.01 12 6 14.01 6 16.5S8.01 21 10.5 21c2.31 0 4.2-1.75 4.45-4H15V6h4V3h-7z"
          />
        </svg>
      );
    }

    if (fileType === "video") {
      return (
        <svg
          viewBox="0 0 24 24"
          style={iconStyle}
          className={`text-red-500 ${className}`}
        >
          <path
            fill="currentColor"
            d="M17 10.5V7a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h12a1 1 0 001-1v-3.5l4 4v-11l-4 4z"
          />
        </svg>
      );
    }

    return (
      <svg viewBox="0 0 24 24" style={iconStyle} className={className}>
        <path
          fill="currentColor"
          d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zM6 4h7v5h5v11H6V4zm8 18v-2h4v2h-4z"
        />
      </svg>
    );
  };

  return getIcon();
};

export default memo(FileIcon);
