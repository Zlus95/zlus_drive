import React, { memo, useState } from "react";
import Button from "../ui-kit/Button/Button";
import api from "../api";
import { useMutation, useQueryClient } from "@tanstack/react-query";

async function changeLimit(storageLimitBytes) {
  const { data } = await api.patch("/user/update", {
    storageLimit: storageLimitBytes,
  });
  return data;
}

const options = [
  { value: 5, label: "5 GB" },
  { value: 10, label: "10 GB" },
  { value: 15, label: "15 GB" },
];

const StorageLimit = ({ onClose }) => {
  const queryClient = useQueryClient();
  const [selectedLimit, setSelectedLimit] = useState(5);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const mutationChange = useMutation({
    mutationFn: ({ storageLimit }) => changeLimit(storageLimit),
    onSuccess: () => {
      queryClient.invalidateQueries(["user"]);
      onClose();
    },
    onError: (error) => {
      console.error("Error updating limit:", error);
      alert("Failed to update storage limit. Please try again");
    },
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const storageLimitBytes = selectedLimit * 1024 * 1024 * 1024;
      await mutationChange.mutateAsync({ storageLimit: storageLimitBytes });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form className="Dialog" onSubmit={handleSubmit}>
      <div className="bgDialog">
        <div className="titleContainer">
          <h2 className="titleDialog">Change Storage Limit</h2>
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
            <select
              value={selectedLimit}
              onChange={(e) => setSelectedLimit(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              {options.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="footerDialog flex justify-end space-x-3">
          <Button
            variant="outlined"
            onClick={onClose}
            type="button"
            className="px-4 py-2"
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={isSubmitting}
            className="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700"
          >
            {isSubmitting ? "Saving..." : "Save Changes"}
          </Button>
        </div>
      </div>
    </form>
  );
};

export default memo(StorageLimit);
