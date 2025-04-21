import React, { useState, useRef } from "react";
import api from "../../api";
import { Link, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";

async function login(user) {
  const { data } = await api.post("/login", user);
  localStorage.setItem("token", data.token);
  return data;
}

const Auth = () => {
  const navigate = useNavigate();
  const emailRef = useRef(null);
  const passwordRef = useRef(null);
  const [validForm, setValid] = useState(false);

  const changeInput = () => {
    const email = emailRef.current.value;
    const password = passwordRef.current.value;
    setValid(email.trim() !== "" && password.trim() !== "");
  };

  const mutationLogin = useMutation({
    mutationFn: login,
    onSuccess: () => navigate("/"),
    onError: ({ response }) => alert(response.data),
  });

  const handlerSubmit = (event) => {
    event.preventDefault();
    if (emailRef.current && passwordRef.current) {
      const email = emailRef.current.value;
      const password = passwordRef.current.value;
      mutationLogin.mutate({ email, password });
    }
  };

  return (
    <div className="flex justify-center items-center h-screen w-full">
      <form className="w-full max-w-xs" onSubmit={handlerSubmit}>
        <div className="h-80 w-72 bg-black/50 border-2 border-neutral-400 rounded">
          <div className="text-primary flex justify-center pt-4 text-lg">
            Drive
          </div>
          <div className="flex flex-col gap-2 p-4">
            <label htmlFor="email" className="text-white/50">
              Email:
            </label>
            <input
              placeholder="login"
              id="email"
              onChange={changeInput}
              ref={emailRef}
            />
            <label htmlFor="password" className="text-white/50">
              Password:
            </label>
            <input
              placeholder="password"
              type="password"
              id="password"
              ref={passwordRef}
              onChange={changeInput}
            />
          </div>
          <div className="flex justify-center pt-4">
            <button
              className={validForm ? "text-primary" : "text-primary/50"}
              disabled={!validForm}
            >
              Sign in
            </button>
          </div>
          <div className="p-4 flex gap-2">
            <p className="text-white/50">Don't have an account?</p>
            <Link to="/register" className="text-primary">
              Sign up
            </Link>
          </div>
        </div>
      </form>
    </div>
  );
};

export default Auth;
