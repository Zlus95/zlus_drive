import React from "react";

const ErrorText = ({ error }) => (
  <div className="flex justify-center items-center h-full text-red-500 text-2xl">
    Error: {error.message}
  </div>
);

export default ErrorText;
