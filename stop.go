package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	BackendPort  = 3000
	FrontendPort = 5173
)

func main() {
	fmt.Println("🛑 Deteniendo BOB Chatbot - Sistema Multiagente")
	fmt.Println("================================================")

	// Detener backend (puerto 3000)
	fmt.Println("🔧 Deteniendo Backend (puerto 3000)...")
	if killPort(BackendPort) {
		fmt.Println("✅ Backend detenido")
	} else {
		fmt.Println("ℹ️  Backend no estaba corriendo")
	}

	// Detener frontend (puerto 5173)
	fmt.Println("🎨 Deteniendo Frontend (puerto 5173)...")
	if killPort(FrontendPort) {
		fmt.Println("✅ Frontend detenido")
	} else {
		fmt.Println("ℹ️  Frontend no estaba corriendo")
	}

	// Limpiar procesos específicos por nombre
	killProcessByName()

	fmt.Println()
	fmt.Println("✅ Todos los servicios detenidos")
	fmt.Println("================================================")
}

func killPort(port int) bool {
	switch runtime.GOOS {
	case "windows":
		return killPortWindows(port)
	case "darwin", "linux":
		return killPortUnix(port)
	default:
		fmt.Printf("⚠️  Sistema operativo no soportado: %s\n", runtime.GOOS)
		return false
	}
}

func killPortWindows(port int) bool {
	// Encontrar el PID usando netstat
	cmd := exec.Command("netstat", "-ano", "-p", "TCP")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	killed := false
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(":%d", port)) && strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pid := fields[len(fields)-1]
				err := exec.Command("taskkill", "/F", "/PID", pid).Run()
				if err == nil {
					killed = true
				}
			}
		}
	}
	return killed
}

func killPortUnix(port int) bool {
	// Usar lsof para encontrar procesos
	cmd := exec.Command("lsof", "-ti:"+strconv.Itoa(port))
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	pids := strings.TrimSpace(string(output))
	if pids == "" {
		return false
	}

	killed := false
	for _, pid := range strings.Split(pids, "\n") {
		err := exec.Command("kill", "-9", pid).Run()
		if err == nil {
			killed = true
		}
	}
	return killed
}

func killProcessByName() {
	switch runtime.GOOS {
	case "windows":
		killProcessWindows()
	case "darwin", "linux":
		killProcessUnix()
	}
}

func killProcessWindows() {
	// Matar procesos de Go
	exec.Command("taskkill", "/F", "/IM", "go.exe", "/FI", "WINDOWTITLE eq cmd/server/main.go").Run()

	// Matar procesos de Node/Vite
	exec.Command("taskkill", "/F", "/IM", "node.exe", "/FI", "WINDOWTITLE eq vite").Run()
}

func killProcessUnix() {
	// Matar procesos de Go
	exec.Command("pkill", "-f", "go run cmd/server/main.go").Run()

	// Matar procesos de Vite y npm
	exec.Command("pkill", "-f", "vite").Run()
	exec.Command("pkill", "-f", "npm run dev").Run()
}
