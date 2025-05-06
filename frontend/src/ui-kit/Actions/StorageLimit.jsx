import React from "react";
import Button from "../Button/Button";
import { useDialog } from "../../providers/DialogProvider";

const StorageLimit = () => {
  const { showDialog, DIALOGS } = useDialog();

  return (
    <Button onClick={() => showDialog(DIALOGS.STORAGE_LIMIT)}>
      Change Storage Limit
    </Button>
  );
};

export default StorageLimit;
