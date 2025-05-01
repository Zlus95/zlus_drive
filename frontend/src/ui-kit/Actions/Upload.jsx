import React, { useCallback, useRef } from "react";
import Button from "../Button/Button";
import api from "../../api";
import { useMutation, useQueryClient } from "@tanstack/react-query";

async function uploadFile(file) {
  const formData = new FormData();
  formData.append("file", file);

  const { data } = await api.post("/upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
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
    <>
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        className="hidden"
      />
      <Button onClick={handleClick}>Upload</Button>
    </>
  );
};

export default Upload;
