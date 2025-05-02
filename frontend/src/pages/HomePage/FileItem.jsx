import { useDrag, useDrop } from "react-dnd";
import React, { memo } from "react";
import Button from "../../ui-kit/Button/Button";
import FolderIcon from "../../ui-kit/icons/FolderIcon";
import FileIcon from "../../ui-kit/icons/FileIcon";
import { formatStorage } from "../../ui-kit/Header/Header";
import clsx from "clsx";

const FileItem = ({
  item,
  handleShowFile,
  handleDeleteFile,
  onDrop,
  depth = 0,
}) => {
  const [{ isDragging }, drag] = useDrag(() => ({
    type: "FILE",
    item: { id: item.id },
    collect: (monitor) => ({
      isDragging: !!monitor.isDragging(),
    }),
  }));

  const [{ isOver }, drop] = useDrop(() => ({
    accept: "FILE",
    drop: (draggedItem) => {
      if (item.isFolder && draggedItem.id !== item.id) {
        onDrop(draggedItem.id, item.id);
      }
    },
    collect: (monitor) => ({
      isOver: !!monitor.isOver(),
    }),
  }));

  return (
    <div
      ref={(node) => drag(drop(node))}
      className={`flex items-center gap-4 p-3 border-b border-gray-200 hover:bg-primary 
        ${isDragging ? "opacity-50" : ""} 
        ${isOver ? "bg-blue-50" : ""}`}
      style={{ marginLeft: `${depth * 16}px` }}
    >
      <div
        className={clsx(
          "flex-shrink-0",
          item.isFolder ? "cursor-auto" : "cursor-pointer"
        )}
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
          {new Date(item.createdAt).toLocaleDateString()}
          {!item.isFolder && formatStorage(item.size)}
        </p>
      </div>
      <div className="flex-shrink-0">
        <Button variant="error" onClick={handleDeleteFile(item)}>
          Delete
        </Button>
      </div>
    </div>
  );
};

export default memo(FileItem);
