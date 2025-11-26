import { useNavigate } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import { getMe, logout } from "../api";

export default function HomePage() {
  const navigate = useNavigate();
  const observerRef = useRef(null);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const startInterview = () => {
    if (user) {
      navigate("/chat");
    } else {
      navigate("/login");
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      setUser(null);
      navigate("/");
    } catch (err) {
      console.error("Logout failed:", err);
    }
  };

  useEffect(() => {
    // Check if user is logged in
    const checkAuth = async () => {
      try {
        const userData = await getMe();
        setUser(userData);
      } catch (err) {
        setUser(null);
      } finally {
        setLoading(false);
      }
    };
    checkAuth();
  }, []);

  useEffect(() => {
    // Lightweight scroll animation using Intersection Observer
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("visible");
          }
        });
      },
      { threshold: 0.1, rootMargin: "0px 0px -50px 0px" }
    );

    observerRef.current = observer;

    // Use setTimeout to ensure DOM is ready
    setTimeout(() => {
      const elements = document.querySelectorAll(".fade-in-on-scroll");
      elements.forEach((el) => {
        if (el) observer.observe(el);
      });
    }, 100);

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, []);

  const addToRefs = (el) => {
    if (el && observerRef.current) {
      observerRef.current.observe(el);
    }
  };

  return (
    <div className="bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 min-h-screen">
      {/* Navigation */}
      <nav className="bg-white shadow-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="text-4xl">🤖</span>
            <span className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              AI Interviewer
            </span>
          </div>
          <div className="flex items-center gap-8">
            <a href="#features" className="text-gray-600 hover:text-indigo-600 font-medium transition-colors">Features</a>
            <a href="#how-it-works" className="text-gray-600 hover:text-indigo-600 font-medium transition-colors">How It Works</a>
            <a href="#pricing" className="text-gray-600 hover:text-indigo-600 font-medium transition-colors">Pricing</a>
            {!loading && (
              user ? (
                <div className="flex items-center gap-4">
                  <span className="text-gray-700 font-medium">{user.name || user.email}</span>
                  <button
                    onClick={handleLogout}
                    className="px-6 py-2 bg-gray-200 text-gray-700 rounded-full hover:bg-gray-300 transition-all"
                  >
                    Logout
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => navigate("/login")}
                  className="px-6 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-full hover:shadow-lg transition-all"
                >
                  Sign In
                </button>
              )
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="max-w-7xl mx-auto px-6 py-20">
        <div className="flex items-center justify-between gap-12 flex-wrap">
          <div className="flex-1 min-w-[300px] animate-fade-in-up">
            <div className="inline-flex items-center gap-2 bg-indigo-100 text-indigo-700 px-4 py-2 rounded-full text-sm font-medium mb-6">
              <span>✨</span>
              <span>Powered by Advanced AI</span>
            </div>
            <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6 leading-tight">
              Practice Interviews with Your
              <span className="bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent"> AI Companion</span>
            </h1>
            <p className="text-xl text-gray-600 mb-8 leading-relaxed">
              Get personalized feedback, improve your answers, and ace your next interview with our intelligent AI interviewer that adapts to your needs.
            </p>
            <div className="flex gap-4 flex-wrap">
              <button
                onClick={startInterview}
                className="px-8 py-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-full font-semibold text-lg hover:shadow-2xl transition-all transform hover:scale-105"
              >
                {user ? "Start Interview" : "Start Free Interview"}
              </button>
              <button className="px-8 py-4 bg-white text-gray-700 rounded-full font-semibold text-lg hover:shadow-lg transition-all border-2 border-gray-200">
                Watch Demo
              </button>
            </div>
            <div className="flex items-center gap-8 mt-8 text-sm text-gray-500 flex-wrap">
              <div className="flex items-center gap-2">
                <span className="text-green-500">✓</span>
                <span>No credit card required</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-green-500">✓</span>
                <span>Free 7-day trial</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-green-500">✓</span>
                <span>Cancel anytime</span>
              </div>
            </div>
          </div>
          <div className="flex-1 relative min-w-[300px] animate-fade-in-scale" style={{ animationDelay: "0.2s" }}>
            <div className="bg-white rounded-3xl shadow-2xl p-12 text-center">
              <div className="text-9xl mb-6 animate-bounce">🤖</div>
              <h3 className="text-2xl font-bold text-gray-800 mb-2">Meet Your AI Interviewer</h3>
              <p className="text-gray-500">Ready to help you succeed!</p>
              <div className="mt-6 grid grid-cols-2 gap-4 text-sm">
                <div className="bg-indigo-50 rounded-xl p-4">
                  <div className="text-3xl mb-2">10k+</div>
                  <div className="text-gray-600">Interviews Conducted</div>
                </div>
                <div className="bg-purple-50 rounded-xl p-4">
                  <div className="text-3xl mb-2">95%</div>
                  <div className="text-gray-600">Success Rate</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="bg-white py-20">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">Why Choose AI Interviewer?</h2>
            <p className="text-xl text-gray-600">Everything you need to ace your next interview</p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <div
              ref={addToRefs}
              className="fade-in-on-scroll bg-gradient-to-br from-indigo-50 to-purple-50 rounded-2xl p-8 hover:shadow-xl transition-all"
            >
              <div className="w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center mb-6 text-white text-3xl">
                ⚡
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Instant Feedback</h3>
              <p className="text-gray-600">Get real-time feedback on your answers and suggestions for improvement after every response.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-8 hover:shadow-xl transition-all"
              style={{ transitionDelay: "0.1s" }}
            >
              <div className="w-16 h-16 bg-purple-600 rounded-2xl flex items-center justify-center mb-6 text-white text-3xl">
                💬
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Natural Conversation</h3>
              <p className="text-gray-600">Experience realistic interview scenarios with our AI that understands context and adapts to you.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll bg-gradient-to-br from-pink-50 to-orange-50 rounded-2xl p-8 hover:shadow-xl transition-all"
              style={{ transitionDelay: "0.2s" }}
            >
              <div className="w-16 h-16 bg-pink-600 rounded-2xl flex items-center justify-center mb-6 text-white text-3xl">
                📈
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Track Progress</h3>
              <p className="text-gray-600">Monitor your improvement over time with detailed analytics and performance insights.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll bg-gradient-to-br from-blue-50 to-indigo-50 rounded-2xl p-8 hover:shadow-xl transition-all"
              style={{ transitionDelay: "0.3s" }}
            >
              <div className="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center mb-6 text-white text-3xl">
                ⏰
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">24/7 Availability</h3>
              <p className="text-gray-600">Practice anytime, anywhere. Your AI interviewer is always ready when you are.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-8 hover:shadow-xl transition-all"
              style={{ transitionDelay: "0.4s" }}
            >
              <div className="w-16 h-16 bg-green-600 rounded-2xl flex items-center justify-center mb-6 text-white text-3xl">
                🛡️
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Private & Secure</h3>
              <p className="text-gray-600">Your interviews are completely confidential. We prioritize your privacy and data security.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll bg-gradient-to-br from-yellow-50 to-orange-50 rounded-2xl p-8 hover:shadow-xl transition-all"
              style={{ transitionDelay: "0.5s" }}
            >
              <div className="w-16 h-16 bg-yellow-600 rounded-2xl flex items-center justify-center mb-6 text-white text-3xl">
                ✨
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Personalized</h3>
              <p className="text-gray-600">Tailored questions based on your industry, role, and experience level for maximum relevance.</p>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="py-20">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">How It Works</h2>
            <p className="text-xl text-gray-600">Get started in three simple steps</p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div
              ref={addToRefs}
              className="fade-in-on-scroll text-center"
            >
              <div className="w-24 h-24 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-full flex items-center justify-center text-white text-4xl font-bold mx-auto mb-6 shadow-xl">
                1
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Create Your Profile</h3>
              <p className="text-gray-600">Tell us about your background, the role you're applying for, and your experience level.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll text-center"
              style={{ transitionDelay: "0.2s" }}
            >
              <div className="w-24 h-24 bg-gradient-to-br from-purple-600 to-pink-600 rounded-full flex items-center justify-center text-white text-4xl font-bold mx-auto mb-6 shadow-xl">
                2
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Start Practicing</h3>
              <p className="text-gray-600">Chat with your AI interviewer and answer questions tailored to your specific needs.</p>
            </div>
            <div
              ref={addToRefs}
              className="fade-in-on-scroll text-center"
              style={{ transitionDelay: "0.4s" }}
            >
              <div className="w-24 h-24 bg-gradient-to-br from-pink-600 to-orange-600 rounded-full flex items-center justify-center text-white text-4xl font-bold mx-auto mb-6 shadow-xl">
                3
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-4">Get Better</h3>
              <p className="text-gray-600">Review feedback, track your progress, and improve with every practice session.</p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-gradient-to-r from-indigo-600 to-purple-600 py-20">
        <div className="max-w-4xl mx-auto px-6 text-center">
          <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">Ready to Ace Your Interview?</h2>
          <p className="text-xl text-indigo-100 mb-8">Join thousands of successful candidates who practiced with AI Interviewer</p>
          <button
            onClick={startInterview}
            className="px-12 py-5 bg-white text-indigo-600 rounded-full font-bold text-xl hover:shadow-2xl transition-all transform hover:scale-105"
          >
            Start Your Free Trial Now
          </button>
          <p className="text-indigo-100 mt-6">No credit card required • 7-day free trial • Cancel anytime</p>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 py-12">
        <div className="max-w-7xl mx-auto px-6">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <span className="text-3xl">🤖</span>
                <span className="text-xl font-bold text-white">AI Interviewer</span>
              </div>
              <p className="text-sm">Empowering candidates with AI-powered interview practice.</p>
            </div>
            <div>
              <h4 className="text-white font-semibold mb-4">Product</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#features" className="hover:text-white transition-colors">Features</a></li>
                <li><a href="#pricing" className="hover:text-white transition-colors">Pricing</a></li>
                <li><a href="#" className="hover:text-white transition-colors">FAQ</a></li>
              </ul>
            </div>
            <div>
              <h4 className="text-white font-semibold mb-4">Company</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#" className="hover:text-white transition-colors">About Us</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Blog</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Careers</a></li>
              </ul>
            </div>
            <div>
              <h4 className="text-white font-semibold mb-4">Legal</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#" className="hover:text-white transition-colors">Privacy Policy</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Terms of Service</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contact</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 pt-8 text-center text-sm">
            <p>© 2025 AI Interviewer. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
