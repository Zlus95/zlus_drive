import React, { memo } from "react";
import Button from "../ui-kit/Button/Button";

const CreateFolder = ({ onClose }) => {
  return (
    <form className="Dialog">
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
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Room Password
            </label>
            <input
              type="text"
              //   ref={passwordRef}
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
          <Button type="submit" className="px-4 py-2">
            Create
          </Button>
        </div>
      </div>
    </form>
  );
};

export default memo(CreateFolder);
