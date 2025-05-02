import FileItem from "./FileItem";

const FilesTree = ({ files, handleShowFile, handleDeleteFile, onDrop }) => {
  const roots = files.filter((f) => !f.parent);
  const getChildren = (id) => files.filter((f) => f.parent === id);

  const renderNode = (item, depth = 0) => (
    <div key={item.id}>
      <FileItem
        item={item}
        handleShowFile={handleShowFile}
        handleDeleteFile={handleDeleteFile}
        onDrop={onDrop}
        depth={depth}
      />
      {item.isFolder &&
        getChildren(item.id).map((child) => renderNode(child, depth + 1))}
    </div>
  );

  return <>{roots.map((root) => renderNode(root))}</>;
};

export default FilesTree;
