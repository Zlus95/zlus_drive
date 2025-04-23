import React from "react";

const Loader = ({ text }) => (
  <div className=" flex items-center justify-center bg-opacity-50 z-50">
    <div className="flex flex-col items-center">
      <div className="w-20 h-20 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
      {text && <p className="mt-4 text-white font-medium">{text}</p>}
    </div>
  </div>
);

export default Loader;
