import React from "react";

export default function Login() {
    const LoginWithGoogle = () => {
        window.location.href = "http://localhost:8080/auth/google";
    };

    return (
        <div style={{display:"grid",placeItems:"center",height:"100vh"}}>
            <div>
            <h1>Alto AI - Login</h1>
            <button onClick={LoginWithGoogle}>Continue with Google</button>
            </div>
        </div>
    );
}