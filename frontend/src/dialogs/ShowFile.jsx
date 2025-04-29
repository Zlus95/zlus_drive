import React, { memo } from "react";

const ShowFile = ({ onClose, item }) => {
  return (
    <div className="fixed inset-0 z-50 bg-black bg-opacity-90">
      <div className="h-full flex flex-col">
        <div className="flex justify-between items-center p-4 bg-gray-900 text-white">
          <h2 className="text-xl font-medium truncate max-w-[80%]">
            {item.name}
          </h2>
          <button
            onClick={onClose}
            className="p-2 rounded-full hover:bg-gray-700"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        {item.mimeType.includes("image") ? (
          <img
            src={`${process.env.REACT_APP_API_URL}/${item.path}`}
            alt={item.name}
            className="w-full h-full rounded"
          />
        ) : (
          <iframe
            src={`${process.env.REACT_APP_API_URL}/${item.path}`}
            width="100%"
            height="100%"
            style={{ border: "none" }}
            title={item.name}
          ></iframe>
        )}
      </div>
    </div>
  );
};

export default memo(ShowFile);
