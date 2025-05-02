import React, { memo, useRef, useCallback, useState } from "react";
import Button from "../ui-kit/Button/Button";
import api from "../api";
import { useMutation, useQueryClient } from "@tanstack/react-query";

async function createFolder(name) {
  const { data } = await api.post("/folder", name);
  return data;
}

const CreateFolder = ({ onClose }) => {
  const queryClient = useQueryClient();
  const nameRef = useRef(null);
  const [validForm, setValid] = useState(false);

  const changeInput = () => {
    const name = nameRef.current.value;
    setValid(name.length > 1);
  };

  const mutationCreate = useMutation({
    mutationFn: ({ name }) => createFolder(name),
    onSuccess: () => queryClient.invalidateQueries(["filesList"]),
  });

  const createCallBack = useCallback(
    async (name) => {
      try {
        await mutationCreate.mutateAsync({ name });
      } catch (error) {
        console.error("error", error);
        alert("Failed to create Folder. Please try again");
      }
    },
    [mutationCreate]
  );

  const handleSubmit = (e) => {
    e.preventDefault();
    const name = nameRef.current.value;
    createCallBack({ name });
    onClose();
  };

  return (
    <form className="Dialog" onSubmit={handleSubmit}>
      <div className="bgDialog">
        <div className="titleContainer">
          <h2 className="titleDialog">Create Folder</h2>
          <button
            type="button"
            onClick={onClose}
            className="closeButton"
            aria-label="Close modal"
          >
            &times;
          </button>
        </div>
        <div className="textDialog">
          <div className="mb-4">
            <input
              type="text"
              ref={nameRef}
              onChange={changeInput}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Enter Name"
              autoComplete="current-password"
              autoFocus
              required
            />
          </div>
        </div>

        <div className="footerDialog flex justify-end space-x-3">
          <Button
            variant="outlined"
            onClick={onClose}
            type="button"
            className="px-4 py-2"
          >
            Close
          </Button>
          <Button type="submit" disabled={!validForm} className="px-4 py-2">
            Create
          </Button>
        </div>
      </div>
    </form>
  );
};

export default memo(CreateFolder);
