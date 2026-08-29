export type LeaderboardPeriod = "daily" | "weekly" | "monthly";

export interface User {
  id: string;
  email: string;
  full_name: string;
  created_at: string;
  updated_at?: string;
  user_workout_score: number;
  sex?: "male" | "female";
  weight_grams?: number;
  birth_date?: string;
  height_cm?: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: "Bearer";
}

export interface RegisterUserRequest extends LoginRequest {
  full_name: string;
}

export interface PatchUserRequest {
  sex?: "male" | "female" | null;
  weight_grams?: number | null;
  birth_date?: string | null;
  height_cm?: number | null;
}

export interface LeaderboardUser {
  id: string;
  full_name: string;
}

export interface LeaderboardEntry {
  rank: number | null;
  user: LeaderboardUser;
  score: number;
  completed_workouts: number;
  last_activity_at: string | null;
  is_current_user: boolean;
}

export interface CurrentUserLeaderboardEntry {
  rank: number | null;
  user: LeaderboardUser;
  score: number;
  completed_workouts: number;
  last_activity_at: string | null;
  eligible: boolean;
  in_top: boolean;
}

export interface LeaderboardResponse {
  period: LeaderboardPeriod;
  status: "live" | "published";
  period_start: string;
  period_end: string;
  timezone: string;
  metric: "workout_score_v1";
  generated_at: string;
  next_refresh_at: string;
  items: LeaderboardEntry[];
  current_user: CurrentUserLeaderboardEntry;
}

export interface ApiErrorBody {
  error: string;
  message: string;
}

export type WorkoutStatus = "planned" | "in_progress" | "completed" | "cancelled";
export type ExerciseType = "weight" | "duration";

export interface Workout {
  id: string;
  version: number;
  user_id: string;
  status: WorkoutStatus;
  started_at: string | null;
  completed_at: string | null;
  created_at: string;
  updated_at: string;
  workout_score: number;
  intensity: number | null;
  personal_score_coefficient: number;
}

export interface PatchWorkoutRequest {
  status?: WorkoutStatus;
  started_at?: string | null;
  completed_at?: string | null;
  intensity?: number | null;
}

export interface Exercise {
  id: string;
  version: number;
  name: string;
  description: string;
  difficulty: number;
  created_at: string;
  updated_at?: string | null;
  type: ExerciseType;
}

export interface WorkoutExercise {
  id: string;
  version: number;
  workout_id: string;
  exercise_id: string;
  weight?: number;
  sets?: number;
  reps?: number;
  duration?: number;
  completed: boolean;
  exercise_load: number;
  created_at: string;
  updated_at?: string;
}

export type CreateWorkoutExerciseRequest =
  | {
      exercise_id: string;
      weight: number;
      sets: number;
      reps: number;
      duration?: never;
      completed?: boolean;
    }
  | {
      exercise_id: string;
      duration: number;
      weight?: never;
      sets?: never;
      reps?: never;
      completed?: boolean;
    };

export interface PatchWorkoutExerciseRequest {
  weight?: number | null;
  sets?: number | null;
  reps?: number | null;
  duration?: number | null;
  completed?: boolean;
}
