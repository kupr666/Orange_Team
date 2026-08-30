#!/usr/bin/env python3
"""Generate the Limit less hackathon deck as PNG, PPTX and PDF.

The slides are deliberately rendered as full-slide images. This keeps the
visual result identical in PowerPoint, LibreOffice and PDF viewers while the
talk track stays editable in the adjacent Markdown file.
"""

from __future__ import annotations

import math
import os
import shutil
import subprocess
import tempfile
import time
from pathlib import Path

from PIL import Image, ImageDraw, ImageEnhance, ImageFilter, ImageFont, ImageOps


ROOT = Path(__file__).resolve().parents[1]
PRESENTATION_DIR = ROOT / "presentation"
BUILD_DIR = PRESENTATION_DIR / "build"
SLIDES_DIR = BUILD_DIR / "slides"

W, H = 1920, 1080

BACKGROUND = ROOT / "frontend/public/assets/hyperborea-citadel.png"
WORKOUT_DESKTOP = ROOT / "frontend/implementation-workouts-authenticated.png"
WORKOUT_MOBILE = ROOT / "frontend/implementation-workouts-mobile.png"
AUTH_MOBILE = ROOT / "frontend/implementation-auth-mobile.png"
LEADERBOARD_DESKTOP = ROOT / "frontend/implementation-desktop.png"

FONT_SANS = Path("/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf")
FONT_SANS_BOLD = Path("/usr/share/fonts/truetype/noto/NotoSans-Bold.ttf")
FONT_SERIF = Path("/usr/share/fonts/truetype/noto/NotoSerifDisplay-Regular.ttf")
FONT_SERIF_BOLD = Path("/usr/share/fonts/truetype/noto/NotoSerifDisplay-Bold.ttf")

ICE = "#EAF7F8"
ICE_DIM = "#B9CCD2"
MUTED = "#8199A3"
CYAN = "#88DCE6"
CYAN_BRIGHT = "#B8F5F6"
CYAN_DARK = "#183A46"
GOLD = "#D3AC71"
GOLD_DIM = "#8B6C3F"
INK = "#050B12"
PANEL = "#08151E"
PANEL_LIGHT = "#0C202B"
LINE = "#29444F"
GREEN = "#79D7B0"
RED = "#D98A82"


def font(size: int, *, serif: bool = False, bold: bool = False) -> ImageFont.FreeTypeFont:
    if serif and bold:
        path = FONT_SERIF_BOLD
    elif serif:
        path = FONT_SERIF
    elif bold:
        path = FONT_SANS_BOLD
    else:
        path = FONT_SANS
    return ImageFont.truetype(str(path), size=size)


def rgb(value: str) -> tuple[int, int, int]:
    value = value.lstrip("#")
    return tuple(int(value[index : index + 2], 16) for index in (0, 2, 4))


def rgba(value: str, alpha: int = 255) -> tuple[int, int, int, int]:
    return (*rgb(value), alpha)


def cover(image: Image.Image, size: tuple[int, int]) -> Image.Image:
    return ImageOps.fit(image.convert("RGB"), size, method=Image.Resampling.LANCZOS)


def base_slide(number: int, *, scenic: bool = False) -> Image.Image:
    scene = cover(Image.open(BACKGROUND), (W, H))
    if scenic:
        scene = ImageEnhance.Brightness(scene).enhance(0.68)
        image = scene.convert("RGBA")
        overlay = Image.new("RGBA", (W, H), (0, 0, 0, 0))
        od = ImageDraw.Draw(overlay)
        for x in range(W):
            ratio = x / W
            alpha = int(225 * (1 - ratio) + 68 * ratio)
            od.line((x, 0, x, H), fill=(2, 8, 14, alpha))
        for y in range(H):
            alpha = int(28 + 112 * (y / H) ** 2)
            od.line((0, y, W, y), fill=(0, 5, 9, alpha))
        image = Image.alpha_composite(image, overlay)
    else:
        scene = ImageEnhance.Brightness(scene).enhance(0.24)
        scene = scene.filter(ImageFilter.GaussianBlur(1.2))
        image = scene.convert("RGBA")
        tint = Image.new("RGBA", (W, H), rgba("#030A11", 205))
        image = Image.alpha_composite(image, tint)
        glow = Image.new("RGBA", (W, H), (0, 0, 0, 0))
        gd = ImageDraw.Draw(glow)
        gd.ellipse((1250, -220, 2100, 640), fill=rgba("#0B6877", 42))
        gd.ellipse((-260, 700, 600, 1400), fill=rgba("#174A59", 30))
        glow = glow.filter(ImageFilter.GaussianBlur(95))
        image = Image.alpha_composite(image, glow)

    draw = ImageDraw.Draw(image)
    draw.line((80, 84, W - 80, 84), fill=rgba(LINE, 120), width=1)
    draw.text((86, 34), "Limit less", font=font(29, serif=True, bold=True), fill=ICE)
    draw.text((W - 166, 38), f"{number:02d} / 10", font=font(21, bold=True), fill=MUTED)
    return image


