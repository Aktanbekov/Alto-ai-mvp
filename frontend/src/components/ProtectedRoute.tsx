import React, { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { getMe } from "../api";
import { setupTokenRefresh, stopTokenRefresh } from "../utils/tokenRefresh";

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
    const [ok, setOk] = useState<boolean | null>(null);

    useEffect(() => {
        getMe().then((data) => {
            setOk(!!data);
            if (data) {
                // Setup token refresh when authenticated
                setupTokenRefresh();
            }
        }).catch(() => {
            setOk(false);
            stopTokenRefresh();
        });

        // Cleanup on unmount
        return () => {
            stopTokenRefresh();
        };
    }, [])

    if (ok === null) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 flex items-center justify-center">
                <div className="text-center">
                    <div className="text-6xl mb-4 animate-bounce">🤖</div>
                    <div className="text-gray-600">Loading...</div>
                </div>
            </div>
        );
    }
    if (!ok) {
        return <Navigate to="/login" replace />;
    }
    return <>{children}</>;
}