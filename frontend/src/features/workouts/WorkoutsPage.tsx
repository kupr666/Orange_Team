import { useEffect, useMemo, useState, type FormEvent } from "react";
import { apiClient } from "../../api/client";
import type {
  CreateWorkoutExerciseRequest,
  Exercise,
  PatchWorkoutExerciseRequest,
  User,
  Workout,
  WorkoutExercise,
  WorkoutStatus,
} from "../../api/types";
import "./workouts.css";

type WorkoutFilter = "all" | "planned" | "in_progress" | "completed";

const statusLabels: Record<WorkoutStatus, string> = {
  planned: "Запланирована",
  in_progress: "Идёт сейчас",
  completed: "Завершена",
  cancelled: "Отменена",
};

const filterLabels: Array<{ id: WorkoutFilter; label: string }> = [
  { id: "all", label: "Все" },
  { id: "planned", label: "Планы" },
  { id: "in_progress", label: "Активные" },
  { id: "completed", label: "Завершённые" },
];

function formatDate(value: string | null, withTime = false) {
  if (!value) return "—";
  return new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: withTime ? "2-digit" : undefined,
    minute: withTime ? "2-digit" : undefined,
  }).format(new Date(value));
}

function formatDuration(totalSeconds: number) {
  const minutes = Math.max(1, Math.round(totalSeconds / 60));
  return `${minutes} мин`;
}

function upsertWorkout(items: Workout[], updated: Workout) {
  const exists = items.some((item) => item.id === updated.id);
  const next = exists
    ? items.map((item) => (item.id === updated.id ? updated : item))
    : [updated, ...items];
  return next.sort(
    (first, second) =>
      new Date(second.created_at).getTime() - new Date(first.created_at).getTime(),
  );
}

