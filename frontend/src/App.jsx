import { useEffect, useState } from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { getMe } from "./api";

function Login() {
  const loginWithGoogle = () => {
    window.location.href = "http://localhost:8080/auth/google";
  };
  return (
    <div style={{display:"grid",placeItems:"center",height:"100vh",gap:16}}>
      <h1>Alto AI</h1>
      <button onClick={loginWithGoogle}>Continue with Google</button>
    </div>
  );
}

function ProtectedRoute({ children }) {
  const [ok, setOk] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    getMe().then((d) => setOk(!!d)).catch(() => setOk(false));
  }, []);

  useEffect(() => {
    if (ok === false) navigate("/");
  }, [ok, navigate]);

  if (ok === null) return <div style={{padding:24}}>Loading…</div>;
  return children;
}

function Dashboard() {
  const [session, setSession] = useState(null);
  useEffect(() => { getMe().then(setSession); }, []);
  return (
    <div style={{padding:24}}>
      <h2>Dashboard</h2>
      <pre>{JSON.stringify(session, null, 2)}</pre>
      <a href="http://localhost:8080/logout">Logout</a>
    </div>
  );
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
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
