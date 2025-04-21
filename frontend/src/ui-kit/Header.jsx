import React, { memo } from "react";

const Header = () => {
  const formatStorage = (bytes) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    if (bytes < 1024 * 1024 * 1024)
      return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
  };

  //   const storagePercentage = Math.min(
  //     Math.round((user.usedStorage / user.storageLimit) * 100),
  //     100
  //   );

  return (
    <header className="bg-gray-800 text-white shadow-md">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <svg
            className="w-8 h-8 text-blue-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"
            />
          </svg>
          <span className="text-xl font-bold">CloudDrive</span>
        </div>

        <div className="flex items-center space-x-6">
          <div className="hidden md:block w-48">
            <div className="flex justify-between text-sm mb-1">
              {/* <span>Storage: {formatStorage(user.usedStorage)}</span>
                <span>{formatStorage(user.storageLimit)}</span> */}
            </div>
            <div className="w-full bg-gray-700 rounded-full h-2">
              {/* <div
                  className={`h-2 rounded-full ${
                    storagePercentage > 90 ? "bg-red-500" : "bg-blue-500"
                  }`}
                  style={{ width: `${storagePercentage}%` }}
                ></div> */}
            </div>
          </div>

          <div className="flex items-center space-x-3">
            <div className="relative">
              <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                {/* {user.name.charAt(0)}
                  {user.lastName.charAt(0)} */}
              </div>
              {/* {user.usedStorage / user.storageLimit > 0.9 && (
                  <span className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full border-2 border-gray-800"></span>
                )} */}
            </div>

            {/* Выпадающее меню */}
            <div className="hidden md:block">
              <div className="text-sm font-medium">
                {/* {user.name} {user.lastName} */}
              </div>
              <div className="text-xs text-gray-400 truncate max-w-xs">
                {/* {user.email} */}
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default memo(Header);
