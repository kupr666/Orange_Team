import type {
  ApiErrorBody,
  CreateWorkoutExerciseRequest,
  Exercise,
  LeaderboardPeriod,
  LeaderboardResponse,
  LoginRequest,
  LoginResponse,
  PatchUserRequest,
  PatchWorkoutExerciseRequest,
  PatchWorkoutRequest,
  RegisterUserRequest,
  User,
  Workout,
  WorkoutExercise,
} from "./types";

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? "/api/v1").replace(/\/$/, "");
export const ACCESS_TOKEN_KEY = "limitless_access_token";

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
    public readonly body?: ApiErrorBody,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

interface RequestOptions extends RequestInit {
  token?: string | null;
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { token, headers, ...requestOptions } = options;
  const requestHeaders = new Headers(headers);
  requestHeaders.set("Accept", "application/json");

  if (requestOptions.body && !requestHeaders.has("Content-Type")) {
    requestHeaders.set("Content-Type", "application/json");
  }

  if (token) {
    requestHeaders.set("Authorization", `Bearer ${token}`);
  }

  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...requestOptions,
      headers: requestHeaders,
    });
  } catch {
    throw new ApiError(0, "Не удалось связаться с сервером. Проверьте подключение и повторите попытку.");
  }

  const isJson = response.headers.get("content-type")?.includes("application/json");
  const payload = isJson ? ((await response.json()) as unknown) : undefined;

  if (!response.ok) {
    const body = payload as ApiErrorBody | undefined;
    throw new ApiError(response.status, body?.message ?? "Сервер не смог обработать запрос.", body);
  }

  return payload as T;
}

export const apiClient = {
  register(payload: RegisterUserRequest): Promise<User> {
    return request<User>("/register", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  login(payload: LoginRequest): Promise<LoginResponse> {
    return request<LoginResponse>("/login", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  getCurrentUser(token: string): Promise<User> {
    return request<User>("/users/me", { token });
  },

  patchCurrentUser(payload: PatchUserRequest, token: string): Promise<User> {
    return request<User>("/users/me", {
      method: "PATCH",
      token,
      body: JSON.stringify(payload),
    });
  },

  getLeaderboard(
    period: LeaderboardPeriod,
    userId: string,
    limit = 50,
    token?: string | null,
  ): Promise<LeaderboardResponse> {
    const params = new URLSearchParams({ user_id: userId, limit: String(limit) });
    return request<LeaderboardResponse>(`/leaderboard/${period}?${params.toString()}`, { token });
  },

  getExercises(): Promise<Exercise[]> {
    return request<Exercise[]>("/exercises");
  },

  getWorkouts(token: string): Promise<Workout[]> {
    return request<Workout[]>("/workouts", { token });
  },

  createWorkout(token: string): Promise<Workout> {
    return request<Workout>("/workouts", { method: "POST", token });
  },

  getWorkout(workoutId: string, token: string): Promise<Workout> {
    return request<Workout>(`/workouts/${encodeURIComponent(workoutId)}`, { token });
  },

  patchWorkout(
    workoutId: string,
    payload: PatchWorkoutRequest,
    token: string,
  ): Promise<Workout> {
    return request<Workout>(`/workouts/${encodeURIComponent(workoutId)}`, {
      method: "PATCH",
      token,
      body: JSON.stringify(payload),
    });
  },

  deleteWorkout(workoutId: string, token: string): Promise<void> {
    return request<void>(`/workouts/${encodeURIComponent(workoutId)}`, {
      method: "DELETE",
      token,
    });
  },

  getWorkoutExercises(workoutId: string, token: string): Promise<WorkoutExercise[]> {
    return request<WorkoutExercise[]>(
      `/workouts/${encodeURIComponent(workoutId)}/exercises`,
      { token },
    );
  },

  createWorkoutExercise(
    workoutId: string,
    payload: CreateWorkoutExerciseRequest,
    token: string,
  ): Promise<WorkoutExercise> {
    return request<WorkoutExercise>(
      `/workouts/${encodeURIComponent(workoutId)}/exercises`,
      { method: "POST", token, body: JSON.stringify(payload) },
    );
  },

  patchWorkoutExercise(
    workoutId: string,
    workoutExerciseId: string,
    payload: PatchWorkoutExerciseRequest,
    token: string,
  ): Promise<WorkoutExercise> {
    return request<WorkoutExercise>(
      `/workouts/${encodeURIComponent(workoutId)}/exercises/${encodeURIComponent(workoutExerciseId)}`,
      { method: "PATCH", token, body: JSON.stringify(payload) },
    );
  },

  deleteWorkoutExercise(
    workoutId: string,
    workoutExerciseId: string,
    token: string,
  ): Promise<void> {
    return request<void>(
      `/workouts/${encodeURIComponent(workoutId)}/exercises/${encodeURIComponent(workoutExerciseId)}`,
      { method: "DELETE", token },
    );
  },
};
