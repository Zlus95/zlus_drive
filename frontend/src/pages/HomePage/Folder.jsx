import React from "react";
import Button from "../../ui-kit/Button/Button";
import { useDialog } from "../../providers/DialogProvider";

const Folder = () => {
  const { showDialog, DIALOGS } = useDialog();

  return (
    <Button onClick={() => showDialog(DIALOGS.CREATE_FOLDER)}>
      Create Folder
    </Button>
  );
};

export default Folder;
