import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Auth from "../pages/Auth/Auth";
import Registration from "../pages/Auth/Registration";

const ProtectedRoute = ({ element, title }) => {
  const navigate = useNavigate();
  const token = localStorage.getItem("token");

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }
  }, [token, navigate]);

  // return <Layout title={title}>{element}</Layout>;
};

const Navigation = () => {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Auth />} />
        <Route path="/register" element={<Registration />} />
      </Routes>
    </Router>
  );
};

export default Navigation;
