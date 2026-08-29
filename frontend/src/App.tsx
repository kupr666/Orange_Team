import { useEffect, useMemo, useState, type FormEvent } from "react";
import { ACCESS_TOKEN_KEY, ApiError, apiClient } from "./api/client";
import type {
  LeaderboardEntry,
  LeaderboardPeriod,
  LeaderboardResponse,
  User,
} from "./api/types";
import { getDemoLeaderboard } from "./demo";
import { WorkoutsPage } from "./features/workouts/WorkoutsPage";

type Page = "habits" | "workouts" | "leaderboard" | "faq";
type AuthMode = "login" | "register";

const navigation: Array<{ id: Page; label: string }> = [
  { id: "habits", label: "Привычки" },
  { id: "workouts", label: "Тренировки" },
  { id: "leaderboard", label: "Лидерборд" },
  { id: "faq", label: "FAQ" },
];

const periods: Array<{ id: LeaderboardPeriod; label: string }> = [
  { id: "daily", label: "День" },
  { id: "weekly", label: "Неделя" },
  { id: "monthly", label: "Месяц" },
];

function getInitialToken() {
  try {
    return window.localStorage.getItem(ACCESS_TOKEN_KEY);
  } catch {
    return null;
  }
}

function getInitials(name: string) {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toLocaleUpperCase("ru"))
    .join("");
}

function completedLabel(value: number) {
  const lastTwo = value % 100;
  const last = value % 10;
  if (lastTwo >= 11 && lastTwo <= 14) return "тренировок";
  if (last === 1) return "тренировка";
  if (last >= 2 && last <= 4) return "тренировки";
  return "тренировок";
}

function periodRange(data: LeaderboardResponse | null, period: LeaderboardPeriod) {
  if (!data) {
    return period === "daily"
      ? "Сегодня"
      : period === "weekly"
        ? "Текущая неделя"
        : "Текущий месяц";
  }

  const formatter = new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "long",
    timeZone: data.timezone,
  });
  const start = formatter.format(new Date(data.period_start));
  const end = formatter.format(new Date(data.period_end));
  return start === end ? start : `${start} — ${end}`;
}

export function App() {
  const [page, setPage] = useState<Page>("leaderboard");
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [authMode, setAuthMode] = useState<AuthMode | null>(null);
  const [token, setToken] = useState<string | null>(getInitialToken);
  const [user, setUser] = useState<User | null>(null);
  const [profileLoading, setProfileLoading] = useState(Boolean(token));

  useEffect(() => {
    if (!token) {
      setUser(null);
      setProfileLoading(false);
      return;
    }

    let active = true;
    setProfileLoading(true);
    apiClient
      .getCurrentUser(token)
      .then((profile) => {
        if (active) setUser(profile);
      })
      .catch((error: unknown) => {
        if (!active) return;
        if (error instanceof ApiError && error.status === 401) {
          try {
            window.localStorage.removeItem(ACCESS_TOKEN_KEY);
          } catch {
            // The in-memory session can still be cleared if storage is unavailable.
          }
          setToken(null);
        }
        setUser(null);
      })
      .finally(() => {
        if (active) setProfileLoading(false);
      });

    return () => {
      active = false;
    };
  }, [token]);

  const navigate = (nextPage: Page) => {
    setPage(nextPage);
    setMobileMenuOpen(false);
  };

  const openAuth = (mode: AuthMode) => {
    setAuthMode(mode);
    setMobileMenuOpen(false);
  };

  const handleAuthenticated = (accessToken: string) => {
    try {
      window.localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    } catch {
      // Keeping the token in state still supports the current session.
    }
    setToken(accessToken);
    setAuthMode(null);
  };

  const logout = () => {
    try {
      window.localStorage.removeItem(ACCESS_TOKEN_KEY);
    } catch {
      // The in-memory session is always cleared below.
    }
    setToken(null);
    setUser(null);
    setMobileMenuOpen(false);
  };

  return (
    <div className="app-shell">
      <aside className="sidebar" aria-label="Панель пользователя">
        <div className="sidebar-top">
          <Brand />
          {(user || profileLoading) && (
            <UserPanel
              user={user}
              loading={profileLoading}
              onLogin={() => openAuth("login")}
              onRegister={() => openAuth("register")}
              onLogout={logout}
            />
          )}
        </div>
        {!user && !profileLoading && (
          <UserPanel
            user={null}
            loading={false}
            onLogin={() => openAuth("login")}
            onRegister={() => openAuth("register")}
            onLogout={logout}
          />
        )}
      </aside>

      <div className="workspace">
        <header className="desktop-header">
          <Navigation current={page} onNavigate={navigate} />
        </header>

        <header className="mobile-header">
          <Brand compact />
          <button
            className="mobile-menu-button"
            type="button"
            aria-expanded={mobileMenuOpen}
            aria-controls="mobile-menu"
            onClick={() => setMobileMenuOpen((open) => !open)}
          >
            {mobileMenuOpen ? "Закрыть" : "Меню"}
          </button>
        </header>

        {mobileMenuOpen && (
          <div className="mobile-menu-backdrop" onClick={() => setMobileMenuOpen(false)}>
            <section
              id="mobile-menu"
              className="mobile-menu"
              aria-label="Мобильное меню"
              onClick={(event) => event.stopPropagation()}
            >
              <Navigation current={page} onNavigate={navigate} vertical />
              <UserPanel
                user={user}
                loading={profileLoading}
                onLogin={() => openAuth("login")}
                onRegister={() => openAuth("register")}
                onLogout={logout}
                compact
              />
            </section>
          </div>
        )}

        <main className="content">
          {page === "leaderboard" ? (
            <LeaderboardPage user={user} token={token} onLogin={() => openAuth("login")} />
          ) : page === "workouts" ? (
            <WorkoutsPage user={user} token={token} onLogin={() => openAuth("login")} />
          ) : (
            <PlaceholderPage page={page} />
          )}
        </main>
      </div>

      {authMode && (
        <AuthDialog
          mode={authMode}
          onModeChange={setAuthMode}
          onClose={() => setAuthMode(null)}
          onAuthenticated={handleAuthenticated}
        />
      )}
    </div>
  );
}

