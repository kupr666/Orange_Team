import type { LeaderboardPeriod, LeaderboardResponse, User } from "./api/types";

export const demoUser: User = {
  id: "b3f6eac7-1a62-40cc-a26b-58aecfa52f4d",
  email: "varvara@limitless.test",
  full_name: "Варвара Соколова",
  created_at: "2026-01-03T09:30:00Z",
  user_workout_score: 4280,
};

const athletes = [
  ["e1c00c5f-7a61-4acd-8d54-030f28c25911", "Ярослав Северный"],
  ["ed213b95-8ca6-41bf-82a8-357e570c2871", "Милана Морозова"],
  ["d9bc56de-1ec6-4280-b146-ea46dd471264", "Святослав Орлов"],
  ["f14aa61e-e628-48ac-9549-44b4270cab08", "Алиса Ветрова"],
  ["2b2a9e8f-7a35-478b-968c-d3c369c40a51", "Ратмир Волков"],
  [demoUser.id, demoUser.full_name],
  ["69b96a75-d4ac-45c4-a602-c3ee7ea53454", "Кира Белова"],
  ["54061185-e5e3-45b0-a8dd-8d4bd5d62643", "Владислав Ледов"],
] as const;

const completedByPeriod: Record<LeaderboardPeriod, number[]> = {
  daily: [8, 7, 6, 6, 5, 4, 3, 2],
  weekly: [28, 25, 23, 20, 19, 17, 15, 13],
  monthly: [72, 68, 61, 58, 55, 49, 45, 41],
};

export function getDemoLeaderboard(period: LeaderboardPeriod): LeaderboardResponse {
  const completed = completedByPeriod[period];
  const now = "2026-08-29T10:00:00+03:00";
  const items = athletes.map(([id, fullName], index) => ({
    rank: index + 1,
    user: { id, full_name: fullName },
    score: completed[index] * 100,
    completed_workouts: completed[index],
    last_activity_at: now,
    is_current_user: id === demoUser.id,
  }));
  const current = items.find((item) => item.user.id === demoUser.id)!;

  return {
    period,
    status: period === "daily" ? "live" : "published",
    period_start: "2026-08-01T00:00:00+03:00",
    period_end: "2026-08-31T23:59:59+03:00",
    timezone: "Europe/Moscow",
    metric: "workout_score_v1",
    generated_at: now,
    next_refresh_at: "2026-08-29T10:05:00+03:00",
    items,
    current_user: {
      rank: current.rank,
      user: current.user,
      score: current.score,
      completed_workouts: current.completed_workouts,
      last_activity_at: current.last_activity_at,
      eligible: true,
      in_top: true,
    },
  };
}
