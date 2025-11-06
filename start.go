package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	BackendPort  = 3000
	FrontendPort = 5173
)

func main() {
	fmt.Println("🚀 Iniciando BOB Chatbot - Sistema Multiagente")
	fmt.Println("================================================")

	// Verificar requisitos
	checkGo()
	checkNode()

	// Limpiar puertos
	cleanupPorts()

	// Verificar configuración
	checkEnv()

	// Instalar dependencias (opcional)
	if askYesNo("¿Instalar/actualizar dependencias? (y/N): ") {
		installDependencies()
	}

	// Iniciar servicios
	backendCmd := startBackend()
	frontendCmd := startFrontend()

	// Verificar que todo esté corriendo
	verifyServices()

	// Mostrar información
	showInfo()

	// Configurar manejo de señales
	setupSignalHandler(backendCmd, frontendCmd)

	// Mantener el programa corriendo
	fmt.Println("⏳ Sistema corriendo... (Presiona Ctrl+C para detener)")
	select {} // Bloquear indefinidamente
}

func checkGo() {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Go no está instalado. Instálalo desde https://go.dev/dl/")
		os.Exit(1)
	}
	fmt.Printf("✅ Go detectado: %s", string(output))
}

func checkNode() {
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Node.js no está instalado. Instálalo desde https://nodejs.org/")
		os.Exit(1)
	}
	fmt.Printf("✅ Node.js detectado: %s", string(output))

	cmd = exec.Command("npm", "--version")
	output, err = cmd.Output()
	if err == nil {
		fmt.Printf("✅ npm detectado: %s", string(output))
	}
}

func cleanupPorts() {
	fmt.Println()
	fmt.Println("🧹 Limpiando puertos anteriores...")
	killPort(BackendPort)
	killPort(FrontendPort)
	fmt.Println("✅ Puertos liberados")
}

func killPort(port int) {
	switch runtime.GOOS {
	case "windows":
		killPortWindows(port)
	case "darwin", "linux":
		killPortUnix(port)
	default:
		fmt.Printf("⚠️  Sistema operativo no soportado: %s\n", runtime.GOOS)
	}
}

func killPortWindows(port int) {
	// Encontrar el PID usando netstat
	cmd := exec.Command("netstat", "-ano", "-p", "TCP")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(":%d", port)) && strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pid := fields[len(fields)-1]
				exec.Command("taskkill", "/F", "/PID", pid).Run()
			}
		}
	}
}

func killPortUnix(port int) {
	// Usar lsof para encontrar procesos
	cmd := exec.Command("lsof", "-ti:"+strconv.Itoa(port))
	output, err := cmd.Output()
	if err != nil {
		return
	}

	pids := strings.TrimSpace(string(output))
	if pids != "" {
		for _, pid := range strings.Split(pids, "\n") {
			exec.Command("kill", "-9", pid).Run()
		}
	}
}

func checkEnv() {
	fmt.Println()
	fmt.Println("🔍 Verificando archivos de configuración...")

	// Verificar backend/.env
	if _, err := os.Stat("backend/.env"); os.IsNotExist(err) {
		fmt.Println("❌ backend/.env no existe")
		fmt.Println("💡 Copia backend/.env.example a backend/.env y configura tu API key")
		os.Exit(1)
	}

	// Verificar/crear frontend/.env
	if _, err := os.Stat("frontend/.env"); os.IsNotExist(err) {
		fmt.Println("⚠️  frontend/.env no existe, creando desde .env.example...")
		input, err := os.ReadFile("frontend/.env.example")
		if err == nil {
			os.WriteFile("frontend/.env", input, 0644)
		} else {
			fmt.Println("ℹ️  No hay .env.example en frontend")
		}
	}

	fmt.Println("✅ Archivos de configuración OK")
}

