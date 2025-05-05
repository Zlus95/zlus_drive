import React, { useCallback } from "react";
import Header from "../../ui-kit/Header/Header";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useDialog } from "../../providers/DialogProvider";
import api from "../../api";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import FilesTree from "./FilesTree";

async function deleteFile(id) {
  const { data } = await api.delete(`/file/${id}`);
  return data;
}

async function moveFile(id, parent) {
  try {
    const { data } = await api.patch(`/file/${id}`, { parent });
    return data;
  } catch (error) {
    console.error("Move file error:", error.response?.data || error.message);
    throw error; // Передаём ошибку в mutation
  }
}

async function deleteFolder(id) {
  const { data } = await api.delete(`/folder/${id}`);
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

  const mutationDeleteFolder = useMutation({
    mutationFn: ({ id }) => deleteFolder(id),
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

  const deleteCallBackFolder = useCallback(
    async (id) => {
      try {
        await mutationDeleteFolder.mutateAsync({ id });
      } catch (error) {
        console.error("error", error);
        alert("Failed to delete file Please try again");
      }
    },
    [mutationDeleteFolder]
  );

  const handleShowFile = useCallback(
    (item) => {
      if (!item.isFolder) return () => showDialog(DIALOGS.SHOW_FILE, { item });
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
          onClick: () =>
            item.isFolder
              ? deleteCallBackFolder(item.id)
              : deleteCallBack(item.id),
        });
    },
    [DIALOGS.CONFIRMATION, showDialog, deleteCallBack, deleteCallBackFolder]
  );

  const mutationMove = useMutation({
    mutationFn: ({ id, parent }) => moveFile(id, parent),
    onSuccess: () => queryClient.invalidateQueries(["filesList"]),
    onError: (error) => {
      console.error("Mutation error:", error.response?.data || error.message);
    },
  });

  const handleDrop = async (dragged, targetParent) => {
    if (!dragged.isFolder && !targetParent.isFolder) return;
    if (dragged.isFolder && !targetParent.isFolder) return;

    try {
      await mutationMove.mutateAsync({
        id: dragged.id,
        parent: targetParent.id,
      });
    } catch (error) {
      console.error("Handle drop error:", error);
    }
  };

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="h-full overflow-auto">
        <Header onChangeSort={onChangeSort} sort={sort} />
        <div className="px-2 space-y-4">
          <FilesTree
            files={files}
            handleShowFile={handleShowFile}
            handleDeleteFile={handleDeleteFile}
            onDrop={handleDrop}
          />
        </div>
      </div>
    </DndProvider>
  );
};

export default Home;
