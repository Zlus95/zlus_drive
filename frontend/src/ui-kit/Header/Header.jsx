import React, { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import { useDialog } from "../../providers/DialogProvider";
import api from "../../api";
import Button from "../Button/Button";

function useGetUser() {
  return useQuery({
    queryKey: ["user"],
    staleTime: 1000 * 60 * 5,
    queryFn: async () => {
      const response = await api.get("/user");
      return response.data;
    },
  });
}

export const formatStorage = (bytes) => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
};

const formatStorageUsage = (used, total) => {
  return `${formatStorage(used)} / ${formatStorage(total)}`;
};

const Header = () => {
  const { isLoading, data } = useGetUser();
  const { showDialog, DIALOGS } = useDialog();
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const logout = useCallback(() => {
    localStorage.removeItem("token");
    queryClient.clear();
    navigate("/login");
  }, [navigate, queryClient]);

  if (isLoading || !data) {
    return (
      <header className="bg-gray-800 text-white shadow-md">
        <div className="container mx-auto px-4 py-3 flex items-center justify-between">
          <div className="animate-pulse h-8 w-32 bg-gray-700 rounded"></div>
          <div className="animate-pulse h-8 w-24 bg-gray-700 rounded"></div>
        </div>
      </header>
    );
  }

  const storagePercentage = Math.min(
    Math.round((data.usedStorage / data.storageLimit) * 100),
    100
  );

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
          <div className="flex flex-col items-end">
            <div className="text-sm mb-1">
              {formatStorageUsage(data.usedStorage, data.storageLimit)}
            </div>
            <div className="w-48 bg-gray-700 rounded-full h-2">
              <div
                className={`h-2 rounded-full ${
                  storagePercentage > 90 ? "bg-red-500" : "bg-blue-500"
                }`}
                style={{ width: `${storagePercentage}%` }}
              ></div>
            </div>
          </div>

          <div className="flex items-center space-x-3">
            <div className="text-right">
              <div className="text-sm font-medium capitalize">
                {data.name} {data.lastName}
              </div>
              <div className="text-xs text-gray-400 truncate max-w-[120px]">
                {data.email}
              </div>
            </div>

            <div className="relative">
              <div className="w-10 h-10 uppercase rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                {data.name.charAt(0)}
                {data.lastName.charAt(0)}
              </div>
              {data.usedStorage / data.storageLimit > 0.9 && (
                <span className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full border-2 border-gray-800"></span>
              )}
            </div>
          </div>
          <Button
            onClick={() =>
              showDialog(DIALOGS.CONFIRMATION, {
                text: "Are you sure you want to log out?",
                title: "Log out",
                submitButton: "Log out",
                onClick: logout,
              })
            }
          >
            Logout
          </Button>
        </div>
      </div>
    </header>
  );
};

export default Header;
