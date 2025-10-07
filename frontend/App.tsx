import React from "react";
import { Routes, Route } from "react-router-dom";
import Login from "./src/pages/Login";
import Dashboard from "./src/pages/Dashboard";
import ProtectedRoute from "./src/components/ProtectedRoute";
import AuthPage from "./src/pages/AuthPage";

export default function App() {
    return (
        <Routes>
            <Route path="/" element={<Login />} />
            <Route path="/auth" element={<AuthPage />} />
            <Route
                path="/dashboard"
                element={
                    <ProtectedRoute>
                        <Dashboard />
                    </ProtectedRoute>
                }
            />
        </Routes>
    );
}