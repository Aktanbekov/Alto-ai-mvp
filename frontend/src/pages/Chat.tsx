import React, { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { sendChatMessage, logout } from "../api";

interface Message {
  id: number;
  text: string;
  sender: "user" | "ai";
  timestamp?: Date;
}

type EmojiState = "happy" | "thinking" | "excited" | "listening" | "impressed" | "proud";

const emojiStates: Record<EmojiState, string> = {
  happy: "🤖",
  thinking: "🤔",
  excited: "🤩",
  listening: "👂",
  impressed: "😮",
  proud: "🎉",
};

const statusTexts: Record<EmojiState, string> = {
  happy: "Active & Listening",
  thinking: "Analyzing Response...",
  excited: "Great Answer!",
  listening: "Waiting for Response",
  impressed: "Impressive!",
  proud: "Excellent Work!",
};


export default function Chat() {
  const navigate = useNavigate();
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 1,
      sender: "ai",
      text: "Hello from AI! I'm here to chat with you. How can I help you today?",
      timestamp: new Date(),
    },
  ]);
  const [inputValue, setInputValue] = useState("");
  const [emojiState, setEmojiState] = useState<EmojiState>("happy");
  const [isTyping, setIsTyping] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const timeoutRefs = useRef<NodeJS.Timeout[]>([]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, isTyping]);

  // Cleanup timeouts on unmount
  useEffect(() => {
    return () => {
      timeoutRefs.current.forEach((timeout) => clearTimeout(timeout));
      timeoutRefs.current = [];
    };
  }, []);

  const progress = Math.min((messages.length / 15) * 100, 100);
  const messageCount = messages.length;
  const timeElapsed = "12m"; // You can calculate this based on start time

  const changeEmoji = (state: EmojiState) => {
    setEmojiState(state);
  };

  const handleSend = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    if (!inputValue.trim()) return;

    // Clear any existing timeouts to prevent double updates
    timeoutRefs.current.forEach((timeout) => clearTimeout(timeout));
    timeoutRefs.current = [];

    const messageText = inputValue.trim();

    // Add user message
    const newUserMessage: Message = {
      id: messages.length + 1,
      text: messageText,
      sender: "user",
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, newUserMessage]);
    setInputValue("");
    
    // Immediately change emoji to thinking
    changeEmoji("thinking");
    setIsTyping(true);

    try {
      // Build conversation history for GPT
      const conversationHistory = [
        ...messages.map((msg) => ({
          role: msg.sender === "user" ? "user" : "assistant",
          content: msg.text,
        })),
        {
          role: "user",
          content: messageText,
        },
      ];

      // Send to GPT API
      const responseText = await sendChatMessage(conversationHistory);
      
      setIsTyping(false);
      changeEmoji("excited");
      
      const newAiMessage: Message = {
        id: messages.length + 2,
        text: responseText,
        sender: "ai",
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, newAiMessage]);

      // Change back to happy after showing excited
      const happyTimeout = setTimeout(() => {
        changeEmoji("happy");
      }, 1500);
      timeoutRefs.current.push(happyTimeout);
    } catch (error) {
      setIsTyping(false);
      changeEmoji("happy");
      
      const errorMessage: Message = {
        id: messages.length + 2,
        text: `Sorry, I encountered an error: ${error instanceof Error ? error.message : "Unknown error"}`,
        sender: "ai",
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleEmojiClick = () => {
    const emojiElement = document.getElementById("aiEmoji");
    if (emojiElement) {
      emojiElement.style.transform = "scale(1.2) rotate(5deg)";
      setTimeout(() => {
        emojiElement.style.transform = "scale(1)";
      }, 200);
    }
  };

  const handleEndSession = async () => {
    if (window.confirm("Are you sure you want to end this interview session?")) {
      try {
        await logout();
      } catch (err) {
        console.error("Logout error:", err);
      }
      navigate("/");
    }
  };

  const handleLogout = async () => {
    if (window.confirm("Are you sure you want to logout?")) {
      try {
        await logout();
        navigate("/");
      } catch (err) {
        console.error("Logout error:", err);
        navigate("/");
      }
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50">
      {/* Top Navigation Bar */}
      <nav className="bg-white shadow-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-2 sm:py-3 flex items-center justify-between">
          <div className="flex items-center gap-2 sm:gap-3">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="lg:hidden p-2 text-gray-600 hover:text-indigo-600 transition-colors min-w-[44px] min-h-[44px] flex items-center justify-center"
              aria-label="Toggle sidebar"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <span className="text-2xl sm:text-3xl">🤖</span>
            <span className="text-lg sm:text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              AI Interviewer
            </span>
          </div>
          <div className="flex items-center gap-2 sm:gap-4">
            <button className="hidden sm:block px-3 sm:px-4 py-2 text-gray-600 hover:text-indigo-600 font-medium transition-colors text-sm sm:text-base min-h-[44px]">
              Save Interview
            </button>
            <button
              onClick={handleEndSession}
              className="hidden sm:block px-3 sm:px-4 py-2 text-gray-600 hover:text-indigo-600 font-medium transition-colors text-sm sm:text-base min-h-[44px]"
            >
              End Session
            </button>
            <button
              onClick={handleLogout}
              className="px-3 sm:px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-full hover:shadow-lg transition-all text-sm sm:text-base min-h-[44px]"
            >
              <span className="hidden sm:inline">Logout</span>
              <span className="sm:hidden">Out</span>
            </button>
          </div>
        </div>
      </nav>

      {/* Main Chat Container */}
      <div className="max-w-7xl mx-auto p-3 sm:p-4 md:p-6">
        <div className="flex gap-4 sm:gap-6 items-start relative">
          {/* Large AI Character Sidebar */}
          <div className={`${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0 fixed lg:sticky top-16 lg:top-24 left-0 z-40 lg:z-auto h-[calc(100vh-4rem)] lg:h-auto overflow-y-auto flex-shrink-0 w-80 sm:w-96 bg-white rounded-r-3xl lg:rounded-3xl shadow-2xl p-6 sm:p-8 flex flex-col items-center transition-transform duration-300 ease-in-out`}>
            {/* Close button for mobile */}
            <button
              onClick={() => setSidebarOpen(false)}
              className="lg:hidden absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600 transition-colors min-w-[44px] min-h-[44px] flex items-center justify-center"
              aria-label="Close sidebar"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            {/* AI Emoji */}
            <div
              id="aiEmoji"
              onClick={handleEmojiClick}
              key={emojiState}
              className="text-6xl sm:text-7xl md:text-8xl lg:text-9xl mb-4 sm:mb-6 cursor-pointer animate-fade-in-scale"
            >
              {emojiStates[emojiState]}
            </div>

            {/* AI Info */}
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-800 mb-2">AI Interviewer</h2>
            <p className="text-gray-500 text-center text-xs sm:text-sm mb-4 sm:mb-6">
              I'm here to chat with you and learn about your experiences!
            </p>

            {/* Status Indicator */}
            <div className="flex items-center gap-2 mb-4 sm:mb-6 px-3 sm:px-4 py-1.5 sm:py-2 bg-green-50 rounded-full">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span className="text-xs sm:text-sm text-green-700 font-medium">{statusTexts[emojiState]}</span>
            </div>

            {/* Progress Section */}
            <div className="w-full mb-4 sm:mb-6">
              <div className="flex justify-between text-xs text-gray-500 mb-2">
                <span>Interview Progress</span>
                <span>{Math.round(progress)}%</span>
              </div>
              <div className="bg-gray-200 rounded-full h-2 sm:h-3 w-full overflow-hidden">
                <div
                  className="bg-gradient-to-r from-green-400 to-emerald-500 rounded-full h-2 sm:h-3 transition-all duration-500"
                  style={{ width: `${progress}%` }}
                ></div>
              </div>
            </div>

            {/* Stats */}
            <div className="w-full grid grid-cols-2 gap-2 sm:gap-3 mb-4 sm:mb-6">
              <div className="bg-indigo-50 rounded-xl p-3 sm:p-4 text-center">
                <div className="text-xl sm:text-2xl font-bold text-indigo-600">{messageCount}</div>
                <div className="text-xs text-gray-600">Messages</div>
              </div>
              <div className="bg-purple-50 rounded-xl p-3 sm:p-4 text-center">
                <div className="text-xl sm:text-2xl font-bold text-purple-600">{timeElapsed}</div>
                <div className="text-xs text-gray-600">Time</div>
              </div>
            </div>

            {/* Tips Section */}
            <div className="w-full bg-gradient-to-br from-yellow-50 to-orange-50 rounded-xl p-3 sm:p-4 border-2 border-yellow-200">
              <div className="flex items-start gap-2">
                <span className="text-lg sm:text-xl">💡</span>
                <div>
                  <h4 className="font-semibold text-gray-800 text-xs sm:text-sm mb-1">Pro Tip</h4>
                  <p className="text-xs text-gray-600">Take your time to think before answering. Quality over speed!</p>
                </div>
              </div>
            </div>
          </div>

          {/* Overlay for mobile sidebar */}
          {sidebarOpen && (
            <div
              className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-30"
              onClick={() => setSidebarOpen(false)}
            />
          )}

          {/* Chat Area */}
          <div className="flex-1 w-full bg-white rounded-2xl sm:rounded-3xl shadow-2xl overflow-hidden flex flex-col" style={{ height: "calc(100vh - 80px)", minHeight: "500px" }}>
            {/* Chat Header */}
            <div className="bg-gradient-to-r from-indigo-600 to-purple-600 p-4 sm:p-6 text-white">
              <div className="flex items-center justify-between">
                <div className="flex-1 min-w-0">
                  <h1 className="text-lg sm:text-xl md:text-2xl font-bold truncate">Interview Session</h1>
                  <p className="text-indigo-100 text-xs sm:text-sm truncate">Software Engineer Position • Technical Round</p>
                </div>
                <div className="flex gap-1 sm:gap-2 ml-2">
                  <button
                    className="p-2 bg-white bg-opacity-20 rounded-lg hover:bg-opacity-30 transition-all min-w-[44px] min-h-[44px] flex items-center justify-center"
                    title="Hints"
                  >
                    💡
                  </button>
                  <button
                    className="p-2 bg-white bg-opacity-20 rounded-lg hover:bg-opacity-30 transition-all min-w-[44px] min-h-[44px] flex items-center justify-center"
                    title="Settings"
                  >
                    ⚙️
                  </button>
                </div>
              </div>
            </div>

            {/* Messages Container */}
            <div className="flex-1 overflow-y-auto p-3 sm:p-4 md:p-6 space-y-3 sm:space-y-4 bg-gradient-to-b from-gray-50 to-white">
              {messages.map((message, index) => (
                <div
                  key={message.id}
                  className={`flex ${message.sender === "user" ? "justify-end" : "justify-start"} animate-fade-in-up`}
                  style={{ animationDelay: `${index * 0.05}s` }}
                >
                  <div
                    className={`max-w-[85%] sm:max-w-[80%] rounded-xl sm:rounded-2xl px-4 sm:px-6 py-3 sm:py-4 shadow-md transition-all hover:shadow-lg ${
                      message.sender === "ai"
                        ? "bg-white text-gray-800 border-2 border-indigo-100"
                        : "bg-gradient-to-r from-indigo-600 to-purple-600 text-white"
                    }`}
                  >
                    {message.sender === "ai" && (
                      <div className="flex items-center gap-2 mb-1 sm:mb-2">
                        <span className="text-lg sm:text-xl">🤖</span>
                        <span className="text-xs font-semibold text-indigo-600">AI Interviewer</span>
                      </div>
                    )}
                    <p className="text-xs sm:text-sm leading-relaxed break-words">{message.text}</p>
                  </div>
                </div>
              ))}

              {/* Typing Indicator */}
              {isTyping && (
                <div className="flex items-center gap-2 sm:gap-3 max-w-[85%] sm:max-w-[80%]">
                  <div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full flex items-center justify-center text-lg sm:text-2xl bg-indigo-100 flex-shrink-0">
                    🤖
                  </div>
                  <div className="bg-white border-2 border-indigo-100 rounded-xl sm:rounded-2xl px-4 sm:px-6 py-2 sm:py-3">
                    <div className="flex gap-1">
                      <div
                        className="w-2 h-2 bg-indigo-400 rounded-full animate-bounce"
                        style={{ animationDelay: "0s" }}
                      ></div>
                      <div
                        className="w-2 h-2 bg-indigo-400 rounded-full animate-bounce"
                        style={{ animationDelay: "0.2s" }}
                      ></div>
                      <div
                        className="w-2 h-2 bg-indigo-400 rounded-full animate-bounce"
                        style={{ animationDelay: "0.4s" }}
                      ></div>
                    </div>
                  </div>
                </div>
              )}

              <div ref={messagesEndRef} />
            </div>

            {/* Input Area */}
            <div className="border-t-2 border-gray-100 p-3 sm:p-4 md:p-6 bg-white">
              <form onSubmit={handleSend} className="flex gap-2 sm:gap-3">
                <input
                  ref={inputRef}
                  type="text"
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onKeyPress={handleKeyPress}
                  placeholder="Type your response here..."
                  className="flex-1 px-4 sm:px-6 py-3 sm:py-4 border-2 border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm sm:text-base transition-all text-gray-900 min-h-[44px]"
                />
                <button
                  type="submit"
                  className="bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-4 sm:px-6 md:px-8 py-3 sm:py-4 rounded-full hover:from-indigo-700 hover:to-purple-700 transition-all shadow-lg hover:shadow-xl transform hover:scale-105 font-medium text-sm sm:text-base min-h-[44px] min-w-[80px] sm:min-w-[100px]"
                >
                  <span className="hidden sm:inline">Send ➤</span>
                  <span className="sm:hidden">➤</span>
                </button>
              </form>
              <div className="hidden sm:flex items-center gap-4 mt-2 sm:mt-3 text-xs text-gray-400">
                <span>💡 Press Enter to send</span>
                <span>•</span>
                <span>Shift + Enter for new line</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
