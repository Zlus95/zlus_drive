import React, { useCallback } from "react";
import Header from "../../ui-kit/Header/Header";
import Button from "../../ui-kit/Button/Button";
import { formatStorage } from "../../ui-kit/Header/Header";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useDialog } from "../../providers/DialogProvider";
import api from "../../api";
import FileIcon from "./FileIcon";

async function deleteFile(id) {
  const { data } = await api.delete(`/file/${id}`);
  return data;
}

// async function showFile(id) {
//   const { data } = await api.get(`/file/${id}`);
//   return data;
// }

const Home = ({ data }) => {
  const { data: files } = data;
  const queryClient = useQueryClient();
  const { showDialog, DIALOGS } = useDialog();

  const mutationDelete = useMutation({
    mutationFn: ({ id }) => deleteFile(id),
    onSuccess: () => queryClient.invalidateQueries(["filesList"]),
  });

  const deleteCallBack = useCallback(
    async (id) => {
      try {
        await mutationDelete.mutateAsync({ id });
      } catch (error) {
        console.error("error", error);
        alert("Failed to delete file Please try again");
      }
    },
    [mutationDelete]
  );

  // const mutationShow = useMutation({
  //   mutationFn: ({ id }) => showFile(id),
  // });

  // const showCallBack = useCallback(
  //   async (id) => {
  //     try {
  //       await mutationShow.mutateAsync({ id });
  //     } catch (error) {
  //       console.error("error", error);
  //       alert("Failed to show file Please try again");
  //     }
  //   },
  //   [mutationShow]
  // );

  return (
    <div className="h-full overflow-auto">
      <Header />
      <div className="px-2 space-y-4">
        {files.map((item) => (
          <div
            key={item.id}
            className="flex items-center gap-4 p-3 border-b border-gray-200 hover:bg-primary"
          >
            <div className="flex-shrink-0">
              <FileIcon type={item.mimeType} name={item.name} />
              {/* <img
                src={`${process.env.REACT_APP_API_URL}/${item.path}`}
                alt={item.name}
                className="w-10 h-10 object-cover rounded"
              /> */}
            </div>

            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-white/50 truncate">
                {item.name}
              </p>
              <p className="text-xs text-gray-500">
                {new Date(item.createdAt).toLocaleDateString()} •{" "}
                {formatStorage(item.size)}
              </p>
            </div>

            <div className="flex-shrink-0">
              <Button
                variant="error"
                onClick={() =>
                  showDialog(DIALOGS.CONFIRMATION, {
                    text: `Are you sure you want to delete the file ${item.name}?`,
                    title: "Delete file",
                    submitButton: "delete",
                    onClick: () => deleteCallBack(item.id),
                  })
                }
              >
                Delete
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Home;