def wrap_lines(draw: ImageDraw.ImageDraw, text: str, face: ImageFont.FreeTypeFont, max_width: int) -> list[str]:
    paragraphs = text.split("\n")
    lines: list[str] = []
    for paragraph in paragraphs:
        if not paragraph:
            lines.append("")
            continue
        words = paragraph.split()
        current = words[0]
        for word in words[1:]:
            candidate = f"{current} {word}"
            if draw.textbbox((0, 0), candidate, font=face)[2] <= max_width:
                current = candidate
            else:
                lines.append(current)
                current = word
        lines.append(current)
    return lines


def text_block(
    image: Image.Image,
    xy: tuple[int, int],
    text: str,
    *,
    face: ImageFont.FreeTypeFont,
    fill: str = ICE,
    max_width: int,
    line_gap: int = 10,
    anchor: str = "la",
) -> int:
    draw = ImageDraw.Draw(image)
    lines = wrap_lines(draw, text, face, max_width)
    x, y = xy
    ascent, descent = face.getmetrics()
    line_height = ascent + descent + line_gap
    for line in lines:
        draw.text((x, y), line, font=face, fill=fill, anchor=anchor)
        y += line_height
    return y


def title_block(
    image: Image.Image,
    eyebrow: str,
    title: str,
    *,
    subtitle: str | None = None,
    title_size: int = 65,
    max_width: int = 1500,
) -> int:
    draw = ImageDraw.Draw(image)
    draw.text((88, 120), eyebrow.upper(), font=font(20, bold=True), fill=GOLD)
    y = text_block(
        image,
        (84, 164),
        title,
        face=font(title_size, serif=True, bold=True),
        fill=ICE,
        max_width=max_width,
        line_gap=0,
    )
    if subtitle:
        y += 12
        y = text_block(
            image,
            (88, y),
            subtitle,
            face=font(28),
            fill=ICE_DIM,
            max_width=max_width,
            line_gap=8,
        )
    return y


def rounded_panel(
    image: Image.Image,
    box: tuple[int, int, int, int],
    *,
    fill: str = PANEL,
    alpha: int = 224,
    outline: str = LINE,
    radius: int = 24,
    shadow: bool = True,
) -> None:
    x1, y1, x2, y2 = box
    if shadow:
        layer = Image.new("RGBA", image.size, (0, 0, 0, 0))
        ld = ImageDraw.Draw(layer)
        ld.rounded_rectangle((x1 + 10, y1 + 16, x2 + 10, y2 + 16), radius=radius, fill=(0, 0, 0, 115))
        layer = layer.filter(ImageFilter.GaussianBlur(18))
        image.alpha_composite(layer)
    draw = ImageDraw.Draw(image)
    draw.rounded_rectangle(box, radius=radius, fill=rgba(fill, alpha), outline=rgba(outline, 210), width=2)
    draw.line((x1 + 28, y1 + 1, x2 - 28, y1 + 1), fill=rgba("#FFFFFF", 20), width=1)


