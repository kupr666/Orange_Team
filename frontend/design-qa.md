# Design QA — Limit less

## Сравниваемые материалы

- Source visual truth: `/home/ilias/Orange_Team/Screenshot 2026-08-29 at 12.42.46 AM.png`
- Implementation screenshot: `/home/ilias/Orange_Team/frontend/implementation-desktop-source-size.png`
- Full-view comparison: `/home/ilias/Orange_Team/frontend/design-qa-comparison.png`
- Focused sidebar comparison: `/home/ilias/Orange_Team/frontend/design-qa-sidebar-focus.png`
- Mobile evidence: `/home/ilias/Orange_Team/frontend/implementation-mobile.png`
- Mobile menu evidence: `/home/ilias/Orange_Team/frontend/implementation-mobile-menu.png`
- Auth evidence: `/home/ilias/Orange_Team/frontend/implementation-auth-mobile.png`
- Workouts guest evidence: `/home/ilias/Orange_Team/frontend/implementation-workouts-guest.png`
- Workouts authenticated evidence: `/home/ilias/Orange_Team/frontend/implementation-workouts-authenticated.png`
- Workouts mobile evidence: `/home/ilias/Orange_Team/frontend/implementation-workouts-mobile.png`

## Нормализация

- Source pixels: 2248 × 1342, density 1×.
- Implementation pixels: 2248 × 1342, CSS viewport 2248 × 1342, `devicePixelRatio: 1`.
- Desktop state: гость, активная страница «Лидерборд», период «Неделя», демонстрационные данные.
- Mobile capture: фактический Firefox CSS viewport 500 × 758 при outer window 500 × 844, density 1×. Проверено отсутствие горизонтального overflow.

## Findings

Существенных P0/P1/P2-расхождений после финальной итерации нет.

- Fonts and typography: декоративная Georgia используется только в бренде, заголовках, навигации и крупных числах; основной интерфейс остаётся на системном sans-serif. Кириллица читаема, переносы и усечения не ломают контент.
- Spacing and layout rhythm: сохранены ключевые пропорции вайрфрейма — отдельная левая пользовательская панель, верхняя навигация и большая основная область. Карточка рейтинга и фильтры добавлены в предусмотренную пустую область согласно MVP-бррифу.
- Colors and visual tokens: тёмный камень, ледяной циан, серебристый текст и бронзовые акценты выдержаны последовательно; контраст основных действий и текста достаточен на фоновой иллюстрации.
- Image quality and asset fidelity: используется отдельный растровый фон `hyperborea-citadel.png` в нужной северной архитектурной стилистике; изображение резкое, без растяжения и заглушек.
- Copy and content: весь пользовательский текст на русском, кроме утверждённого имени продукта `Limit less` и общеупотребимого `FAQ`. Таблица показывает место, участника и завершённые тренировки.

## Comparison history

1. Первый desktop-pass выявил P2: при viewport 2248 px боковая панель занимала около 12% ширины вместо примерно 15% в вайрфрейме.
2. Исправление: desktop grid изменён на `clamp(272px, 15vw, 340px)`, сохранив удобную ширину на обычных ноутбуках и приблизив крупный viewport к референсу.
3. Повторный browser capture 2248 × 1342 и focused sidebar comparison подтвердили совпадение пропорции. P2 закрыт.

## Primary interactions tested

- Переключение «День» / «Неделя» / «Месяц»; месячный список обновился до 8 строк, первое значение — 72 тренировки.
- Открытие и закрытие мобильного меню.
- Переход на заглушку «Привычки» и обратно в «Лидерборд»; активная навигация обновилась.
- Открытие формы входа; присутствуют `role="dialog"`, заголовок и два обязательных поля.
- Адаптивность таблицы и отсутствие горизонтального переполнения на мобильном viewport.
- Vite error overlay отсутствует; TypeScript, production build и Sites worker test прошли без ошибок.
- Авторизованный раздел тренировок: загрузились профиль, две тренировки, каталог и состав активной тренировки.
- Полный workout-flow через API: создание плана, добавление упражнения, запуск, отметка упражнения выполненным и завершение с интенсивностью 7.
- Мобильный workout viewport 500 × 844: горизонтального overflow нет, список прокручивается по горизонтали, детали остаются в вертикальном потоке.

## Follow-up polish

- P3: проверить интерфейс дополнительно на физическом экране шириной 360–390 px; установленный Firefox ограничивает headless viewport минимумом 500 px, хотя CSS breakpoint и отсутствие overflow проверены программно.

final result: passed
