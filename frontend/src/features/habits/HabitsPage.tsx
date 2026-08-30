import { useEffect, useMemo, useState, type FormEvent } from "react";
import { apiClient } from "../../api/client";
import type { Habit, User } from "../../api/types";
import "./habits.css";

function streakLabel(value: number) {
  const lastTwo = value % 100;
  const last = value % 10;
  if (lastTwo >= 11 && lastTwo <= 14) return "дней";
  if (last === 1) return "день";
  if (last >= 2 && last <= 4) return "дня";
  return "дней";
}

export function HabitsPage({
  user,
  token,
  onLogin,
}: {
  user: User | null;
  token: string | null;
  onLogin: () => void;
}) {
  const [habits, setHabits] = useState<Habit[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [creating, setCreating] = useState(false);
  const [busyHabitId, setBusyHabitId] = useState<string | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!token || !user) {
      setHabits([]);
      setLoading(false);
      setLoadError(null);
      return;
    }

    let active = true;
    setLoading(true);
    setLoadError(null);
    apiClient
      .getHabits(token)
      .then((items) => {
        if (active) setHabits(items);
      })
      .catch((requestError: unknown) => {
        if (!active) return;
        setLoadError(
          requestError instanceof Error
            ? requestError.message
            : "Не удалось загрузить привычки.",
        );
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [reloadKey, token, user]);

  const completedToday = useMemo(
    () => habits.filter((habit) => habit.completed_today).length,
    [habits],
  );
  const strongestStreak = useMemo(
    () => habits.reduce((highest, habit) => Math.max(highest, habit.current_streak), 0),
    [habits],
  );

  const createHabit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token) return;

    setCreating(true);
    setActionError(null);
    try {
      const created = await apiClient.createHabit({ name, description }, token);
      setHabits((items) => [created, ...items]);
      setName("");
      setDescription("");
    } catch (requestError: unknown) {
      setActionError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось создать привычку.",
      );
    } finally {
      setCreating(false);
    }
  };

  const completeHabit = async (habitId: string) => {
    if (!token) return;
    setBusyHabitId(habitId);
    setActionError(null);
    try {
      const updated = await apiClient.completeHabit(habitId, token);
      setHabits((items) => items.map((habit) => (habit.id === updated.id ? updated : habit)));
    } catch (requestError: unknown) {
      setActionError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось отметить привычку.",
      );
    } finally {
      setBusyHabitId(null);
    }
  };

  const deleteHabit = async (habit: Habit) => {
    if (!token || !window.confirm(`Удалить привычку «${habit.name}»? Текущая серия будет потеряна.`)) {
      return;
    }

    setBusyHabitId(habit.id);
    setActionError(null);
    try {
      await apiClient.deleteHabit(habit.id, token);
      setHabits((items) => items.filter((item) => item.id !== habit.id));
    } catch (requestError: unknown) {
      setActionError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось удалить привычку.",
      );
    } finally {
      setBusyHabitId(null);
    }
  };

  return (
    <section className="habits-page">
      <div className="page-heading habits-heading">
        <div>
          <span className="eyebrow">Ежедневный ритм</span>
          <h1>Привычки</h1>
          <p>Отмечайте выполнение один раз в день и не прерывайте свою серию.</p>
        </div>
        {user && (
          <div className="habits-summary" aria-label="Итоги привычек">
            <div>
              <span>Сегодня</span>
              <strong>{completedToday} / {habits.length}</strong>
            </div>
            <div>
              <span>Макс. сейчас</span>
              <strong>{strongestStreak} {streakLabel(strongestStreak)}</strong>
            </div>
          </div>
        )}
      </div>

      {!user || !token ? (
        <div className="habits-card habits-state">
          <h2>Начните свою серию</h2>
          <p>Войдите, чтобы создать ежедневные привычки и сохранять прогресс.</p>
          <button className="button button--primary" type="button" onClick={onLogin}>Войти</button>
        </div>
      ) : (
        <div className="habits-layout">
          <aside className="habits-card habit-create-card">
            <span className="row-kicker">Новая привычка</span>
            <h2>Что повторяем?</h2>
            <p>Выберите простое действие, которое реально выполнить каждый день.</p>
            <form className="habit-form" onSubmit={createHabit}>
              <label>
                <span>Название</span>
                <input
                  name="name"
                  minLength={2}
                  maxLength={80}
                  required
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="Например, читать 10 страниц"
                />
              </label>
              <label>
                <span>Описание</span>
                <textarea
                  name="description"
                  maxLength={500}
                  rows={4}
                  value={description}
                  onChange={(event) => setDescription(event.target.value)}
                  placeholder="Короткое напоминание для себя"
                />
              </label>
              <button className="button button--primary button--wide" type="submit" disabled={creating}>
                {creating ? "Создаём…" : "Добавить привычку"}
              </button>
            </form>
          </aside>

          <div className="habits-column">
            {actionError && <p className="form-error habits-action-error" role="alert">{actionError}</p>}
            {loading ? (
              <div className="habits-card habits-state" aria-live="polite">
                <span className="profile-skeleton" />
                <span className="profile-skeleton profile-skeleton--short" />
              </div>
            ) : loadError ? (
              <div className="habits-card habits-state">
                <h2>Не удалось загрузить привычки</h2>
                <p>{loadError}</p>
                <button className="button button--primary" type="button" onClick={() => setReloadKey((value) => value + 1)}>
                  Повторить
                </button>
              </div>
            ) : habits.length === 0 ? (
              <div className="habits-card habits-state">
                <h2>Первая серия ждёт</h2>
                <p>Создайте привычку слева и отметьте её выполненной сегодня.</p>
              </div>
            ) : (
              <div className="habit-list">
                {habits.map((habit) => (
                  <article
                    className={`habits-card habit-item ${habit.completed_today ? "habit-item--completed" : ""}`}
                    key={habit.id}
                  >
                    <div className="habit-item-copy">
                      <span className="row-kicker">
                        {habit.completed_today ? "Сегодня выполнено" : "Ждёт выполнения"}
                      </span>
                      <h2>{habit.name}</h2>
                      {habit.description && <p>{habit.description}</p>}
                    </div>
                    <div className="habit-streak">
                      <strong>{habit.current_streak}</strong>
                      <span>{streakLabel(habit.current_streak)} подряд</span>
                    </div>
                    <div className="habit-actions">
                      <button
                        className="button button--primary habit-complete-button"
                        type="button"
                        disabled={habit.completed_today || busyHabitId === habit.id}
                        onClick={() => void completeHabit(habit.id)}
                      >
                        {habit.completed_today
                          ? "Выполнено сегодня"
                          : busyHabitId === habit.id
                            ? "Отмечаем…"
                            : "Выполнить сегодня"}
                      </button>
                      <button
                        className="text-button habit-delete-button"
                        type="button"
                        disabled={busyHabitId === habit.id}
                        onClick={() => void deleteHabit(habit)}
                      >
                        Удалить
                      </button>
                    </div>
                  </article>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </section>
  );
}
