import React, { useEffect, useState } from "react";
import { getMe } from "../api";

interface User {
    email: string;
    name: string;
    picture: string | null;
}

export default function Dashboard() {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        getMe()
            .then(data => {
                if (!data) {
                    throw new Error("Failed to fetch user data");
                }
                setUser(data);
            })
            .catch(err => {
                setError(err.message);
                console.error("Dashboard error:", err);
            })
            .finally(() => setLoading(false));
    }, []);

    const handleLogout = () => {
        document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        window.location.href = "/";
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-[#0b0d10] text-slate-200 flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-sky-500"></div>
            </div>
        );
    }

    if (error || !user) {
        return (
            <div className="min-h-screen bg-[#0b0d10] text-slate-200 flex flex-col items-center justify-center gap-4">
                <p className="text-red-400">Failed to load dashboard</p>
                <button
                    onClick={() => window.location.reload()}
                    className="px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600 transition"
                >
                    Try Again
                </button>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#0b0d10] text-slate-200">
            {/* Header */}
            <header className="bg-[#121418] border-b border-white/10">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <svg width="34" height="34" viewBox="0 0 24 24" className="text-sky-400">
                            <path fill="currentColor" d="M12 2l4.2 2.5v5L12 12l-4.2-2.5v-5z"/>
                            <path fill="currentColor" opacity=".6" d="M7.8 9.5 12 12v5l-4.2-2.5zM16.2 9.5 12 12v5l4.2-2.5z"/>
                        </svg>
                        <span className="text-xl font-semibold tracking-wide">Alto</span>
                    </div>
                    <button
                        onClick={handleLogout}
                        className="px-4 py-2 text-sm font-medium text-slate-300 hover:text-white transition"
                    >
                        Logout
                    </button>
                </div>
            </header>

            {/* Main Content */}
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="bg-[#121418]/95 rounded-2xl shadow-xl ring-1 ring-white/5 p-6">
                    <div className="flex items-center gap-4">
                        {user.picture && (
                            <img
                                src={user.picture}
                                alt={user.name}
                                className="w-16 h-16 rounded-full ring-2 ring-sky-500/20"
                                referrerPolicy="no-referrer"
                            />
                        )}
                        <div>
                            <h1 className="text-2xl font-semibold text-white">
                                Welcome, {user.name}!
                            </h1>
                            <p className="text-slate-400 mt-1">{user.email}</p>
                        </div>
                    </div>

                    {/* Dashboard Content */}
                    <div className="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {/* Stats Card */}
                        <div className="bg-[#0f1115] rounded-xl p-6 ring-1 ring-white/5">
                            <h3 className="text-lg font-medium text-slate-300">Your Stats</h3>
                            <p className="mt-2 text-sm text-slate-400">
                                View your activity and statistics here.
                            </p>
                        </div>

                        {/* Settings Card */}
                        <div className="bg-[#0f1115] rounded-xl p-6 ring-1 ring-white/5">
                            <h3 className="text-lg font-medium text-slate-300">Settings</h3>
                            <p className="mt-2 text-sm text-slate-400">
                                Configure your account preferences.
                            </p>
                        </div>

                        {/* Help Card */}
                        <div className="bg-[#0f1115] rounded-xl p-6 ring-1 ring-white/5">
                            <h3 className="text-lg font-medium text-slate-300">Help & Support</h3>
                            <p className="mt-2 text-sm text-slate-400">
                                Get help or contact support.
                            </p>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}
