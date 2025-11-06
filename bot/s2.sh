#!/usr/bin/env bash
set -euo pipefail

# s2.sh — Apaga todo (stop clean + sweep opcional)

# === Paths (mismos que s1.sh) ===
ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$ROOT_DIR/bin"
RUN_DIR="$ROOT_DIR/.run"
LOG_DIR="$ROOT_DIR/logs"

SERVER_NAME="whserver"
BOT_NAME="whbot"

SERVER_BIN="$BIN_DIR/$SERVER_NAME"
BOT_BIN="$BIN_DIR/$BOT_NAME"

SERVER_PID="$RUN_DIR/$SERVER_NAME.pid"
BOT_PID="$RUN_DIR/$BOT_NAME.pid"

SERVER_LOG="$LOG_DIR/$SERVER_NAME.log"
BOT_LOG="$LOG_DIR/$BOT_NAME.log"

# Puerto del server para sweep por socket (override: export WHSERVER_PORT=XXXX)
WHSERVER_PORT="${WHSERVER_PORT:-9000}"

# Flags
FORCE_KILL=0   # --force: mata rezagos con pkill/puerto
PURGE_FILES=0  # --purge: limpia pidfiles (y opcionalmente logs)

# === Helpers ===
usage() {
  cat <<EOF
Uso: $0 [--force] [--purge]
  --force  : intenta matar procesos rezagados por nombre/puerto (:${WHSERVER_PORT})
  --purge  : elimina pidfiles tras detener; no borra logs por defecto
EOF
}

ensure_dirs() {
  mkdir -p "$BIN_DIR" "$RUN_DIR" "$LOG_DIR"
}

is_running() {
  # $1 = pidfile, $2 = expected binary name (whserver|whbot)
  [[ -f "$1" ]] || return 1
  local pid; pid="$(cat "$1" 2>/dev/null || true)"
  [[ -n "${pid:-}" ]] || return 1
  if ps -p "$pid" -o cmd= >/dev/null 2>&1; then
    ps -p "$pid" -o cmd= | grep -q "$2" && return 0 || return 1
  fi
  return 1
}

stop_one() {
  # $1 = pidfile, $2 = name (también esperado por is_running)
  if is_running "$1" "$2"; then
    local pid; pid="$(cat "$1")"
    echo "→ Parando $2 (pid=$pid)…"
    kill -TERM "$pid" 2>/dev/null || true
    # espera suave hasta 6s
    for _ in {1..12}; do
      sleep 0.5
      if ! ps -p "$pid" >/dev/null 2>&1; then
        break
      fi
    done
    if ps -p "$pid" >/dev/null 2>&1; then
      echo "⚠ $2 no cerró a tiempo; enviando SIGKILL"
      kill -KILL "$pid" 2>/dev/null || true
    fi
    echo "✔ $2 detenido."
  else
    echo "ℹ $2 no estaba corriendo."
  fi
}

sweep_leftovers() {
  # Mata por nombre exacto del bin si quedara algo huérfano
  echo "→ Sweep de rezagos (--force)…"
  # pkill por ruta/nombre (silencioso si no existe)
  pkill -f "$SERVER_BIN" 2>/dev/null || true
  pkill -f "$BOT_BIN" 2>/dev/null || true
  # Barrido por nombre, por si se ejecutó fuera de BIN_DIR
  pkill -x "$SERVER_NAME" 2>/dev/null || true
  pkill -x "$BOT_NAME" 2>/dev/null || true
  # Barrido por puerto del server (si disponible: ss o lsof)
  if command -v ss >/dev/null 2>&1; then
    mapfile -t pids < <(ss -lptn 2>/dev/null | awk -v p=":${WHSERVER_PORT}" '$4 ~ p {print $0}' | sed -n 's/.*pid=\([0-9]\+\).*/\1/p' | sort -u)
    for p in "${pids[@]:-}"; do
      [[ -n "$p" ]] && kill -KILL "$p" 2>/dev/null || true
    done
  elif command -v lsof >/dev/null 2>&1; then
    mapfile -t pids < <(lsof -iTCP:"$WHSERVER_PORT" -sTCP:LISTEN -t 2>/dev/null | sort -u)
    for p in "${pids[@]:-}"; do
      [[ -n "$p" ]] && kill -KILL "$p" 2>/dev/null || true
    done
  fi
  echo "✔ Sweep completado."
}

purge_files() {
  echo "→ Limpiando pidfiles…"
  rm -f "$SERVER_PID" "$BOT_PID" 2>/dev/null || true
  echo "✔ PIDs eliminados."
}

status_all() {
  if is_running "$SERVER_PID" "$SERVER_NAME"; then
    echo "✖ $SERVER_NAME aún corre (pid=$(cat "$SERVER_PID"))"
  else
    echo "✔ $SERVER_NAME detenido"
  fi
  if is_running "$BOT_PID" "$BOT_NAME"; then
    echo "✖ $BOT_NAME aún corre (pid=$(cat "$BOT_PID"))"
  else
    echo "✔ $BOT_NAME detenido"
  fi
}

# === Parse args ===
while [[ $# -gt 0 ]]; do
  case "${1}" in
    --force) FORCE_KILL=1; shift ;;
    --purge) PURGE_FILES=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Arg desconocido: $1"; usage; exit 1 ;;
  esac
done

# === Main ===
ensure_dirs
stop_one "$BOT_PID" "$BOT_NAME"
stop_one "$SERVER_PID" "$SERVER_NAME"

if (( FORCE_KILL == 1 )); then
  sweep_leftovers
fi

if (( PURGE_FILES == 1 )); then
  purge_files
fi

status_all
echo "✅ Shutdown completo."
