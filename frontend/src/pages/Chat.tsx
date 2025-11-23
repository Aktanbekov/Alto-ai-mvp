import React, { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";

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

const aiResponses = [
  "That's interesting! Tell me more about that experience.",
  "Great answer! Can you elaborate on the challenges you faced?",
  "I see. How did that experience shape your approach to problem-solving?",
  "Excellent! What did you learn from that situation?",
  "That's impressive! How do you think that applies to this role?",
];

export default function Chat() {
  const navigate = useNavigate();
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 1,
      sender: "ai",
      text: "Hello! I'm your AI interviewer today. I'm here to learn more about you and your experiences. Let's start with an easy question: Can you tell me a bit about yourself?",
      timestamp: new Date(),
    },
  ]);
  const [inputValue, setInputValue] = useState("");
  const [emojiState, setEmojiState] = useState<EmojiState>("happy");
  const [isTyping, setIsTyping] = useState(false);
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

  const handleSend = (e?: React.FormEvent) => {
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

    // Simulate AI response after delay
    const typingTimeout = setTimeout(() => {
      setIsTyping(false);
      changeEmoji("excited");
      
      const randomResponse = aiResponses[Math.floor(Math.random() * aiResponses.length)];
      const newAiMessage: Message = {
        id: messages.length + 2,
        text: randomResponse,
        sender: "ai",
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, newAiMessage]);

      // Change back to happy after showing excited
      const happyTimeout = setTimeout(() => {
        changeEmoji("happy");
      }, 1500);
      timeoutRefs.current.push(happyTimeout);
    }, 1500);
    timeoutRefs.current.push(typingTimeout);
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

  const handleEndSession = () => {
    if (window.confirm("Are you sure you want to end this interview session?")) {
      navigate("/");
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50">
      {/* Top Navigation Bar */}
      <nav className="bg-white shadow-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="text-3xl">🤖</span>
            <span className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              AI Interviewer
            </span>
          </div>
          <div className="flex items-center gap-4">
            <button className="px-4 py-2 text-gray-600 hover:text-indigo-600 font-medium transition-colors">
              Save Interview
            </button>
            <button
              onClick={handleEndSession}
              className="px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-full hover:shadow-lg transition-all"
            >
              End Session
            </button>
          </div>
        </div>
      </nav>

      {/* Main Chat Container */}
      <div className="max-w-7xl mx-auto p-6">
        <div className="flex gap-6 items-start">
          {/* Large AI Character Sidebar */}
          <div className="flex-shrink-0 w-96 bg-white rounded-3xl shadow-2xl p-8 flex flex-col items-center sticky top-24">
            {/* AI Emoji */}
            <div
              id="aiEmoji"
              onClick={handleEmojiClick}
              key={emojiState}
              className="text-9xl mb-6 cursor-pointer animate-fade-in-scale"
            >
              {emojiStates[emojiState]}
            </div>

            {/* AI Info */}
            <h2 className="text-3xl font-bold text-gray-800 mb-2">AI Interviewer</h2>
            <p className="text-gray-500 text-center text-sm mb-6">
              I'm here to chat with you and learn about your experiences!
            </p>

            {/* Status Indicator */}
            <div className="flex items-center gap-2 mb-6 px-4 py-2 bg-green-50 rounded-full">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span className="text-sm text-green-700 font-medium">{statusTexts[emojiState]}</span>
            </div>

            {/* Progress Section */}
            <div className="w-full mb-6">
              <div className="flex justify-between text-xs text-gray-500 mb-2">
                <span>Interview Progress</span>
                <span>{Math.round(progress)}%</span>
              </div>
              <div className="bg-gray-200 rounded-full h-3 w-full overflow-hidden">
                <div
                  className="bg-gradient-to-r from-green-400 to-emerald-500 rounded-full h-3 transition-all duration-500"
                  style={{ width: `${progress}%` }}
                ></div>
              </div>
            </div>

            {/* Stats */}
            <div className="w-full grid grid-cols-2 gap-3 mb-6">
              <div className="bg-indigo-50 rounded-xl p-4 text-center">
                <div className="text-2xl font-bold text-indigo-600">{messageCount}</div>
                <div className="text-xs text-gray-600">Messages</div>
              </div>
              <div className="bg-purple-50 rounded-xl p-4 text-center">
                <div className="text-2xl font-bold text-purple-600">{timeElapsed}</div>
                <div className="text-xs text-gray-600">Time</div>
              </div>
            </div>

            {/* Tips Section */}
            <div className="w-full bg-gradient-to-br from-yellow-50 to-orange-50 rounded-xl p-4 border-2 border-yellow-200">
              <div className="flex items-start gap-2">
                <span className="text-xl">💡</span>
                <div>
                  <h4 className="font-semibold text-gray-800 text-sm mb-1">Pro Tip</h4>
                  <p className="text-xs text-gray-600">Take your time to think before answering. Quality over speed!</p>
                </div>
              </div>
            </div>
          </div>

          {/* Chat Area */}
          <div className="flex-1 bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col" style={{ height: "calc(100vh - 140px)" }}>
            {/* Chat Header */}
            <div className="bg-gradient-to-r from-indigo-600 to-purple-600 p-6 text-white">
              <div className="flex items-center justify-between">
                <div>
                  <h1 className="text-2xl font-bold">Interview Session</h1>
                  <p className="text-indigo-100 text-sm">Software Engineer Position • Technical Round</p>
                </div>
                <div className="flex gap-2">
                  <button
                    className="p-2 bg-white bg-opacity-20 rounded-lg hover:bg-opacity-30 transition-all"
                    title="Hints"
                  >
                    💡
                  </button>
                  <button
                    className="p-2 bg-white bg-opacity-20 rounded-lg hover:bg-opacity-30 transition-all"
                    title="Settings"
                  >
                    ⚙️
                  </button>
                </div>
              </div>
            </div>

            {/* Messages Container */}
            <div className="flex-1 overflow-y-auto p-6 space-y-4 bg-gradient-to-b from-gray-50 to-white">
              {messages.map((message, index) => (
                <div
                  key={message.id}
                  className={`flex ${message.sender === "user" ? "justify-end" : "justify-start"} animate-fade-in-up`}
                  style={{ animationDelay: `${index * 0.05}s` }}
                >
                  <div
                    className={`max-w-[80%] rounded-2xl px-6 py-4 shadow-md transition-all hover:shadow-lg ${
                      message.sender === "ai"
                        ? "bg-white text-gray-800 border-2 border-indigo-100"
                        : "bg-gradient-to-r from-indigo-600 to-purple-600 text-white"
                    }`}
                  >
                    {message.sender === "ai" && (
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-xl">🤖</span>
                        <span className="text-xs font-semibold text-indigo-600">AI Interviewer</span>
                      </div>
                    )}
                    <p className="text-sm leading-relaxed">{message.text}</p>
                  </div>
                </div>
              ))}

              {/* Typing Indicator */}
              {isTyping && (
                <div className="flex items-center gap-3 max-w-[80%]">
                  <div className="w-10 h-10 rounded-full flex items-center justify-center text-2xl bg-indigo-100">
                    🤖
                  </div>
                  <div className="bg-white border-2 border-indigo-100 rounded-2xl px-6 py-3">
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
            <div className="border-t-2 border-gray-100 p-6 bg-white">
              <form onSubmit={handleSend} className="flex gap-3">
                <input
                  ref={inputRef}
                  type="text"
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onKeyPress={handleKeyPress}
                  placeholder="Type your response here..."
                  className="flex-1 px-6 py-4 border-2 border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm transition-all text-gray-900"
                />
                <button
                  type="submit"
                  className="bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-8 py-4 rounded-full hover:from-indigo-700 hover:to-purple-700 transition-all shadow-lg hover:shadow-xl transform hover:scale-105 font-medium"
                >
                  Send ➤
                </button>
              </form>
              <div className="flex items-center gap-4 mt-3 text-xs text-gray-400">
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