export function WorkoutsPage({
  user,
  token,
  onLogin,
}: {
  user: User | null;
  token: string | null;
  onLogin: () => void;
}) {
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [catalog, setCatalog] = useState<Exercise[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [workoutExercises, setWorkoutExercises] = useState<WorkoutExercise[]>([]);
  const [filter, setFilter] = useState<WorkoutFilter>("all");
  const [loadingList, setLoadingList] = useState(false);
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [busyAction, setBusyAction] = useState<string | null>(null);
  const [showAddExercise, setShowAddExercise] = useState(false);
  const [intensity, setIntensity] = useState(7);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!user || !token) {
      setWorkouts([]);
      setCatalog([]);
      setSelectedId(null);
      setWorkoutExercises([]);
      setLoadError(null);
      return;
    }

    let active = true;
    setLoadingList(true);
    setLoadError(null);
    Promise.all([apiClient.getWorkouts(token), apiClient.getExercises()])
      .then(([workoutItems, exerciseItems]) => {
        if (!active) return;
        const sorted = [...workoutItems].sort(
          (first, second) =>
            new Date(second.created_at).getTime() - new Date(first.created_at).getTime(),
        );
        setWorkouts(sorted);
        setCatalog(exerciseItems);
        setSelectedId((current) => {
          if (current && sorted.some((workout) => workout.id === current)) return current;
          return (
            sorted.find((workout) => workout.status === "in_progress")?.id ??
            sorted.find((workout) => workout.status === "planned")?.id ??
            sorted[0]?.id ??
            null
          );
        });
      })
      .catch((error: unknown) => {
        if (active) {
          setLoadError(
            error instanceof Error
              ? error.message
              : "Не удалось загрузить тренировки.",
          );
        }
      })
      .finally(() => {
        if (active) setLoadingList(false);
      });

    return () => {
      active = false;
    };
  }, [reloadKey, token, user]);

  useEffect(() => {
    if (!selectedId || !token) {
      setWorkoutExercises([]);
      return;
    }

    let active = true;
    setLoadingDetail(true);
    setActionError(null);
    setShowAddExercise(false);
    Promise.all([
      apiClient.getWorkout(selectedId, token),
      apiClient.getWorkoutExercises(selectedId, token),
    ])
      .then(([workout, exerciseItems]) => {
        if (!active) return;
        setWorkouts((items) => upsertWorkout(items, workout));
        setWorkoutExercises(exerciseItems);
        if (workout.intensity) setIntensity(workout.intensity);
      })
      .catch((error: unknown) => {
        if (active) {
          setActionError(
            error instanceof Error
              ? error.message
              : "Не удалось загрузить состав тренировки.",
          );
        }
      })
      .finally(() => {
        if (active) setLoadingDetail(false);
      });

    return () => {
      active = false;
    };
  }, [selectedId, token]);

  const selectedWorkout = useMemo(
    () => workouts.find((workout) => workout.id === selectedId) ?? null,
    [selectedId, workouts],
  );

  const visibleWorkouts = useMemo(() => {
    if (filter === "all") return workouts;
    return workouts.filter((workout) => workout.status === filter);
  }, [filter, workouts]);

  const completedExercises = workoutExercises.filter((exercise) => exercise.completed).length;
  const canModifyExercises =
    selectedWorkout?.status === "planned" || selectedWorkout?.status === "in_progress";

  const refreshWorkoutSummary = async (workoutId: string) => {
    if (!token) return;
    const updated = await apiClient.getWorkout(workoutId, token);
    setWorkouts((items) => upsertWorkout(items, updated));
  };

  const createWorkout = async () => {
    if (!token) return;
    setBusyAction("create");
    setActionError(null);
    try {
      const created = await apiClient.createWorkout(token);
      setWorkouts((items) => upsertWorkout(items, created));
      setSelectedId(created.id);
      setWorkoutExercises([]);
      setFilter("all");
    } catch (error: unknown) {
      setActionError(error instanceof Error ? error.message : "Не удалось создать тренировку.");
    } finally {
      setBusyAction(null);
    }
  };

  const patchWorkoutStatus = async (
    status: Extract<WorkoutStatus, "in_progress" | "completed" | "cancelled">,
  ) => {
    if (!selectedWorkout || !token) return;
    setBusyAction(status);
    setActionError(null);
    try {
      const payload = status === "completed" ? { status, intensity } : { status };
      const updated = await apiClient.patchWorkout(selectedWorkout.id, payload, token);
      setWorkouts((items) => upsertWorkout(items, updated));
      setShowAddExercise(false);
    } catch (error: unknown) {
      setActionError(
        error instanceof Error ? error.message : "Не удалось изменить статус тренировки.",
      );
    } finally {
      setBusyAction(null);
    }
  };

  const deleteWorkout = async () => {
    if (!selectedWorkout || !token) return;
    if (!window.confirm("Удалить эту тренировку без возможности восстановления?")) return;
    setBusyAction("delete-workout");
    setActionError(null);
    try {
      await apiClient.deleteWorkout(selectedWorkout.id, token);
      const remaining = workouts.filter((workout) => workout.id !== selectedWorkout.id);
      setWorkouts(remaining);
      setSelectedId(remaining[0]?.id ?? null);
      setWorkoutExercises([]);
    } catch (error: unknown) {
      setActionError(error instanceof Error ? error.message : "Не удалось удалить тренировку.");
    } finally {
      setBusyAction(null);
    }
  };

  const updateCompletedIntensity = async () => {
    if (!selectedWorkout || !token) return;
    setBusyAction("update-intensity");
    setActionError(null);
    try {
      const updated = await apiClient.patchWorkout(
        selectedWorkout.id,
        { intensity },
        token,
      );
      setWorkouts((items) => upsertWorkout(items, updated));
    } catch (error: unknown) {
      setActionError(
        error instanceof Error ? error.message : "Не удалось обновить интенсивность.",
      );
    } finally {
      setBusyAction(null);
    }
  };

  if (!user || !token) {
    return (
      <section className="workouts-guest">
        <span className="eyebrow">Путь силы</span>
        <h1>Тренировки</h1>
        <p>
          Собирайте тренировку из упражнений, отмечайте выполненные подходы и
          завершайте занятие — результат сразу попадёт в ваш прогресс и рейтинг.
        </p>
        <div className="workouts-guest-features">
          <span>Планы и активные занятия</span>
          <span>Силовые и временные упражнения</span>
          <span>Интенсивность и личный счёт</span>
        </div>
        <button className="button button--primary" type="button" onClick={onLogin}>
          Войти и начать
        </button>
      </section>
    );
  }

  return (
    <section className="workouts-page">
      <div className="workouts-heading">
        <div>
          <span className="eyebrow">Путь силы</span>
          <h1>Тренировки</h1>
          <p>Планируйте нагрузку, выполняйте упражнения и фиксируйте результат.</p>
        </div>
        <button
          className="button button--primary"
          type="button"
          disabled={busyAction === "create"}
          onClick={createWorkout}
        >
          {busyAction === "create" ? "Создаём…" : "Новая тренировка"}
        </button>
      </div>

      {loadError && (
        <div className="workout-error" role="alert">
          <span>{loadError}</span>
          <button type="button" onClick={() => setReloadKey((value) => value + 1)}>
            Повторить
          </button>
        </div>
      )}

      <div className="workouts-layout">
        <aside className="workout-index" aria-label="Список тренировок">
          <div className="workout-filters" aria-label="Фильтр тренировок">
            {filterLabels.map((item) => (
              <button
                key={item.id}
                type="button"
                className={filter === item.id ? "active" : ""}
                onClick={() => setFilter(item.id)}
              >
                {item.label}
              </button>
            ))}
          </div>

          {loadingList ? (
            <WorkoutIndexLoading />
          ) : visibleWorkouts.length === 0 ? (
            <div className="workout-index-empty">
              <strong>{workouts.length === 0 ? "Пока нет тренировок" : "Здесь пусто"}</strong>
              <span>
                {workouts.length === 0
                  ? "Создайте первую тренировку."
                  : "В этом фильтре нет записей."}
              </span>
            </div>
          ) : (
            <div className="workout-index-list">
              {visibleWorkouts.map((workout, index) => (
                <button
                  key={workout.id}
                  type="button"
                  className={`workout-index-card ${selectedId === workout.id ? "active" : ""}`}
                  onClick={() => setSelectedId(workout.id)}
                >
                  <span className={`workout-status workout-status--${workout.status}`}>
                    {statusLabels[workout.status]}
                  </span>
                  <strong>Тренировка {workouts.length - index}</strong>
                  <span>{formatDate(workout.created_at)}</span>
                  <span className="workout-index-score">{workout.workout_score} очков</span>
                </button>
              ))}
            </div>
          )}
        </aside>

        <div className="workout-detail">
          {!selectedWorkout ? (
            <div className="workout-detail-empty">
              <h2>Выберите тренировку</h2>
              <p>Откройте существующую запись слева или создайте новую.</p>
            </div>
          ) : (
            <>
              <header className="workout-detail-header">
                <div>
                  <span className={`workout-status workout-status--${selectedWorkout.status}`}>
                    {statusLabels[selectedWorkout.status]}
                  </span>
                  <h2>Тренировка от {formatDate(selectedWorkout.created_at)}</h2>
                  <p>
                    {selectedWorkout.started_at
                      ? `Начата ${formatDate(selectedWorkout.started_at, true)}`
                      : "Готова к наполнению и запуску"}
                  </p>
                </div>
                <button
                  className="workout-delete-button"
                  type="button"
                  disabled={busyAction === "delete-workout"}
                  onClick={deleteWorkout}
                >
                  Удалить
                </button>
              </header>

              <div className="workout-metrics">
                <div><span>Упражнения</span><strong>{workoutExercises.length}</strong></div>
                <div><span>Выполнено</span><strong>{completedExercises}</strong></div>
                <div><span>Счёт</span><strong>{selectedWorkout.workout_score}</strong></div>
                <div>
                  <span>Интенсивность</span>
                  <strong>{selectedWorkout.intensity ?? "—"}</strong>
                </div>
              </div>

              {actionError && <div className="workout-error" role="alert">{actionError}</div>}

              <div className="workout-actions">
                {selectedWorkout.status === "planned" && (
                  <>
                    <button
                      className="button button--primary"
                      type="button"
                      disabled={Boolean(busyAction)}
                      onClick={() => patchWorkoutStatus("in_progress")}
                    >
                      {busyAction === "in_progress" ? "Запускаем…" : "Начать тренировку"}
                    </button>
                    <button
                      className="button button--secondary"
                      type="button"
                      disabled={Boolean(busyAction)}
                      onClick={() => patchWorkoutStatus("cancelled")}
                    >
                      Отменить план
                    </button>
                  </>
                )}
                {selectedWorkout.status === "in_progress" && (
                  <div className="complete-workout-control">
                    <label>
                      <span>Интенсивность: {intensity}</span>
                      <input
                        type="range"
                        min="1"
                        max="10"
                        value={intensity}
                        onChange={(event) => setIntensity(Number(event.target.value))}
                      />
                    </label>
                    <button
                      className="button button--primary"
                      type="button"
                      disabled={Boolean(busyAction) || completedExercises === 0}
                      title={completedExercises === 0 ? "Сначала выполните хотя бы одно упражнение" : undefined}
                      onClick={() => patchWorkoutStatus("completed")}
                    >
                      {busyAction === "completed" ? "Завершаем…" : "Завершить тренировку"}
                    </button>
                  </div>
                )}
                {selectedWorkout.status === "completed" && (
                  <div className="complete-workout-control">
                    <label>
                      <span>Интенсивность: {intensity}</span>
                      <input
                        type="range"
                        min="1"
                        max="10"
                        value={intensity}
                        onChange={(event) => setIntensity(Number(event.target.value))}
                      />
                    </label>
                    <button
                      className="button button--secondary"
                      type="button"
                      disabled={Boolean(busyAction) || intensity === selectedWorkout.intensity}
                      onClick={updateCompletedIntensity}
                    >
                      {busyAction === "update-intensity" ? "Сохраняем…" : "Обновить интенсивность"}
                    </button>
                  </div>
                )}
                {canModifyExercises && (
                  <button
                    className="workout-add-button"
                    type="button"
                    onClick={() => setShowAddExercise((visible) => !visible)}
                  >
                    {showAddExercise ? "Закрыть форму" : "Добавить упражнение"}
                  </button>
                )}
              </div>

              {showAddExercise && canModifyExercises && (
                <AddExerciseForm
                  workoutId={selectedWorkout.id}
                  catalog={catalog}
                  token={token}
                  onCancel={() => setShowAddExercise(false)}
                  onCreated={(created) => {
                    setWorkoutExercises((items) => [...items, created]);
                    setShowAddExercise(false);
                    void refreshWorkoutSummary(selectedWorkout.id);
                  }}
                />
              )}

              <section className="workout-exercises-section">
                <div className="workout-section-heading">
                  <div>
                    <span className="eyebrow">Состав занятия</span>
                    <h3>Упражнения</h3>
                  </div>
                  <span>{completedExercises} из {workoutExercises.length} выполнено</span>
                </div>

                {loadingDetail ? (
                  <WorkoutExerciseLoading />
                ) : workoutExercises.length === 0 ? (
                  <div className="workout-exercises-empty">
                    <strong>Тренировка пока пуста</strong>
                    <p>Добавьте упражнение из каталога, чтобы рассчитать нагрузку и счёт.</p>
                  </div>
                ) : (
                  <div className="workout-exercise-list">
                    {workoutExercises.map((item) => (
                      <WorkoutExerciseRow
                        key={item.id}
                        item={item}
                        exercise={catalog.find((exercise) => exercise.id === item.exercise_id)}
                        workoutId={selectedWorkout.id}
                        token={token}
                        editable={Boolean(canModifyExercises)}
                        onChanged={(updated) => {
                          setWorkoutExercises((items) =>
                            items.map((entry) => (entry.id === updated.id ? updated : entry)),
                          );
                          void refreshWorkoutSummary(selectedWorkout.id);
                        }}
                        onDeleted={(deletedId) => {
                          setWorkoutExercises((items) =>
                            items.filter((entry) => entry.id !== deletedId),
                          );
                          void refreshWorkoutSummary(selectedWorkout.id);
                        }}
                      />
                    ))}
                  </div>
                )}
              </section>
            </>
          )}
        </div>
      </div>
    </section>
  );
}

