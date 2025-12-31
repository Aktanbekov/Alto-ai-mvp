import React from "react";

interface AnalysisScores {
  migration_intent: number;
  goal_understanding: number;
  answer_length: number;
  total_score: number;
}

interface FeedbackByCriterion {
  migration_intent: string;
  goal_understanding: string;
  answer_length: string;
}

interface StructuredFeedback {
  overall: string;
  by_criterion: FeedbackByCriterion;
  improvements: string[];
}

interface ChatAnalysis {
  scores: AnalysisScores;
  classification: string;
  feedback: StructuredFeedback;
}

interface AnswerFeedbackCardProps {
  analysis: ChatAnalysis;
  questionNumber?: number;
}

const AnswerFeedbackCard: React.FC<AnswerFeedbackCardProps> = ({
  analysis,
  questionNumber,
}) => {
  if (!analysis || !analysis.scores) {
    return null;
  }

  const { scores, classification, feedback } = analysis;
  const totalScore = scores.total_score || 0;
  const percentage = ((totalScore - 3) / 12) * 100;
  
  // Helper to check if feedback has content
  const hasFeedback = feedback && (
    feedback.by_criterion?.migration_intent ||
    feedback.by_criterion?.goal_understanding ||
    feedback.by_criterion?.answer_length ||
    (feedback.improvements && feedback.improvements.length > 0)
  );

  const getClassificationStyle = () => {
    const lowerClass = classification?.toLowerCase() || "";
    if (lowerClass.includes("excellent")) {
      return {
        gradient: "from-green-500 to-emerald-600",
        bg: "bg-green-50",
        border: "border-green-200",
        emoji: "😇",
        badgeBg: "bg-green-100",
        badgeText: "text-green-800",
        progressBar: "bg-green-500",
      };
    }
    if (lowerClass.includes("good")) {
      return {
        gradient: "from-blue-500 to-indigo-600",
        bg: "bg-blue-50",
        border: "border-blue-200",
        emoji: "☺️",
        badgeBg: "bg-blue-100",
        badgeText: "text-blue-800",
        progressBar: "bg-blue-500",
      };
    }
    if (lowerClass.includes("average")) {
      return {
        gradient: "from-yellow-500 to-orange-500",
        bg: "bg-yellow-50",
        border: "border-yellow-200",
        emoji: "😕",
        badgeBg: "bg-yellow-100",
        badgeText: "text-yellow-800",
        progressBar: "bg-yellow-500",
      };
    }
    if (lowerClass.includes("weak")) {
      return {
        gradient: "from-orange-500 to-red-500",
        bg: "bg-orange-50",
        border: "border-orange-200",
        emoji: "😟",
        badgeBg: "bg-orange-100",
        badgeText: "text-orange-800",
        progressBar: "bg-orange-500",
      };
    }
    return {
      gradient: "from-red-500 to-red-700",
      bg: "bg-red-50",
      border: "border-red-200",
      emoji: "❌",
      badgeBg: "bg-red-100",
      badgeText: "text-red-800",
      progressBar: "bg-red-500",
    };
  };

  const style = getClassificationStyle();

  const getScoreColor = (score: number) => {
    if (score >= 4) return "text-green-600 bg-green-50 border-green-300";
    if (score === 3) return "text-yellow-600 bg-yellow-50 border-yellow-300";
    return "text-red-600 bg-red-50 border-red-300";
  };

  const criteriaLabels: Record<
    keyof AnalysisScores,
    { label: string; icon: string }
  > = {
    migration_intent: { label: "Intent", icon: "🏠" },
    goal_understanding: { label: "Goal", icon: "🎯" },
    answer_length: { label: "Length", icon: "📏" },
    total_score: { label: "Total", icon: "📊" },
  };

  return (
    <div className="my-3 animate-slide-in">
      <div
        className={`${style.bg} ${style.border} border-2 rounded-2xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow`}
      >
        {/* Header */}
        <div className={`bg-gradient-to-r ${style.gradient} p-4 text-white`}>
          <div className="flex items-center gap-3">
            <span className="text-3xl">{style.emoji}</span>
            <div>
              <div className="flex items-center gap-2">
                {questionNumber && (
                  <span className="bg-white bg-opacity-30 px-2 py-0.5 rounded-full text-xs font-bold">
                    Q{questionNumber}
                  </span>
                )}
                <span
                  className={`${style.badgeBg} ${style.badgeText} px-3 py-1 rounded-full text-xs font-bold`}
                >
                  {classification}
                </span>
              </div>
              <p className="text-sm opacity-90 mt-1">Answer Analysis</p>
            </div>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="bg-gray-200 h-2">
          <div
            className={`${style.progressBar} h-full transition-all duration-1000`}
            style={{ width: `${percentage}%` }}
          />
        </div>

        {/* Content */}
        <div className="p-4">
          {/* Score Breakdown */}
          <div className="grid grid-cols-3 gap-2 mb-4">
            {(
              [
                "migration_intent",
                "goal_understanding",
                "answer_length",
              ] as (keyof AnalysisScores)[]
            ).map((key) => {
              const meta = criteriaLabels[key];
              const score = scores[key] || 0;
              return (
                <div
                  key={key}
                  className={`${getScoreColor(
                    score
                  )} rounded-lg p-2 text-center border-2 transition-transform hover:scale-105`}
                >
                  <div className="text-xl mb-1">{meta.icon}</div>
                  <div className="text-xs font-medium">{meta.label}</div>
                </div>
              );
            })}
          </div>

          {/* Feedback */}
          {hasFeedback && (
            <div className="space-y-3">
              {/* Feedback by Criterion */}
              {(feedback.by_criterion?.migration_intent || 
                feedback.by_criterion?.goal_understanding || 
                feedback.by_criterion?.answer_length) && (
                <div className="bg-white border-l-4 border-blue-500 p-3 rounded-lg shadow-sm">
                  <p className="text-xs font-semibold text-gray-500 mb-2">
                    DETAILED FEEDBACK
                  </p>
                  <div className="space-y-2 text-sm text-gray-700">
                    {feedback.by_criterion.migration_intent && (
                      <div>
                        <span className="font-medium">Intent: </span>
                        {feedback.by_criterion.migration_intent}
                      </div>
                    )}
                    {feedback.by_criterion.goal_understanding && (
                      <div>
                        <span className="font-medium">Goal: </span>
                        {feedback.by_criterion.goal_understanding}
                      </div>
                    )}
                    {feedback.by_criterion.answer_length && (
                      <div>
                        <span className="font-medium">Length: </span>
                        {feedback.by_criterion.answer_length}
                      </div>
                    )}
                  </div>
                </div>
              )}

              {/* Improvements */}
              {feedback.improvements && feedback.improvements.length > 0 && (
                <div className="bg-white border-l-4 border-green-500 p-3 rounded-lg shadow-sm">
                  <p className="text-xs font-semibold text-gray-500 mb-2">
                    SUGGESTIONS FOR IMPROVEMENT
                  </p>
                  <ul className="list-disc list-inside space-y-1 text-sm text-gray-700">
                    {feedback.improvements.map((improvement, idx) => (
                      <li key={idx}>{improvement}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <style>{`
        @keyframes slide-in {
          from {
            transform: translateY(-10px);
            opacity: 0;
          }
          to {
            transform: translateY(0);
            opacity: 1;
          }
        }
        .animate-slide-in {
          animation: slide-in 0.4s ease-out;
        }
      `}</style>
    </div>
  );
};

export default AnswerFeedbackCard;



