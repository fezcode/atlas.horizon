# atlas.horizon 🛰️

![Banner](banner-image.png)

**atlas.horizon** is a high-fidelity environmental and weather dashboard for the **Atlas Suite**. It provides real-time atmospheric monitoring, including temperature, humidity, wind patterns, and UV index—which we track as simulated "Radiation Levels" for the wasteland survivor.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

## 🔧 Features
- **Real-time Weather:** Powered by Open-Meteo for high accuracy.
- **Rad-Meter (UV Index):** Tracks UV exposure with a "caution" threshold.
- **Wind Radar:** Aesthetic ASCII wind direction and speed monitoring.
- **Location Auto-Detection:** Automatically finds your coordinates via IP.
- **Retro-Future TUI:** Polished Onyx & Gold aesthetic with animations.

## 🚀 Installation
Ensure you have the [Atlas Hub](https://github.com/fezcode/atlas.hub) installed, then:
```bash
atlas.hub install atlas.horizon
```

Or build from source:
```bash
gobake build
```

## ⌨️ Controls
- **'r':** Refresh data manually.
- **'q':** Exit dashboard.
- **'h':** Toggle detailed atmospheric stats.

---
Built with ❤️ by [fezcode](https://github.com/fezcode)