def chip(
    image: Image.Image,
    xy: tuple[int, int],
    label: str,
    *,
    color: str = CYAN,
    fill: str = PANEL_LIGHT,
    height: int = 44,
    size: int = 18,
) -> int:
    draw = ImageDraw.Draw(image)
    face = font(size, bold=True)
    width = draw.textbbox((0, 0), label, font=face)[2] + 36
    x, y = xy
    draw.rounded_rectangle((x, y, x + width, y + height), radius=height // 2, fill=rgba(fill, 220), outline=rgba(color, 160), width=2)
    draw.text((x + width // 2, y + height // 2 - 1), label, font=face, fill=color, anchor="mm")
    return width


def footer(image: Image.Image, criterion: str, timing: str) -> None:
    draw = ImageDraw.Draw(image)
    draw.line((84, H - 74, W - 84, H - 74), fill=rgba(LINE, 145), width=1)
    draw.text((88, H - 52), criterion.upper(), font=font(18, bold=True), fill=CYAN)
    draw.text((W - 88, H - 52), timing.replace("—", "-"), font=font(18, bold=True), fill=MUTED, anchor="ra")


def fit_image(source: Image.Image, size: tuple[int, int], mode: str = "cover") -> Image.Image:
    if mode == "contain":
        result = Image.new("RGBA", size, rgba(INK))
        fitted = ImageOps.contain(source.convert("RGBA"), size, method=Image.Resampling.LANCZOS)
        x = (size[0] - fitted.width) // 2
        y = (size[1] - fitted.height) // 2
        result.alpha_composite(fitted, (x, y))
        return result
    return ImageOps.fit(source.convert("RGBA"), size, method=Image.Resampling.LANCZOS)


def screen(
    image: Image.Image,
    source_path: Path,
    box: tuple[int, int, int, int],
    *,
    crop_box: tuple[int, int, int, int] | None = None,
    mode: str = "cover",
    radius: int = 22,
    border: str = "#42616B",
) -> None:
    source = Image.open(source_path).convert("RGBA")
    if crop_box:
        source = source.crop(crop_box)
    x1, y1, x2, y2 = box
    width, height = x2 - x1, y2 - y1
    fitted = fit_image(source, (width, height), mode=mode)
    mask = Image.new("L", (width, height), 0)
    md = ImageDraw.Draw(mask)
    md.rounded_rectangle((0, 0, width - 1, height - 1), radius=radius, fill=255)

    shadow_layer = Image.new("RGBA", image.size, (0, 0, 0, 0))
    sd = ImageDraw.Draw(shadow_layer)
    sd.rounded_rectangle((x1 + 12, y1 + 18, x2 + 12, y2 + 18), radius=radius, fill=(0, 0, 0, 145))
    shadow_layer = shadow_layer.filter(ImageFilter.GaussianBlur(20))
    image.alpha_composite(shadow_layer)
    image.paste(fitted, (x1, y1), mask)
    draw = ImageDraw.Draw(image)
    draw.rounded_rectangle(box, radius=radius, outline=rgba(border, 235), width=2)


def label_number(image: Image.Image, x: int, y: int, number: str, label: str, *, width: int = 470) -> None:
    draw = ImageDraw.Draw(image)
    draw.ellipse((x, y, x + 56, y + 56), fill=rgba(CYAN_DARK, 245), outline=rgba(CYAN, 190), width=2)
    draw.text((x + 28, y + 28), number, font=font(19, bold=True), fill=CYAN_BRIGHT, anchor="mm")
    text_block(image, (x + 78, y + 4), label, face=font(25, bold=True), fill=ICE, max_width=width - 78, line_gap=4)


def arrow(image: Image.Image, start: tuple[int, int], end: tuple[int, int], *, color: str = CYAN) -> None:
    draw = ImageDraw.Draw(image)
    draw.line((*start, *end), fill=rgba(color, 205), width=4)
    angle = math.atan2(end[1] - start[1], end[0] - start[0])
    length = 17
    for delta in (2.55, -2.55):
        point = (end[0] + length * math.cos(angle + delta), end[1] + length * math.sin(angle + delta))
        draw.line((*end, *point), fill=rgba(color, 205), width=4)


def slide_01() -> Image.Image:
    image = base_slide(1, scenic=True)
    draw = ImageDraw.Draw(image)
    draw.text((104, 146), "HACKATHON · ORANGE TEAM · 10 МИНУТ", font=font(21, bold=True), fill=GOLD)
    draw.text((96, 270), "Limit less", font=font(124, serif=True, bold=True), fill=ICE)
    text_block(
        image,
        (104, 436),
        "Тренировка превращается в измеримый прогресс —\nи в причину вернуться завтра.",
        face=font(43),
        fill=ICE_DIM,
        max_width=1060,
        line_gap=9,
    )
    x = 104
    for label in ("ТРЕНИРОВКА", "ПЕРСОНАЛЬНЫЙ SCORE", "ЛИДЕРБОРД"):
        x += chip(image, (x, 620), label, color=CYAN, fill="#0B222C", height=48, size=18) + 14
    rounded_panel(image, (1370, 706, 1780, 900), fill="#07151D", alpha=210, outline=GOLD_DIM, radius=22)
    draw.text((1410, 744), "РАБОЧИЙ MVP", font=font(18, bold=True), fill=GOLD)
    draw.text((1410, 795), "1 сквозной\nсценарий", font=font(34, serif=True, bold=True), fill=ICE, spacing=10)
    draw.text((105, 980), "Планируй · выполняй · сравнивай свой прогресс", font=font(23, bold=True), fill=CYAN_BRIGHT)
    return image


def slide_02() -> Image.Image:
    image = base_slide(2)
    title_block(
        image,
        "Проблема / необходимость",
        "Обратная связь теряется между намерением и результатом",
        subtitle="Пользователю нужен не ещё один список упражнений, а понятный повод продолжать.",
        max_width=1550,
    )
    cards = [
        ("01", "Разрозненно", "План, упражнения и результат живут в разных местах."),
        ("02", "Непонятно", "Вес и повторы сами по себе не складываются в общий прогресс."),
        ("03", "Одиноко", "Без видимого места среди других сложнее поддерживать ритм."),
    ]
    x_positions = (84, 688, 1292)
    for x, (number, heading, body) in zip(x_positions, cards):
        rounded_panel(image, (x, 420, x + 544, 756), fill=PANEL, alpha=230)
        draw = ImageDraw.Draw(image)
        draw.text((x + 38, 454), number, font=font(22, bold=True), fill=GOLD)
        draw.text((x + 38, 514), heading, font=font(39, serif=True, bold=True), fill=ICE)
        text_block(image, (x + 38, 584), body, face=font(26), fill=ICE_DIM, max_width=468, line_gap=8)

    rounded_panel(image, (84, 806, 1836, 954), fill="#0A2029", alpha=222, outline=CYAN_DARK, radius=20, shadow=False)
    draw = ImageDraw.Draw(image)
    draw.text((126, 844), "НУЖЕН ОДИН КОРОТКИЙ ЦИКЛ", font=font(18, bold=True), fill=CYAN)
    draw.text((126, 889), "действие", font=font(30, bold=True), fill=ICE)
    arrow(image, (292, 906), (430, 906))
    draw.text((468, 889), "понятный результат", font=font(30, bold=True), fill=ICE)
    arrow(image, (808, 906), (946, 906))
    draw.text((984, 889), "мотивация продолжить", font=font(30, bold=True), fill=ICE)
    footer(image, "Необходимость проекта", "00:25 — 01:10")
    return image


def slide_03() -> Image.Image:
    image = base_slide(3)
    title_block(
        image,
        "Решение задачи",
        "Limit less замыкает цикл прогресса",
        subtitle="Каждое действие пользователя сразу связано со следующим полезным результатом.",
    )
    nodes = [
        ("01", "Планирую", "Создаю тренировку"),
        ("02", "Выполняю", "Подходы или время"),
        ("03", "Получаю score", "Сервер считает результат"),
        ("04", "Вижу место", "День · неделя · месяц"),
    ]
    xs = [96, 536, 976, 1416]
    y1, y2 = 435, 716
    for idx, (x, item) in enumerate(zip(xs, nodes)):
        number, heading, body = item
        rounded_panel(image, (x, y1, x + 380, y2), fill=PANEL, alpha=232, outline=LINE, radius=22)
        draw = ImageDraw.Draw(image)
        draw.ellipse((x + 32, y1 + 32, x + 92, y1 + 92), fill=rgba(CYAN_DARK), outline=rgba(CYAN, 170), width=2)
        draw.text((x + 62, y1 + 62), number, font=font(19, bold=True), fill=CYAN_BRIGHT, anchor="mm")
        draw.text((x + 32, y1 + 124), heading, font=font(31, serif=True, bold=True), fill=ICE)
        text_block(image, (x + 32, y1 + 182), body, face=font(23), fill=ICE_DIM, max_width=316, line_gap=5)
        if idx < 3:
            arrow(image, (x + 385, (y1 + y2) // 2), (xs[idx + 1] - 10, (y1 + y2) // 2))

    rounded_panel(image, (340, 780, 1580, 930), fill="#0B222C", alpha=226, outline=CYAN_DARK, radius=22, shadow=False)
    draw = ImageDraw.Draw(image)
    draw.text((384, 815), "ЦЕННОСТЬ", font=font(18, bold=True), fill=GOLD)
    draw.text((384, 858), "Пользователь видит не архив записей, а движение вперёд после каждого занятия.", font=font(29, bold=True), fill=ICE)
    footer(image, "Решение задачи", "01:10 — 01:55")
    return image


def slide_04() -> Image.Image:
    image = base_slide(4)
    title_block(
        image,
        "Демонстрация работы / пункт 5.2",
        "От новой тренировки — до места в рейтинге",
        max_width=1500,
    )

    rounded_panel(image, (82, 367, 456, 897), fill="#050D14", alpha=244, radius=30)
    screen(image, AUTH_MOBILE, (114, 391, 424, 824), crop_box=(18, 130, 470, 725), mode="cover", radius=18)
    draw = ImageDraw.Draw(image)
    draw.text((104, 848), "01  ВХОД", font=font(19, bold=True), fill=GOLD)
    draw.text((104, 878), "Аккаунт и профиль", font=font(22, bold=True), fill=ICE)

    screen(image, WORKOUT_DESKTOP, (490, 367, 1210, 826), crop_box=(268, 80, 1430, 960), mode="cover", radius=20)
    draw.text((512, 848), "02  ТРЕНИРОВКА", font=font(19, bold=True), fill=GOLD)
    draw.text((512, 878), "Нагрузка / выполнение / score", font=font(22, bold=True), fill=ICE)

    # The saved leaderboard capture predates score labels. Crop out the old
    # right-hand metric and use it only to prove rank/period/navigation.
    screen(image, LEADERBOARD_DESKTOP, (1244, 367, 1838, 826), crop_box=(270, 90, 1110, 914), mode="cover", radius=20)
    draw.text((1266, 848), "03  РЕЙТИНГ", font=font(19, bold=True), fill=GOLD)
    draw.text((1266, 878), "Позиция за период", font=font(22, bold=True), fill=ICE)

    draw.rounded_rectangle((1575, 126, 1838, 225), radius=20, fill=rgba("#0B2C36", 238), outline=rgba(CYAN, 170), width=2)
    draw.text((1706, 158), "LIVE", font=font(20, bold=True), fill=CYAN, anchor="mm")
    draw.text((1706, 197), "02:30", font=font(34, bold=True), fill=ICE, anchor="mm")
    footer(image, "Демонстрация работы", "01:55 — 04:25")
    return image


def slide_05() -> Image.Image:
    image = base_slide(5)
    title_block(
        image,
        "Работоспособность",
        "Рабочее ядро уже собрано",
        subtitle="Не кликабельный макет: интерфейс связан с Go API и PostgreSQL.",
        max_width=1380,
    )

    screen(image, WORKOUT_DESKTOP, (84, 360, 1230, 898), crop_box=(0, 70, 1440, 960), mode="cover", radius=24)
    rounded_panel(image, (1100, 440, 1374, 904), fill="#03090E", alpha=248, outline=CYAN_DARK, radius=35)
    screen(image, WORKOUT_MOBILE, (1125, 460, 1349, 870), crop_box=(0, 0, 500, 844), mode="cover", radius=27, border=CYAN_DARK)

    rounded_panel(image, (1410, 360, 1838, 898), fill=PANEL, alpha=232, radius=24)
    draw = ImageDraw.Draw(image)
    draw.text((1450, 396), "РЕАЛИЗОВАНО", font=font(18, bold=True), fill=GOLD)
    features = [
        "Регистрация, вход и JWT",
        "Профиль и личные параметры",
        "Жизненный цикл тренировки",
        "Силовые и временные упражнения",
        "Автопересчёт score",
        "Рейтинг: день / неделя / месяц",
        "Адаптивный интерфейс",
    ]
    y = 452
    for feature in features:
        draw.ellipse((1450, y + 3, 1473, y + 26), fill=rgba(CYAN_DARK), outline=rgba(CYAN, 180), width=2)
        draw.line((1457, y + 15, 1462, y + 21), fill=rgba(CYAN_BRIGHT), width=3)
        draw.line((1462, y + 21, 1470, y + 10), fill=rgba(CYAN_BRIGHT), width=3)
        text_block(image, (1492, y), feature, face=font(22, bold=True), fill=ICE, max_width=302, line_gap=4)
        y += 61
    footer(image, "Работоспособность", "04:25 — 05:10")
    return image


def slide_06() -> Image.Image:
    image = base_slide(6)
    title_block(
        image,
        "Инновационность",
        "Сравниваем усилие, а не только записи",
        subtitle="score v1 пересчитывается на сервере после изменения выполненного упражнения.",
    )

    rounded_panel(image, (82, 365, 1838, 650), fill=PANEL, alpha=234, radius=26)
    draw = ImageDraw.Draw(image)
    stages = [
        ("ВЫПОЛНЕНО", "вес × подходы × повторы\nили длительность"),
        ("СЛОЖНОСТЬ", "коэффициент\nупражнения 1–10"),
        ("ПРОФИЛЬ", "персональный\nкоэффициент 1–10"),
        ("SCORE", "нелинейная шкала\nс убывающей отдачей"),
        ("РЕЙТИНГ", "live день + snapshots\nнеделя / месяц"),
    ]
    stage_x = [120, 472, 824, 1176, 1528]
    for idx, ((label, body), x) in enumerate(zip(stages, stage_x)):
        draw.ellipse((x, 410, x + 84, 494), fill=rgba(CYAN_DARK), outline=rgba(CYAN, 180), width=2)
        draw.text((x + 42, 452), str(idx + 1), font=font(25, bold=True), fill=CYAN_BRIGHT, anchor="mm")
        draw.text((x, 520), label, font=font(17, bold=True), fill=GOLD)
        text_block(image, (x, 554), body, face=font(20, bold=True), fill=ICE, max_width=250, line_gap=3)
        if idx < len(stages) - 1:
            arrow(image, (x + 100, 451), (stage_x[idx + 1] - 18, 451))

    rounded_panel(image, (82, 698, 1138, 932), fill="#0A2029", alpha=226, outline=CYAN_DARK, radius=24)
    draw.text((122, 730), "НЕЛИНЕЙНАЯ ШКАЛА", font=font(18, bold=True), fill=GOLD)
    draw.text((122, 785), "score: x / (K + x)", font=font(40, bold=True), fill=ICE)
    draw.text((122, 846), "x = эффективная нагрузка × персональный коэффициент", font=font(23), fill=ICE_DIM)
    draw.text((122, 887), "Насыщение ограничивает влияние экстремальных значений.", font=font(21, bold=True), fill=CYAN_BRIGHT)

    rounded_panel(image, (1174, 698, 1838, 932), fill=PANEL, alpha=232, radius=24)
    draw.text((1216, 730), "ПОЧЕМУ ЭТО ВАЖНО", font=font(18, bold=True), fill=GOLD)
    reasons = ["Считается только выполненная работа", "Результат персонализирован", "Периоды фиксируются снимками"]
    y = 784
    for reason in reasons:
        draw.ellipse((1216, y + 5, 1234, y + 23), fill=rgba(CYAN, 220))
        text_block(image, (1254, y), reason, face=font(22, bold=True), fill=ICE, max_width=520, line_gap=2)
        y += 50
    footer(image, "Инновационность", "05:10 — 06:05")
    return image


def slide_07() -> Image.Image:
    image = base_slide(7)
    title_block(
        image,
        "Стек технологий / пункт 5.2",
        "Сквозная система: интерфейс, API и данные",
        subtitle="Разделение слоёв позволяет независимо развивать продукт и бизнес-логику.",
    )

    y1, y2 = 405, 690
    architecture = [
        (110, 430, "WEB", "React 19\nTypeScript\nVite"),
        (530, 430, "HTTP API", "Go 1.25\nnet/http\nOpenAPI 3"),
        (950, 430, "CORE", "transport\nservice\nrepository"),
        (1370, 430, "DATA", "PostgreSQL 18\npgx\nmigrations"),
    ]
    for idx, (x, y, label, body) in enumerate(architecture):
        rounded_panel(image, (x, y, x + 340, y + 250), fill=PANEL, alpha=234, radius=24)
        draw = ImageDraw.Draw(image)
        draw.text((x + 32, y + 28), label, font=font(18, bold=True), fill=GOLD)
        text_block(image, (x + 32, y + 78), body, face=font(29, bold=True), fill=ICE, max_width=270, line_gap=6)
        if idx < len(architecture) - 1:
            arrow(image, (x + 350, y + 125), (architecture[idx + 1][0] - 10, y + 125))

    rounded_panel(image, (110, 748, 1710, 910), fill="#0A2029", alpha=226, outline=CYAN_DARK, radius=22, shadow=False)
    draw = ImageDraw.Draw(image)
    stack_items = [
        ("AUTH", "JWT + роли"),
        ("BACKGROUND", "worker snapshots"),
        ("INFRA", "Docker Compose"),
        ("OBSERVABILITY", "Loki + Alloy"),
    ]
    x = 148
    for label, value in stack_items:
        draw.text((x, 780), label, font=font(16, bold=True), fill=MUTED)
        draw.text((x, 822), value, font=font(24, bold=True), fill=ICE)
        x += 386
    footer(image, "Стек технологий", "06:05 — 06:55")
    return image


def slide_08() -> Image.Image:
    image = base_slide(8)
    title_block(
        image,
        "Доказательства работоспособности",
        "Надёжность заложена в реализацию",
        subtitle="Проверяем то, что реально подключено к основному приложению — без завышенных обещаний.",
    )

    tests = [
        ("FRONTEND", "PASS", "TypeScript typecheck\nVite production build\nSites worker test"),
        ("GO CORE", "PASS", "cmd/api + подключённые модули\nleaderboard service/transport tests"),
        ("API CONTRACT", "v1", "OpenAPI 3\n6 миграций\nограничения БД"),
    ]
    xs = [84, 680, 1276]
    for x, (label, status, body) in zip(xs, tests):
        rounded_panel(image, (x, 385, x + 560, 650), fill=PANEL, alpha=234, radius=24)
        draw = ImageDraw.Draw(image)
        draw.text((x + 36, 420), label, font=font(18, bold=True), fill=GOLD)
        status_color = GREEN if status == "PASS" else CYAN
        draw.rounded_rectangle((x + 410, 408, x + 524, 455), radius=20, fill=rgba("#123328", 230), outline=rgba(status_color, 150), width=2)
        draw.text((x + 467, 431), status, font=font(18, bold=True), fill=status_color, anchor="mm")
        text_block(image, (x + 36, 494), body, face=font(24, bold=True), fill=ICE, max_width=482, line_gap=8)

    rounded_panel(image, (84, 705, 1836, 936), fill="#0A2029", alpha=226, outline=CYAN_DARK, radius=24)
    draw = ImageDraw.Draw(image)
    draw.text((124, 742), "ЗАЩИТА И ЦЕЛОСТНОСТЬ", font=font(18, bold=True), fill=GOLD)
    safeguards = [
        ("JWT", "аутентификация и роли"),
        ("OWNERSHIP", "ID владельца из токена"),
        ("SQL", "только параметры $1…$n"),
        ("TIMEOUTS", "границы операций БД"),
        ("VERSION", "optimistic locking"),
    ]
    x = 124
    for label, body in safeguards:
        draw.text((x, 800), label, font=font(17, bold=True), fill=CYAN)
        text_block(image, (x, 836), body, face=font(20, bold=True), fill=ICE, max_width=280, line_gap=2)
        x += 338
    footer(image, "Работоспособность", "06:55 — 07:45")
    return image


def slide_09() -> Image.Image:
    image = base_slide(9)
    title_block(
        image,
        "Развитие проекта",
        "Следующий шаг — доказать, что цикл удерживает",
        subtitle="MVP проверяет механику. Пилот должен проверить ценность для реальных пользователей.",
    )

    rounded_panel(image, (84, 380, 820, 900), fill=PANEL, alpha=232, radius=24)
    draw = ImageDraw.Draw(image)
    draw.text((126, 418), "СЕЙЧАС · MVP", font=font(18, bold=True), fill=GOLD)
    now_items = [
        "Профиль и авторизация",
        "Тренировка от плана до завершения",
        "score v1 и личный прогресс",
        "Рейтинги за три периода",
        "Desktop + mobile",
    ]
    y = 486
    for item in now_items:
        draw.ellipse((128, y + 4, 150, y + 26), fill=rgba(GREEN, 210))
        text_block(image, (174, y), item, face=font(25, bold=True), fill=ICE, max_width=580, line_gap=3)
        y += 76

    rounded_panel(image, (860, 380, 1836, 900), fill="#0A2029", alpha=230, outline=CYAN_DARK, radius=24)
    draw.text((904, 418), "ПИЛОТ · 4 НЕДЕЛИ", font=font(18, bold=True), fill=GOLD)
    draw.text((904, 470), "20–30 пользователей", font=font(38, serif=True, bold=True), fill=ICE)
    draw.text((904, 527), "Цель — измерить повторные тренировки, а не установки.", font=font(23), fill=ICE_DIM)
    pilot_items = [
        ("01", "Валидировать score v1 с тренером"),
        ("02", "Считать возврат к следующей тренировке"),
        ("03", "Добавить привычки, серии и командные вызовы"),
        ("04", "Подключить носимые устройства после проверки спроса"),
    ]
    y = 604
    for number, item in pilot_items:
        draw.text((904, y), number, font=font(19, bold=True), fill=CYAN)
        text_block(image, (958, y - 2), item, face=font(23, bold=True), fill=ICE, max_width=800, line_gap=4)
        y += 67

    footer(image, "Необходимость + решение", "07:45 — 08:35")
    return image


def slide_10() -> Image.Image:
    image = base_slide(10, scenic=True)
    draw = ImageDraw.Draw(image)
    draw.text((102, 150), "ФИНАЛ", font=font(21, bold=True), fill=GOLD)
    draw.text((96, 238), "Limit less", font=font(104, serif=True, bold=True), fill=ICE)
    text_block(
        image,
        (104, 390),
        "Мы не просто сохраняем тренировку.\nМы превращаем её в следующий повод продолжить.",
        face=font(41),
        fill=ICE_DIM,
        max_width=1160,
        line_gap=10,
    )

    summary = [
        ("РАБОТАЕТ", "сквозной MVP"),
        ("РЕШАЕТ", "разрыв обратной связи"),
        ("ОТЛИЧАЕТСЯ", "score v1 + snapshots"),
    ]
    x = 104
    for heading, body in summary:
        rounded_panel(image, (x, 610, x + 430, 758), fill="#07151D", alpha=218, outline=CYAN_DARK, radius=20, shadow=False)
        draw.text((x + 28, 642), heading, font=font(17, bold=True), fill=GOLD)
        draw.text((x + 28, 692), body, font=font(23, bold=True), fill=ICE)
        x += 454

    rounded_panel(image, (1360, 820, 1818, 944), fill="#07151D", alpha=214, outline=GOLD_DIM, radius=20, shadow=False)
    draw.text((1589, 855), "ORANGE TEAM", font=font(19, bold=True), fill=GOLD, anchor="mm")
    draw.text((1589, 903), "Спасибо · вопросы", font=font(25, bold=True), fill=ICE, anchor="mm")
    draw.text((104, 972), "Limit less · превзойди предел", font=font(23, bold=True), fill=CYAN_BRIGHT)
    return image


SLIDE_BUILDERS = [
    slide_01,
    slide_02,
    slide_03,
    slide_04,
    slide_05,
    slide_06,
    slide_07,
    slide_08,
    slide_09,
    slide_10,
]


def generate_slide_pngs() -> list[Path]:
    SLIDES_DIR.mkdir(parents=True, exist_ok=True)
    result: list[Path] = []
    for index, builder in enumerate(SLIDE_BUILDERS, start=1):
        image = builder().convert("RGB")
        path = SLIDES_DIR / f"{index:02d}.png"
        image.save(path, format="PNG", optimize=True)
        result.append(path)
        print(f"rendered {path.relative_to(ROOT)}")
    return result


def uno_property(name: str, value: object):
    from com.sun.star.beans import PropertyValue  # type: ignore

    prop = PropertyValue()
    prop.Name = name
    prop.Value = value
    return prop


def build_office_files(slide_paths: list[Path]) -> tuple[Path, Path]:
    """Create PPTX/PDF via a private headless LibreOffice instance."""

    import uno  # type: ignore
    from com.sun.star.awt import Point, Size  # type: ignore

    pptx_path = PRESENTATION_DIR / "limit-less-hackathon-10min.pptx"
    pdf_path = PRESENTATION_DIR / "limit-less-hackathon-10min.pdf"

    profile_dir = Path(tempfile.mkdtemp(prefix="limitless-impress-", dir="/tmp"))
    runtime_dir = profile_dir / "runtime"
    config_dir = profile_dir / "config"
    cache_dir = profile_dir / "cache"
    runtime_dir.mkdir(mode=0o700)
    config_dir.mkdir()
    cache_dir.mkdir()
    port = 2083
    accept = f"socket,host=127.0.0.1,port={port};urp;StarOffice.ComponentContext"
    office_env = os.environ.copy()
    office_env.update(
        {
            "XDG_RUNTIME_DIR": str(runtime_dir),
            "XDG_CONFIG_HOME": str(config_dir),
            "XDG_CACHE_HOME": str(cache_dir),
            "SAL_USE_VCLPLUGIN": "svp",
        }
    )
    process = subprocess.Popen(
        [
            shutil.which("soffice") or "soffice",
            "--headless",
            "--nologo",
            "--nodefault",
            "--nofirststartwizard",
            "--norestore",
            f"-env:UserInstallation={profile_dir.as_uri()}",
            f"--accept={accept}",
        ],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.PIPE,
        text=True,
        env=office_env,
    )

    document = None
    try:
        local_context = uno.getComponentContext()
        resolver = local_context.ServiceManager.createInstanceWithContext(
            "com.sun.star.bridge.UnoUrlResolver", local_context
        )
        context = None
        for _ in range(80):
            try:
                context = resolver.resolve(
                    f"uno:socket,host=127.0.0.1,port={port};urp;StarOffice.ComponentContext"
                )
                break
            except Exception:
                if process.poll() is not None:
                    stderr = process.stderr.read() if process.stderr else ""
                    raise RuntimeError(f"LibreOffice exited before connection: {stderr}")
                time.sleep(0.15)
        if context is None:
            raise RuntimeError("Timed out connecting to LibreOffice")

        service_manager = context.ServiceManager
        desktop = service_manager.createInstanceWithContext("com.sun.star.frame.Desktop", context)
        document = desktop.loadComponentFromURL(
            "private:factory/simpress",
            "_blank",
            0,
            (uno_property("Hidden", True),),
        )

        pages = document.getDrawPages()
        provider = service_manager.createInstanceWithContext(
            "com.sun.star.graphic.GraphicProvider", context
        )
        slide_width = 33867  # 13.333 in in 1/100 mm
        slide_height = 19050  # 7.5 in in 1/100 mm

        for index, slide_path in enumerate(slide_paths):
            if index == 0:
                page = pages.getByIndex(0)
            else:
                pages.insertNewByIndex(index)
                page = pages.getByIndex(index)
            page.Width = slide_width
            page.Height = slide_height
            try:
                page.Layout = 20  # blank
            except Exception:
                pass
            while page.getCount():
                page.remove(page.getByIndex(0))

            graphic_url = uno.systemPathToFileUrl(str(slide_path.resolve()))
            graphic = provider.queryGraphic((uno_property("URL", graphic_url),))
            shape = document.createInstance("com.sun.star.drawing.GraphicObjectShape")
            shape.Position = Point(0, 0)
            shape.Size = Size(slide_width, slide_height)
            shape.Graphic = graphic
            page.add(shape)

        pptx_path.unlink(missing_ok=True)
        pdf_path.unlink(missing_ok=True)
        document.storeAsURL(
            uno.systemPathToFileUrl(str(pptx_path.resolve())),
            (
                uno_property("FilterName", "Impress MS PowerPoint 2007 XML"),
                uno_property("Overwrite", True),
            ),
        )
        document.storeToURL(
            uno.systemPathToFileUrl(str(pdf_path.resolve())),
            (
                uno_property("FilterName", "impress_pdf_Export"),
                uno_property("Overwrite", True),
            ),
        )
    finally:
        if document is not None:
            try:
                document.close(True)
            except Exception:
                pass
        process.terminate()
        try:
            process.wait(timeout=10)
        except subprocess.TimeoutExpired:
            process.kill()
        shutil.rmtree(profile_dir, ignore_errors=True)

    print(f"created {pptx_path.relative_to(ROOT)}")
    print(f"created {pdf_path.relative_to(ROOT)}")
    return pptx_path, pdf_path


def main() -> None:
    required = [BACKGROUND, WORKOUT_DESKTOP, WORKOUT_MOBILE, AUTH_MOBILE, LEADERBOARD_DESKTOP]
    missing = [str(path) for path in required if not path.exists()]
    if missing:
        raise SystemExit("Missing assets:\n" + "\n".join(missing))
    slide_paths = generate_slide_pngs()
    build_office_files(slide_paths)


if __name__ == "__main__":
    main()
