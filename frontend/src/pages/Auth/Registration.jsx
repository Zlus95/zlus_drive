import React, { useState, useRef } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "../../api";

async function registration(user) {
  const { data } = await api.post("/register", user);
  return data;
}

const Registration = () => {
  const navigate = useNavigate();
  const emailRef = useRef(null);
  const passwordRef = useRef(null);
  const nameRef = useRef(null);
  const lastNameRef = useRef(null);
  const [validForm, setValid] = useState(false);

  const changeInput = () => {
    const email = emailRef.current.value;
    const password = passwordRef.current.value;
    const name = nameRef.current.value;
    const lastName = lastNameRef.current.value;
    setValid(
      email.trim() !== "" &&
        password.trim() !== "" &&
        name.trim() !== "" &&
        lastName.trim() !== ""
    );
  };

  const mutationReg = useMutation({
    mutationFn: registration,
    onSuccess: () => navigate("/login"),
    onError: ({ response }) => alert(response.data.error),
  });

  const handlerSubmit = (event) => {
    event.preventDefault();
    if (
      emailRef.current &&
      passwordRef.current &&
      nameRef.current &&
      lastNameRef.current
    ) {
      const email = emailRef.current.value;
      const password = passwordRef.current.value;
      const name = nameRef.current.value;
      const lastName = lastNameRef.current.value;
      mutationReg.mutate({ email, password, name, lastName });
    }
  };

  const inputFields = [
    {
      id: "Name",
      label: "Name:",
      placeholder: "Name",
      ref: nameRef,
    },
    {
      id: "Last Name",
      label: "Last Name:",
      placeholder: "Last Name",
      ref: lastNameRef,
    },
    {
      id: "email",
      label: "Email:",
      placeholder: "login",
      ref: emailRef,
    },
    {
      id: "password",
      label: "Password:",
      placeholder: "password",
      type: "password",
      ref: passwordRef,
    },
  ];

  return (
    <div className="flex justify-center items-center h-screen w-full">
      <form className="w-full max-w-xs" onSubmit={handlerSubmit}>
        <div className="h-96 w-72 bg-black/50 border-2 border-neutral-400 rounded">
          <p className="text-primary flex justify-center pt-4 text-lg">Drive</p>
          <div className="flex flex-col gap-1 px-2">
            {inputFields.map((field) => (
              <React.Fragment key={field.id}>
                <label htmlFor={field.id} className="text-white/50">
                  {field.label}
                </label>
                <input
                  placeholder={field.placeholder}
                  id={field.id}
                  onChange={changeInput}
                  ref={field.ref}
                  type={field.type || "text"}
                />
              </React.Fragment>
            ))}
          </div>
          <div className="flex justify-center pt-4">
            <button
              className={validForm ? "text-primary" : "text-primary/50"}
              disabled={!validForm}
            >
              Sign up
            </button>
          </div>
          <div className="p-4 flex gap-2">
            <p className="text-white/50">Already have an account?</p>
            <Link to="/login" className="text-primary">
              Sign in
            </Link>
          </div>
        </div>
      </form>
    </div>
  );
};

export default Registration;
