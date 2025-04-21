import React, { memo } from "react";
import Button from "../ui-kit/Button/Button";
import "./Dialog.css";

const Confirmation = ({ text, title, onClose, submitButton, onClick }) => {
  return (
    <div className="Dialog">
      <div className="bgDialog">
        <div className="titleContainer">
          <h2 className="titleDialog">{title}</h2>
          <button onClick={onClose} className="closeButton">
            &times;
          </button>
        </div>
        <div className="textDialog">{text}</div>
        <div className="footerDialog">
          <Button variant="outlined" onClick={onClose}>
            Close
          </Button>
          <Button
            onClick={() => {
              onClick && onClick();
              onClose();
            }}
          >
            {submitButton || "Submit"}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default memo(Confirmation);
