import React from 'react';
import { ArrowRight } from 'lucide-react';

const LevelSelection = () => {
  const levels = [
    {
      id: 'easy',
      name: 'Easy',
      questions: 4,
      color: 'from-green-400 to-green-500',
      bgColor: 'bg-green-50',
      borderColor: 'border-green-200',
      textColor: 'text-green-700',
      description: 'Perfect for first-time practice',
      icon: '🌱'
    },
    {
      id: 'medium',
      name: 'Medium',
      questions: 7,
      color: 'from-blue-400 to-blue-500',
      bgColor: 'bg-blue-50',
      borderColor: 'border-blue-200',
      textColor: 'text-blue-700',
      description: 'Build your confidence',
      icon: '🎯'
    },
    {
      id: 'hard',
      name: 'Hard',
      questions: 12,
      color: 'from-purple-400 to-purple-500',
      bgColor: 'bg-purple-50',
      borderColor: 'border-purple-200',
      textColor: 'text-purple-700',
      description: 'Master the interview',
      icon: '🏆'
    }
  ];

  const handleLevelSelect = (level) => {
    console.log('Selected level:', level);
    // Add your navigation logic here
  };

  return (
    <div className="w-full h-screen flex flex-col bg-gradient-to-b from-green-50 to-blue-50">
      {/* Header */}
      <div className="bg-gradient-to-r from-green-500 to-blue-500 text-white p-6 shadow-lg">
        <div className="max-w-4xl mx-auto text-center">
          <div className="text-5xl mb-3">🎓</div>
          <h1 className="text-3xl font-bold mb-2">F1 Visa Interview Practice</h1>
          <p className="text-green-100">Choose your difficulty level to begin</p>
        </div>
      </div>

      {/* Level Cards */}
      <div className="flex-1 flex items-center justify-center p-6">
        <div className="max-w-5xl w-full grid grid-cols-1 md:grid-cols-3 gap-6">
          {levels.map((level) => (
            <div
              key={level.id}
              className={`${level.bgColor} border-2 ${level.borderColor} rounded-3xl p-8 cursor-pointer transform transition-all hover:scale-105 hover:shadow-2xl`}
              onClick={() => handleLevelSelect(level)}
            >
              <div className="text-center">
                <div className="text-6xl mb-4">{level.icon}</div>
                <h2 className="text-2xl font-bold text-gray-800 mb-2">{level.name}</h2>
                <div className={`inline-block px-4 py-2 rounded-full bg-white ${level.textColor} font-semibold mb-4`}>
                  {level.questions} Questions
                </div>
                <p className="text-gray-600 mb-6">{level.description}</p>
                <button className={`w-full bg-gradient-to-r ${level.color} text-white py-3 rounded-full font-semibold flex items-center justify-center gap-2 hover:shadow-lg transition-all`}>
                  Start Practice
                  <ArrowRight size={20} />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Footer Tip */}
      <div className="bg-white border-t p-4 text-center text-gray-500 text-sm">
        💡 Start with Easy if this is your first time practicing
      </div>
    </div>
  );
};

export default LevelSelection;