function Brand({ compact = false }: { compact?: boolean }) {
  return (
    <div className={`brand ${compact ? "brand--compact" : ""}`} aria-label="Limit less">
      <span className="brand-name">Limit less</span>
      {!compact && <span className="brand-tagline">Превзойди предел</span>}
    </div>
  );
}

function Navigation({
  current,
  onNavigate,
  vertical = false,
}: {
  current: Page;
  onNavigate: (page: Page) => void;
  vertical?: boolean;
}) {
  return (
    <nav className={`navigation ${vertical ? "navigation--vertical" : ""}`} aria-label="Основная навигация">
      {navigation.map((item) => (
        <button
          key={item.id}
          className={`nav-link ${current === item.id ? "nav-link--active" : ""}`}
          type="button"
          aria-current={current === item.id ? "page" : undefined}
          onClick={() => onNavigate(item.id)}
        >
          {item.label}
        </button>
      ))}
    </nav>
  );
}

function UserPanel({
  user,
  loading,
  compact = false,
  onLogin,
  onRegister,
  onLogout,
}: {
  user: User | null;
  loading: boolean;
  compact?: boolean;
  onLogin: () => void;
  onRegister: () => void;
  onLogout: () => void;
}) {
  if (loading) {
    return (
      <div className={`user-panel ${compact ? "user-panel--compact" : ""}`} aria-label="Загрузка профиля">
        <div className="profile-skeleton" />
        <div className="profile-skeleton profile-skeleton--short" />
      </div>
    );
  }

  if (user) {
    return (
      <div className={`user-panel ${compact ? "user-panel--compact" : ""}`}>
        <div className="profile">
          <span className="avatar" aria-hidden="true">{getInitials(user.full_name)}</span>
          <div>
            <span className="profile-label">Участник</span>
            <strong className="profile-name">{user.full_name}</strong>
          </div>
        </div>
        <button className="text-button" type="button" onClick={onLogout}>Выйти</button>
      </div>
    );
  }

  return (
    <div className={`user-panel user-panel--guest ${compact ? "user-panel--compact" : ""}`}>
      <p className="guest-copy">Сохраняйте прогресс и находите своё место среди сильнейших.</p>
      <button className="button button--primary" type="button" onClick={onRegister}>Регистрация</button>
      <button className="button button--secondary" type="button" onClick={onLogin}>Войти</button>
    </div>
  );
}

