import React from "react";
import Button from "../../ui-kit/Button/Button";

const SortButton = ({ sort, onChangeSort }) => (
  <Button onClick={onChangeSort}>
    {sort ? (
      <>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
          <path d="M20 6h-8l-2-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm0 12H4V8h16v10z" />
        </svg>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
          <path d="M10 17l5-5-5-5v10z" />
        </svg>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
          <path d="M14 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V8l-6-6zM6 20V4h7v5h5v11H6z" />
        </svg>
      </>
    ) : (
      <>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
          <path d="M14 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V8l-6-6zM6 20V4h7v5h5v11H6z" />
        </svg>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
          <path d="M10 17l5-5-5-5v10z" />
        </svg>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
          <path d="M20 6h-8l-2-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm0 12H4V8h16v10z" />
        </svg>
      </>
    )}
  </Button>
);

export default SortButton;
