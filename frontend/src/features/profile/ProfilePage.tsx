import { useEffect, useState, type FormEvent } from "react";
import { apiClient } from "../../api/client";
import type { PatchUserRequest, User } from "../../api/types";
import "./profile.css";

function getInitials(name: string) {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toLocaleUpperCase("ru"))
    .join("");
}

function currentLocalDate() {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, "0");
  const day = String(today.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export function ProfilePage({
  user,
  token,
  loading,
  onLogin,
  onUserUpdated,
}: {
  user: User | null;
  token: string | null;
  loading: boolean;
  onLogin: () => void;
  onUserUpdated: (user: User) => void;
}) {
  const [sex, setSex] = useState("");
  const [weightKg, setWeightKg] = useState("");
  const [birthDate, setBirthDate] = useState("");
  const [heightCM, setHeightCM] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    setSex(user?.sex ?? "");
    setWeightKg(user?.weight_grams == null ? "" : String(user.weight_grams / 1000));
    setBirthDate(user?.birth_date?.slice(0, 10) ?? "");
    setHeightCM(user?.height_cm == null ? "" : String(user.height_cm));
  }, [user]);

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token) return;

    const weight = weightKg === "" ? null : Number(weightKg);
    const height = heightCM === "" ? null : Number(heightCM);
    if ((weight !== null && !Number.isFinite(weight)) || (height !== null && !Number.isFinite(height))) {
      setError("Проверьте числовые значения веса и роста.");
      return;
    }

    const payload: PatchUserRequest = {
      sex: sex === "male" || sex === "female" ? sex : null,
      weight_grams: weight === null ? null : Math.round(weight * 1000),
      birth_date: birthDate || null,
      height_cm: height === null ? null : Math.round(height),
    };

    setSubmitting(true);
    setError(null);
    setSuccess(null);
    try {
      const updated = await apiClient.patchCurrentUser(payload, token);
      onUserUpdated(updated);
      setSuccess("Профиль сохранён.");
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось сохранить профиль. Повторите попытку.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="profile-page">
      <div className="page-heading">
        <div>
          <span className="eyebrow">Личные параметры</span>
          <h1>Профиль</h1>
          <p>Заполните данные, которые используются для вашего личного прогресса.</p>
        </div>
      </div>

      {loading ? (
        <div className="profile-card profile-state" aria-live="polite">
          <span className="profile-skeleton" />
          <span className="profile-skeleton profile-skeleton--short" />
        </div>
      ) : !user || !token ? (
        <div className="profile-card profile-state">
          <h2>Сначала войдите</h2>
          <p>Профиль и личные параметры доступны после авторизации.</p>
          <button className="button button--primary" type="button" onClick={onLogin}>
            Войти
          </button>
        </div>
      ) : (
        <div className="profile-layout">
          <aside className="profile-card profile-summary-card">
            <span className="profile-avatar" aria-hidden="true">{getInitials(user.full_name)}</span>
            <div>
              <span className="profile-label">Участник</span>
              <h2>{user.full_name}</h2>
              <p>{user.email}</p>
            </div>
            <dl className="profile-stats">
              <div>
                <dt>Общий счёт</dt>
                <dd>{user.user_workout_score}</dd>
              </div>
              <div>
                <dt>В проекте с</dt>
                <dd>{new Intl.DateTimeFormat("ru-RU").format(new Date(user.created_at))}</dd>
              </div>
            </dl>
          </aside>

          <div className="profile-card profile-form-card">
            <div className="profile-form-heading">
              <h2>Параметры</h2>
              <p>Пустое поле удалит сохранённое значение.</p>
            </div>
            <form className="profile-form" onSubmit={submit}>
              <label>
                <span>Пол</span>
                <select name="sex" value={sex} onChange={(event) => setSex(event.target.value)}>
                  <option value="">Не указан</option>
                  <option value="male">Мужской</option>
                  <option value="female">Женский</option>
                </select>
              </label>
              <label>
                <span>Вес, кг</span>
                <input
                  name="weight"
                  type="number"
                  inputMode="decimal"
                  min="20"
                  max="300"
                  step="0.1"
                  value={weightKg}
                  onChange={(event) => setWeightKg(event.target.value)}
                  placeholder="Например, 80"
                />
              </label>
              <label>
                <span>Дата рождения</span>
                <input
                  name="birth_date"
                  type="date"
                  min="1900-01-01"
                  max={currentLocalDate()}
                  value={birthDate}
                  onChange={(event) => setBirthDate(event.target.value)}
                />
              </label>
              <label>
                <span>Рост, см</span>
                <input
                  name="height"
                  type="number"
                  inputMode="numeric"
                  min="100"
                  max="250"
                  step="1"
                  value={heightCM}
                  onChange={(event) => setHeightCM(event.target.value)}
                  placeholder="Например, 180"
                />
              </label>

              <div className="profile-form-actions">
                {error && <p className="form-error" role="alert">{error}</p>}
                {success && <p className="form-success" role="status">{success}</p>}
                <button className="button button--primary" type="submit" disabled={submitting}>
                  {submitting ? "Сохраняем…" : "Сохранить профиль"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </section>
  );
}
