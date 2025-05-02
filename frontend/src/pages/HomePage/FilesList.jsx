import React, { memo } from "react";
import Button from "../../ui-kit/Button/Button";
import FolderIcon from "../../ui-kit/icons/FolderIcon";
import FileIcon from "../../ui-kit/icons/FileIcon";
import { formatStorage } from "../../ui-kit/Header/Header";

const FilesList = ({ files, handleShowFile, handleDeleteFile }) => (
  <>
    {files.map((item) => (
      <div
        key={item.id}
        className="flex items-center gap-4 p-3 border-b border-gray-200 hover:bg-primary"
      >
        <div
          className="flex-shrink-0 cursor-pointer"
          onClick={handleShowFile(item)}
        >
          {item.isFolder ? (
            <FolderIcon />
          ) : (
            <FileIcon type={item.mimeType} name={item.name} />
          )}
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-white/50 truncate">
            {item.name}
          </p>
          <p className="text-xs text-gray-500">
            {new Date(item.createdAt).toLocaleDateString()}{" "}
            {!item.isFolder && formatStorage(item.size)}
          </p>
        </div>
        <div className="flex-shrink-0">
          <Button variant="error" onClick={handleDeleteFile(item)}>
            Delete
          </Button>
        </div>
      </div>
    ))}
  </>
);

export default memo(FilesList);
