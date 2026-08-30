# Материалы защиты Limit less

- `limit-less-hackathon-10min.pptx` — презентация 16:9 для PowerPoint/LibreOffice.
- `limit-less-hackathon-10min.pdf` — офлайн-резерв с тем же визуальным результатом.
- `pitch-script.md` — готовый текст выступления и хронометраж.
- `demo-checklist.md` — подготовка живого демо, известные риски и план B.
- `generate_deck.py` — воспроизводимый генератор слайдов.
- `build/slides/*.png` — полноразмерные 1920×1080 исходники слайдов.

Пересборка:

```bash
python3 presentation/generate_deck.py
```

Генератор использует Pillow, PyUNO и установленный LibreOffice. Слайды в PPTX намеренно вставляются как полноэкранные изображения: это сохраняет типографику и композицию одинаковыми на разных компьютерах. Текст выступления редактируется отдельно в `pitch-script.md`.

