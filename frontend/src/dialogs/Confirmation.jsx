import React, { memo } from "react";
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
          <button variant="outlined" onClick={onClose}>
            Close
          </button>
          <button
            onClick={() => {
              onClick && onClick();
              onClose();
            }}
          >
            {submitButton || "Submit"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default memo(Confirmation);