function WorkoutIndexLoading() {
  return (
    <div className="workout-index-loading" aria-label="Загрузка тренировок">
      {[0, 1, 2].map((item) => <span key={item} />)}
    </div>
  );
}

function WorkoutExerciseLoading() {
  return (
    <div className="workout-exercise-loading" aria-label="Загрузка упражнений">
      {[0, 1].map((item) => <span key={item} />)}
    </div>
  );
}

function AddExerciseForm({
  workoutId,
  catalog,
  token,
  onCancel,
  onCreated,
}: {
  workoutId: string;
  catalog: Exercise[];
  token: string;
  onCancel: () => void;
  onCreated: (exercise: WorkoutExercise) => void;
}) {
  const [exerciseId, setExerciseId] = useState(catalog[0]?.id ?? "");
  const [weight, setWeight] = useState("40");
  const [sets, setSets] = useState("3");
  const [reps, setReps] = useState("10");
  const [durationMinutes, setDurationMinutes] = useState("10");
  const [completed, setCompleted] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const selected = catalog.find((exercise) => exercise.id === exerciseId) ?? catalog[0];

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selected) return;
    setSubmitting(true);
    setError(null);
    try {
      let payload: CreateWorkoutExerciseRequest;
      if (selected.type === "weight") {
        payload = {
          exercise_id: selected.id,
          weight: Number(weight),
          sets: Number(sets),
          reps: Number(reps),
          completed,
        };
      } else {
        payload = {
          exercise_id: selected.id,
          duration: Math.round(Number(durationMinutes) * 60),
          completed,
        };
      }
      const created = await apiClient.createWorkoutExercise(workoutId, payload, token);
      onCreated(created);
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось добавить упражнение.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  if (!selected) {
    return (
      <div className="add-exercise-form add-exercise-form--empty">
        <strong>Каталог упражнений пуст</strong>
        <p>Сначала администратор должен добавить упражнения в общий каталог.</p>
        <button className="button button--secondary" type="button" onClick={onCancel}>Закрыть</button>
      </div>
    );
  }

  return (
    <form className="add-exercise-form" onSubmit={submit}>
      <div className="add-exercise-title">
        <div>
          <span className="eyebrow">Новое упражнение</span>
          <h3>Настройте нагрузку</h3>
        </div>
        <button type="button" onClick={onCancel}>Отмена</button>
      </div>

      <label className="workout-field workout-field--wide">
        <span>Упражнение из каталога</span>
        <select value={selected.id} onChange={(event) => setExerciseId(event.target.value)}>
          {catalog.map((exercise) => (
            <option key={exercise.id} value={exercise.id}>
              {exercise.name} · сложность {exercise.difficulty}/10
            </option>
          ))}
        </select>
      </label>

      <div className="exercise-input-grid">
        {selected.type === "weight" ? (
          <>
            <NumberField label="Вес, кг" value={weight} onChange={setWeight} min={1} />
            <NumberField label="Подходы" value={sets} onChange={setSets} min={1} />
            <NumberField label="Повторения" value={reps} onChange={setReps} min={1} />
          </>
        ) : (
          <NumberField
            label="Длительность, мин"
            value={durationMinutes}
            onChange={setDurationMinutes}
            min={1}
          />
        )}
      </div>

      <label className="workout-check">
        <input
          type="checkbox"
          checked={completed}
          onChange={(event) => setCompleted(event.target.checked)}
        />
        <span>Уже выполнено</span>
      </label>

      {error && <div className="workout-error" role="alert">{error}</div>}
      <button className="button button--primary" type="submit" disabled={submitting}>
        {submitting ? "Добавляем…" : "Добавить в тренировку"}
      </button>
    </form>
  );
}