function LeaderboardPage({
  user,
  token,
  onLogin,
}: {
  user: User | null;
  token: string | null;
  onLogin: () => void;
}) {
  const [period, setPeriod] = useState<LeaderboardPeriod>("weekly");
  const [data, setData] = useState<LeaderboardResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(null);

    if (!user) {
      const timer = window.setTimeout(() => {
        if (!active) return;
        setData(getDemoLeaderboard(period));
        setLoading(false);
      }, 420);
      return () => {
        active = false;
        window.clearTimeout(timer);
      };
    }

    apiClient
      .getLeaderboard(period, user.id, 50, token)
      .then((leaderboard) => {
        if (active) setData(leaderboard);
      })
      .catch((requestError: unknown) => {
        if (!active) return;
        setData(null);
        setError(
          requestError instanceof Error
            ? requestError.message
            : "Не удалось загрузить рейтинг. Попробуйте ещё раз.",
        );
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [period, refreshKey, token, user]);

  const visibleItems = useMemo(() => data?.items ?? [], [data]);
  const currentOutsideTop = Boolean(
    user && data?.current_user && !data.current_user.in_top,
  );

  return (
    <section className="leaderboard-page">
      <div className="page-heading">
        <div>
          <span className="eyebrow">Чертог достижений</span>
          <h1>Лидерборд</h1>
          <p>Рейтинг по количеству завершённых тренировок.</p>
        </div>
        <div className="period-summary">
          <span>Период</span>
          <strong>{periodRange(data, period)}</strong>
        </div>
      </div>

      <div className="leaderboard-toolbar">
        <div className="period-tabs" role="tablist" aria-label="Период рейтинга">
          {periods.map((item) => (
            <button
              key={item.id}
              type="button"
              role="tab"
              aria-selected={period === item.id}
              className={`period-tab ${period === item.id ? "period-tab--active" : ""}`}
              onClick={() => setPeriod(item.id)}
            >
              {item.label}
            </button>
          ))}
        </div>
        {!user && (
          <button className="login-hint" type="button" onClick={onLogin}>
            Войдите, чтобы увидеть своё место
          </button>
        )}
      </div>

      <div className="leaderboard-card" aria-busy={loading}>
        <div className="table-heading" aria-hidden="true">
          <span>Место</span>
          <span>Участник</span>
          <span>Завершено</span>
        </div>

        {loading ? (
          <LeaderboardLoading />
        ) : error ? (
          <StatePanel
            title="Лёд скрыл результаты"
            description={error}
            action="Повторить"
            onAction={() => setRefreshKey((value) => value + 1)}
          />
        ) : visibleItems.length === 0 ? (
          <StatePanel
            title="Пока здесь тихо"
            description="За выбранный период ещё никто не завершил тренировку. Станьте первым участником рейтинга."
          />
        ) : (
          <ol className="leaderboard-list">
            {visibleItems.map((entry, index) => (
              <LeaderboardRow
                key={entry.user.id}
                entry={entry}
                fallbackRank={index + 1}
                currentUserId={user?.id}
              />
            ))}
          </ol>
        )}

        {currentOutsideTop && data?.current_user && (
          <div className="current-user-summary">
            <span className="rank-number">{data.current_user.rank ?? "—"}</span>
            <div>
              <span className="row-kicker">Ваше место</span>
              <strong>{data.current_user.user.full_name}</strong>
            </div>
            <strong className="workout-count">
              {data.current_user.completed_workouts}
              <span>{completedLabel(data.current_user.completed_workouts)}</span>
            </strong>
          </div>
        )}

        {!loading && !error && visibleItems.length > 0 && (
          <footer className="leaderboard-footer">
            <span>{user ? "Ваш результат выделен ледяным сиянием" : "Показан демонстрационный рейтинг"}</span>
            <span>До 50 участников</span>
          </footer>
        )}
      </div>
    </section>
  );
}

function LeaderboardRow({
  entry,
  fallbackRank,
  currentUserId,
}: {
  entry: LeaderboardEntry;
  fallbackRank: number;
  currentUserId?: string;
}) {
  const rank = entry.rank ?? fallbackRank;
  const current = entry.is_current_user || entry.user.id === currentUserId;
  return (
    <li className={`leaderboard-row ${current ? "leaderboard-row--current" : ""}`}>
      <span className={`rank-number ${rank <= 3 ? `rank-number--${rank}` : ""}`}>{rank}</span>
      <div className="athlete">
        <span className="athlete-avatar" aria-hidden="true">{getInitials(entry.user.full_name)}</span>
        <div>
          {current && <span className="row-kicker">Это вы</span>}
          <strong>{entry.user.full_name}</strong>
        </div>
      </div>
      <strong className="workout-count">
        {entry.completed_workouts}
        <span>{completedLabel(entry.completed_workouts)}</span>
      </strong>
    </li>
  );
}

function LeaderboardLoading() {
  return (
    <div className="leaderboard-loading" aria-live="polite">
      <span className="visually-hidden">Загрузка рейтинга</span>
      {[0, 1, 2, 3, 4].map((item) => (
        <div className="loading-row" key={item}>
          <span className="loading-block loading-block--rank" />
          <span className="loading-block loading-block--name" />
          <span className="loading-block loading-block--score" />
        </div>
      ))}
    </div>
  );
}

function StatePanel({
  title,
  description,
  action,
  onAction,
}: {
  title: string;
  description: string;
  action?: string;
  onAction?: () => void;
}) {
  return (
    <div className="state-panel" role="status">
      <h2>{title}</h2>
      <p>{description}</p>
      {action && onAction && (
        <button className="button button--primary" type="button" onClick={onAction}>{action}</button>
      )}
    </div>
  );
}

function PlaceholderPage({ page }: { page: Exclude<Page, "leaderboard" | "workouts"> }) {
  const content = {
    habits: {
      eyebrow: "Ритм дня",
      title: "Привычки",
      text: "Раздел привычек появится в следующем обновлении. Здесь можно будет собирать ежедневные серии и следить за устойчивостью ритуалов.",
    },
    faq: {
      eyebrow: "Северные записи",
      title: "Частые вопросы",
      text: "Мы собираем короткие ответы о тренировках, рейтинге и прогрессе. Раздел будет открыт в следующей версии.",
    },
  }[page];

  return (
    <section className="placeholder-page">
      <span className="eyebrow">{content.eyebrow}</span>
      <h1>{content.title}</h1>
      <p>{content.text}</p>
      <span className="coming-soon">Скоро</span>
    </section>
  );
}

function AuthDialog({
  mode,
  onModeChange,
  onClose,
  onAuthenticated,
}: {
  mode: AuthMode;
  onModeChange: (mode: AuthMode) => void;
  onClose: () => void;
  onAuthenticated: (token: string) => void;
}) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") onClose();
    };
    window.addEventListener("keydown", handleEscape);
    return () => window.removeEventListener("keydown", handleEscape);
  }, [onClose]);

  const switchMode = (nextMode: AuthMode) => {
    setError(null);
    onModeChange(nextMode);
  };

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      if (mode === "register") {
        await apiClient.register({ email, password, full_name: fullName });
      }
      const result = await apiClient.login({ email, password });
      onAuthenticated(result.access_token);
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Не удалось продолжить. Проверьте данные и повторите попытку.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="dialog-backdrop" role="presentation" onMouseDown={onClose}>
      <section
        className="auth-dialog"
        role="dialog"
        aria-modal="true"
        aria-labelledby="auth-title"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <button className="dialog-close" type="button" onClick={onClose} aria-label="Закрыть окно">Закрыть</button>
        <span className="eyebrow">Врата Limit less</span>
        <h2 id="auth-title">{mode === "login" ? "С возвращением" : "Начните свой путь"}</h2>
        <p className="dialog-copy">
          {mode === "login"
            ? "Войдите, чтобы увидеть личное место в рейтинге."
            : "Создайте профиль — после регистрации вход выполнится автоматически."}
        </p>

        <form className="auth-form" onSubmit={submit}>
          {mode === "register" && (
            <label>
              <span>Имя</span>
              <input
                autoFocus
                name="full_name"
                autoComplete="name"
                minLength={2}
                maxLength={50}
                required
                value={fullName}
                onChange={(event) => setFullName(event.target.value)}
                placeholder="Как к вам обращаться"
              />
            </label>
          )}
          <label>
            <span>Электронная почта</span>
            <input
              autoFocus={mode === "login"}
              name="email"
              type="email"
              autoComplete="email"
              minLength={5}
              maxLength={30}
              required
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="name@example.ru"
            />
          </label>
          <label>
            <span>Пароль</span>
            <input
              name="password"
              type="password"
              autoComplete={mode === "login" ? "current-password" : "new-password"}
              minLength={8}
              maxLength={72}
              required
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="Не менее 8 символов"
            />
          </label>
          {error && <p className="form-error" role="alert">{error}</p>}
          <button className="button button--primary button--wide" type="submit" disabled={submitting}>
            {submitting ? "Открываем врата…" : mode === "login" ? "Войти" : "Зарегистрироваться"}
          </button>
        </form>

        <p className="auth-switch">
          {mode === "login" ? "Впервые здесь?" : "Уже есть профиль?"}{" "}
          <button type="button" onClick={() => switchMode(mode === "login" ? "register" : "login")}>
            {mode === "login" ? "Регистрация" : "Войти"}
          </button>
        </p>
      </section>
    </div>
  );
}
