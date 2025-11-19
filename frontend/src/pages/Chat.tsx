import React, { useState, useRef, useEffect } from "react";

interface Message {
    id: number;
    text: string;
    sender: "user" | "bot";
    timestamp: Date;
}

export default function Chat() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: 1,
            text: "test",
            sender: "bot",
            timestamp: new Date(),
        },
    ]);
    const [inputValue, setInputValue] = useState("");
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const handleSend = (e: React.FormEvent) => {
        e.preventDefault();
        if (!inputValue.trim()) return;

        const newMessage: Message = {
            id: messages.length + 1,
            text: inputValue.trim(),
            sender: "user",
            timestamp: new Date(),
        };

        setMessages([...messages, newMessage]);
        setInputValue("");

        // Simple echo response (you can modify this logic)
        setTimeout(() => {
            const botResponse: Message = {
                id: messages.length + 2,
                text: `You said: ${inputValue.trim()}`,
                sender: "bot",
                timestamp: new Date(),
            };
            setMessages((prev) => [...prev, botResponse]);
        }, 500);
    };

    const formatTime = (date: Date) => {
        return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    };

    return (
        <div className="min-h-screen bg-[#0a0a0a] flex flex-col">
            {/* Header */}
            <header className="bg-[#1a1a1a] border-b border-white/10 px-4 py-3 flex items-center justify-between">
                <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                        C
                    </div>
                    <div>
                        <h1 className="text-white font-semibold text-lg">Chat</h1>
                        <p className="text-xs text-gray-400">Online</p>
                    </div>
                </div>
                <button className="text-gray-400 hover:text-white transition">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                    >
                        <circle cx="12" cy="12" r="1" />
                        <circle cx="19" cy="12" r="1" />
                        <circle cx="5" cy="12" r="1" />
                    </svg>
                </button>
            </header>

            {/* Messages Container */}
            <div className="flex-1 overflow-y-auto px-4 py-6 space-y-4">
                {messages.map((message) => (
                    <div
                        key={message.id}
                        className={`flex ${message.sender === "user" ? "justify-end" : "justify-start"
                            }`}
                    >
                        <div
                            className={`max-w-[75%] rounded-2xl px-4 py-2 ${message.sender === "user"
                                    ? "bg-blue-500 text-white rounded-tr-sm"
                                    : "bg-[#1f1f1f] text-gray-100 rounded-tl-sm"
                                }`}
                        >
                            <p className="text-sm leading-relaxed">{message.text}</p>
                            <p
                                className={`text-xs mt-1 ${message.sender === "user"
                                        ? "text-blue-100"
                                        : "text-gray-400"
                                    }`}
                            >
                                {formatTime(message.timestamp)}
                            </p>
                        </div>
                    </div>
                ))}
                <div ref={messagesEndRef} />
            </div>

            {/* Input Area */}
            <div className="bg-[#1a1a1a] border-t border-white/10 px-4 py-3">
                <form onSubmit={handleSend} className="flex items-end gap-2">
                    <div className="flex-1 relative">
                        <input
                            type="text"
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            placeholder="Type a message..."
                            className="w-full bg-[#0f0f0f] text-white placeholder-gray-500 rounded-full px-4 py-3 pr-12 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition"
                        />
                        <button
                            type="button"
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white transition"
                        >
                            <svg
                                width="20"
                                height="20"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                strokeWidth="2"
                            >
                                <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
                                <circle cx="12" cy="10" r="3" />
                            </svg>
                        </button>
                    </div>
                    <button
                        type="submit"
                        disabled={!inputValue.trim()}
                        className="bg-blue-500 text-white rounded-full p-3 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition"
                    >
                        <svg
                            width="20"
                            height="20"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                        >
                            <line x1="22" y1="2" x2="11" y2="13" />
                            <polygon points="22 2 15 22 11 13 2 9 22 2" />
                        </svg>
                    </button>
                </form>
            </div>
        </div>
    );
}

