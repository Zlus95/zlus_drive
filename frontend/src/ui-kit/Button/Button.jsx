import React from "react";
import clsx from "clsx";

import "./Button.css";

const DefaultSize = "xs";
const DefaultVariant = "contained";

const Button = ({
  children,
  size = DefaultSize,
  variant = DefaultVariant,
  disabled = false,
  fullWidth = false,
  onClick,
  type = "button",
}) => {
  const sizes = ["xs", "sm", "md", "lg"];
  const variants = ["regular", "outlined", "contained", "error", "edit"];

  const className = clsx(
    "Button",
    size && sizes.includes(size) ? size : DefaultSize,
    variant && variants.includes(variant) ? variant : DefaultVariant,
    disabled && "disabled",
    fullWidth && "w-full"
  );

  return (
    <button
      className={className}
      type={type}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  );
};

export default Button;
