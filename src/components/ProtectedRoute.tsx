import React, { useEffect, useState } from "react";
import { getMe } from "../api";


export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
    const [ok, setOk] = useState<boolean | null>(null);

    useEffect(() => {
        getMe().then((data) => setOk(!!data)).catch(() => setOk(false));
    }, [])

    if (ok === null) return <div>Loading...</div>;
    if (!ok) {
        window.location.href = "/";
        return null;
    }
    return <>{children}</>;
}