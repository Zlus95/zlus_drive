import React, { createContext, useContext, useState, useCallback } from "react";
import Confirmation from "../dialogs/Confirmation";
import ShowFile from "../dialogs/ShowFile";

const Context = createContext();

export const DIALOGS = {
  CONFIRMATION: "CONFIRMATION",
  SHOW_FILE: "SHOW_FILE",
};

const COMPONENTS = {
  CONFIRMATION: Confirmation,
  SHOW_FILE: ShowFile,
};

export const DialogProvider = ({ children }) => {
  const [dialogs, setDialogs] = useState([]);

  const showDialog = useCallback((type, initialData = {}) => {
    setDialogs([
      {
        id: `${type}-${Date.now()}`,
        Component: COMPONENTS[type],
        initialData,
      },
    ]);
  }, []);

  const addDialog = useCallback((type, initialData = {}) => {
    setDialogs((prev) => [
      ...prev,
      {
        id: `${type}-${Date.now()}`,
        Component: COMPONENTS[type],
        initialData,
      },
    ]);
  }, []);

  const hideDialog = useCallback(() => {
    setDialogs((prev) => prev.slice(0, -1));
  }, []);

  return (
    <Context.Provider
      value={{
        showDialog,
        addDialog,
        hideDialog,
        DIALOGS,
      }}
    >
      {children}
      {dialogs.map((dialog) => (
        <dialog.Component
          key={dialog.id}
          {...dialog.initialData}
          onClose={hideDialog}
          isOpen={true}
        />
      ))}
    </Context.Provider>
  );
};

export const useDialog = () => {
  const context = useContext(Context);
  if (!context) {
    throw new Error("useDialog must be used within a DialogProvider");
  }
  return context;
};
