import React, { useCallback, useRef } from "react";
import Button from "../../ui-kit/Button/Button";
import api from "../../api";
import { useMutation, useQueryClient } from "@tanstack/react-query";

async function uploadFile(file) {
  const formData = new FormData();
  formData.append("file", file); // Добавляем файл в FormData

  // Важно: передаём formData в запрос
  const { data } = await api.post("/upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data", // Обязательный заголовок для файлов
    },
  });
  return data;
}

const Upload = () => {
  const fileInputRef = useRef(null);

  const queryClient = useQueryClient();

  const mutationUpload = useMutation({
    mutationFn: (file) => uploadFile(file),
    onSuccess: () => queryClient.invalidateQueries(["filesList"]),
  });

  const handleFileChange = useCallback(
    async (e) => {
      const file = e.target.files?.[0];
      if (!file) return;

      try {
        await mutationUpload.mutateAsync(file);
        e.target.value = null;
      } catch (error) {
        console.error("Upload error:", error);
        alert("Failed to upload file. Please try again");
      }
    },
    [mutationUpload]
  );

  const handleClick = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  return (
    <label className="cursor-pointer">
      <input type="file" ref={fileInputRef} onChange={handleFileChange} />
      <Button onClick={handleClick}>Upload</Button>
    </label>
  );
};

export default Upload;