function NumberField({
  label,
  value,
  onChange,
  min,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  min: number;
}) {
  return (
    <label className="workout-field">
      <span>{label}</span>
      <input
        type="number"
        min={min}
        step="1"
        required
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
    </label>
  );
}

function WorkoutExerciseRow({
  item,
  exercise,
  workoutId,
  token,
  editable,
  onChanged,
  onDeleted,
}: {
  item: WorkoutExercise;
  exercise?: Exercise;
  workoutId: string;
  token: string;
  editable: boolean;
  onChanged: (exercise: WorkoutExercise) => void;
  onDeleted: (id: string) => void;
}) {
  const type = exercise?.type ?? (item.duration ? "duration" : "weight");
  const [weight, setWeight] = useState(String(item.weight ?? 1));
  const [sets, setSets] = useState(String(item.sets ?? 1));
  const [reps, setReps] = useState(String(item.reps ?? 1));
  const [durationMinutes, setDurationMinutes] = useState(
    String(Math.max(1, Math.round((item.duration ?? 60) / 60))),
  );
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setWeight(String(item.weight ?? 1));
    setSets(String(item.sets ?? 1));
    setReps(String(item.reps ?? 1));
    setDurationMinutes(String(Math.max(1, Math.round((item.duration ?? 60) / 60))));
  }, [item.duration, item.reps, item.sets, item.weight]);

  const patch = async (payload: PatchWorkoutExerciseRequest) => {
    setSaving(true);
    setError(null);
    try {
      const updated = await apiClient.patchWorkoutExercise(workoutId, item.id, payload, token);
      onChanged(updated);
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось обновить упражнение.",
      );
    } finally {
      setSaving(false);
    }
  };

  const saveLoad = () => {
    if (type === "weight") {
      void patch({ weight: Number(weight), sets: Number(sets), reps: Number(reps) });
    } else {
      void patch({ duration: Math.round(Number(durationMinutes) * 60) });
    }
  };

  const deleteExercise = async () => {
    if (!window.confirm("Удалить упражнение из этой тренировки?")) return;
    setSaving(true);
    setError(null);
    try {
      await apiClient.deleteWorkoutExercise(workoutId, item.id, token);
      onDeleted(item.id);
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось удалить упражнение.",
      );
      setSaving(false);
    }
  };

  return (
    <article className={`workout-exercise ${item.completed ? "workout-exercise--completed" : ""}`}>
      <div className="workout-exercise-main">
        <label className="exercise-complete-check">
          <input
            type="checkbox"
            checked={item.completed}
            disabled={!editable || saving}
            onChange={(event) => void patch({ completed: event.target.checked })}
          />
          <span className="visually-hidden">Отметить выполнение</span>
        </label>
        <div className="workout-exercise-copy">
          <span>{type === "weight" ? "Силовое" : "На время"}</span>
          <h4>{exercise?.name ?? "Упражнение"}</h4>
          {exercise?.description && <p>{exercise.description}</p>}
        </div>
        <div className="exercise-load">
          <span>Нагрузка</span>
          <strong>{item.exercise_load}</strong>
        </div>
      </div>

      {editable ? (
        <div className="workout-exercise-editor">
          {type === "weight" ? (
            <>
              <NumberField label="Вес, кг" value={weight} onChange={setWeight} min={1} />
              <NumberField label="Подходы" value={sets} onChange={setSets} min={1} />
              <NumberField label="Повторы" value={reps} onChange={setReps} min={1} />
            </>
          ) : (
            <NumberField
              label="Минуты"
              value={durationMinutes}
              onChange={setDurationMinutes}
              min={1}
            />
          )}
          <div className="exercise-row-actions">
            <button type="button" disabled={saving} onClick={saveLoad}>
              {saving ? "Сохраняем…" : "Сохранить"}
            </button>
            <button type="button" disabled={saving} onClick={deleteExercise}>Удалить</button>
          </div>
        </div>
      ) : (
        <div className="workout-exercise-summary">
          {type === "weight"
            ? `${item.weight} кг · ${item.sets} × ${item.reps}`
            : formatDuration(item.duration ?? 0)}
        </div>
      )}

      {error && <div className="exercise-row-error" role="alert">{error}</div>}
    </article>
  );
}
