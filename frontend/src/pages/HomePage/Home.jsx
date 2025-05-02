import React, { useCallback } from "react";
import Header from "../../ui-kit/Header/Header";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useDialog } from "../../providers/DialogProvider";
import api from "../../api";
import FilesList from "./FilesList";

async function deleteFile(id) {
  const { data } = await api.delete(`/file/${id}`);
  return data;
}

const Home = ({ data, onChangeSort, sort }) => {
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

  const handleShowFile = useCallback(
    (item) => {
      return () => showDialog(DIALOGS.SHOW_FILE, { item });
    },
    [DIALOGS.SHOW_FILE, showDialog]
  );

  const handleDeleteFile = useCallback(
    (item) => {
      return () =>
        showDialog(DIALOGS.CONFIRMATION, {
          text: `Are you sure you want to delete the file ${item.name}?`,
          title: "Delete file",
          submitButton: "delete",
          onClick: () => deleteCallBack(item.id),
        });
    },
    [DIALOGS.CONFIRMATION, showDialog, deleteCallBack]
  );

  return (
    <div className="h-full overflow-auto">
      <Header onChangeSort={onChangeSort} sort={sort} />
      <div className="px-2 space-y-4">
        <FilesList
          files={files}
          handleShowFile={handleShowFile}
          handleDeleteFile={handleDeleteFile}
        />
      </div>
    </div>
  );
};

export default Home;