func installDependencies() {
	fmt.Println()
	fmt.Println("📦 Instalando dependencias...")

	// Backend Go
	fmt.Println("📦 Backend Go dependencies...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = "backend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Frontend React
	if _, err := os.Stat("frontend/package.json"); err == nil {
		fmt.Println("📦 Frontend React dependencies...")
		cmd := exec.Command("npm", "install", "--silent")
		cmd.Dir = "frontend"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	fmt.Println("✅ Dependencias instaladas")
}

func startBackend() *exec.Cmd {
	fmt.Println()
	fmt.Println("🔧 Iniciando Backend Go (puerto 3000)...")

	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Dir = "backend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("❌ Error iniciando backend: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Backend iniciado (PID: %d)\n", cmd.Process.Pid)
	time.Sleep(3 * time.Second)

	return cmd
}

func startFrontend() *exec.Cmd {
	if _, err := os.Stat("frontend/package.json"); os.IsNotExist(err) {
		fmt.Println("ℹ️  No hay frontend configurado, solo backend corriendo")
		return nil
	}

	fmt.Println()
	fmt.Println("🎨 Iniciando Frontend React (puerto 5173)...")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "npm", "run", "dev")
	} else {
		cmd = exec.Command("npm", "run", "dev")
	}

	cmd.Dir = "frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("❌ Error iniciando frontend: %v\n", err)
		return nil
	}

	fmt.Printf("✅ Frontend iniciado (PID: %d)\n", cmd.Process.Pid)
	time.Sleep(2 * time.Second)

	return cmd
}

func verifyServices() {
	fmt.Println()
	fmt.Println("🔍 Verificando servicios...")

	time.Sleep(2 * time.Second)

	// Verificar backend
	resp, err := http.Get("http://localhost:3000/health")
	if err == nil && resp.StatusCode == 200 {
		fmt.Println("✅ Backend OK: http://localhost:3000")
		resp.Body.Close()
	} else {
		fmt.Println("❌ Backend no responde en http://localhost:3000")
		fmt.Println("⚠️  Verifica los logs arriba")
	}

	// Verificar frontend
	if _, err := os.Stat("frontend/package.json"); err == nil {
		resp, err := http.Get("http://localhost:5173")
		if err == nil {
			fmt.Println("✅ Frontend OK: http://localhost:5173")
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		} else {
			fmt.Println("⚠️  Frontend puede tardar unos segundos en iniciar...")
		}
	}
}

func showInfo() {
	fmt.Println()
	fmt.Println("================================================")
	fmt.Println("✅ SISTEMA BOB CHATBOT INICIADO")
	fmt.Println("================================================")
	fmt.Println()
	fmt.Println("📍 Endpoints disponibles:")
	fmt.Println("   Backend:  http://localhost:3000")
	fmt.Println("   Health:   http://localhost:3000/health")
	fmt.Println("   API Docs: http://localhost:3000/")

	if _, err := os.Stat("frontend/package.json"); err == nil {
		fmt.Println("   Frontend: http://localhost:5173")
	}

	fmt.Println()
	fmt.Println("📊 Sistema Multiagente:")
	fmt.Println("   Orchestrator: ✅ Activo (spam, routing, intención)")
	fmt.Println("   FAQ Agent:    ✅ Activo (preguntas frecuentes)")
	fmt.Println("   Auction Agent:✅ Activo (búsqueda vehículos)")
	fmt.Println("   Scoring Agent:✅ Activo (7 dimensiones, 0-100 pts)")
	fmt.Println()
	fmt.Println("🔌 Endpoint principal para WhatsApp:")
	fmt.Println("   POST http://localhost:3000/api/chat/message")
	fmt.Println()
	fmt.Println("📝 Para detener todo:")
	fmt.Println("   Ctrl+C en esta terminal o ejecuta: go run stop.go")
	fmt.Println()
	fmt.Println("📋 Logs:")
	fmt.Println("   Los logs aparecerán debajo de este mensaje...")
	fmt.Println("================================================")
}

func setupSignalHandler(backendCmd, frontendCmd *exec.Cmd) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println()
		fmt.Println("🛑 Deteniendo servicios...")

		// Detener procesos
		if backendCmd != nil && backendCmd.Process != nil {
			backendCmd.Process.Kill()
		}
		if frontendCmd != nil && frontendCmd.Process != nil {
			frontendCmd.Process.Kill()
		}

		// Limpiar puertos
		killPort(BackendPort)
		killPort(FrontendPort)

		fmt.Println("✅ Servicios detenidos")
		os.Exit(0)
	}()
}

func askYesNo(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
